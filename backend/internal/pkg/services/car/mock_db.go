package car

import (
    "context"

    "github.com/jackc/pgconn"
    "github.com/jackc/pgx/v4"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
)

// MockDB is a mock implementation of the DB interface
type MockDB struct {
    QueryFunc func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
    ExecFunc  func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
    QueryRowFunc func(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
}

// Implement the Query method
func (m *MockDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
    return m.QueryFunc(ctx, sql, args...)
}

// Implement the Exec method
func (m *MockDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
    return m.ExecFunc(ctx, sql, args...)
}

// Implement the QueryRow method
func (m *MockDB) QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row {
    return m.QueryRowFunc(ctx, sql, optionsAndArgs...)
}

// MockRows simulates pgx.Rows for testing
type MockRows struct {
    rows []models.Car
    idx  int
}

// Implement pgx.Rows interface
func (r *MockRows) Next() bool {
    return r.idx < len(r.rows)
}

func (r *MockRows) Scan(dest ...interface{}) error {
    if r.idx >= len(r.rows) {
        return nil
    }
    car := r.rows[r.idx]
    *dest[0].(*int) = car.CarID
    *dest[1].(*string) = car.LicensePlate
    *dest[2].(*string) = car.Make
    *dest[3].(*string) = car.Model
    *dest[4].(*string) = car.Color
    r.idx++
    return nil
}

func (r *MockRows) CommandTag() pgconn.CommandTag {
    return pgconn.CommandTag{}
}

func (r *MockRows) Close() {}

func (r *MockRows) Err() error {
    return nil
}

func (r *MockRows) Values() ([]interface{}, error) {
    return nil, nil
}