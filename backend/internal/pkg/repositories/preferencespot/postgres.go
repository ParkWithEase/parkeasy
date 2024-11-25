package preferencespot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmodels"
	// "github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/aarondl/opt/omit"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/sm"
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
	dbmodels.Parkingspot
	Preferencespotid int64 `db:"preferencespotid"`
}

func (p *PostgresRepository) Create(ctx context.Context, userID int64, spotID int64) error {
	_, err := dbmodels.Preferencespots.Insert(ctx, p.db, &dbmodels.PreferencespotSetter{
		Userid:        omit.From(userID),
		Parkingspotid: omit.From(spotID),
	})

	if err != nil {
		// Handle duplicate error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation && pgErr.ConstraintName == "preferencespotidx" {
				return ErrDuplicatedPreference
			}
		}
		return fmt.Errorf("could not execute insert: %w", err)
	}

	return nil
}

func (p *PostgresRepository) GetMany(ctx context.Context, userID int64, limit int, after omit.Val[Cursor]) ([]Entry, error) {
	where := dbmodels.SelectWhere.Preferencespots.Userid.EQ(userID)
	if cursor, ok := after.Get(); ok {
		where = psql.WhereAnd(where, dbmodels.SelectWhere.Preferencespots.Preferencespotid.GT(cursor.ID))
		fmt.Printf("YES")
	}

	// cursor, err := dbmodels.Parkingspots.Query(
	// 	ctx, p.db,
	// 	sm.Columns(dbmodels.Parkingspots.Columns()),
	// 	psql.WhereAnd(
	// 		dbmodels.SelectWhere.Preferencespots.Userid.EQ(userID),
	// 	),
	// 	sm.InnerJoin(dbmodels.SelectJoins.Preferencespots.InnerJoin.ParkingspotidParkingspot(ctx)),
	// 	sm.Limit(limit),
	// ).Cursor()

	smods := []bob.Mod[*dialect.SelectQuery]{
		sm.Columns(dbmodels.Parkingspots.Columns()),                                  // Select columns from Parkingspots
		sm.Columns(dbmodels.PreferencespotColumns.Preferencespotid),
		sm.From(dbmodels.Preferencespots.Name(ctx)),                                  // Preferencespots as the main table
		dbmodels.SelectJoins.Preferencespots.InnerJoin.ParkingspotidParkingspot(ctx), // Join with Preferencespots
		sm.Limit(limit), // Limit modifier
		where,
	}

	// Build the query
	query := psql.Select(smods...)

	sqlBuilder := strings.Builder{}
	query.WriteSQL(&sqlBuilder, nil, 0)
	fmt.Printf("Generated SQL: %s\n", sqlBuilder.String())

	// Execute the query and map results to getManyResult
	entryCursor, err := bob.Cursor(ctx, p.db, query, scan.StructMapper[Entry]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Entry{}, nil // No rows found
		}
		return nil, err // Other errors
	}
	defer entryCursor.Close()

	result := make([]Entry, 0, 8)
	for entryCursor.Next() {
		entry, err := entryCursor.Get()
		if err != nil { // if there's an error, just return what we already have
			break
		}
		result = append(result, entry)
	}

	return result, nil
}

func (p *PostgresRepository) Delete(ctx context.Context, userID int64, spotID int64) error {
	rowsAffected, err := dbmodels.Preferencespots.DeleteQ(
		ctx, p.db,
		dbmodels.DeleteWhere.Preferencespots.Userid.EQ(userID),
		dbmodels.DeleteWhere.Preferencespots.Parkingspotid.EQ(spotID),
	).Exec()
	if err != nil {
		return fmt.Errorf("could not execute delete: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// func entryFromDB(model getManyResult) (Entry, error) {
// 	lat, ok := model.Latitude.Float64()
// 	if !ok {
// 		return Entry{}, fmt.Errorf("could not convert %v to float64", model.Latitude)
// 	}
// 	lon, ok := model.Longitude.Float64()
// 	if !ok {
// 		return Entry{}, fmt.Errorf("could not convert %v to float64", model.Longitude)
// 	}
// 	price, ok := model.Priceperhour.Float64()
// 	if !ok {
// 		return Entry{}, fmt.Errorf("could not convert %v to float64", model.Priceperhour)
// 	}

// 	return Entry{
// 		ParkingSpot: models.ParkingSpot{
// 			Location: models.ParkingSpotLocation{
// 				PostalCode:    model.Postalcode,
// 				CountryCode:   model.Countrycode,
// 				City:          model.City,
// 				State:         model.State,
// 				StreetAddress: model.Streetaddress,
// 				Longitude:     lon,
// 				Latitude:      lat,
// 			},
// 			Features: models.ParkingSpotFeatures{
// 				Shelter:         model.Hasshelter,
// 				PlugIn:          model.Hasplugin,
// 				ChargingStation: model.Haschargingstation,
// 			},
// 			PricePerHour: price,
// 			ID:           model.Parkingspotuuid,
// 		},
// 	}, nil
// }
