package car

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmodels"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
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

func (p *PostgresRepository) Create(ctx context.Context, userID64 int64, car *models.CarCreationInput) (int64, Entry, error) {
	if userID64 > math.MaxInt32 || userID64 < 0 {
		return -1, Entry{}, fmt.Errorf("userID out of range: %v", userID64)
	}
	userID := int32(userID64)

	inserted, err := dbmodels.Cars.Insert(ctx, p.db, &dbmodels.CarSetter{
		Userid:       omit.From(userID),
		Licenseplate: omit.From(car.LicensePlate),
		Make:         omit.From(car.Make),
		Model:        omit.From(car.Model),
		Color:        omit.From(car.Color),
	})
	if err != nil {
		return -1, Entry{}, fmt.Errorf("could not commit transaction: %w", err)
	}
	return int64(inserted.Carid), Entry{}, nil
}

func (p *PostgresRepository) DeleteByUUID(ctx context.Context, carID uuid.UUID) error {
	query := psql.Delete(
		dm.From(dbmodels.Cars.Name(ctx)),
		dbmodels.DeleteWhere.Cars.Caruuid.EQ(carID),
	)

	// Execute the query
	_, err := bob.Exec(ctx, p.db, query)
	if err != nil {
		return fmt.Errorf("could not execute delete: %w", err)
	}

	return nil
}

func (p *PostgresRepository) UpdateByUUID(ctx context.Context, carID uuid.UUID, car *models.CarCreationInput) (Entry, error) {
	query := psql.Update(
		um.From(dbmodels.Cars.Name(ctx)),
		dbmodels.UpdateWhere.Cars.Caruuid.EQ(carID),
		um.Set(dbmodels.CarColumns.Licenseplate, psql.Arg(car.LicensePlate)),
		um.Set(dbmodels.CarColumns.Make, psql.Arg(car.Make)),
		um.Set(dbmodels.CarColumns.Model, psql.Arg(car.Model)),
		um.Set(dbmodels.CarColumns.Color, psql.Arg(car.Color)),
	)

	// Execute the query
	_, err := bob.Exec(ctx, p.db, query)
	if err != nil {
		return Entry{}, fmt.Errorf("could not execute update: %w", err)
	}

	return Entry{}, nil
}

func (p *PostgresRepository) GetByUUID(ctx context.Context, carID uuid.UUID) (Entry, error) {
	query := psql.Select(
		sm.Columns(
			dbmodels.CarColumns.Licenseplate,
			dbmodels.CarColumns.Make,
			dbmodels.CarColumns.Model,
			dbmodels.CarColumns.Color,
			dbmodels.CarColumns.Carid,
			dbmodels.CarColumns.Userid,
		),
		sm.From(dbmodels.Cars.Name(ctx)),
		dbmodels.SelectWhere.Cars.Caruuid.EQ(carID),
	)
	result, err := bob.One(ctx, p.db, query, scan.StructMapper[dbmodels.Car]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return Entry{}, err
	}

	details := models.CarDetails{
		LicensePlate: result.Licenseplate,
		Make:         result.Make,
		Model:        result.Model,
		Color:        result.Color,
	}

	car := models.Car{
		Details: details,
	}

	return Entry{
		Car:        car,
		InternalID: int64(result.Carid),
		OwnerID:    int64(result.Userid),
	}, err
}

func (p *PostgresRepository) GetOwnerByUUID(ctx context.Context, carID uuid.UUID) (int64, error) {
	query := psql.Select(
		sm.Columns(dbmodels.CarColumns.Userid),
		sm.From(dbmodels.Cars.Name(ctx)),
		dbmodels.SelectWhere.Cars.Caruuid.EQ(carID),
	)
	result, err := bob.One(ctx, p.db, query, scan.SingleColumnMapper[int32])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return -1, err
	}

	return int64(result), err
}
