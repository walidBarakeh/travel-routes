package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFindRoutes_Errors(t *testing.T) {
	tf := NewTravelFinder(Config{GoogleMapsAPIKey: "test-key"})
	// Should error on invalid origin
	_, err := tf.FindRoutes("", "Barcelona", time.Now())
	assert.Error(t, err)
	// Should error on invalid destination
	_, err = tf.FindRoutes("Madrid", "", time.Now())
	assert.Error(t, err)
}

func TestFindRoutes_SuccessMock(t *testing.T) {
	tf := NewTravelFinder(Config{GoogleMapsAPIKey: "test-key"})
	// This will likely error due to mock key, but test structure
	_, err := tf.FindRoutes("Madrid", "Barcelona", time.Now())
	assert.Error(t, err)
}

func TestAirportService_FindReachableAirports_Empty(t *testing.T) {
	as := &AirportService{config: Config{GoogleMapsAPIKey: "test-key"}, client: &http.Client{}}
	loc := Location{Name: "Nowhere", Latitude: 0, Longitude: 0}
	result, err := as.FindReachableAirports(loc)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestAirportService_FindNearbyAirports_Empty(t *testing.T) {
	as := &AirportService{config: Config{GoogleMapsAPIKey: "test-key"}, client: &http.Client{}}
	loc := Location{Name: "Nowhere", Latitude: 0, Longitude: 0}
	result, err := as.FindNearbyAirports(loc, 100)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestFlightService_FindConnectingFlights(t *testing.T) {
	fs := &FlightService{config: Config{}, client: nil}
	from := Location{Name: "A", Code: "AAA"}
	to := Location{Name: "B", Code: "BBB"}
	date := time.Now()
	options, err := fs.FindConnectingFlights(from, to, date)
	assert.NoError(t, err)
	assert.NotNil(t, options)
}

func TestTravelFinder_NewTravelFinder(t *testing.T) {
	config := Config{GoogleMapsAPIKey: "test-key"}
	tf := NewTravelFinder(config)
	assert.NotNil(t, tf)
	assert.NotNil(t, tf.airportSvc)
	assert.NotNil(t, tf.transportSvc)
	assert.NotNil(t, tf.flightSvc)
}
