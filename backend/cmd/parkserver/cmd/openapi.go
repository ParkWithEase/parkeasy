package cmd

import (
	"fmt"

	"github.com/ParkWithEase/parkeasy/backend/internal/app/parkserver"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// openapiCmd represents the openapi command
var openapiCmd = &cobra.Command{
	Use:   "openapi",
	Short: "Print the OpenAPI spec",
	Run: func(cmd *cobra.Command, args []string) {
		var config parkserver.Config
		api := config.NewHumaApi()
		oapi, err := api.OpenAPI().YAML()
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot export OpenAPI spec")
		}
		fmt.Println(string(oapi))
	},
}

func init() {
	rootCmd.AddCommand(openapiCmd)
}
