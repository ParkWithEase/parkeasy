package booking

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmodels"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbtype"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/rs/zerolog"
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

type getManyResult struct {
	dbmodels.Booking
	Parkingspotuuid uuid.UUID `db:"parkingspotuuid" `
	Caruuid         uuid.UUID `db:"caruuid" `
}

func (p *PostgresRepository) Create(ctx context.Context, userID int64, spotID int64, carID int64, booking *models.BookingCreationDBInput) (EntryWithTimes, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return EntryWithTimes{}, fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback if commit is not done

	paidAmount, err := decimal.NewFromFloat64(booking.PaidAmount)
	if err != nil {
		return EntryWithTimes{}, ErrInvalidPaidAmount
	}

	inserted, err := dbmodels.Bookings.Insert(ctx, tx, &dbmodels.BookingSetter{
		Userid:        omit.From(userID),
		Parkingspotid: omit.From(spotID),
		Carid:         omit.From(carID),
		Paidamount:    omit.From(paidAmount),
	})
	if err != nil {
		return EntryWithTimes{}, fmt.Errorf("could not execute insert: %w", err)
	}

	//--------Update the corresponding time slots--------
	// Check if all the time slots are available
	umods := []bob.Mod[*dialect.UpdateQuery]{&dbmodels.TimeunitSetter{
		Bookingid: omitnull.From(inserted.Bookingid),
	}}

	//Variable for OR clauses in where, used for timeslots timeranges
	var whereOrMods []mods.Where[*dialect.UpdateQuery]
	for _, time := range booking.BookingInfo.BookedTimes {
		whereOrMods = append(whereOrMods, um.Where(
			dbmodels.TimeunitColumns.Timerange.OP(
				"&&",
				psql.Arg(dbtype.Tstzrange{
					Start: time.StartTime,
					End:   time.EndTime,
				}),
			),
		))
	}

	//Variable for AND clauses in where
	var whereMods []mods.Where[*dialect.UpdateQuery]
	whereMods = append(
		whereMods,
		dbmodels.UpdateWhere.Timeunits.Bookingid.IsNull(),
		dbmodels.UpdateWhere.Timeunits.Parkingspotid.EQ(spotID),
		psql.WhereOr(whereOrMods...),
	)

	umods = append(umods, psql.WhereAnd(whereMods...), um.Returning(dbmodels.Timeunits.Columns()))
	query := dbmodels.Timeunits.UpdateQ(ctx, tx, umods...)

	updatedCursor, err := query.Cursor()
	if err != nil {
		return EntryWithTimes{}, fmt.Errorf("could not update time units: %w", err)
	}
	defer updatedCursor.Close()

	bookedSlots := make([]models.TimeUnit, 0, len(booking.BookingInfo.BookedTimes))
	for updatedCursor.Next() {
		dbtime, err := updatedCursor.Get()
		if err != nil {
			return EntryWithTimes{}, fmt.Errorf("could not get time units: %w", err)
		}

		bookedSlots = append(bookedSlots, timeUnitFromDB(dbtime))
	}
	// Check if the count matches the expected number of booked times
	if len(bookedSlots) != len(booking.BookingInfo.BookedTimes) {
		return EntryWithTimes{}, ErrTimeAlreadyBooked
	}
	//------------------------------------------------

	err = tx.Commit()
	if err != nil {
		return EntryWithTimes{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	amount, ok := inserted.Paidamount.Float64()
	if !ok {
		return EntryWithTimes{}, fmt.Errorf("could not convert %v to float64", inserted.Paidamount)
	}

	entry := EntryWithTimes{
		Entry: formEntry(amount,
			inserted.Bookinguuid,
			inserted.Bookingid,
			booking.BookingInfo.ParkingSpotID,
			booking.BookingInfo.CarID,
			inserted.Createdat,
			inserted.Userid,
		),
		BookedTimes: bookedSlots,
	}

	return entry, nil
}

func (p *PostgresRepository) GetByUUID(ctx context.Context, bookingID uuid.UUID) (EntryWithTimes, error) {
	log := zerolog.Ctx(ctx).
		With().
		Str("component", "booking.Postgres").
		Logger()

	// Build select mods
	smods := []bob.Mod[*dialect.SelectQuery]{
		sm.Columns(dbmodels.Bookings.Columns()),
		sm.Columns(dbmodels.ParkingspotColumns.Parkingspotuuid),
		sm.Columns(dbmodels.CarColumns.Caruuid),
		sm.From(dbmodels.Bookings.Name(ctx)),
		dbmodels.SelectJoins.Bookings.InnerJoin.ParkingspotidParkingspot(ctx),
		dbmodels.SelectJoins.Bookings.InnerJoin.CaridCar(ctx),
		dbmodels.SelectWhere.Bookings.Bookinguuid.EQ(bookingID),
		sm.Limit(1),
	}

	// Build the query
	query := psql.Select(smods...)

	// Execute the query
	cursor, err := bob.Cursor(ctx, p.db, query, scan.StructMapper[getManyResult]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return EntryWithTimes{}, err
	}
	defer cursor.Close()

	var bookingResult getManyResult

	if cursor.Next() {
		bookingResult, err = cursor.Get()
		if err != nil {
			log.Err(err).Msg("error retrieving record from cursor")
			return EntryWithTimes{}, err
		}
	} else {
		// No results found in cursor
		log.Info().Msg("no results found in cursor")
		return EntryWithTimes{}, ErrNotFound
	}

	timeResult, err := dbmodels.Timeunits.Query(
		ctx, p.db,
		sm.Columns(dbmodels.TimeunitColumns.Timerange),
		sm.Columns(dbmodels.TimeunitColumns.Bookingid),
		psql.WhereAnd(
			dbmodels.SelectWhere.Timeunits.Bookingid.EQ(bookingResult.Bookingid),
		),
		sm.OrderBy(psql.F("lower", dbmodels.TimeunitColumns.Timerange)),
	).All()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return EntryWithTimes{}, err
	}

	amount, ok := bookingResult.Paidamount.Float64()
	if !ok {
		return EntryWithTimes{}, fmt.Errorf("could not convert %v to float64", bookingResult.Paidamount)
	}

	entry := EntryWithTimes{
		Entry: formEntry(amount,
			bookingResult.Bookinguuid,
			bookingResult.Bookingid,
			bookingResult.Parkingspotuuid,
			bookingResult.Caruuid,
			bookingResult.Createdat,
			bookingResult.Userid,
		),
		BookedTimes: timeUnitsFromDB(timeResult),
	}

	return entry, nil
}

func (p *PostgresRepository) GetBookedTimesByUUID(ctx context.Context, bookingUUID uuid.UUID) ([]models.TimeUnit, error) {
	// Query timeunits
	result, err := dbmodels.Timeunits.Query(
		ctx, p.db,
		sm.Columns(dbmodels.TimeunitColumns.Timerange),
		sm.Columns(dbmodels.TimeunitColumns.Bookingid),
		psql.WhereAnd(
			dbmodels.SelectWhere.Bookings.Bookinguuid.EQ(bookingUUID),
		),
		dbmodels.SelectJoins.Timeunits.InnerJoin.BookingidBooking(ctx),
		sm.OrderBy(psql.F("lower", dbmodels.TimeunitColumns.Timerange)),
	).All()
	if err != nil {
		return nil, err
	}

	// If no rows found
	if len(result) == 0 {
		// Ignore errors here, just treat it as not existing

		// Check if the booking exists
		exists, _ := dbmodels.Bookings.Query(
			ctx, p.db,
			sm.Columns(1),
			dbmodels.SelectWhere.Bookings.Bookinguuid.EQ(bookingUUID),
		).Exists()

		if !exists {
			return nil, ErrNotFound
		}

		// No time units is not an error
		return []models.TimeUnit{}, nil
	}

	return timeUnitsFromDB(result), nil
}

func (p *PostgresRepository) GetBookerByUUID(ctx context.Context, bookingID uuid.UUID) (int64, error) {
	result, err := dbmodels.Bookings.Query(
		ctx, p.db,
		sm.Columns(
			dbmodels.BookingColumns.Userid,
		),
		dbmodels.SelectWhere.Bookings.Bookinguuid.EQ(bookingID),
	).One()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return 0, err
	}

	return result.Userid, nil
}

