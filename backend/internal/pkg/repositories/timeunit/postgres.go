package timeunit

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

func (p *PostgresRepository) Create(ctx context.Context, timeslots []models.TimeSlot, listingID int64) (Entry, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return Entry{}, fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback in case of panic

	var entry Entry
	entry.TimeSlots = make([]models.TimeSlot, 0, len(timeslots)) // Initialize slice

	for _, timeslot := range timeslots {
		var newUnits []int16 // Collect units for the current timeslot

		for _, unit := range timeslot.Units {
			// Insert each unit
			inserted, err := dbmodels.Timeunits.Insert(ctx, p.db, &dbmodels.TimeunitSetter{
				Unitnum:   omit.From(unit),
				Date:      omit.From(timeslot.Date),
				Listingid: omit.From(listingID),
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

			// Populate newUnits with the current unit
			newUnits = append(newUnits, inserted.Unitnum)
		}

		// After processing all units for the current day, create a new TimeSlot
		newTimeSlot := models.TimeSlot{
			Date:  timeslot.Date,
			Units: newUnits,
		}

		entry.TimeSlots = append(entry.TimeSlots, newTimeSlot)
	}

	entry.ListingId = listingID

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return Entry{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	return entry, nil
}

func (p *PostgresRepository) GetByListingID(ctx context.Context, listingID int64) (Entry, error) {
	// Query to get all time units for the specified listing ID
	result, err := dbmodels.Timeunits.Query(
		ctx, p.db,
		sm.Columns(
			dbmodels.TimeunitColumns.Date,
			dbmodels.TimeunitColumns.Unitnum,
		),
		dbmodels.SelectWhere.Timeunits.Listingid.EQ(listingID),
		sm.OrderBy(dbmodels.TimeunitColumns.Date),
		sm.OrderBy(dbmodels.TimeunitColumns.Unitnum),
	).All()

	if (err != nil && errors.Is(err, sql.ErrNoRows)) || len(result) == 0 {
		err = ErrNotFound
		return Entry{}, err
	}

	// Create a map to group units by date
	timeslotMap := make(map[time.Time][]int16)

	// Iterate through the query results to populate the timeslotMap
	for _, row := range result {
		date := row.Date
		unitNum := row.Unitnum

		// Append the unit number to the corresponding date
		timeslotMap[date] = append(timeslotMap[date], unitNum)
	}

	// Construct the Entry struct
	var entry Entry
	entry.TimeSlots = make([]models.TimeSlot, 0, len(timeslotMap))

	// Populate Timeslots from the timeslotMap
	for date, units := range timeslotMap {
		entry.TimeSlots = append(entry.TimeSlots, models.TimeSlot{
			Date:  date,
			Units: units,
		})
	}

	return entry, nil
}

func (p *PostgresRepository) GetUnbookedByListingID(ctx context.Context, listingID int64) (Entry, error) {
	// Query to get all time units for the specified listing ID
	result, err := dbmodels.Timeunits.Query(
		ctx, p.db,
		sm.Columns(
			dbmodels.TimeunitColumns.Date,
			dbmodels.TimeunitColumns.Unitnum,
		),
		dbmodels.SelectWhere.Timeunits.Listingid.EQ(listingID),
		dbmodels.SelectWhere.Timeunits.Bookingid.IsNull(),
	).All()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return Entry{}, err
	}

	// Create a map to group units by date
	timeslotMap := make(map[time.Time][]int16)

	// Iterate through the query results to populate the timeslotMap
	for _, row := range result {
		date := row.Date
		unitNum := row.Unitnum

		// Append the unit number to the corresponding date
		timeslotMap[date] = append(timeslotMap[date], unitNum)
	}

	// Construct the Entry struct
	var entry Entry
	entry.TimeSlots = make([]models.TimeSlot, 0, len(timeslotMap))

	// Populate Timeslots from the timeslotMap
	for date, units := range timeslotMap {
		entry.TimeSlots = append(entry.TimeSlots, models.TimeSlot{
			Date:  date,
			Units: units,
		})
	}

	return entry, nil
}

func (p *PostgresRepository) GetByBookingID(ctx context.Context, bookingID int64) (Entry, error) {
	// Query to get all time units for the specified listing ID
	dateUnitNum, err := dbmodels.Timeunits.Query(
		ctx, p.db,
		sm.Columns(
			dbmodels.TimeunitColumns.Date,
			dbmodels.TimeunitColumns.Unitnum,
		),
		dbmodels.SelectWhere.Timeunits.Bookingid.EQ(bookingID),
	).All()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return Entry{}, err
	}

	// Query to get the listingID of the time units
	listingID, err := dbmodels.Timeunits.Query(
		ctx, p.db,
		sm.Columns(
			dbmodels.TimeunitColumns.Listingid,
		),
		dbmodels.SelectWhere.Timeunits.Bookingid.EQ(bookingID),
	).One()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return Entry{}, err
	}

	// Create a map to group units by date
	timeslotMap := make(map[time.Time][]int16)

	// Iterate through the query results to populate the timeslotMap
	for _, row := range dateUnitNum {
		date := row.Date
		unitNum := row.Unitnum

		// Append the unit number to the corresponding date
		timeslotMap[date] = append(timeslotMap[date], unitNum)
	}

	// Construct the Entry struct
	var entry Entry
	entry.TimeSlots = make([]models.TimeSlot, 0, len(timeslotMap))

	// Populate Timeslots from the timeslotMap
	for date, units := range timeslotMap {
		entry.TimeSlots = append(entry.TimeSlots, models.TimeSlot{
			Date:  date,
			Units: units,
		})
	}

	entry.BookingId = bookingID
	entry.ListingId = listingID.Listingid

	return entry, nil
}

