package standardbooking

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
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
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

func (p *PostgresRepository) Create(ctx context.Context, userID int64, listingID int64, booking *models.StandardBookingCreationInput,) (Entry, error) {
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

	updated, err := dbmodels.Timeunits.UpdateQ(
		ctx, p.db,
		dbmodels.UpdateWhere.Timeunits.Listingid.EQ(listingID),
		dbmodels.UpdateWhere.Timeunits.Date.EQ(booking.Date),
		dbmodels.UpdateWhere.Timeunits.Unitnum.GTE(booking.StartUnitNum),
		dbmodels.UpdateWhere.Timeunits.Unitnum.LTE(booking.EndUnitNum),
		&dbmodels.TimeunitSetter{
			Bookingid: omitnull.From(bookingInserted.Bookingid),
		},
		um.Returning(dbmodels.Timeunits.Columns()),
	).All()


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

	var units []int16

	for _, row := range updated {
		units = append(units, row.Unitnum)
	}

	timeslot := models.TimeSlot{
		Date:  booking.Date,
		Units: units,
	}

	entry := Entry{
		StandardBooking: standardBooking,
		TimeSlot: 		 timeslot,
		InternalID:      standardBookingInserted.Bookingid,
		BookingID:       bookingInserted.Bookingid,
		OwnerID:         userID,
		ListingID:       listingID,
	}

	return entry, nil
}

func (p *PostgresRepository) GetByUUID(ctx context.Context, bookingID uuid.UUID) (Entry, error) {
	type Result struct {
		Stdbook dbmodels.Standardbooking
		Booking dbmodels.Booking
	}
	query := psql.Select(
		sm.From(dbmodels.Standardbookings.Name(ctx)),
		sm.Columns(
			dbmodels.Standardbookings.Columns().WithPrefix("stdbook."),
			dbmodels.Bookings.Columns().WithPrefix("booking."),
		),
		dbmodels.SelectJoins.Standardbookings.InnerJoin.BookingidBooking(ctx),
		dbmodels.SelectWhere.Standardbookings.Standardbookinguuid.EQ(bookingID),
	)
	result, err := bob.One(ctx, p.db, query, scan.StructMapper[Result]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return Entry{}, err
	}

	details := models.StandardBookingDetails{
		StartUnitNum: result.Stdbook.Startunitnum,
		EndUnitNum:   result.Stdbook.Endunitnum,
		Date:         result.Stdbook.Date,
		PaidAmount:   result.Booking.Paidamount,
	}

	standardBooking := models.StandardBooking{
		Details: details,
		ID:      result.Stdbook.Standardbookinguuid,
	}

	entry := Entry{
		StandardBooking: standardBooking,
		InternalID:      result.Stdbook.Bookingid,
		BookingID:       result.Stdbook.Bookingid,
		OwnerID:         result.Booking.Buyeruserid,
		ListingID:       result.Stdbook.Listingid,
	}

	return entry, nil
}
