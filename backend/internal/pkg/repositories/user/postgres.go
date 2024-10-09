package user

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
	"github.com/stephenafamo/bob/dialect/psql"
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

func (p *PostgresRepository) Create(ctx context.Context, id uuid.UUID, profile models.UserProfile) (int64, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback if commit is not done
	inserted, err := dbmodels.Users.Insert(ctx, tx, &dbmodels.UserSetter{
		Fullname: omit.From(profile.FullName),
		Email:    omit.From(profile.Email),
		Authuuid: omit.From(id),
	})
	if err != nil {
		// Handle duplicate error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				err = ErrProfileExists
			}
		}
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("could not commit transaction: %w", err)
	}
	return inserted.Userid, nil
}

// GetProfileById implements Repository.
func (p *PostgresRepository) GetProfileByID(ctx context.Context, id int64) (Profile, error) {
	query := psql.Select(
		sm.Columns(dbmodels.UserColumns.Email, dbmodels.UserColumns.Fullname, dbmodels.UserColumns.Authuuid, dbmodels.UserColumns.Userid),
		sm.From(dbmodels.Users.Name(ctx)),
		dbmodels.SelectWhere.Users.Userid.EQ(id),
	)
	result, err := bob.One(ctx, p.db, query, scan.StructMapper[dbmodels.User]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrUnknownID
		}
		return Profile{}, err
	}
	return Profile{
		UserProfile: models.UserProfile{
			FullName: result.Fullname,
			Email:    result.Email,
		},
		Auth: result.Authuuid,
		ID:   result.Userid,
	}, nil
}

// GetProfileByAuth implements Repository.
func (p *PostgresRepository) GetProfileByAuth(ctx context.Context, id uuid.UUID) (Profile, error) {
	query := psql.Select(
		sm.Columns(
			dbmodels.UserColumns.Email,
			dbmodels.UserColumns.Fullname,
			dbmodels.UserColumns.Authuuid,
			dbmodels.UserColumns.Userid,
		),
		sm.From(dbmodels.Users.Name(ctx)),
		dbmodels.SelectWhere.Users.Authuuid.EQ(id),
	)
	result, err := bob.One(ctx, p.db, query, scan.StructMapper[dbmodels.User]())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrUnknownID
		}
		return Profile{}, err
	}
	return Profile{
		UserProfile: models.UserProfile{
			FullName: result.Fullname,
			Email:    result.Email,
		},
		Auth: result.Authuuid,
		ID:   result.Userid,
	}, nil
}