func (p *PostgresRepository) DeleteByListingID(ctx context.Context, listingID int64, timeslots []models.TimeSlot) error {

	totalRowsAffected := int64(0)

	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback in case of panic

	for _, timeslot := range timeslots {

		rowsAffected, err := dbmodels.Timeunits.DeleteQ(
			ctx, p.db,
			dbmodels.DeleteWhere.Timeunits.Listingid.EQ(listingID),
			dbmodels.DeleteWhere.Timeunits.Date.EQ(timeslot.Date),
			dbmodels.DeleteWhere.Timeunits.Unitnum.In(timeslot.Units...),
			dbmodels.DeleteWhere.Timeunits.Bookingid.IsNull(),
		).Exec()
		if err != nil {
			return fmt.Errorf("could not execute delete: %w", err)
		}

		totalRowsAffected = totalRowsAffected + rowsAffected
	}

	if totalRowsAffected == 0 {
		return ErrNotFound
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

// Not sure this will be needed (can simply delete and create instead)
func (p *PostgresRepository) UpdateByListingID(ctx context.Context, listingID int64, timeslots []models.TimeSlot) (Entry, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return Entry{}, fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback in case of panic

	var entry Entry
	entry.TimeSlots = make([]models.TimeSlot, 0, len(timeslots)) // Initialize slice

	for _, timeslot := range timeslots {
		var newUnits []int16 // Collect units for the current timeslot

		for _, unit := range timeslot.Units {
			// Insert each unit
			inserted, err := dbmodels.Timeunits.Insert(ctx, p.db, &dbmodels.TimeunitSetter{
				Unitnum:   omit.From(unit),
				Date:      omit.From(timeslot.Date),
				Listingid: omit.From(listingID),
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

			// Populate newUnits with the current unit
			newUnits = append(newUnits, inserted.Unitnum)
		}

		// After processing all units for the current day, create a new TimeSlot
		newTimeSlot := models.TimeSlot{
			Date:  timeslot.Date,
			Units: newUnits,
		}

		entry.TimeSlots = append(entry.TimeSlots, newTimeSlot)
	}

	entry.ListingId = listingID

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return Entry{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	return entry, nil
}
