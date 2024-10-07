package auth

import (
	"testing"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/testutils"
)

func TestPostgresCreateIntegration(t *testing.T) {
	testutils.Integration(t)
}
