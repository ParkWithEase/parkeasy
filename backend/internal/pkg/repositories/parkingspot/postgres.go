package parkingspot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmodels"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/fm"
	"github.com/stephenafamo/bob/dialect/psql/sm"
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

func (p *PostgresRepository) Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (Entry, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return Entry{}, fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback if commit is not done

	longitude := float32(spot.Location.Longitude)
	latitude := float32(spot.Location.Latitude)

	inserted, err := dbmodels.Parkingspots.Insert(ctx, p.db, &dbmodels.ParkingspotSetter{
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

	availability := make([]models.TimeUnit, 0, len(spot.Availability)) // Initialize slice

	for _, timeslot := range spot.Availability {
		// Insert each unit
		inserted, err := dbmodels.Timeunits.Insert(ctx, p.db, &dbmodels.TimeunitSetter{
			Starttime:       omit.From(timeslot.StartTime.UTC()),
			Endtime:         omit.From(timeslot.EndTime.UTC()),
			Parkingspotuuid: omit.From(inserted.Parkingspotuuid),
			Status:          omit.From(timeslot.Status),
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					err = ErrDuplicatedTimeUnit
				}
			}
			return Entry{}, err
		}

		availability = append(availability, models.TimeUnit{
			StartTime: inserted.Starttime.UTC(),
			EndTime:   inserted.Endtime.UTC(),
			Status:    inserted.Status,
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
		StreetAddress: inserted.Streetaddress,
		State:         inserted.State,
		Longitude:     float64(inserted.Longitude),
		Latitude:      float64(inserted.Latitude),
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

func (p *PostgresRepository) GetByUUID(ctx context.Context, spotID uuid.UUID, startDate time.Time, endDate time.Time) (Entry, error) {
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

	//Initialize error variable
	var timeError error
	availability, err := p.GetAvalByUUID(ctx, spotID, startDate, endDate)
	if err != nil && !errors.Is(err, ErrTimeUnitNotFound) {
		timeError = err
	}

	location := models.ParkingSpotLocation{
		PostalCode:    spotResult.Postalcode,
		CountryCode:   spotResult.Countrycode,
		City:          spotResult.City,
		StreetAddress: spotResult.Streetaddress,
		State:         spotResult.State,
		Longitude:     float64(spotResult.Longitude),
		Latitude:      float64(spotResult.Latitude),
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
		Availability: availability,
		PricePerHour: spotResult.Priceperhour,
	}

	return Entry{
		ParkingSpot: parkingspot,
		InternalID:  spotResult.Parkingspotid,
		OwnerID:     spotResult.Userid,
	}, timeError
}

func (p *PostgresRepository) GetAvalByUUID(ctx context.Context, spotID uuid.UUID, startDate time.Time, endDate time.Time) ([]models.TimeUnit, error) {
	// startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	// endDate := startDate.AddDate(0, 0, 7)

	timeUnitResult, err := dbmodels.Timeunits.Query(
		ctx, p.db,
		sm.Columns(
			dbmodels.TimeunitColumns.Starttime,
			dbmodels.TimeunitColumns.Endtime,
			dbmodels.TimeunitColumns.Status,
		),
		dbmodels.SelectWhere.Timeunits.Parkingspotuuid.EQ(spotID),
		dbmodels.SelectWhere.Timeunits.Starttime.GTE(startDate),
		dbmodels.SelectWhere.Timeunits.Endtime.LTE(endDate),
		sm.OrderBy(dbmodels.TimeunitColumns.Starttime),
	).All()

	if (err != nil && errors.Is(err, sql.ErrNoRows)) || len(timeUnitResult) == 0 {
		err = ErrTimeUnitNotFound
		return make([]models.TimeUnit, 0, 1), err
	}

	availability := make([]models.TimeUnit, 0, len(timeUnitResult)) // Initialize slice

	for _, timeslot := range timeUnitResult {
		availability = append(availability, models.TimeUnit{
			StartTime: timeslot.Starttime.UTC(),
			EndTime:   timeslot.Endtime.UTC(),
			Status:    timeslot.Status,
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

func (p *PostgresRepository) GetMany(ctx context.Context, limit int, after omit.Val[Cursor], longitude float64, latitude float64, distance int32, startDate time.Time, endDate time.Time) ([]Entry, error) {

	type Result struct {
		dbmodels.Parkingspot
		DistanceToOrigin float64 `db:"distance_to_origin"`
	}

	centre := psql.F("ll_to_earth", psql.Arg(latitude), psql.Arg(longitude))
	spotPosition := psql.F("ll_to_earth", dbmodels.ParkingspotColumns.Latitude, dbmodels.ParkingspotColumns.Longitude)

	query := psql.Select(
		sm.Columns(
			dbmodels.Parkingspots.Columns(),
			psql.F("earth_distance", centre, spotPosition)(fm.As("distance_to_origin")),
		),
		sm.From(dbmodels.Parkingspots.Name(ctx)),
		sm.OrderBy("distance_to_origin").Asc(),
		sm.Where(
			psql.F(
				"earth_box",
				centre,
				psql.Arg(distance),
			)().OP(
				"@>",
				spotPosition,
			),
		),
		sm.Limit(limit),
	)

	entryCursor, err := bob.Cursor(ctx, p.db, query, scan.StructMapper[Result]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Entry{}, nil
		}
		return nil, err
	}
	defer entryCursor.Close()

	result := make([]Entry, 0, 8)
	for entryCursor.Next() {
		dbSpot, err := entryCursor.Get()
		if err != nil { // if there's an error, just return what we already have
			break
		}

		//Get availability
		availability, _ := p.GetAvalByUUID(ctx, dbSpot.Parkingspotuuid, startDate, endDate)
		// FIXME: How to handle errors apart from nothing found

		// TODO: Move dbmodels.Parkingspot -> Entry to a proc or something.
		result = append(result, Entry{
			ParkingSpot: models.ParkingSpot{
				Location: models.ParkingSpotLocation{
					PostalCode:    dbSpot.Postalcode,
					CountryCode:   dbSpot.Countrycode,
					City:          dbSpot.City,
					State:         dbSpot.State,
					StreetAddress: dbSpot.Streetaddress,
					Longitude:     float64(dbSpot.Longitude),
					Latitude:      float64(dbSpot.Latitude),
				},
				Features: models.ParkingSpotFeatures{
					Shelter:         dbSpot.Hasshelter,
					PlugIn:          dbSpot.Hasplugin,
					ChargingStation: dbSpot.Haschargingstation,
				},
				PricePerHour: dbSpot.Priceperhour,
				ID:           dbSpot.Parkingspotuuid,
				Availability: availability,
			},
			InternalID: dbSpot.Parkingspotid,
			OwnerID:    dbSpot.Userid,
		})
	}
	return result, nil
}

func formEntry(dbSpot dbmodels.Parkingspot, availability []models.TimeUnit) Entry {
	return Entry{
		ParkingSpot: models.ParkingSpot{
			Location: models.ParkingSpotLocation{
				PostalCode:    dbSpot.Postalcode,
				CountryCode:   dbSpot.Countrycode,
				City:          dbSpot.City,
				State:         dbSpot.State,
				StreetAddress: dbSpot.Streetaddress,
				Longitude:     float64(dbSpot.Longitude),
				Latitude:      float64(dbSpot.Latitude),
			},
			Features: models.ParkingSpotFeatures{
				Shelter:         dbSpot.Hasshelter,
				PlugIn:          dbSpot.Hasplugin,
				ChargingStation: dbSpot.Haschargingstation,
			},
			PricePerHour: dbSpot.Priceperhour,
			ID:           dbSpot.Parkingspotuuid,
			Availability: availability,
		},
		InternalID: dbSpot.Parkingspotid,
		OwnerID:    dbSpot.Userid,
	}
}
