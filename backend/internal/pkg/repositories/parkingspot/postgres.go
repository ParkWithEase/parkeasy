package parkingspot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmodels"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
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
		Province:           omit.From(spot.Location.State),
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
			Starttime:     omit.From(timeslot.StartTime),
			Endtime:       omit.From(timeslot.EndTime),
			Parkingspotid: omit.From(inserted.Parkingspotuuid),
			Status:        omit.From(timeslot.Status),
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
			StartTime: inserted.Starttime,
			EndTime:   inserted.Endtime,
			Status:    inserted.Status})
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

func (p *PostgresRepository) GetByUUID(ctx context.Context, spotID uuid.UUID) (Entry, error) {
	spotResult, err := dbmodels.Parkingspots.Query(
		ctx, p.db,
		sm.Columns(
			dbmodels.ParkingspotColumns.Postalcode,
			dbmodels.ParkingspotColumns.Countrycode,
			dbmodels.ParkingspotColumns.City,
			dbmodels.ParkingspotColumns.Province,
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

	timeUnitResult, err := dbmodels.Timeunits.Query(
		ctx, p.db,
		sm.Columns(
			dbmodels.TimeunitColumns.Starttime,
			dbmodels.TimeunitColumns.Endtime,
			dbmodels.TimeunitColumns.Status,
		),
		dbmodels.SelectWhere.Timeunits.Parkingspotid.EQ(spotID),
		sm.OrderBy(dbmodels.TimeunitColumns.Starttime),
	).All()

	availability := make([]models.TimeUnit, 0, len(timeUnitResult)) // Initialize slice

	for _, timeslot := range timeUnitResult {
		availability = append(availability, models.TimeUnit{
			StartTime: timeslot.Starttime,
			EndTime:   timeslot.Endtime,
			Status:    timeslot.Status})
	}

	location := models.ParkingSpotLocation{
		PostalCode:    spotResult.Postalcode,
		CountryCode:   spotResult.Countrycode,
		City:          spotResult.City,
		StreetAddress: spotResult.Streetaddress,
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
	}

	return Entry{
		ParkingSpot: parkingspot,
		InternalID:  spotResult.Parkingspotid,
		OwnerID:     spotResult.Userid,
	}, err
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

func (p *PostgresRepository) GetMany(ctx context.Context, limit int, after omit.Val[Cursor], longitude float64, latitude float64) ([]Entry, error) {
	where := dbmodels.SelectWhere.Parkingspots.
	
	if cursor, ok := after.Get(); ok {
		where = psql.WhereAnd(where, dbmodels.SelectWhere.Cars.Carid.GT(cursor.ID))
	}

	entryCursor, err := dbmodels.Cars.Query(
		ctx, p.db,
		sm.Columns(dbmodels.Cars.Columns()),
		sm.OrderBy(dbmodels.CarColumns.Carid),
		where,
		sm.Limit(limit),
	).Cursor()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Entry{}, nil
		}
		return nil, err
	}
	defer entryCursor.Close()

	result := make([]Entry, 0, 8)
	for entryCursor.Next() {
		dbCar, err := entryCursor.Get()
		if err != nil { // if there's an error, just return what we already have
			break
		}
		result = append(result, Entry{
			Car: models.Car{
				Details: models.CarDetails{
					LicensePlate: dbCar.Licenseplate,
					Make:         dbCar.Make,
					Model:        dbCar.Model,
					Color:        dbCar.Color,
				},
				ID: dbCar.Caruuid,
			},
			InternalID: dbCar.Carid,
			OwnerID:    dbCar.Userid,
		})
	}
	return result, nil
}
