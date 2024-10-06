package car

import (
	"context"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/scan"
)

type CarRepository struct {
	DBPool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *CarRepository {
	return &CarRepository{
		DBPool: pool,
	}
}

type Entry struct {
	models.Car
	InternalID int64 // The internal ID of this car
	OwnerID    int64 // The user id owning this car
}

var ErrNotFound = errors.New("no car found")

type Repository interface {
	Create(ctx context.Context, userID int64, car *models.CarCreationInput) (int64, Entry, error)
	GetByUUID(ctx context.Context, carID uuid.UUID) (Entry, error)
	GetOwnerByUUID(ctx context.Context, carID uuid.UUID) (int64, error)
	DeleteByUUID(ctx context.Context, carID uuid.UUID) error
	UpdateByUUID(ctx context.Context, carID uuid.UUID, car *models.CarCreationInput) (Entry, error)
}

func (r *CarRepository) Create(ctx context.Context, userID int64, car *models.CarCreationInput) (int64, Entry, error) {
	db, err := bob.Open("postgres", viper.GetString("db.url"))
	if err != nil {
		panic(err)
	}

	query := psql.Insert(
		im.Into("Car"),
		im.Values(psql.Arg(uuid.New(), car.LicensePlate, car.Make, car.Model, car.Color, userID)),
	)

	// Executing the query and mapping the result using scan.StructMapper
	entry, err := bob.One(ctx, db, query, scan.StructMapper[Entry]())
	if err != nil {
		return 0, Entry{}, err
	}

	return entry.InternalID, entry, nil
}

func (r *CarRepository) DeleteByUUID(ctx context.Context, carID uuid.UUID) error {
	db, err := bob.Open("postgres", viper.GetString("db.url"))
	if err != nil {
		panic(err)
	}

	query := psql.Delete(
		dm.From("Car"),
		dm.Where(psql.Quote("UUID").EQ(psql.Arg(carID))),
	)

	_, err = bob.Exec(ctx, db, query)
	if err != nil {
		panic(err)
	}

	return err
}

func (r *CarRepository) UpdateByUUID(ctx context.Context, carID uuid.UUID, car *models.CarCreationInput) (Entry, error) {
	db, err := bob.Open("postgres", viper.GetString("db.url"))
	if err != nil {
		panic(err)
	}

	query := psql.Update(
		um.Table("Car"),
		um.SetCol("LicensePlate").ToArg(car.LicensePlate),
		um.SetCol("Make").ToArg(car.Make),
		um.SetCol("Model").ToArg(car.Model),
		um.SetCol("Color").ToArg(car.Color),
		um.Where(psql.Quote("UUID").EQ(psql.Arg(carID))),
	)

	// Executing the query and mapping the result using scan.StructMapper
	entry, err := bob.One(ctx, db, query, scan.StructMapper[Entry]())
	if err != nil {
		return Entry{}, err
	}

	return entry, nil
}

func (r *CarRepository) GetByUUID(ctx context.Context, carID uuid.UUID) (Entry, error) {
	db, err := bob.Open("postgres", viper.GetString("db.url"))
	if err != nil {
		panic(err)
	}

	query := psql.Select(
		sm.From("Car"),
		sm.Where(psql.Quote("UUID").EQ(psql.Arg(carID))),
	)

	// Executing the query and mapping the result using scan.StructMapper
	entry, err := bob.One(ctx, db, query, scan.StructMapper[Entry]())
	if err != nil {
		return Entry{}, err
	}

	return entry, nil
}

func (r *CarRepository) GetOwnerByUUID(ctx context.Context, carID uuid.UUID) (int64, error) {
	db, err := bob.Open("postgres", viper.GetString("db.url"))
	if err != nil {
		panic(err)
	}

	query := psql.Select(
		sm.Columns("UserId"),
		sm.From("Car"),
		sm.Where(psql.Quote("UUID").EQ(psql.Arg(carID))),
	)

	// Executing the query and mapping the result using scan.StructMapper
	entry, err := bob.One(ctx, db, query, scan.StructMapper[Entry]())
	if err != nil {
		return 0, err
	}

	return entry.OwnerID, nil
}
