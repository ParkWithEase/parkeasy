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
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
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
		PricePerHour:  inserted.Priceperhour.GetOrZero(),
		ParkingSpotID: inserted.Parkingspotid,
		IsActive:      inserted.Isactive,
	}

	return inserted.Parkingspotid, entry, nil
}

// func (p *PostgresRepository) GetByUUID(ctx context.Context, listingID uuid.UUID) (Entry, error) {

// }
