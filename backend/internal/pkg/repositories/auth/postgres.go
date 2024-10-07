package auth

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

// Create implements Repository.
func (p *PostgresRepository) Create(ctx context.Context, email string, passwordHash models.HashedPassword) (uuid.UUID, error) {
	db := bob.NewDB(p.db)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback if commit is not done
	inserted, err := dbmodels.Auths.Insert(ctx, db, &dbmodels.AuthSetter{
		Email:        omit.From(email),
		Passwordhash: omit.From(string(passwordHash)),
	})
	if err != nil {
		// TODO: Handle duplicate error
		return uuid.UUID{}, err
	}
	err = tx.Commit()
	if err != nil {
		// TODO: Handle duplicate error
		return uuid.UUID{}, fmt.Errorf("could not commit transaction: %w", err)
	}
	return inserted.Authuuid, nil
}

// Get implements Repository.
func (p *PostgresRepository) Get(ctx context.Context, id uuid.UUID) (Identity, error) {
	db := bob.NewDB(p.db)
	query := psql.Select(
		sm.Columns(dbmodels.AuthColumns.Email, dbmodels.AuthColumns.Passwordhash, dbmodels.AuthColumns.Authuuid),
		sm.From(dbmodels.Auths.Name(ctx)),
		dbmodels.SelectWhere.Auths.Authuuid.EQ(id),
	)
	result, err := bob.One(ctx, db, query, scan.StructMapper[dbmodels.Auth]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrIdentityNotFound
		}
		return Identity{}, err
	}
	return Identity{
		Email:        result.Email,
		PasswordHash: models.HashedPassword(result.Passwordhash),
		ID:           result.Authuuid,
	}, err
}

// GetByEmail implements Repository.
func (p *PostgresRepository) GetByEmail(_ context.Context, _ string) (Identity, error) {
	panic("unimplemented")
}

// UpdatePassword implements Repository.
func (p *PostgresRepository) UpdatePassword(_ context.Context, _ uuid.UUID, _ models.HashedPassword) error {
	panic("unimplemented")
}