func (p *PostgresRepository) GetManyForBuyer(ctx context.Context, limit int, after omit.Val[Cursor], userID int64, filter *Filter) ([]Entry, error) {
	log := zerolog.Ctx(ctx).
		With().
		Str("component", "booking.Postgres").
		Logger()

	smods := []bob.Mod[*dialect.SelectQuery]{
		sm.Columns(dbmodels.Bookings.Columns()),
		sm.Columns(dbmodels.ParkingspotColumns.Parkingspotuuid),
		sm.Columns(dbmodels.CarColumns.Caruuid),
	}
	var whereMods []mods.Where[*dialect.SelectQuery]

	whereMods = append(whereMods, dbmodels.SelectWhere.Bookings.Userid.EQ(userID))
	if cursor, ok := after.Get(); ok {
		whereMods = append(whereMods, dbmodels.SelectWhere.Bookings.Bookingid.LT(cursor.ID))
	}

	if filter.SpotID != 0 {
		whereMods = append(whereMods, dbmodels.SelectWhere.Bookings.Parkingspotid.EQ(filter.SpotID))
	}

	smods = append(
		smods,
		sm.From(dbmodels.Bookings.Name(ctx)),
		dbmodels.SelectJoins.Bookings.InnerJoin.ParkingspotidParkingspot(ctx),
		dbmodels.SelectJoins.Bookings.InnerJoin.CaridCar(ctx),
		sm.Limit(limit),
		sm.OrderBy(dbmodels.BookingColumns.Bookingid).Desc(),
	)

	smods = append(smods, psql.WhereAnd(whereMods...))
	query := psql.Select(smods...)
	str, _, err := query.Build()
	fmt.Println(str)

	entryCursor, err := bob.Cursor(ctx, p.db, query, scan.StructMapper[getManyResult]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Entry{}, nil
		}
		return nil, err
	}
	defer entryCursor.Close()

	result := make([]Entry, 0, 8)

	for entryCursor.Next() {
		get, err := entryCursor.Get()
		if err != nil {
			return []Entry{}, fmt.Errorf("could not fetch booking details: %w", err)
		}

		amount, ok := get.Paidamount.Float64()
		if !ok {
			// if there is an error converting amount to float
			log.Err(err).Msg("error while converting amount to float")
			continue
		}

		entry := formEntry(amount,
			get.Bookinguuid,
			get.Bookingid,
			get.Parkingspotuuid,
			get.Caruuid,
			get.Createdat,
			get.Userid,
		)

		result = append(result, entry)
	}

	return result, nil
}

