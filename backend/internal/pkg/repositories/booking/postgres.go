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
	"github.com/govalues/decimal"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
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
		Userid:     omit.From(userID),
		Paidamount: omit.From(paidAmount),
	})
	if err != nil {
		return Entry{}, fmt.Errorf("could not execute insert: %w", err)
	}

	result := make([]models.TimeUnit, 0, len(booking.BookedTimes))

	umods := []bob.Mod[*dialect.UpdateQuery]{
		um.Table(dbmodels.Timeunits.Name(ctx)),
		um.SetCol(dbmodels.ColumnNames.Timeunits.Bookingid).ToArg(inserted.Bookingid),
	}
	var whereMods []mods.Where[*dialect.UpdateQuery]

	for _, time := range booking.BookedTimes {
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

		entryCursor, err := bob.Cursor(ctx, p.db, query, scan.StructMapper[*dbmodels.Timeunit]())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return Entry{}, nil
			}
		}
		defer entryCursor.Close()

		if entryCursor.Next() {
			get, _ := entryCursor.Get()
			result = append(result, timeUnitsFromDB(get))
		}
	
		if err != nil {
			return Entry{}, fmt.Errorf("could not execute update: %w", err)
		}
	}

	entry := Entry{
		Booking: models.Booking{
			Details: models.BookingDetails{
				PaidAmount:  float64(inserted.Bookingid),
				BookedTimes: result,
			},
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

func timeUnitsFromDB(model *dbmodels.Timeunit) models.TimeUnit {
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
