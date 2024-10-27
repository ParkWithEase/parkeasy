package geocoding

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var geocodioBaseURL = url.URL{
	Scheme: "https",
	Host:   "api.geocod.io",
	Path:   "/v1.7",
}

type Geocodio struct {
	client *http.Client
	apiKey string
}

func NewGeocodio(client *http.Client, apiKey string) *Geocodio {
	return &Geocodio{
		client: client,
		apiKey: apiKey,
	}
}

func (g *Geocodio) Geocode(address Address) ([]Result, error) {
	reqURL := geocodioBaseURL.JoinPath("geocode")
	queryParams := make(url.Values)
	queryParams.Set("api_key", g.apiKey)
	queryParams.Set("limit", "5")
	queryParams.Set("street", address.Street)
	queryParams.Set("city", address.City)
	queryParams.Set("state", address.State)
	queryParams.Set("postal_code", address.PostalCode)
	queryParams.Set("country", address.Country)
	reqURL.RawQuery = queryParams.Encode()

	resp, err := g.client.Get(reqURL.String())
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		var gerr GeocodioError
		err := json.NewDecoder(resp.Body).Decode(&gerr)
		if err != nil {
			gerr.Message = resp.Status
		}
		return nil, err
	}

	var apiResult struct {
		Results []struct {
			AddressComponents struct {
				Number          string `json:"number"`
				FormattedStreet string `json:"formatted_street"`
				City            string `json:"city"`
				State           string `json:"state"`
				Zip             string `json:"zip"`
				Country         string `json:"country"`
			} `json:"address_components"`
			FormattedAddress string `json:"formatted_address"`
			Location         struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			Accuracy float32 `json:"accuracy"`
		} `json:"results"`
	}
	err = json.NewDecoder(resp.Body).Decode(&apiResult)
	if err != nil {
		return nil, fmt.Errorf("could not decode result: %w", err)
	}

	result := make([]Result, 0, len(apiResult.Results))
	for _, r := range apiResult.Results {
		// geocod.io split street name and number, so rejoin them here
		street := r.AddressComponents.Number
		if street != "" {
			street += " " + r.AddressComponents.FormattedStreet
		}

		result = append(result, Result{
			Address: Address{
				Street:     street,
				City:       r.AddressComponents.City,
				State:      r.AddressComponents.State,
				PostalCode: r.AddressComponents.Zip,
				Country:    r.AddressComponents.Country,
			},
			FormattedAddress: r.FormattedAddress,
			Latitude:         r.Location.Lat,
			Longitude:        r.Location.Lng,
			Accuracy:         r.Accuracy,
		})
	}
	return result, nil
}

type GeocodioError struct {
	Message string `json:"error"`
}

func (e GeocodioError) Error() string {
	return e.Message
}
