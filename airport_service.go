package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

type AirportService struct {
	config Config
	client *http.Client
}

func NewAirportService(config Config, client *http.Client) *AirportService {
	return &AirportService{
		config: config,
		client: client,
	}
}

// GeocodeLocation converts a location name to coordinates
func (as *AirportService) GeocodeLocation(locationName string) (Location, error) {
	baseURL := "https://maps.googleapis.com/maps/api/geocode/json"
	params := url.Values{}
	params.Add("address", locationName)
	params.Add("key", as.config.GoogleMapsAPIKey)

	resp, err := as.client.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Location{}, err
	}

	var geocodeResp GoogleGeocodingResponse
	if err := json.Unmarshal(body, &geocodeResp); err != nil {
		return Location{}, err
	}

	if geocodeResp.Status != "OK" || len(geocodeResp.Results) == 0 {
		return Location{}, fmt.Errorf("geocoding failed for %s", locationName)
	}

	result := geocodeResp.Results[0]

	// Extract country from address components
	country := ""
	for _, component := range result.AddressComponents {
		for _, componentType := range component.Types {
			if componentType == "country" {
				country = component.LongName
				break
			}
		}
	}

	return Location{
		Name:      locationName,
		Latitude:  result.Geometry.Location.Lat,
		Longitude: result.Geometry.Location.Lng,
		Type:      "city",
		Country:   country,
	}, nil
}

// FindReachableAirports finds airports reachable from a given location
func (as *AirportService) FindReachableAirports(origin Location) ([]Location, error) {
	airports, err := as.FindNearbyAirports(origin, as.config.DefaultRadius)
	if err != nil {
		return nil, fmt.Errorf("failed to search nearby airports: %v", err)
	}

	var reachableAirports []AirportDistance
	for _, airport := range airports {
		distance := CalculateDistance(origin.Latitude, origin.Longitude, airport.Latitude, airport.Longitude)

		if distance <= as.config.MaxDistance {
			reachableAirports = append(reachableAirports, AirportDistance{
				Airport:  airport,
				Distance: distance,
			})
		}
	}

	// Sort by distance (closest first)
	sort.Slice(reachableAirports, func(i, j int) bool {
		return reachableAirports[i].Distance < reachableAirports[j].Distance
	})

	// Return top airports
	var result []Location
	maxAirports := as.config.MaxAirports
	if len(reachableAirports) < maxAirports {
		maxAirports = len(reachableAirports)
	}

	for i := 0; i < maxAirports; i++ {
		result = append(result, reachableAirports[i].Airport)
	}

	log.Printf("Found %d reachable airports from %s", len(result), origin.Name)
	for i, airport := range result {
		if i < len(reachableAirports) {
			log.Printf("  %d. %s (%s) - %.1f km", i+1, airport.Name, airport.Code, reachableAirports[i].Distance)
		}
	}

	return result, nil
}

// FindNearbyAirports searches for airports near a location using Google Places API
func (as *AirportService) FindNearbyAirports(origin Location, radiusMeters int) ([]Location, error) {
	baseURL := "https://maps.googleapis.com/maps/api/place/nearbysearch/json"
	params := url.Values{}
	params.Add("location", fmt.Sprintf("%f,%f", origin.Latitude, origin.Longitude))
	params.Add("radius", strconv.Itoa(radiusMeters))
	params.Add("type", "airport")
	params.Add("key", as.config.GoogleMapsAPIKey)

	resp, err := as.client.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var placesResp GooglePlacesResponse
	if err := json.Unmarshal(body, &placesResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if placesResp.Status != "OK" {
		log.Printf("Google Places API returned status: %s", placesResp.Status)
		return []Location{}, nil
	}

	var airports []Location
	for _, place := range placesResp.Results {
		iataCode := ExtractIATACode(place.Name)

		airport := Location{
			Name:      place.Name,
			Latitude:  place.Geometry.Location.Lat,
			Longitude: place.Geometry.Location.Lng,
			Type:      "airport",
			Code:      iataCode,
			PlaceID:   place.PlaceID,
		}
		airports = append(airports, airport)
	}

	return airports, nil
}
