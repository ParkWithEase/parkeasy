package standardbooking

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmodels"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/aarondl/opt/omit"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
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



func (p *PostgresRepository) Create(ctx context.Context, userID int64, listingID int64, booking models.StandardBookingCreationInput) (Entry, err) {
	interted, err = dbmodels.Bookings.Insert(ctx, p.db, &dbmodels.BookingSetter{
		
	})
}





func (p *PostgresRepository) GetByUUID(ctx context.Context, bookingID uuid.UUID) (Entry, err) {

}




