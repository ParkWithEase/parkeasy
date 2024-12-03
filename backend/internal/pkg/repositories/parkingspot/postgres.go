package parkingspot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmodels"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbtype"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/fm"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/bob/mods"
	"github.com/stephenafamo/scan"
)

type PostgresRepository struct {
	db bob.DB
}

func NewPostgres(db bob.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

type getManyResult struct {
	dbmodels.Parkingspot
	DistanceToOrigin float64 `db:"distance_to_origin"`
}

func (p *PostgresRepository) Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (Entry, []models.TimeUnit, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return Entry{}, nil, fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback if commit is not done

	spotSetter, unitSetters, err := settersFromCreateInput(userID, spot)
	if err != nil {
		return Entry{}, nil, err
	}
	inserted, err := dbmodels.Parkingspots.Insert(ctx, tx, &spotSetter)
	if err != nil {
		// Handle duplicate error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.ExclusionViolation && pgErr.ConstraintName == "latlon_overlap_exclude" {
				err = ErrDuplicatedAddress
			}
		}
		return Entry{}, nil, err
	}

	err = inserted.InsertParkingspotidTimeunits(ctx, tx, unitSetters...)
	if err != nil {
		return Entry{}, nil, err
	}

	err = tx.Commit()
	if err != nil {
		return Entry{}, nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	entry, err := entryFromDB(inserted)
	if err != nil {
		return Entry{}, nil, fmt.Errorf("could not adapt dbmodels.Parkingspot: %w", err)
	}
	availability := timeUnitsFromDB(inserted.R.ParkingspotidTimeunits)

	return entry, availability, nil
}

func (p *PostgresRepository) UpdateByUUID(ctx context.Context, spotID uuid.UUID, updateSpot *models.ParkingSpotUpdateInput) (Entry, []models.TimeUnit, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return Entry{}, nil, fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback if commit is not done

	spotSetter, err := spotSetterFromUpdateInput(updateSpot)
	if err != nil {
		return Entry{}, nil, err
	}

	// Update the spot details
	updated, err := dbmodels.Parkingspots.UpdateQ(
		ctx, p.db,
		dbmodels.UpdateWhere.Parkingspots.Parkingspotuuid.EQ(spotID),
		spotSetter,
		um.Returning(dbmodels.Parkingspots.Columns()),
	).One()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Entry{}, []models.TimeUnit{}, ErrNotFound
		}
		return Entry{}, []models.TimeUnit{}, fmt.Errorf("could not execute update: %w", err)
	}

	// Delete unbooked times
	_, err = dbmodels.Timeunits.DeleteQ(
		ctx, p.db,
		dbmodels.DeleteWhere.Timeunits.Parkingspotid.EQ(updated.Parkingspotid),
		dbmodels.DeleteWhere.Timeunits.Bookingid.IsNull(),
	).Exec()
	if err != nil {
		return Entry{}, []models.TimeUnit{}, fmt.Errorf("could not execute update: %w", err)
	}

	// Get remaining booked times
	bookedTimes, err := dbmodels.Timeunits.Query(
		ctx, p.db,
		sm.Columns(dbmodels.TimeunitColumns.Timerange),
		psql.WhereAnd(
			dbmodels.SelectWhere.Parkingspots.Parkingspotuuid.EQ(spotID),
		),
	).All()
	if err != nil {
		return Entry{}, []models.TimeUnit{}, fmt.Errorf("could not execute update: %w", err)
	}

	// Determine unbooked times to be inserted
	newTimeUnits := timeUnitSetMinus(updateSpot.Availability, timeUnitsFromDB(bookedTimes))

	// Insert unbooked times
	newTimesSetter := timeSetterFromInput(newTimeUnits, updated.Parkingspotid)
	newTimes := make(dbmodels.TimeunitSlice, 0, 8)
	for _, setter := range newTimesSetter {
		insertedTimes, err := dbmodels.Timeunits.Insert(ctx, p.db, setter)
		if err != nil {
			return Entry{}, []models.TimeUnit{}, fmt.Errorf("could not execute update: %w", err)
		}
		newTimes = append(newTimes, insertedTimes)
	}

	// Get times from DB after update
	bookedTimes, err = dbmodels.Timeunits.Query(
		ctx, p.db,
		sm.Columns(dbmodels.TimeunitColumns.Timerange),
		psql.WhereAnd(
			dbmodels.SelectWhere.Timeunits.Parkingspotid.EQ(updated.Parkingspotid),
		),
	).All()
	if err != nil {
		return Entry{}, []models.TimeUnit{}, fmt.Errorf("could not execute update: %w", err)
	}

	// Compare state of times in DB after update with times from input
	audit := timeUnitSetMinus(updateSpot.Availability, timeUnitsFromDB(bookedTimes))
	if len(audit) != 0 {
		return Entry{}, []models.TimeUnit{}, errors.New("could not execute update: audit failed")
	}

	// Commit since audit passed
	err = tx.Commit()
	if err != nil {
		return Entry{}, nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	entry, err := entryFromDB(updated)
	if err != nil {
		return Entry{}, nil, fmt.Errorf("could not adapt dbmodels.Parkingspot: %w", err)
	}
	availability := timeUnitsFromDB(newTimes)

	return entry, availability, nil
}

