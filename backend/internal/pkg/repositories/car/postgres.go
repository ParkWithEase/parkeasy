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
	"github.com/stephenafamo/scan"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (p *PostgresRepository) Create(ctx context.Context, userID int64, car *models.CarCreationInput) (uuid.UUID, error) {
	db := bob.NewDB(p.db)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	defer tx.Rollback()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("could not start a transaction: %v", err)
	}
	inserted, err := dbmodels.Cars.Insert(ctx, db, &dbmodels.CarSetter{
		Userid:             omit.From(userID),
		Licenseplate:       omit.From(car.LicensePlate),
		Make: 				omit.From(car.Make),
		Model:				omit.From(car.Model),
		Color:				omit.From(car.Color),
	})
	err = tx.Commit()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("could not commit transaction: %v", err)
	}
	return inserted.Caruuid, nil
}

func (p *PostgresRepository) DeleteByUUID(ctx context.Context, carID uuid.UUID) error {
	panic("unimplemented")

}

func (p *PostgresRepository) UpdateByUUID(ctx context.Context, carID uuid.UUID, car *models.CarCreationInput) (Entry, error) {
	panic("unimplemented")

}

func (p *PostgresRepository) GetByUUID(ctx context.Context, carID uuid.UUID) (Entry, error) {
	db := bob.NewDB(p.db)
	query := psql.Select(
		sm.Columns(dbmodels.CarColumns.Licenseplate, dbmodels.CarColumns.Make, dbmodels.CarColumns.Model, dbmodels.CarColumns.Color, dbmodels.CarColumns.Carid, dbmodels.CarColumns.Userid),
		sm.From(dbmodels.Auths.Name(ctx)),
		dbmodels.SelectWhere.Cars.Caruuid.EQ(carID),
	)
	result, err := bob.One(ctx, db, query, scan.StructMapper[dbmodels.Car]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return Entry{}, err
	}

	details := models.CarDetails{
		LicensePlate:   result.Licenseplate,
		Make: 			result.Make,
		Model:			result.Model,
		Color: 			result.Color,
	}

	car := models.Car{
		Details:		details,
	}

	return Entry{
		models.Car: car,
		InternalID: result.Carid,
		OwnerID: result.Userid,
	}, err
}

func (p *PostgresRepository) GetOwnerByUUID(ctx context.Context, carID uuid.UUID) (int64, error) {
	db := bob.NewDB(p.db)
	query := psql.Select(
		sm.Columns(dbmodels.CarColumns.Licenseplate, dbmodels.CarColumns.Make, dbmodels.CarColumns.Model, dbmodels.CarColumns.Color, dbmodels.CarColumns.Carid, dbmodels.CarColumns.Userid),
		sm.From(dbmodels.Auths.Name(ctx)),
		dbmodels.SelectWhere.Cars.Caruuid.EQ(carID),
	)
	result, err := bob.One(ctx, db, query, scan.StructMapper[dbmodels.Car]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return -1, err
	}
	
	return result.Userid, err
}