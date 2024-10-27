package cmd

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"path"
	"strconv"

	"github.com/ParkWithEase/parkeasy/backend/internal/app/parkserver"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type DBUrlBuilder struct {
	Host     string `env:"HOST" help:"Postgres database host."`
	User     string `env:"USER" help:"Postgres database username."`
	Password string `env:"PASSWORD" help:"Postgres database password."`
	Name     string `env:"NAME" help:"Postgres database name."`
	Port     uint16 `env:"PORT" placeholder:"PORT" help:"Postgres database port."`
}

// Returns the composed connection string
func (ub *DBUrlBuilder) String() string {
	host := ub.Host
	if ub.Port != 0 {
		host = net.JoinHostPort(host, strconv.Itoa(int(ub.Port)))
	}
	var user *url.Userinfo
	if ub.Password != "" {
		user = url.UserPassword(ub.User, ub.Password)
	} else if ub.User != "" {
		user = url.User(ub.User)
	}
	u := url.URL{
		Scheme: "postgres",
		User:   user,
		Host:   host,
		Path:   path.Join("/", ub.Name),
	}
	return u.String()
}

type DBConfig struct {
	URL        *url.URL     `env:"URL" help:"Postgres database connection URL, preferred over components if provided."`
	Components DBUrlBuilder `embed:""`
}

// Returns the connection string
func (c *DBConfig) String() string {
	if c.URL == nil {
		return c.Components.String()
	}
	return c.URL.String()
}

type ServeCmd struct {
	APIPrefix  *url.URL `env:"API_PREFIX" placeholder:"PREFIX" help:"Specify the base prefix of the API server (example: http://localhost:8080/). If not specified, will be set to localhost at serve port."`
	CorsOrigin string   `placeholder:"ORIGIN" env:"CORS_ORIGIN" help:"Allow pages from ORIGIN to access the API server."`
	DB         DBConfig `embed:"" group:"db" prefix:"db-" envprefix:"DB_"`
	Port       uint16   `short:"p" placeholder:"PORT" env:"PORT" default:"8080" help:"Port to serve the server on (default: ${default})."`
	Insecure   bool     `env:"INSECURE" help:"Run in insecure mode for development (ie. CORS allow-all, HTTP cookies)."`
}

func (s *ServeCmd) getAPIPrefix() string {
	if s.APIPrefix == nil {
		prefix := net.JoinHostPort("localhost", strconv.Itoa(int(s.Port)))
		prefix = "http://" + prefix
		return prefix
	}
	if s.APIPrefix.Scheme == "" {
		s.APIPrefix.Scheme = "https"
	}
	return s.APIPrefix.String()
}

func (s *ServeCmd) Run(ctx context.Context, l *zerolog.Logger, globals *Globals) error {
	log := globals.ConfigureZerolog(l).
		With().
		Str("command", "serve").
		Logger()

	ctx = log.WithContext(ctx)

	pool, err := pgxpool.New(ctx, s.DB.String())
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	defer pool.Close()

	config := parkserver.Config{
		DBPool:     pool,
		APIPrefix:  s.getAPIPrefix(),
		Addr:       net.JoinHostPort("", strconv.Itoa(int(s.Port))),
		Insecure:   s.Insecure,
		CorsOrigin: s.CorsOrigin,
	}

	log.Info().Msg("running migrations")
	err = config.RunMigrations(ctx)
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	if config.Insecure {
		log.Warn().Msg("running in insecure mode")
	} else {
		log.Debug().Str("origin", s.CorsOrigin).Msg("allowing cross-origin requests")
	}
	log.Info().Uint16("port", s.Port).Msg("server started")
	err = config.ListenAndServe(ctx)
	if err != nil {
		return err
	}

	return nil
}
