package auth

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

// Create implements Repository.
func (p *PostgresRepository) Create(ctx context.Context, email string, passwordHash models.HashedPassword) (uuid.UUID, error) {
	db := bob.NewDB(p.db)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	defer tx.Rollback()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("could not start a transaction: %v", err)
	}
	inserted, err := dbmodels.Auths.Insert(ctx, db, &dbmodels.AuthSetter{
		Email:        omit.From(email),
		Passwordhash: omit.From(string(passwordHash)),
	})
	err = tx.Commit()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("could not commit transaction: %v", err)
	}
	return inserted.Authuuid, nil
}

// Get implements Repository.
func (p *PostgresRepository) Get(ctx context.Context, id uuid.UUID) (Identity, error) {
	panic("unimplemented")
}

// GetByEmail implements Repository.
func (p *PostgresRepository) GetByEmail(ctx context.Context, email string) (Identity, error) {
	panic("unimplemented")
}

// UpdatePassword implements Repository.
func (p *PostgresRepository) UpdatePassword(ctx context.Context, authID uuid.UUID, newPassword models.HashedPassword) error {
	panic("unimplemented")
}