func timeUnitSetMinus(firstSet, secondSet []models.TimeUnit) []models.TimeUnit {
	// Create a map to store booked TimeUnits for quick lookup
	bookedSet := make(map[time.Time]map[time.Time]struct{})
	for _, booked := range secondSet {
		start := booked.StartTime.UTC().Truncate(time.Second)
		end := booked.EndTime.UTC().Truncate(time.Second)

		if _, exists := bookedSet[start]; !exists {
			bookedSet[start] = make(map[time.Time]struct{})
		}
		bookedSet[start][end] = struct{}{}
	}

	var result []models.TimeUnit
	for _, update := range firstSet {
		start := update.StartTime.UTC().Truncate(time.Second)
		end := update.EndTime.UTC().Truncate(time.Second)

		if _, startFound := bookedSet[start]; startFound {
			if _, endFound := bookedSet[start][end]; endFound {
				continue
			}
		}
		result = append(result, models.TimeUnit{
			StartTime: start,
			EndTime:   end,
			Status:    update.Status,
		})
	}
	return result
}

func (p *PostgresRepository) GetByUUID(ctx context.Context, spotID uuid.UUID) (Entry, error) {
	spotResult, err := dbmodels.Parkingspots.Query(
		ctx, p.db,
		dbmodels.SelectWhere.Parkingspots.Parkingspotuuid.EQ(spotID),
	).One()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return Entry{}, err
	}

	entry, err := entryFromDB(spotResult)
	if err != nil {
		return Entry{}, fmt.Errorf("could not adapt dbmodels.Parkingspot: %w", err)
	}
	return entry, nil
}

func (p *PostgresRepository) GetAvailByUUID(ctx context.Context, spotID uuid.UUID, startDate, endDate time.Time) ([]models.TimeUnit, error) {
	result, err := dbmodels.Timeunits.Query(
		ctx, p.db,
		sm.Columns(dbmodels.TimeunitColumns.Timerange),
		psql.WhereAnd(
			dbmodels.SelectWhere.Parkingspots.Parkingspotuuid.EQ(spotID),
			sm.Where(dbmodels.TimeunitColumns.Timerange.OP("&&", psql.Arg(dbtype.Tstzrange{
				Start: startDate,
				End:   endDate,
			}))),
		),
		dbmodels.SelectJoins.Timeunits.InnerJoin.ParkingspotidParkingspot(ctx),
		sm.OrderBy(psql.F("lower", dbmodels.TimeunitColumns.Timerange)),
	).All()
	if err != nil {
		return nil, err
	}

	// If no rows found
	if len(result) == 0 {
		// Ignore errors here, just treat it as not existing
		exists, _ := dbmodels.Parkingspots.Query(
			ctx, p.db,
			sm.Columns(1),
			dbmodels.SelectWhere.Parkingspots.Parkingspotuuid.EQ(spotID),
		).Exists()

		if !exists {
			return nil, ErrNotFound
		}

		// No time units is not an error
		return []models.TimeUnit{}, nil
	}

	return timeUnitsFromDB(result), nil
}

