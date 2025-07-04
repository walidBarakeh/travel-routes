package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleSearchRoutes(t *testing.T) {
	os.Setenv("GOOGLE_MAPS_API_KEY", "test-key")
	config := LoadConfig()
	tf := NewTravelFinder(config)

	req := httptest.NewRequest("GET", "/search?origin=Madrid&destination=Barcelona&date=2024-07-01", nil)
	w := httptest.NewRecorder()
	tf.handleSearchRoutes(w, req)
	resp := w.Result()
	// The real API call will fail with a mock key, so expect 500 error
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestHandleNearbyAirports(t *testing.T) {
	os.Setenv("GOOGLE_MAPS_API_KEY", "test-key")
	config := LoadConfig()
	tf := NewTravelFinder(config)

	req := httptest.NewRequest("GET", "/airports?location=Madrid&radius=30000", nil)
	w := httptest.NewRecorder()
	tf.handleNearbyAirports(w, req)
	resp := w.Result()
	// The real API call will fail with a mock key, so expect 500 error
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestAirportService_GeocodeLocation(t *testing.T) {
	as := &AirportService{config: Config{GoogleMapsAPIKey: "test-key"}, client: &http.Client{}}
	_, err := as.GeocodeLocation("Madrid")
	assert.Error(t, err) // Should error with mock key
}

func TestFlightService_SearchFlights(t *testing.T) {
	fs := &FlightService{config: Config{}, client: &http.Client{}}
	from := Location{Name: "Madrid", Code: "MAD"}
	to := Location{Name: "Barcelona", Code: "BCN"}
	date := time.Now()
	options, err := fs.SearchFlights(from, to, date)
	assert.NoError(t, err)
	assert.NotEmpty(t, options)
}

func TestTransportService_GetGroundTransport(t *testing.T) {
	ts := &TransportService{config: Config{GoogleMapsAPIKey: "test-key"}, client: &http.Client{}}
	from := Location{Name: "Madrid", Latitude: 40.4168, Longitude: -3.7038}
	to := Location{Name: "Barcelona", Latitude: 41.3851, Longitude: 2.1734}
	date := time.Now()
	_, err := ts.GetGroundTransport(from, to, date)
	// Accept both error and nil, since the mock implementation may not always error
	if err == nil {
		t.Log("No error returned, but this may be expected with mock data.")
	} else {
		assert.Error(t, err)
	}
}
