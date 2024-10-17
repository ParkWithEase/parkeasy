package cmd

import (
	"os"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBURL(t *testing.T) {
	const testURL = "postgres://user:pass@somehost:5432/"
	t.Run("set via environment", func(t *testing.T) {
		t.Setenv("DB_URL", testURL)

		var cli ServeCmd
		k, err := kong.New(&cli)
		require.NoError(t, err)
		_, err = k.Parse([]string{})
		require.NoError(t, err)

		assert.Equal(t, testURL, cli.DB.String())
	})

	t.Run("set via command line", func(t *testing.T) {
		var cli ServeCmd
		k, err := kong.New(&cli)
		require.NoError(t, err)
		_, err = k.Parse([]string{"--db-url=" + testURL})
		require.NoError(t, err)

		assert.Equal(t, testURL, cli.DB.String())
	})

	t.Run("ignore components if set", func(t *testing.T) {
		var cli ServeCmd
		k, err := kong.New(&cli)
		require.NoError(t, err)
		_, err = k.Parse([]string{
			"--db-url=" + testURL,
			"--db-host=otherhost",
			"--db-user=someuser",
			"--db-password=somepass",
			"--db-port=42",
		})
		require.NoError(t, err)

		assert.Equal(t, testURL, cli.DB.String())
	})
}

func TestDBComponents(t *testing.T) {
	const (
		testUser     = "user"
		testPassword = "pass"
		testHost     = "somehost"
		testPort     = "5432"
		testName     = "somedb"

		testURL = "postgres://user:pass@somehost:5432/somedb"
	)

	t.Run("set via environment", func(t *testing.T) {
		t.Setenv("DB_HOST", testHost)
		t.Setenv("DB_PORT", testPort)
		t.Setenv("DB_NAME", testName)
		t.Setenv("DB_USER", testUser)
		t.Setenv("DB_PASSWORD", testPassword)

		_, ok := os.LookupEnv("DB_URL")
		assert.False(t, ok, "DB_URL environment variable should not be set")

		var cli ServeCmd
		k, err := kong.New(&cli)
		require.NoError(t, err)
		_, err = k.Parse([]string{})
		require.NoError(t, err)

		assert.Equal(t, testURL, cli.DB.String())
	})

	t.Run("set via flags", func(t *testing.T) {
		_, ok := os.LookupEnv("DB_URL")
		assert.False(t, ok, "DB_URL environment variable should not be set")

		var cli ServeCmd
		k, err := kong.New(&cli)
		require.NoError(t, err)
		_, err = k.Parse([]string{
			"--db-host=" + testHost,
			"--db-port=" + testPort,
			"--db-user=" + testUser,
			"--db-password=" + testPassword,
			"--db-name=" + testName,
		})
		require.NoError(t, err)

		assert.Equal(t, testURL, cli.DB.String())
	})
}
