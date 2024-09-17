package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGreetName(t *testing.T) {
	var greeter SimpleGreeting
	msg := greeter.Greet("Someone")

	require.Equal(t, "Hello, Someone!", msg.Message)
}
