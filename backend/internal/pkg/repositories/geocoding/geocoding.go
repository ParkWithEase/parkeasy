package geocoding

// An address, splitted into components
type Address struct {
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

// Geocoding result
type Result struct {
	// Address as resolved by the provider
	Address Address
	// The formatted form of the address above
	FormattedAddress string
	Latitude         float64
	Longitude        float64
	// The accuracy of the result, ranging from 0 to 1.
	//
	// Results with accuracy >= 0.8 is considered accurate.
	Accuracy float32
}

type Geocoder interface {
	// Resolve an address into real location
	Geocode(Address) ([]Result, error)
}
