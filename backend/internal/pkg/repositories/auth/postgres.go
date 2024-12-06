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
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

type PostgresRepository struct {
	db bob.DB
}

func NewPostgres(db bob.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// Create implements Repository.
func (p *PostgresRepository) Create(ctx context.Context, email string, passwordHash models.HashedPassword) (uuid.UUID, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback if commit is not done
	inserted, err := dbmodels.Auths.Insert(&dbmodels.AuthSetter{
		Email:        omit.From(email),
		Passwordhash: omit.From(string(passwordHash)),
	}).One(ctx, tx)
	if err != nil {
		// Handle duplicate error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				err = ErrDuplicateIdentity
			}
		}
		return uuid.UUID{}, err
	}
	err = tx.Commit()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("could not commit transaction: %w", err)
	}
	return inserted.Authuuid, nil
}

// Get implements Repository.
func (p *PostgresRepository) Get(ctx context.Context, id uuid.UUID) (Identity, error) {
	result, err := dbmodels.Auths.Query(
		sm.Columns(dbmodels.AuthColumns.Email, dbmodels.AuthColumns.Passwordhash, dbmodels.AuthColumns.Authuuid),
		dbmodels.SelectWhere.Auths.Authuuid.EQ(id),
	).One(ctx, p.db)
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
	}, nil
}

// GetByEmail implements Repository.
func (p *PostgresRepository) GetByEmail(ctx context.Context, email string) (Identity, error) {
	result, err := dbmodels.Auths.Query(
		sm.Columns(dbmodels.AuthColumns.Email, dbmodels.AuthColumns.Passwordhash, dbmodels.AuthColumns.Authuuid),
		dbmodels.SelectWhere.Auths.Email.EQ(email),
	).One(ctx, p.db)
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
	}, nil
}

// UpdatePassword implements Repository.
func (p *PostgresRepository) UpdatePassword(ctx context.Context, authID uuid.UUID, newPassword models.HashedPassword) error {
	// Execute the query
	rowsAffected, err := dbmodels.Auths.Update(
		dbmodels.UpdateWhere.Auths.Authuuid.EQ(authID),
		dbmodels.AuthSetter{
			Passwordhash: omit.From(string(newPassword)),
		}.UpdateMod(),
	).Exec(ctx, p.db)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrIdentityNotFound
	}

	return nil
}
