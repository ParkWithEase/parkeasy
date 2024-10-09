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
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
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

func (p *PostgresRepository) Create(ctx context.Context, userID int64, spot *models.ParkingSpotCreationInput) (int64, Entry, error) {
	longitude := float32(spot.Location.Longitude)
	latitude := float32(spot.Location.Latitude)

	inserted, err := dbmodels.Parkingspots.Insert(ctx, p.db, &dbmodels.ParkingspotSetter{
		Userid:             omit.From(userID),
		Postalcode:         omit.From(spot.Location.PostalCode),
		Countrycode:        omit.From(spot.Location.CountryCode),
		City:               omit.From(spot.Location.City),
		Streetaddress:      omit.From(spot.Location.StreetAddress),
		Longitude:          omit.From(longitude),
		Latitude:           omit.From(latitude),
		Hasshelter:         omit.From(spot.Features.Shelter),
		Hasplugin:          omit.From(spot.Features.PlugIn),
		Haschargingstation: omit.From(spot.Features.ChargingStation),
	})

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
		IsPublic:    inserted.Ispublic,
	}

	if err != nil {
		return -1, Entry{}, fmt.Errorf("could not commit transaction: %w", err)
	}
	return inserted.Parkingspotid, entry, nil
}

func (p *PostgresRepository) DeleteByUUID(ctx context.Context, spotID uuid.UUID) error {
	query := psql.Delete(
		dm.From(dbmodels.Parkingspots.Name(ctx)),
		dbmodels.DeleteWhere.Parkingspots.Parkingspotuuid.EQ(spotID),
	)

	// Execute the query
	_, err := bob.Exec(ctx, p.db, query)
	if err != nil {
		return fmt.Errorf("could not execute delete: %w", err)
	}

	return nil
}

func (p *PostgresRepository) GetByUUID(ctx context.Context, spotID uuid.UUID) (Entry, error) {
	query := psql.Select(
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
		sm.From(dbmodels.Parkingspots.Name(ctx)),
		dbmodels.SelectWhere.Parkingspots.Parkingspotuuid.EQ(spotID),
	)
	result, err := bob.One(ctx, p.db, query, scan.StructMapper[dbmodels.Parkingspot]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrParkingSpotNotFound
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
	query := psql.Select(
		sm.Columns(dbmodels.ParkingspotColumns.Userid),
		sm.From(dbmodels.Parkingspots.Name(ctx)),
		dbmodels.SelectWhere.Parkingspots.Parkingspotuuid.EQ(spotID),
	)
	result, err := bob.One(ctx, p.db, query, scan.SingleColumnMapper[int32])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrParkingSpotNotFound
		}
		return -1, err
	}

	return int64(result), err
}
