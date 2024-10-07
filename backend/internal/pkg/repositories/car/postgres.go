package car

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmodels"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
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
	panic("unimplemented")

}

func (p *PostgresRepository) GetOwnerByUUID(ctx context.Context, carID uuid.UUID) (int64, error) {
	panic("unimplemented")

}
