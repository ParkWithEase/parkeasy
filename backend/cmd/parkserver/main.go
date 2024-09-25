package main

import (
	"os"

	"github.com/ParkWithEase/parkeasy/backend/cmd/parkserver/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cmd.Execute()
}
