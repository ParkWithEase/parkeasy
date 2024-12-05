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

type getManyWithLocation struct {
	dbmodels.Booking
	Caruuid         uuid.UUID       `db:"caruuid" `
	Parkingspotuuid uuid.UUID       `db:"parkingspotuuid" `
	Postalcode      string          `db:"postalcode" `
	Countrycode     string          `db:"countrycode" `
	City            string          `db:"city" `
	State           string          `db:"state" `
	Streetaddress   string          `db:"streetaddress" `
	Longitude       decimal.Decimal `db:"longitude" `
	Latitude        decimal.Decimal `db:"latitude" `
}

func (p *PostgresRepository) Create(ctx context.Context, booking *CreateInput) (EntryWithTimes, error) {
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
		Userid:        omit.From(booking.UserID),
		Parkingspotid: omit.From(booking.SpotID),
		Carid:         omit.From(booking.CarID),
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

	// Variable for OR clauses in where, used for timeslots timeranges
	timeslots := timeSlotsToSQLExpr(booking.BookedTimes)

	// Variable for AND clauses in where
	var whereMods []mods.Where[*dialect.UpdateQuery]
	whereMods = append(
		whereMods,
		dbmodels.UpdateWhere.Timeunits.Bookingid.IsNull(),
		dbmodels.UpdateWhere.Timeunits.Parkingspotid.EQ(booking.SpotID),
		um.Where(timeslots),
	)

	umods = append(umods, psql.WhereAnd(whereMods...), um.Returning(dbmodels.Timeunits.Columns()))
	query := dbmodels.Timeunits.UpdateQ(ctx, tx, umods...)

	updatedCursor, err := query.Cursor()
	if err != nil {
		return EntryWithTimes{}, fmt.Errorf("could not update time units: %w", err)
	}
	defer updatedCursor.Close()

	bookedSlots := make([]models.TimeUnit, 0, len(booking.BookedTimes))
	for updatedCursor.Next() {
		dbtime, err := updatedCursor.Get()
		if err != nil {
			return EntryWithTimes{}, fmt.Errorf("could not get time units: %w", err)
		}

		bookedSlots = append(bookedSlots, timeUnitFromDB(dbtime))
	}
	// Check if the count matches the expected number of booked times
	if len(bookedSlots) != len(booking.BookedTimes) {
		return EntryWithTimes{}, ErrTimeAlreadyBooked
	}
	//------------------------------------------------

	type bookingUUIDs struct {
		Parkingspotuuid uuid.UUID
		Caruuid         uuid.UUID
	}
	uuidQuery := psql.Select(
		sm.Columns(dbmodels.ParkingspotColumns.Parkingspotuuid, dbmodels.CarColumns.Caruuid),
		sm.From(dbmodels.Bookings.Name(ctx)),
		dbmodels.SelectJoins.Bookings.InnerJoin.ParkingspotidParkingspot(ctx),
		dbmodels.SelectJoins.Bookings.InnerJoin.CaridCar(ctx),
		dbmodels.SelectWhere.Bookings.Bookingid.EQ(inserted.Bookingid),
	)
	ids, err := bob.One(ctx, tx, uuidQuery, scan.StructMapper[bookingUUIDs]())
	if err != nil {
		return EntryWithTimes{}, fmt.Errorf("could not get car and spot uuids: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return EntryWithTimes{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	entry := EntryWithTimes{
		Entry:       formEntry(inserted, ids.Parkingspotuuid, ids.Caruuid),
		BookedTimes: bookedSlots,
	}

	return entry, nil
}

func timeSlotsToSQLExpr(units []models.TimeUnit) dialect.Expression {
	var expression dialect.Expression
	for _, bookTime := range units {
		test := dbmodels.TimeunitColumns.Timerange.OP(
			"&&",
			psql.Arg(dbtype.Tstzrange{
				Start: bookTime.StartTime,
				End:   bookTime.EndTime,
			}),
		)
		if expression.Base == nil {
			expression = test
		} else {
			expression = expression.Or(test)
		}
	}
	return expression
}

func (p *PostgresRepository) GetByUUID(ctx context.Context, bookingID uuid.UUID) (EntryWithTimes, error) {
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
	var bookingResult getManyResult
	bookingResult, err := bob.One(ctx, p.db, query, scan.StructMapper[getManyResult]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return EntryWithTimes{}, err
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

	entry := EntryWithTimes{
		Entry:       formEntry(&bookingResult.Booking, bookingResult.Parkingspotuuid, bookingResult.Caruuid),
		BookedTimes: timeUnitsFromDB(timeResult),
	}

	return entry, nil
}

func (p *PostgresRepository) GetManyForBuyer(ctx context.Context, limit int, after omit.Val[Cursor], userID int64, filter *Filter) ([]EntryWithLocation, error) {
	log := zerolog.Ctx(ctx).
		With().
		Str("component", "booking.Postgres").
		Logger()

	smods := []bob.Mod[*dialect.SelectQuery]{
		sm.Columns(dbmodels.Bookings.Columns()),
		sm.Columns(dbmodels.ParkingspotColumns.Parkingspotuuid),
		sm.Columns(dbmodels.ParkingspotColumns.City),
		sm.Columns(dbmodels.ParkingspotColumns.Countrycode),
		sm.Columns(dbmodels.ParkingspotColumns.Latitude),
		sm.Columns(dbmodels.ParkingspotColumns.Longitude),
		sm.Columns(dbmodels.ParkingspotColumns.Postalcode),
		sm.Columns(dbmodels.ParkingspotColumns.State),
		sm.Columns(dbmodels.ParkingspotColumns.Streetaddress),
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
		psql.WhereAnd(whereMods...),
	)

	query := psql.Select(smods...)

	entryCursor, err := bob.Cursor(ctx, p.db, query, scan.StructMapper[getManyWithLocation]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []EntryWithLocation{}, nil
		}
		return nil, err
	}
	defer entryCursor.Close()

	result := make([]EntryWithLocation, 0, 8)

	for entryCursor.Next() {
		get, err := entryCursor.Get()
		if err != nil {
			log.Err(err).Msg("error iterating get many cursor")
			break
		}

		// Convert lat and long from deciaml to float
		lat, _ := get.Latitude.Float64()
		long, _ := get.Longitude.Float64()

		result = append(result, EntryWithLocation{
			Entry: formEntry(&get.Booking, get.Parkingspotuuid, get.Caruuid),
			ParkingSpotLocation: models.ParkingSpotLocation{
				PostalCode:    get.Postalcode,
				CountryCode:   get.Countrycode,
				StreetAddress: get.Streetaddress,
				State:         get.State,
				City:          get.City,
				Latitude:      lat,
				Longitude:     long,
			},
		})
	}

	return result, nil
}

func (p *PostgresRepository) GetManyForOwner(ctx context.Context, limit int, after omit.Val[Cursor], userID int64, filter *Filter) ([]EntryWithLocation, error) {
	log := zerolog.Ctx(ctx).
		With().
		Str("component", "booking.Postgres").
		Logger()

	smods := []bob.Mod[*dialect.SelectQuery]{
		sm.Columns(dbmodels.Bookings.Columns()),
		sm.Columns(dbmodels.ParkingspotColumns.Parkingspotuuid),
		sm.Columns(dbmodels.ParkingspotColumns.City),
		sm.Columns(dbmodels.ParkingspotColumns.Countrycode),
		sm.Columns(dbmodels.ParkingspotColumns.Latitude),
		sm.Columns(dbmodels.ParkingspotColumns.Longitude),
		sm.Columns(dbmodels.ParkingspotColumns.Postalcode),
		sm.Columns(dbmodels.ParkingspotColumns.State),
		sm.Columns(dbmodels.ParkingspotColumns.Streetaddress),
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
		psql.WhereAnd(whereMods...),
	)

	query := psql.Select(smods...)

	entryCursor, err := bob.Cursor(ctx, p.db, query, scan.StructMapper[getManyWithLocation]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []EntryWithLocation{}, nil
		}
		return nil, err
	}
	defer entryCursor.Close()

	result := make([]EntryWithLocation, 0, 8)

	for entryCursor.Next() {
		get, err := entryCursor.Get()
		if err != nil {
			log.Err(err).Msg("error iterating get many cursor")
			break
		}

		// Convert lat and long from deciaml to float
		lat, _ := get.Latitude.Float64()
		long, _ := get.Longitude.Float64()

		result = append(result, EntryWithLocation{
			Entry: formEntry(&get.Booking, get.Parkingspotuuid, get.Caruuid),
			ParkingSpotLocation: models.ParkingSpotLocation{
				PostalCode:    get.Postalcode,
				CountryCode:   get.Countrycode,
				StreetAddress: get.Streetaddress,
				State:         get.State,
				City:          get.City,
				Latitude:      lat,
				Longitude:     long,
			},
		})
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

func formEntry(entry *dbmodels.Booking, spotUUID, carUUID uuid.UUID) Entry {
	amount, _ := entry.Paidamount.Float64()

	return Entry{
		Booking: models.Booking{
			CreatedAt:     entry.Createdat,
			PaidAmount:    amount,
			ID:            entry.Bookinguuid,
			ParkingSpotID: spotUUID,
			CarID:         carUUID,
		},
		InternalID: entry.Bookingid,
		BookerID:   entry.Userid,
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
