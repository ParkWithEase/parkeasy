package listing

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmodels"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
)

type PostgresRepository struct {
	db bob.DB
}

func NewPostgres(db bob.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (p *PostgresRepository) Create(ctx context.Context, parkingSpotID int64, listing *models.ListingCreationInput) (int64, Entry, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return -1, Entry{}, fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback if commit is not done

	inserted, err := dbmodels.Listings.Insert(ctx, p.db, &dbmodels.ListingSetter{
		Parkingspotid: omit.From(parkingSpotID),
		Priceperhour:  omitnull.From(listing.PricePerHour),
		Isactive:      omit.From(listing.MakePublic),
	})
	if err != nil {
		// Handle duplicate error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				err = ErrDuplicatedListing
			}
		}
		return -1, Entry{}, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, Entry{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	entry := Entry{
		ID:            inserted.Listinguuid,
		InternalID:    inserted.Listingid,
		PricePerHour:  inserted.Priceperhour.GetOr(-1),
		ParkingSpotID: inserted.Parkingspotid,
		IsActive:      inserted.Isactive,
	}

	return inserted.Listingid, entry, nil
}

func (p *PostgresRepository) GetByUUID(ctx context.Context, listingID uuid.UUID) (Entry, error) {
	result, err := dbmodels.Listings.Query(
		ctx, p.db,
		sm.Columns(
			dbmodels.ListingColumns.Listinguuid,
			dbmodels.ListingColumns.Priceperhour,
			dbmodels.ListingColumns.Parkingspotid,
			dbmodels.ListingColumns.Isactive,
		),
		dbmodels.SelectWhere.Listings.Listinguuid.EQ(listingID),
	).One()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return Entry{}, err
	}

	return Entry{
		ID:            result.Listinguuid,
		InternalID:    result.Listingid,
		PricePerHour:  result.Priceperhour.GetOr(-1),
		IsActive:      result.Isactive,
		ParkingSpotID: result.Parkingspotid,
	}, nil
}

func (p *PostgresRepository) GetSpotByUUID(ctx context.Context, listingID uuid.UUID) (int64, error) {
	result, err := dbmodels.Listings.Query(
		ctx, p.db,
		sm.Columns(
			dbmodels.ListingColumns.Parkingspotid,
		),
		dbmodels.SelectWhere.Listings.Listinguuid.EQ(listingID),
	).One()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return -1, err
	}

	return result.Parkingspotid, nil
}

func (p *PostgresRepository) UnlistByUUID(ctx context.Context, listingID uuid.UUID) (Entry, error) {
	result, err := dbmodels.Listings.UpdateQ(
		ctx, p.db,
		dbmodels.UpdateWhere.Listings.Listinguuid.EQ(listingID),
		&dbmodels.ListingSetter{
			Isactive: omit.From(false),
		},
		um.Returning(dbmodels.Listings.Columns()),
	).One()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Entry{}, ErrNotFound
		}
		return Entry{}, fmt.Errorf("could not execute update: %w", err)
	}

	return Entry{
		ID:            result.Listinguuid,
		InternalID:    result.Listingid,
		PricePerHour:  result.Priceperhour.GetOr(-1),
		ParkingSpotID: result.Parkingspotid,
		IsActive:      result.Isactive,
	}, nil
}

func (p *PostgresRepository) UpdateByUUID(ctx context.Context, listingID uuid.UUID, listing *models.ListingCreationInput) (Entry, error) {
	result, err := dbmodels.Listings.UpdateQ(
		ctx, p.db,
		dbmodels.UpdateWhere.Listings.Listinguuid.EQ(listingID),
		&dbmodels.ListingSetter{
			Priceperhour: omitnull.From(listing.PricePerHour),
			Isactive:     omit.From(listing.MakePublic),
		},
		um.Returning(dbmodels.Listings.Columns()),
	).One()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Entry{}, ErrNotFound
		}
		return Entry{}, fmt.Errorf("could not execute update: %w", err)
	}

	return Entry{
		ID:            result.Listinguuid,
		InternalID:    result.Listingid,
		PricePerHour:  result.Priceperhour.GetOr(-1),
		ParkingSpotID: result.Parkingspotid,
		IsActive:      result.Isactive,
	}, nil
}
