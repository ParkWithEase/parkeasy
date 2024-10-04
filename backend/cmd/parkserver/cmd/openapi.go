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
	Run: func(_ *cobra.Command, _ []string) {
		var config parkserver.Config
		api := config.NewHumaAPI()  
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
