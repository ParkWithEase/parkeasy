package cmd

import (
	"fmt"
	"log"

	"github.com/ParkWithEase/parkeasy/backend/internal/app/parkserver"
	"github.com/spf13/cobra"
)

// openapiCmd represents the openapi command
var openapiCmd = &cobra.Command{
	Use:   "openapi",
	Short: "Print the OpenAPI spec",
	Run: func(cmd *cobra.Command, args []string) {
		api := parkserver.NewHumaApi()
		oapi, err := api.OpenAPI().YAML()
		if err != nil {
			log.Fatalf("Cannot export OpenAPI spec: %v", err)
		}
		fmt.Println(string(oapi))
	},
}

func init() {
	rootCmd.AddCommand(openapiCmd)
}
