package car

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmodels"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
)

type PostgresRepository struct {
	db bob.DB
}

func NewPostgres(db bob.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (p *PostgresRepository) Create(ctx context.Context, userID int64, car *models.CarCreationInput) (int64, Entry, error) {
	inserted, err := dbmodels.Cars.Insert(&dbmodels.CarSetter{
		Userid:       omit.From(userID),
		Licenseplate: omit.From(car.LicensePlate),
		Make:         omit.From(car.Make),
		Model:        omit.From(car.Model),
		Color:        omit.From(car.Color),
	}).One(ctx, p.db)
	if err != nil {
		return -1, Entry{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	details := models.CarDetails{
		LicensePlate: inserted.Licenseplate,
		Make:         inserted.Make,
		Model:        inserted.Model,
		Color:        inserted.Color,
	}

	insertedCar := models.Car{
		Details: details,
		ID:      inserted.Caruuid,
	}

	entry := Entry{
		Car:        insertedCar,
		InternalID: inserted.Carid,
		OwnerID:    inserted.Userid,
	}

	return inserted.Carid, entry, nil
}

func (p *PostgresRepository) DeleteByUUID(ctx context.Context, carID uuid.UUID) error {
	rowsAffected, err := dbmodels.Cars.Delete(
		dbmodels.DeleteWhere.Cars.Caruuid.EQ(carID),
	).Exec(ctx, p.db)
	if err != nil {
		return fmt.Errorf("could not execute delete: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (p *PostgresRepository) UpdateByUUID(ctx context.Context, carID uuid.UUID, car *models.CarCreationInput) (Entry, error) {
	result, err := dbmodels.Cars.Update(
		dbmodels.UpdateWhere.Cars.Caruuid.EQ(carID),
		dbmodels.CarSetter{
			Licenseplate: omit.From(car.LicensePlate),
			Make:         omit.From(car.Make),
			Model:        omit.From(car.Model),
			Color:        omit.From(car.Color),
		}.UpdateMod(),
		um.Returning(dbmodels.Cars.Columns()),
	).One(ctx, p.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Entry{}, ErrNotFound
		}
		return Entry{}, fmt.Errorf("could not execute update: %w", err)
	}

	details := models.CarDetails{
		LicensePlate: result.Licenseplate,
		Make:         result.Make,
		Model:        result.Model,
		Color:        result.Color,
	}

	return Entry{
		Car: models.Car{
			Details: details,
			ID:      carID,
		},
		InternalID: result.Carid,
		OwnerID:    result.Userid,
	}, err
}

func (p *PostgresRepository) GetByUUID(ctx context.Context, carID uuid.UUID) (Entry, error) {
	result, err := dbmodels.Cars.Query(
		sm.Columns(
			dbmodels.CarColumns.Licenseplate,
			dbmodels.CarColumns.Make,
			dbmodels.CarColumns.Model,
			dbmodels.CarColumns.Color,
			dbmodels.CarColumns.Carid,
			dbmodels.CarColumns.Userid,
		),
		dbmodels.SelectWhere.Cars.Caruuid.EQ(carID),
	).One(ctx, p.db)
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
		ID:      carID,
	}

	return Entry{
		Car:        car,
		InternalID: result.Carid,
		OwnerID:    result.Userid,
	}, err
}

func (p *PostgresRepository) GetOwnerByUUID(ctx context.Context, carID uuid.UUID) (int64, error) {
	result, err := dbmodels.Cars.Query(
		sm.Columns(
			dbmodels.CarColumns.Userid,
		),
		dbmodels.SelectWhere.Cars.Caruuid.EQ(carID),
	).One(ctx, p.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return -1, err
	}

	return result.Userid, err
}

func (p *PostgresRepository) GetMany(ctx context.Context, userID int64, limit int, after omit.Val[Cursor]) ([]Entry, error) {
	where := dbmodels.SelectWhere.Cars.Userid.EQ(userID)
	if cursor, ok := after.Get(); ok {
		where = psql.WhereAnd(where, dbmodels.SelectWhere.Cars.Carid.GT(cursor.ID))
	}

	entryCursor, err := dbmodels.Cars.Query(
		sm.Columns(dbmodels.Cars.Columns()),
		sm.OrderBy(dbmodels.CarColumns.Carid),
		where,
		sm.Limit(limit),
	).Cursor(ctx, p.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Entry{}, nil
		}
		return nil, err
	}
	defer entryCursor.Close()

	result := make([]Entry, 0, 8)
	for entryCursor.Next() {
		dbCar, err := entryCursor.Get()
		if err != nil { // if there's an error, just return what we already have
			break
		}
		result = append(result, Entry{
			Car: models.Car{
				Details: models.CarDetails{
					LicensePlate: dbCar.Licenseplate,
					Make:         dbCar.Make,
					Model:        dbCar.Model,
					Color:        dbCar.Color,
				},
				ID: dbCar.Caruuid,
			},
			InternalID: dbCar.Carid,
			OwnerID:    dbCar.Userid,
		})
	}
	return result, nil
}