func (p *PostgresRepository) GetManyForSeller(ctx context.Context, limit int, after omit.Val[Cursor], userID int64, filter *Filter) ([]Entry, error) {
	log := zerolog.Ctx(ctx).
		With().
		Str("component", "booking.Postgres").
		Logger()

	smods := []bob.Mod[*dialect.SelectQuery]{
		sm.Columns(dbmodels.Bookings.Columns()),
		sm.Columns(dbmodels.ParkingspotColumns.Parkingspotuuid),
		sm.Columns(dbmodels.CarColumns.Caruuid),
	}
	var whereMods []mods.Where[*dialect.SelectQuery]

	whereMods = append(whereMods, dbmodels.SelectWhere.Parkingspots.Userid.EQ(userID))
	if cursor, ok := after.Get(); ok {
		whereMods = append(whereMods, dbmodels.SelectWhere.Bookings.Bookingid.LT(cursor.ID))
	}

	if filter.SpotID != 0 {
		whereMods = append(whereMods, dbmodels.SelectWhere.Bookings.Parkingspotid.EQ(filter.SpotID))
	}

	smods = append(
		smods,
		sm.From(dbmodels.Bookings.Name(ctx)),
		dbmodels.SelectJoins.Bookings.InnerJoin.ParkingspotidParkingspot(ctx),
		dbmodels.SelectJoins.Bookings.InnerJoin.CaridCar(ctx),
		sm.Limit(limit),
		sm.OrderBy(dbmodels.BookingColumns.Bookingid).Desc(),
	)

	smods = append(smods, psql.WhereAnd(whereMods...))
	query := psql.Select(smods...)

	entryCursor, err := bob.Cursor(ctx, p.db, query, scan.StructMapper[getManyResult]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Entry{}, nil
		}
		return nil, err
	}
	defer entryCursor.Close()

	result := make([]Entry, 0, 8)

	for entryCursor.Next() {
		get, err := entryCursor.Get()
		if err != nil {
			return []Entry{}, fmt.Errorf("could not fetch booking details: %w", err)
		}

		amount, ok := get.Paidamount.Float64()
		if !ok {
			// if there is an error converting amount to float
			log.Err(err).Msg("error while converting amount to float")
			continue
		}

		entry := formEntry(
			amount,
			get.Bookinguuid,
			get.Bookingid,
			get.Parkingspotuuid,
			get.Caruuid,
			get.Createdat,
			get.Userid,
		)

		result = append(result, entry)
	}

	return result, nil
}

func timeUnitsFromDB(model []*dbmodels.Timeunit) []models.TimeUnit {
	result := make([]models.TimeUnit, 0, len(model))
	for _, unit := range model {
		result = append(result, timeUnitFromDB(unit))
	}
	return result
}

func formEntry(amount float64, bookingUUID uuid.UUID, bookingID int64, parkingSpotUUID uuid.UUID, carID uuid.UUID, createdAt time.Time, bookerID int64) Entry {
	return Entry{
		Booking: models.Booking{
			PaidAmount:    amount,
			ID:            bookingUUID,
			ParkingSpotID: parkingSpotUUID,
			CarID:         carID,
			CreatedAt:     createdAt,
		},
		InternalID: bookingID,
		BookerID:   bookerID,
	}
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
