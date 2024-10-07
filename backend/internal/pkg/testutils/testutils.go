package testutils

import (
	"os"
	"strconv"
	"testing"
)

// Mark the current test as integration and skip it when INTEGRATION environment is not set
func Integration(tb testing.TB) {
	tb.Helper()
	if ok, _ := strconv.ParseBool(os.Getenv("INTEGRATION")); !ok {
		tb.Skip("Integration test disabled, set INTEGRATION=1 environment variable to enable")
	}
}
