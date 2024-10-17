package standardbooking

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
	"github.com/stephenafamo/bob/dialect/psql/sm"
	// "github.com/stephenafamo/bob/dialect/psql"
)

type PostgresRepository struct {
	db bob.DB
}

func NewPostgres(db bob.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (p *PostgresRepository) Create(ctx context.Context, userID int64, listingID int64, booking models.StandardBookingCreationInput) (Entry, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return Entry{}, fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback in case of panic

	bookingInserted, err := dbmodels.Bookings.Insert(ctx, p.db, &dbmodels.BookingSetter{
		Buyeruserid: omit.From(userID),
		Paidamount:  omit.From(booking.PaidAmount),
	})

	standardBookingInserted, err := dbmodels.Standardbookings.Insert(ctx, p.db, &dbmodels.StandardbookingSetter{
		Bookingid:    omit.From(bookingInserted.Bookingid),
		Listingid:    omit.From(listingID),
		Startunitnum: omit.From(booking.StartUnitNum),
		Endunitnum:   omit.From(booking.EndUnitNum),
		Date:         omit.From(booking.Date),
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				err = ErrDuplicatedStandardBooking
			}
		}
		return Entry{}, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return Entry{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	details := models.StandardBookingDetails{
		StartUnitNum: standardBookingInserted.Startunitnum,
		EndUnitNum:   standardBookingInserted.Endunitnum,
		Date:         standardBookingInserted.Date,
		PaidAmount:   booking.PaidAmount,
	}

	standardBooking := models.StandardBooking{
		Details: details,
		ID:      standardBookingInserted.Standardbookinguuid,
	}

	entry := Entry{
		StandardBooking: standardBooking,
		InternalID:      standardBookingInserted.Bookingid,
		BookingID:       bookingInserted.Bookingid,
		OwnerID:         userID,
		ListingID:       listingID,
	}

	return entry, nil
}

func (p *PostgresRepository) GetByUUID(ctx context.Context, bookingID uuid.UUID) (Entry, error) {
	result, err := dbmodels.Standardbookings.Query(
		ctx, p.db,
		sm.Columns(
			dbmodels.StandardbookingColumns.Startunitnum,
			dbmodels.StandardbookingColumns.Endunitnum,
			dbmodels.StandardbookingColumns.Date,
			dbmodels.BookingColumns.Paidamount,
			dbmodels.StandardbookingColumns.Bookingid,
			dbmodels.BookingColumns.Bookingid,
			dbmodels.BookingColumns.Buyeruserid,
			dbmodels.StandardbookingColumns.Listingid,
		),
		dbmodels.SelectJoins.Standardbookings.InnerJoin.BookingidBooking(ctx),
		dbmodels.SelectWhere.Standardbookings.Standardbookinguuid.EQ(bookingID),
	).One()

	// query := dbmodels.Standardbookings.Query(
	// 	ctx, p.db,
	// 	sm.Columns(
	// 		dbmodels.StandardbookingColumns.Startunitnum,
	// 		dbmodels.StandardbookingColumns.Endunitnum,
	// 		dbmodels.StandardbookingColumns.Date,
	// 		dbmodels.BookingColumns.Paidamount,
	// 		dbmodels.StandardbookingColumns.Bookingid,
	// 		dbmodels.BookingColumns.Bookingid,
	// 		dbmodels.BookingColumns.Buyeruserid,
	// 		dbmodels.StandardbookingColumns.Listingid,
	// 	),
	// 	sm.InnerJoin(
	// 		psql.On(dbmodels.StandardbookingColumns.Bookingid.EQ(dbmodels.BookingColumns.Bookingid)),
	// 	),
	// )

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return Entry{}, err
	}

	details := models.StandardBookingDetails{
		StartUnitNum: result.Startunitnum,
		EndUnitNum:   result.Endunitnum,
		Date:         result.Date,
		PaidAmount:   result.R.BookingidBooking.Paidamount,
	}

	standardBooking := models.StandardBooking{
		Details: details,
		ID:      result.Standardbookinguuid,
	}

	entry := Entry{
		StandardBooking: standardBooking,
		InternalID:      result.Bookingid,
		BookingID:       result.Bookingid,
		OwnerID:         result.R.BookingidBooking.Buyeruserid,
		ListingID:       result.Listingid,
	}

	return entry, nil
}
