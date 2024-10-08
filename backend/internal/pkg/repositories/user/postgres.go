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

func (p *PostgresRepository) Create(ctx context.Context, id uuid.UUID, profile models.UserProfile) (int64, error) {
	db := bob.NewDB(p.db)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback if commit is not done
	inserted, err := dbmodels.Users.Insert(ctx, db, &dbmodels.UserSetter{
		Fullname: omit.From(profile.FullName),
		Email:    omit.From(profile.Email),
		Authuuid: omit.From(id),
	})
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		// TODO: Handle duplicate error
		return 0, fmt.Errorf("could not commit transaction: %w", err)
	}
	return int64(inserted.Userid), nil
}

// GetProfileById implements Repository.
func (p *PostgresRepository) GetProfileByID(ctx context.Context, id int64) (Profile, error) {
	db := bob.NewDB(p.db)
	query := psql.Select(
		sm.Columns(dbmodels.UserColumns.Email, dbmodels.UserColumns.Fullname, dbmodels.UserColumns.Authuuid, dbmodels.UserColumns.Userid),
		sm.From(dbmodels.Auths.Name(ctx)),
		dbmodels.SelectWhere.Users.Userid.EQ(id),
	)
	result, err := bob.One(ctx, db, query, scan.StructMapper[dbmodels.User]())
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
		ID:   int64(result.Userid),
	}, err
}

// GetProfileByAuth implements Repository.
func (p *PostgresRepository) GetProfileByAuth(ctx context.Context, id uuid.UUID) (Profile, error) {
	db := bob.NewDB(p.db)
	query := psql.Select(
		sm.Columns(dbmodels.UserColumns.Email, dbmodels.UserColumns.Fullname, dbmodels.UserColumns.Authuuid, dbmodels.UserColumns.Userid),
		sm.From(dbmodels.Auths.Name(ctx)),
		dbmodels.SelectWhere.Users.Authuuid.EQ(id),
	)
	result, err := bob.One(ctx, db, query, scan.StructMapper[dbmodels.User]())
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
		ID:   int64(result.Userid),
	}, err
}
