package cmd

import (
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerbosity(t *testing.T) {
	t.Run("set via environment", func(t *testing.T) {
		t.Setenv("VERBOSITY", "2")
		var cli RootCmd
		k, err := kong.New(&cli)
		require.NoError(t, err)
		_, err = k.Parse([]string{})
		require.NoError(t, err)

		assert.Equal(t, 2, cli.Verbose)
	})

	t.Run("set via flag", func(t *testing.T) {
		var cli RootCmd
		k, err := kong.New(&cli)
		require.NoError(t, err)
		_, err = k.Parse([]string{"-vvvv"})
		require.NoError(t, err)

		assert.Equal(t, 4, cli.Verbose)
	})
}
