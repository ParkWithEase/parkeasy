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
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/fm"
	"github.com/stephenafamo/bob/dialect/psql/sm"
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

type result struct {
	dbmodels.Parkingspot
	DistanceToOrigin float64 `db:"distance_to_origin"`
}

func (p *PostgresRepository) Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (Entry, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return Entry{}, fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback if commit is not done

	longitude := spot.Location.Longitude
	latitude := spot.Location.Latitude

	inserted, err := dbmodels.Parkingspots.Insert(ctx, tx, &dbmodels.ParkingspotSetter{
		Userid:             omit.From(userID),
		Postalcode:         omit.From(spot.Location.PostalCode),
		Countrycode:        omit.From(spot.Location.CountryCode),
		City:               omit.From(spot.Location.City),
		State:              omit.From(spot.Location.State),
		Streetaddress:      omit.From(spot.Location.StreetAddress),
		Longitude:          omit.From(longitude),
		Latitude:           omit.From(latitude),
		Hasshelter:         omit.From(spot.Features.Shelter),
		Hasplugin:          omit.From(spot.Features.PlugIn),
		Haschargingstation: omit.From(spot.Features.ChargingStation),
		Priceperhour:       omit.From(spot.PricePerHour),
	})
	if err != nil {
		// Handle duplicate error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				err = ErrDuplicatedAddress
			}
		}
		return Entry{}, err
	}

	availabilitySetters := make([]*dbmodels.TimeunitSetter, 0, len(spot.Availability))
	for _, timeslot := range spot.Availability {
		availabilitySetters = append(availabilitySetters, &dbmodels.TimeunitSetter{
			Timerange: omit.From(dbtype.Tstzrange{
				Start: timeslot.StartTime,
				End:   timeslot.EndTime,
			}),
			Parkingspotid: omit.From(inserted.Parkingspotid),
		})
	}

	units, err := dbmodels.Timeunits.InsertMany(ctx, tx, availabilitySetters...)
	availability := make([]models.TimeUnit, 0, len(units))

	for _, unit := range units {
		availability = append(availability, models.TimeUnit{
			StartTime: unit.Timerange.Start,
			EndTime:   unit.Timerange.End,
		})
	}

	err = tx.Commit()
	if err != nil {
		return Entry{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	location := models.ParkingSpotLocation{
		PostalCode:    inserted.Postalcode,
		CountryCode:   inserted.Countrycode,
		City:          inserted.City,
		State:         inserted.State,
		StreetAddress: inserted.Streetaddress,
		Longitude:     inserted.Longitude,
		Latitude:      inserted.Latitude,
	}

	features := models.ParkingSpotFeatures{
		Shelter:         inserted.Hasshelter,
		PlugIn:          inserted.Hasplugin,
		ChargingStation: inserted.Haschargingstation,
	}

	parkingspot := models.ParkingSpot{
		Location:     location,
		Features:     features,
		ID:           inserted.Parkingspotuuid,
		PricePerHour: inserted.Priceperhour,
		Availability: availability,
	}

	entry := Entry{
		ParkingSpot: parkingspot,
		InternalID:  inserted.Parkingspotid,
		OwnerID:     inserted.Userid,
	}
	return entry, nil
}

func (p *PostgresRepository) GetByUUID(ctx context.Context, spotID uuid.UUID) (Entry, error) {
	spotResult, err := dbmodels.Parkingspots.Query(
		ctx, p.db,
		sm.Columns(
			dbmodels.ParkingspotColumns.Postalcode,
			dbmodels.ParkingspotColumns.Countrycode,
			dbmodels.ParkingspotColumns.City,
			dbmodels.ParkingspotColumns.State,
			dbmodels.ParkingspotColumns.Streetaddress,
			dbmodels.ParkingspotColumns.Longitude,
			dbmodels.ParkingspotColumns.Latitude,
			dbmodels.ParkingspotColumns.Hasshelter,
			dbmodels.ParkingspotColumns.Hasplugin,
			dbmodels.ParkingspotColumns.Haschargingstation,
			dbmodels.ParkingspotColumns.Parkingspotid,
			dbmodels.ParkingspotColumns.Userid,
			dbmodels.ParkingspotColumns.Priceperhour,
		),
		dbmodels.SelectWhere.Parkingspots.Parkingspotuuid.EQ(spotID),
	).One()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return Entry{}, err
	}

	location := models.ParkingSpotLocation{
		PostalCode:    spotResult.Postalcode,
		CountryCode:   spotResult.Countrycode,
		City:          spotResult.City,
		State:         spotResult.State,
		StreetAddress: spotResult.Streetaddress,
		Longitude:     spotResult.Longitude,
		Latitude:      spotResult.Latitude,
	}

	features := models.ParkingSpotFeatures{
		Shelter:         spotResult.Hasshelter,
		PlugIn:          spotResult.Hasplugin,
		ChargingStation: spotResult.Haschargingstation,
	}

	parkingspot := models.ParkingSpot{
		Location:     location,
		Features:     features,
		ID:           spotID,
		PricePerHour: spotResult.Priceperhour,
	}

	return Entry{
		ParkingSpot: parkingspot,
		InternalID:  spotResult.Parkingspotid,
		OwnerID:     spotResult.Userid,
	}, nil
}

