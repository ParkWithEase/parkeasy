package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type CheckHealthCmd struct {
	Host    string        `arg:"" default:"localhost" help:"The hostname of the API server (default: ${default})."`
	Port    uint16        `env:"PORT" default:"8080" help:"The port of the API server (default: ${default})."`
	Timeout time.Duration `default:"30s" help:"Timeout for the query (default: ${default})."`
}

func (c *CheckHealthCmd) Run(ctx context.Context, l *zerolog.Logger, globals *Globals) error {
	log := globals.ConfigureZerolog(l).
		With().
		Str("command", "check-health").
		Logger()

	// Set a timeout for the operation
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	healthURL := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(c.Host, strconv.Itoa(int(c.Port))),
		Path:   "/healthz",
	}
	log.Debug().
		Stringer("url", &healthURL).
		Msg("querying health endpoint")
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, healthURL.String(), nil)
	if err != nil {
		return fmt.Errorf("could not create a new request: %w", err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("could not send HTTP request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode > 299 {
		e := log.Debug().
			Int("status_code", response.StatusCode)
		// Read at most 1 MiB from the server
		respBody, err := io.ReadAll(io.LimitReader(response.Body, 1024*1024))
		if err != nil {
			e.AnErr("read_error", err)
		}
		e.Bytes("response_body", respBody).Send()

		return errors.New("server is not healthy")
	}

	log.Info().Msg("server is healthy")
	return nil
}
