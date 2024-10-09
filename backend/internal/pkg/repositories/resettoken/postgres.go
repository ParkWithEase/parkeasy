package resettoken

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmodels"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
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

func (p *PostgresRepository) Create(ctx context.Context, authID uuid.UUID, token Token) error {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not start a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // Default to rollback if commit is not done

	// Delete any existing tokens
	delQuery := psql.Delete(
		dm.From(dbmodels.Resettokens.Name(ctx)),
		dbmodels.DeleteWhere.Resettokens.Authuuid.EQ(authID),
	)
	_, err = bob.Exec(ctx, tx, delQuery)
	if err != nil {
		return fmt.Errorf("unable to delete existing token: %w", err)
	}

	// Create a new token
	_, err = dbmodels.Resettokens.Insert(ctx, tx, &dbmodels.ResettokenSetter{
		Token:    omit.From(string(token)),
		Authuuid: omit.From(authID),
		Expiry:   omit.From(time.Now().Add(15 * time.Minute)),
	})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (p *PostgresRepository) Get(ctx context.Context, token Token) (uuid.UUID, error) {
	query := psql.Select(
		sm.Columns(dbmodels.ResettokenColumns.Authuuid),
		sm.From(dbmodels.Resettokens.Name(ctx)),
		dbmodels.SelectWhere.Resettokens.Token.EQ(string(token)),
	)

	result, err := bob.One(ctx, p.db, query, scan.SingleColumnMapper[uuid.UUID])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrInvalidToken
		}
		return uuid.Nil, err
	}

	return result, nil
}

func (p *PostgresRepository) Delete(ctx context.Context, token Token) error {
	query := psql.Delete(
		dm.From(dbmodels.Resettokens.Name(ctx)),
		dbmodels.DeleteWhere.Resettokens.Token.EQ(string(token)),
	)

	_, err := bob.Exec(ctx, p.db, query)
	if err != nil {
		return fmt.Errorf("unable to delete existing token: %w", err)
	}

	return nil
}
