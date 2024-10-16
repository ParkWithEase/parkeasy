package cmd

import (
	"fmt"
	"os"

	"github.com/ParkWithEase/parkeasy/backend/internal/app/parkserver"
)

type OpenAPICmd struct {
	Output    string `placeholder:"FILE" help:"Output to file instead of stdout."`
	JSON      bool   `help:"Output JSON instead of YAML."`
	Downgrade bool   `help:"Downgrade spec to version 3.0."`
}

func (c *OpenAPICmd) Run() error {
	var config parkserver.Config
	api := config.NewHumaAPI()
	var oapi []byte
	var err error
	if c.JSON {
		if c.Downgrade {
			oapi, err = api.OpenAPI().Downgrade()
		} else {
			oapi, err = api.OpenAPI().MarshalJSON()
		}
	} else {
		if c.Downgrade {
			oapi, err = api.OpenAPI().DowngradeYAML()
		} else {
			oapi, err = api.OpenAPI().YAML()
		}
	}
	if err != nil {
		return err
	}
	if c.Output != "" {
		if err := os.WriteFile(c.Output, oapi, 0o666); err != nil {
			return err
		}
	} else {
		fmt.Printf("%s", oapi)
	}
	return nil
}
