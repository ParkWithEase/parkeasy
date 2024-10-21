package timeunit

import (
	"context"
	"testing"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	// "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	// "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	// "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/testutils"
	// "github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stephenafamo/bob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestPostgresIntegration(t *testing.T) {
	t.Parallel()

	testutils.Integration(t)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	container, connString := testutils.CreatePostgresContainer(ctx, t)
	t.Cleanup(func() { _ = container.Terminate(ctx) })
	testutils.RunMigrations(t, connString)

	pool, err := pgxpool.New(ctx, connString)
	require.NoError(t, err, "could not connect to db")
	t.Cleanup(func() { pool.Close() })
	db := bob.NewDB(stdlib.OpenDBFromPool(pool))

	repo := NewPostgres(db)
	// listingrepo = listing.NewPostgres(db)

	var sampleListingId int64 = 1

	pool.Reset()
	snapshotErr := container.Snapshot(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
	require.NoError(t, snapshotErr, "could not snapshot db")

	t.Run("basic add & get & delete", func(t *testing.T) {
		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		sampleTimeSlots := []models.TimeSlot{
			{
				Date:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Units: []int16{1, 2, 3, 4, 5},
			},
			{
				Date:  time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
				Units: []int16{6, 7, 8},
			},
			{
				Date:  time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
				Units: []int16{9, 10},
			},
		}

		// Testing create
		createEntry, err := repo.Create(ctx, sampleTimeSlots, sampleListingId)
		require.NoError(t, err)
		assert.Equal(t, sampleListingId, createEntry.ListingId)
		assert.Equal(t, int64(0), createEntry.BookingId)
		assert.Equal(t, sampleTimeSlots, createEntry.TimeSlots)

		// Testing get
		getEntry, err := repo.GetByListingID(ctx, sampleListingId)
		require.NoError(t, err)
		assert.Equal(t, sampleTimeSlots, getEntry.TimeSlots)

		// Testing delete
		err = repo.DeleteByListingID(ctx, sampleListingId, sampleTimeSlots)
		require.NoError(t, err)
	})

	t.Run("duplicate address creation should fail", func(t *testing.T) {

		t.Cleanup(func() {
			err := container.Restore(ctx, postgres.WithSnapshotName(testutils.PostgresSnapshotName))
			require.NoError(t, err, "could not restore db")

			// clear all idle connections
			// required since Restore() deletes the current DB
			pool.Reset()
		})

		sampleTimeSlots := []models.TimeSlot{
			{
				Date:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Units: []int16{1, 2, 3, 4, 5},
			},
			{
				Date:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Units: []int16{1, 2, 3, 4, 5},
			},
		}

		// Attempt to create duplicate timeunits
		_, err := repo.Create(ctx, sampleTimeSlots, sampleListingId)
		if assert.Error(t, err, "Creating two identical time units should fail") {
			assert.ErrorIs(t, err, ErrDuplicatedTimeUnit)
		}
	})

	t.Run("get non-existent", func(t *testing.T) {
		_, err := repo.GetByListingID(ctx, int64(-1))
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrNotFound)
		}
	})

}