func (p *PostgresRepository) GetOwnerByUUID(ctx context.Context, spotID uuid.UUID) (int64, error) {
	result, err := dbmodels.Parkingspots.Query(
		ctx, p.db,
		sm.Columns(
			dbmodels.ParkingspotColumns.Userid,
		),
		dbmodels.SelectWhere.Parkingspots.Parkingspotuuid.EQ(spotID),
	).One()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return -1, err
	}

	return result.Userid, err
}

func (p *PostgresRepository) GetMany(ctx context.Context, limit int, filter *Filter) ([]GetManyEntry, error) {
	log := zerolog.Ctx(ctx).
		With().
		Str("component", "parkingspot.Postgres").
		Logger()

	smods := []bob.Mod[*dialect.SelectQuery]{sm.Columns(dbmodels.Parkingspots.Columns())}
	var whereMods []mods.Where[*dialect.SelectQuery]

	if userID, ok := filter.UserID.Get(); ok {
		whereMods = append(whereMods, dbmodels.SelectWhere.Parkingspots.Userid.EQ(userID))
	}

	if locFilter, ok := filter.Location.Get(); ok {
		centre := psql.F("ll_to_earth", psql.Arg(locFilter.Latitude), psql.Arg(locFilter.Longitude))
		spotPosition := psql.F("ll_to_earth", dbmodels.ParkingspotColumns.Latitude, dbmodels.ParkingspotColumns.Longitude)
		whereMods = append(whereMods, sm.Where(
			psql.F(
				"earth_box",
				centre,
				psql.Arg(locFilter.Radius),
			)().OP(
				"@>",
				spotPosition,
			),
		))
		smods = append(
			smods,
			sm.Columns(psql.F("earth_distance", centre, spotPosition)(fm.As("distance_to_origin"))),
			sm.OrderBy("distance_to_origin").Asc(),
		)
	}

	if availFilter, ok := filter.Availability.Get(); ok {
		whereMods = append(whereMods, sm.Where(
			dbmodels.TimeunitColumns.Timerange.OP(
				"&&",
				psql.Arg(dbtype.Tstzrange{
					Start: availFilter.Start,
					End:   availFilter.End,
				}),
			),
		))
		smods = append(
			smods,
			dbmodels.SelectJoins.Parkingspots.InnerJoin.ParkingspotidTimeunits(ctx),
		)
	}

	if len(whereMods) == 0 {
		return nil, ErrNoConstraint
	}

	smods = append(
		smods,
		sm.From(dbmodels.Parkingspots.Name(ctx)),
		sm.Limit(limit),
		sm.Distinct(),
		psql.WhereAnd(whereMods...),
	)
	query := psql.Select(smods...)

	entryCursor, err := bob.Cursor(ctx, p.db, query, scan.StructMapper[getManyResult]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []GetManyEntry{}, nil
		}
		return nil, err
	}
	defer entryCursor.Close()

	result := make([]GetManyEntry, 0, 8)
	for entryCursor.Next() {
		r, err := entryCursor.Get()
		if err != nil { // if there's an error, just return what we already have
			break
		}

		res, err := r.ToEntry()
		if err != nil { // if there is an error converting lat, long or price to float, then log and skip this entry
			log.Err(err).Msg("error while converting DB entry")
			continue
		}

		result = append(result, res)
	}
	return result, nil
}

func (r *getManyResult) ToEntry() (GetManyEntry, error) {
	entry, err := entryFromDB(&r.Parkingspot)
	if err != nil {
		return GetManyEntry{}, err
	}

	return GetManyEntry{
		Entry:              entry,
		DistanceToLocation: r.DistanceToOrigin,
	}, nil
}

