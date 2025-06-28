package main

// Google Maps API Response structures
type GoogleDirectionsResponse struct {
	Routes []struct {
		Legs []struct {
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration"`
			Distance struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"distance"`
			StartAddress string `json:"start_address"`
			EndAddress   string `json:"end_address"`
		} `json:"legs"`
	} `json:"routes"`
	Status string `json:"status"`
}

// Google Places API Response structures
type GooglePlacesResponse struct {
	Results []struct {
		Name     string `json:"name"`
		PlaceID  string `json:"place_id"`
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
		Types           []string `json:"types"`
		BusinessStatus  string   `json:"business_status,omitempty"`
		Rating          float64  `json:"rating,omitempty"`
		PriceLevel      int      `json:"price_level,omitempty"`
	} `json:"results"`
	Status string `json:"status"`
}

// Geocoding API Response
type GoogleGeocodingResponse struct {
	Results []struct {
		FormattedAddress string `json:"formatted_address"`
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
	} `json:"results"`
	Status string `json:"status"`
}
