package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ParkWithEase/parkeasy/backend/cmd/parkserver/cmd"
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
)

func run(ctx context.Context, log *zerolog.Logger) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var cli cmd.RootCmd
	kctx := kong.Parse(
		&cli,
		cmd.DefaultKongOptions()...,
	)
	cli.Bind(ctx, kctx, log)
	return kctx.Run()
}

func main() {
	ctx := context.Background()

	logOutput := zerolog.NewConsoleWriter()
	log := zerolog.New(logOutput).
		With().
		Timestamp().
		Logger().
		Level(zerolog.InfoLevel)

	if err := run(ctx, &log); err != nil {
		log.Fatal().Err(err).Send()
	}
}