func entryFromDB(model *dbmodels.Parkingspot) (Entry, error) {
	lat, ok := model.Latitude.Float64()
	if !ok {
		return Entry{}, fmt.Errorf("could not convert %v to float64", model.Latitude)
	}
	lon, ok := model.Longitude.Float64()
	if !ok {
		return Entry{}, fmt.Errorf("could not convert %v to float64", model.Longitude)
	}
	price, ok := model.Priceperhour.Float64()
	if !ok {
		return Entry{}, fmt.Errorf("could not convert %v to float64", model.Priceperhour)
	}

	return Entry{
		ParkingSpot: models.ParkingSpot{
			Location: models.ParkingSpotLocation{
				PostalCode:    model.Postalcode,
				CountryCode:   model.Countrycode,
				City:          model.City,
				State:         model.State,
				StreetAddress: model.Streetaddress,
				Longitude:     lon,
				Latitude:      lat,
			},
			Features: models.ParkingSpotFeatures{
				Shelter:         model.Hasshelter,
				PlugIn:          model.Hasplugin,
				ChargingStation: model.Haschargingstation,
			},
			PricePerHour: price,
			ID:           model.Parkingspotuuid,
		},
		InternalID: model.Parkingspotid,
		OwnerID:    model.Userid,
	}, nil
}

func timeUnitsFromDB(model []*dbmodels.Timeunit) []models.TimeUnit {
	result := make([]models.TimeUnit, 0, len(model))
	for _, unit := range model {
		result = append(result, models.TimeUnit{
			StartTime: unit.Timerange.Start,
			EndTime:   unit.Timerange.End,
			Status:    "available",
		})
	}
	return result
}

func settersFromCreateInput(userID int64, input *models.ParkingSpotCreationInput) (dbmodels.ParkingspotSetter, []*dbmodels.TimeunitSetter, error) {
	lon, err := decimal.NewFromFloat64(input.Location.Longitude)
	if err != nil {
		return dbmodels.ParkingspotSetter{}, nil, ErrInvalidCoordinate
	}
	lat, err := decimal.NewFromFloat64(input.Location.Latitude)
	if err != nil {
		return dbmodels.ParkingspotSetter{}, nil, ErrInvalidCoordinate
	}
	price, err := decimal.NewFromFloat64(input.PricePerHour)
	if err != nil {
		return dbmodels.ParkingspotSetter{}, nil, ErrInvalidPrice
	}
	timeunits := make([]*dbmodels.TimeunitSetter, 0, len(input.Availability))
	for _, timeslot := range input.Availability {
		timeunits = append(timeunits, &dbmodels.TimeunitSetter{
			Timerange: omit.From(dbtype.Tstzrange{
				Start: timeslot.StartTime,
				End:   timeslot.EndTime,
			}),
		})
	}

	return dbmodels.ParkingspotSetter{
		Userid:             omit.From(userID),
		Postalcode:         omit.From(input.Location.PostalCode),
		Countrycode:        omit.From(input.Location.CountryCode),
		City:               omit.From(input.Location.City),
		State:              omit.From(input.Location.State),
		Streetaddress:      omit.From(input.Location.StreetAddress),
		Longitude:          omit.From(lon),
		Latitude:           omit.From(lat),
		Hasshelter:         omit.From(input.Features.Shelter),
		Hasplugin:          omit.From(input.Features.PlugIn),
		Haschargingstation: omit.From(input.Features.ChargingStation),
		Priceperhour:       omit.From(price),
	}, timeunits, nil
}

func spotSetterFromUpdateInput(input *models.ParkingSpotUpdateInput) (dbmodels.ParkingspotSetter, error) {
	price, err := decimal.NewFromFloat64(input.PricePerHour)
	if err != nil {
		return dbmodels.ParkingspotSetter{}, ErrInvalidPrice
	}

	return dbmodels.ParkingspotSetter{
		Hasshelter:         omit.From(input.Features.Shelter),
		Hasplugin:          omit.From(input.Features.PlugIn),
		Haschargingstation: omit.From(input.Features.ChargingStation),
		Priceperhour:       omit.From(price),
	}, nil
}

func timeSetterFromInput(input []models.TimeUnit, spotID int64) []*dbmodels.TimeunitSetter {
	timeunits := make([]*dbmodels.TimeunitSetter, 0, len(input))
	for _, timeslot := range input {
		timeunits = append(timeunits, &dbmodels.TimeunitSetter{
			Timerange: omit.From(dbtype.Tstzrange{
				Start: timeslot.StartTime,
				End:   timeslot.EndTime,
			}),
			Parkingspotid: omit.From(spotID),
		})
	}

	return timeunits
}
