package booking

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmodels"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbtype"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
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

func (p *PostgresRepository) Create(ctx context.Context, userID int64, spotID int64, booking *models.BookingCreationInput) (Entry, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return Entry{}, fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback if commit is not done

	paidAmount, err := decimal.NewFromFloat64(booking.PaidAmount)
	if err != nil {
		return Entry{}, ErrInvalidPaidAmount
	}

	inserted, err := dbmodels.Bookings.Insert(ctx, p.db, &dbmodels.BookingSetter{
		Userid:        omit.From(userID),
		Parkingspotid: omit.From(spotID),
		Paidamount:    omit.From(paidAmount),
	})
	if err != nil {
		return Entry{}, fmt.Errorf("could not execute insert: %w", err)
	}

	updatedTimes := make([]models.TimeUnit, 0, len(booking.BookedTimes))

	umods := []bob.Mod[*dialect.UpdateQuery]{
		um.Table(dbmodels.Timeunits.Name(ctx)),
		um.SetCol(dbmodels.ColumnNames.Timeunits.Bookingid).ToArg(inserted.Bookingid),
	}

	for _, time := range booking.BookedTimes {
		var whereMods []mods.Where[*dialect.UpdateQuery]

		whereMods = append(whereMods, um.Where(
			dbmodels.TimeunitColumns.Timerange.OP(
				"&&",
				psql.Arg(dbtype.Tstzrange{
					Start: time.StartTime,
					End:   time.EndTime,
				}),
			),
		))

		uWhereMod := append(
			umods,
			psql.WhereAnd(whereMods...),
			um.Returning(dbmodels.Timeunits.Columns()),
		)

		query := psql.Update(uWhereMod...)

		updateCursor, err := bob.Cursor(ctx, p.db, query, scan.StructMapper[*dbmodels.Timeunit]())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return Entry{}, err
			}
		}
		defer updateCursor.Close()

		if updateCursor.Next() {
			get, err := updateCursor.Get()
			if err != nil {
				return Entry{}, fmt.Errorf("could not fetch updated time: %w", err)
			}
			updatedTimes = append(updatedTimes, timeUnitFromDB(get))
		}
	}

	amount, ok := inserted.Paidamount.Float64()
	if !ok {
		return Entry{}, fmt.Errorf("could not convert %v to float64", inserted.Paidamount)
	}

	entry := Entry{
		Booking: models.Booking{
			Details: models.BookingDetails{
				PaidAmount:  amount,
				BookedTimes: updatedTimes,
			},
			ID: inserted.Bookinguuid,
		},
		InternalID: inserted.Bookingid,
		OwnerID:    inserted.Userid,
	}

	err = tx.Commit()
	if err != nil {
		return Entry{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	return entry, nil
}

func (p *PostgresRepository) GetByUUID(ctx context.Context, bookingID uuid.UUID) (Entry, error) {
	bookingResult, err := dbmodels.Bookings.Query(
		ctx, p.db,
		sm.Columns(
			dbmodels.BookingColumns.Paidamount,
			dbmodels.BookingColumns.Bookingid,
			dbmodels.BookingColumns.Userid,
		),
		dbmodels.SelectWhere.Bookings.Bookinguuid.EQ(bookingID),
	).One()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return Entry{}, err
	}

	timeResult, err := dbmodels.Timeunits.Query(
		ctx, p.db,
		sm.Columns(dbmodels.TimeunitColumns.Timerange),
		sm.Columns(dbmodels.TimeunitColumns.Bookingid),
		psql.WhereAnd(
			dbmodels.SelectWhere.Timeunits.Bookingid.EQ(bookingResult.Bookingid),
		),
		dbmodels.SelectJoins.Timeunits.InnerJoin.ParkingspotidParkingspot(ctx),
		sm.OrderBy(psql.F("lower", dbmodels.TimeunitColumns.Timerange)),
	).All()
	if err != nil {
		return Entry{}, err
	}

	amount, ok := bookingResult.Paidamount.Float64()
	if !ok {
		return Entry{}, fmt.Errorf("could not convert %v to float64", bookingResult.Paidamount)
	}

	entry := Entry{
		Booking: models.Booking{
			Details: models.BookingDetails{
				PaidAmount:  amount,
				BookedTimes: timeUnitsFromDB(timeResult),
			},
			ID: bookingID,
		},
		InternalID: bookingResult.Bookingid,
		OwnerID:    bookingResult.Userid,
	}

	return entry, nil
}




func timeUnitsFromDB(model []*dbmodels.Timeunit) []models.TimeUnit {
	result := make([]models.TimeUnit, 0, len(model))
	for _, unit := range model {
		result = append(result, timeUnitFromDB(unit))
	}
	return result
}

func timeUnitFromDB(model *dbmodels.Timeunit) models.TimeUnit {
	var status string
	if _, ok := model.Bookingid.Get(); ok {
		status = "booked"
	} else {
		status = "available"
	}

	result := models.TimeUnit{
		StartTime: model.Timerange.Start,
		EndTime:   model.Timerange.End,
		Status:    status,
	}

	return result
}
