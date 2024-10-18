package parkingspot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmodels"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
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

func (p *PostgresRepository) Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (int64, Entry, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return -1, Entry{}, fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback if commit is not done

	longitude := float32(spot.Location.Longitude)
	latitude := float32(spot.Location.Latitude)

	inserted, err := dbmodels.Parkingspots.InsertQ(ctx, p.db,
		im.Columns(
			dbmodels.ParkingspotColumns.Userid,
			dbmodels.ParkingspotColumns.Postalcode,
			dbmodels.ParkingspotColumns.Countrycode,
			dbmodels.ParkingspotColumns.City,
			dbmodels.ParkingspotColumns.Streetaddress,
			dbmodels.ParkingspotColumns.Longitude,
			dbmodels.ParkingspotColumns.Latitude,
			dbmodels.ParkingspotColumns.Coordinates,
			dbmodels.ParkingspotColumns.Hasshelter,
			dbmodels.ParkingspotColumns.Hasplugin,
			dbmodels.ParkingspotColumns.Haschargingstation,
		),
		im.Values(
			psql.Arg(userID),
			psql.Arg(spot.Location.PostalCode),
			psql.Arg(spot.Location.CountryCode),
			psql.Arg(spot.Location.City),
			psql.Arg(spot.Location.StreetAddress),
			psql.Arg(longitude),
			psql.Arg(latitude),
			psql.Arg(psql.F("ST_SetSRID", psql.F("ST_MakePoint", psql.Arg(longitude, latitude)), 4326)),
			psql.Arg(spot.Features.Shelter),
			psql.Arg(spot.Features.PlugIn),
			psql.Arg(spot.Features.ChargingStation),
		),
	).One()

	if err != nil {
		// Handle duplicate error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				err = ErrDuplicatedAddress
			}
		}
		return -1, Entry{}, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, Entry{}, fmt.Errorf("could not commit transaction: %w", err)
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
		Location: location,
		Features: features,
		ID:       inserted.Parkingspotuuid,
	}

	entry := Entry{
		ParkingSpot: parkingspot,
		InternalID:  inserted.Parkingspotid,
		OwnerID:     inserted.Userid,
	}
	return inserted.Parkingspotid, entry, nil
}

func (p *PostgresRepository) DeleteByUUID(ctx context.Context, spotID uuid.UUID) error {
	rowsAffected, err := dbmodels.Parkingspots.DeleteQ(
		ctx, p.db,
		dbmodels.DeleteWhere.Parkingspots.Parkingspotuuid.EQ(spotID),
	).Exec()
	if err != nil {
		return fmt.Errorf("could not execute delete: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (p *PostgresRepository) GetByUUID(ctx context.Context, spotID uuid.UUID) (Entry, error) {
	result, err := dbmodels.Parkingspots.Query(
		ctx, p.db,
		sm.Columns(
			dbmodels.ParkingspotColumns.Postalcode,
			dbmodels.ParkingspotColumns.Countrycode,
			dbmodels.ParkingspotColumns.City,
			dbmodels.ParkingspotColumns.Streetaddress,
			dbmodels.ParkingspotColumns.Longitude,
			dbmodels.ParkingspotColumns.Latitude,
			dbmodels.ParkingspotColumns.Hasshelter,
			dbmodels.ParkingspotColumns.Hasplugin,
			dbmodels.ParkingspotColumns.Haschargingstation,
			dbmodels.ParkingspotColumns.Parkingspotid,
			dbmodels.ParkingspotColumns.Userid,
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
		PostalCode:    result.Postalcode,
		CountryCode:   result.Countrycode,
		City:          result.City,
		StreetAddress: result.Streetaddress,
		Longitude:     float64(result.Longitude),
		Latitude:      float64(result.Latitude),
	}

	features := models.ParkingSpotFeatures{
		Shelter:         result.Hasshelter,
		PlugIn:          result.Hasplugin,
		ChargingStation: result.Haschargingstation,
	}

	parkingspot := models.ParkingSpot{
		Location: location,
		Features: features,
		ID:       spotID,
	}

	return Entry{
		ParkingSpot: parkingspot,
		InternalID:  result.Parkingspotid,
		OwnerID:     result.Userid,
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
