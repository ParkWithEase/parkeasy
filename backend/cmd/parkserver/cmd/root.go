package cmd

import (
	"context"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
)

type Globals struct {
	Verbose int `short:"v" env:"VERBOSITY" type:"counter" help:"Set log verbosity level, can be repeated multiple times."`
}

// Create a child logger with the options in globals applied
func (g *Globals) ConfigureZerolog(log *zerolog.Logger) zerolog.Logger {
	level := zerolog.InfoLevel
	switch {
	case g.Verbose > 1:
		level = zerolog.TraceLevel
	case g.Verbose > 0:
		level = zerolog.DebugLevel
	}
	return log.Level(level)
}

type RootCmd struct {
	OpenAPI OpenAPICmd `cmd:"" name:"openapi" help:"Dump OpenAPI schema."`
	Serve   ServeCmd   `cmd:"" default:"withargs" help:"Run API server."`
	Globals
}

// Returns a list of kong.Option that describes the program
func DefaultKongOptions() []kong.Option {
	return []kong.Option{
		kong.Name("parkserver"),
		kong.Description("The API server for ParkEasy app."),
		kong.ExplicitGroups([]kong.Group{
			{
				Key:         "db",
				Title:       "Database flags",
				Description: "Configures database connection",
			},
		}),
	}
}

// Bind all dependencies to a kong.Context.
//
// This should be used before the parsed result is run.
func (r *RootCmd) Bind(
	ctx context.Context,
	kctx *kong.Context,
	log *zerolog.Logger,
) {
	kctx.BindTo(ctx, (*context.Context)(nil))
	kctx.Bind(&r.Globals)
	kctx.Bind(log)
}
