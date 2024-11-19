package demo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/dbmodels"
	"github.com/stephenafamo/bob"
)

type PostgresRepository struct {
	db bob.DB
}

func NewPostgres(db bob.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (p *PostgresRepository) Get(ctx context.Context) (string, error) {
	demoResult, err := dbmodels.Demos.Query(
		ctx, p.db,
	).One()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		return "", err
	}

	return demoResult.Demostring.GetOrZero(), nil
}