func (p *PostgresRepository) GetAvalByUUID(ctx context.Context, spotID uuid.UUID, startDate time.Time, endDate time.Time) ([]models.TimeUnit, error) {
	timeUnitResult, err := dbmodels.Timeunits.Query(
		ctx, p.db,
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
	).All()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrTimeUnitNotFound
		}
		return []models.TimeUnit{}, err
	}

	// Handle no rows found
	if len(timeUnitResult) == 0 {
		return []models.TimeUnit{}, ErrTimeUnitNotFound
	}

	availability := make([]models.TimeUnit, 0, len(timeUnitResult)) // Initialize slice

	for _, timeslot := range timeUnitResult {
		var status string
		_, err := timeslot.Bookingid.Value()

		if err != nil {
			status = "booked"
		} else {
			status = "available"
		}

		availability = append(availability, models.TimeUnit{
			StartTime: timeslot.Timerange.Start,
			EndTime:   timeslot.Timerange.End,
			Status:    status,
		})
	}

	return availability, nil
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

func (p *PostgresRepository) GetMany(ctx context.Context, limit int, filter Filter) ([]GetManyEntry, error) {
	smods := []bob.Mod[*dialect.SelectQuery]{sm.Columns(dbmodels.Parkingspots.Columns())}
	var where mods.Where[*dialect.SelectQuery]

	if userID, ok := filter.UserID.Get(); ok {
		constraint := dbmodels.SelectWhere.Parkingspots.Userid.EQ(userID)
		if where.E != nil {
			where = psql.WhereAnd(
				where,
				constraint,
			)
		} else {
			where = constraint
		}
	}

	if locFilter, ok := filter.Location.Get(); ok {
		centre := psql.F("ll_to_earth", psql.Arg(locFilter.Latitude), psql.Arg(locFilter.Longitude))
		spotPosition := psql.F("ll_to_earth", dbmodels.ParkingspotColumns.Latitude, dbmodels.ParkingspotColumns.Longitude)
		constraint := sm.Where(
			psql.F(
				"earth_box",
				centre,
				psql.Arg(locFilter.Radius),
			)().OP(
				"@>",
				spotPosition,
			),
		)
		if where.E != nil {
			where = psql.WhereAnd(where, constraint)
		} else {
			where = constraint
		}
		smods = append(
			smods,
			sm.Columns(psql.F("earth_distance", centre, spotPosition)(fm.As("distance_to_origin"))),
			sm.OrderBy("distance_to_origin").Asc(),
		)
	}

	if availFilter, ok := filter.Availability.Get(); ok {
		constraint := sm.Where(
			dbmodels.TimeunitColumns.Timerange.OP(
				"&&",
				psql.Arg(dbtype.Tstzrange{
					Start: availFilter.Start,
					End:   availFilter.End,
				}),
			),
		)
		if where.E != nil {
			where = psql.WhereAnd(where, constraint)
		} else {
			where = constraint
		}
		smods = append(
			smods,
			dbmodels.SelectJoins.Parkingspots.InnerJoin.ParkingspotidTimeunits(ctx),
		)
	}

	if where.E == nil {
		return nil, ErrNoConstraint
	}

	smods = append(
		smods,
		sm.From(dbmodels.Parkingspots.Name(ctx)),
		sm.Limit(limit),
		where,
	)
	query := psql.Select(smods...)

	entryCursor, err := bob.Cursor(ctx, p.db, query, scan.StructMapper[result]())
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

		result = append(result, r.ToEntry())
	}
	return result, nil
}

func (r *result) ToEntry() GetManyEntry {
	return GetManyEntry{
		Entry: Entry{
			ParkingSpot: models.ParkingSpot{
				Location: models.ParkingSpotLocation{
					PostalCode:    r.Postalcode,
					CountryCode:   r.Countrycode,
					City:          r.City,
					State:         r.State,
					StreetAddress: r.Streetaddress,
					Longitude:     r.Longitude,
					Latitude:      r.Latitude,
				},
				Features: models.ParkingSpotFeatures{
					Shelter:         r.Hasshelter,
					PlugIn:          r.Hasplugin,
					ChargingStation: r.Haschargingstation,
				},
				PricePerHour: r.Priceperhour,
				ID:           r.Parkingspotuuid,
			},
			InternalID: r.Parkingspotid,
			OwnerID:    r.Userid,
		},
		DistanceToLocation: r.DistanceToOrigin,
	}
}
