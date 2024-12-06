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
	"github.com/stephenafamo/bob/dialect/psql/dm"
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
	inserted, err := dbmodels.Parkingspots.Insert(&spotSetter).One(ctx, tx)
	if err != nil {
		// Handle duplicate error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ExclusionViolation {
			err = ErrDuplicatedAddress
		}
		return Entry{}, nil, err
	}

	err = inserted.InsertParkingspotidTimeunits(ctx, tx, unitSetters...)
	if err != nil {
		if dbmodels.ErrUniqueConstraint.Is(err) {
			err = ErrDuplicatedTimeUnit
		}
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

func (p *PostgresRepository) UpdateSpotByUUID(ctx context.Context, spotID uuid.UUID, updateSpot *models.ParkingSpotUpdateInput) (Entry, error) {
	spotSetter, err := spotSetterFromUpdateInput(updateSpot)
	if err != nil {
		return Entry{}, err
	}

	// Update the spot details
	updated, err := dbmodels.Parkingspots.Update(
		dbmodels.UpdateWhere.Parkingspots.Parkingspotuuid.EQ(spotID),
		spotSetter.UpdateMod(),
		um.Returning(dbmodels.Parkingspots.Columns()),
	).One(ctx, p.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Entry{}, ErrNotFound
		}
		return Entry{}, fmt.Errorf("could not execute update: %w", err)
	}

	entry, err := entryFromDB(updated)
	if err != nil {
		return Entry{}, fmt.Errorf("could not adapt dbmodels.Parkingspot: %w", err)
	}

	return entry, nil
}

func (p *PostgresRepository) UpdateAvailByUUID(ctx context.Context, spotID uuid.UUID, updateTimes *models.ParkingSpotAvailUpdateInput) error {
	entry, err := p.GetByUUID(ctx, spotID)
	if err != nil {
		return err
	}

	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback if commit is not done

	err = removeAvailability(ctx, tx, entry.InternalID, updateTimes.RemoveAvailability)
	if err != nil {
		return err
	}

	if len(updateTimes.AddAvailability) > 0 {
		// Insert new availability times
		newTimesSetter := timeSetterFromInput(updateTimes.AddAvailability, entry.InternalID)
		_, err = dbmodels.Timeunits.Insert(bob.ToMods(newTimesSetter...)).Exec(ctx, tx)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					return ErrDuplicatedTimeUnit
				}
			}
			return err
		}
	}

	// Commit since audit passed
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func removeAvailability(ctx context.Context, tx bob.Tx, spotID int64, remove []models.TimeUnit) error {
	if len(remove) == 0 {
		return nil
	}

	timeslots := timeSlotsToSQLExpr(remove)

	var whereMods []mods.Where[*dialect.DeleteQuery]
	whereMods = append(
		whereMods,
		dbmodels.DeleteWhere.Timeunits.Parkingspotid.EQ(spotID),
		dbmodels.DeleteWhere.Timeunits.Bookingid.IsNull(),
		dm.Where(timeslots),
	)

	// Delete unbooked times
	deleted, err := dbmodels.Timeunits.Delete(
		psql.WhereAnd(whereMods...),
	).Exec(ctx, tx)
	if err != nil {
		return err
	}
	// Audit failure must mean update availability input is attempting to removed a booked time unit
	if deleted != int64(len(remove)) {
		return ErrDeleteBookedTimeUnit
	}
	return nil
}

func timeSlotsToSQLExpr(units []models.TimeUnit) dialect.Expression {
	var expression dialect.Expression
	for _, bookTime := range units {
		test := dbmodels.TimeunitColumns.Timerange.OP(
			"&&",
			psql.Arg(dbtype.Tstzrange{
				Start: bookTime.StartTime,
				End:   bookTime.EndTime,
			}),
		)
		if expression.Base == nil {
			expression = test
		} else {
			expression = expression.Or(test)
		}
	}
	return expression
}

func (p *PostgresRepository) GetByUUID(ctx context.Context, spotID uuid.UUID) (Entry, error) {
	spotResult, err := dbmodels.Parkingspots.Query(
		dbmodels.SelectWhere.Parkingspots.Parkingspotuuid.EQ(spotID),
	).One(ctx, p.db)
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
		sm.Columns(dbmodels.TimeunitColumns.Timerange),
		sm.Columns(dbmodels.TimeunitColumns.Bookingid),
		psql.WhereAnd(
			dbmodels.SelectWhere.Parkingspots.Parkingspotuuid.EQ(spotID),
			sm.Where(dbmodels.TimeunitColumns.Timerange.OP("&&", psql.Arg(dbtype.Tstzrange{
				Start: startDate,
				End:   endDate,
			}))),
		),
		dbmodels.SelectJoins.Timeunits.InnerJoin.ParkingspotidParkingspot(ctx),
		sm.OrderBy(psql.F("lower", dbmodels.TimeunitColumns.Timerange)),
	).All(ctx, p.db)
	if err != nil {
		return nil, err
	}

	// If no rows found
	if len(result) == 0 {
		// Ignore errors here, just treat it as not existing
		exists, _ := dbmodels.Parkingspots.Query(
			sm.Columns(1),
			dbmodels.SelectWhere.Parkingspots.Parkingspotuuid.EQ(spotID),
		).Exists(ctx, p.db)

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
		sm.Columns(
			dbmodels.ParkingspotColumns.Userid,
		),
		dbmodels.SelectWhere.Parkingspots.Parkingspotuuid.EQ(spotID),
	).One(ctx, p.db)
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
		sm.From(dbmodels.Parkingspots.Name()),
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
		var status string
		if _, ok := unit.Bookingid.Get(); ok {
			status = "booked"
		} else {
			status = "available"
		}

		result = append(result, models.TimeUnit{
			StartTime: unit.Timerange.Start,
			EndTime:   unit.Timerange.End,
			Status:    status,
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
