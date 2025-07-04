// File: flight_service.go
package main

import (
	"fmt"
	"net/http"
	"time"
)

type FlightService struct {
	config Config
	client *http.Client
}

func NewFlightService(config Config, client *http.Client) *FlightService {
	return &FlightService{
		config: config,
		client: client,
	}
}

// SearchFlights searches for direct flights
func (fs *FlightService) SearchFlights(from, to Location, date time.Time) ([]TransportOption, error) {
	// Check if direct route is likely available
	if !fs.isDirectRouteAvailable(from.Code, to.Code) {
		return []TransportOption{}, fmt.Errorf("no direct flights available")
	}

	// Mock flight data - replace with real API call
	flight := TransportOption{
		Mode:      "flight",
		From:      from,
		To:        to,
		Duration:  4*time.Hour + 30*time.Minute,
		Price:     fs.estimateFlightPrice(from.Code, to.Code, "direct"),
		Currency:  "EUR",
		Departure: date.Add(2 * time.Hour),
		Arrival:   date.Add(6*time.Hour + 30*time.Minute),
		Provider:  "Airlines",
	}

	return []TransportOption{flight}, nil
}

// FindConnectingFlights finds flights with connections
func (fs *FlightService) FindConnectingFlights(origin, destination Location, date time.Time) ([]Route, error) {
	var routes []Route

	// Major European hubs that typically have good connections
	hubs := []string{"LHR", "CDG", "FRA", "AMS", "FCO", "MUC", "VIE", "ZUR", "IST"}

	for _, hubCode := range hubs[:3] { // Limit to 3 hubs for demo
		route, err := fs.createConnectingRoute(origin, destination, hubCode, date)
		if err == nil {
			routes = append(routes, route)
		}
	}

	return routes, nil
}

func (fs *FlightService) createConnectingRoute(origin, destination Location, hubCode string, date time.Time) (Route, error) {
	hubAirport := Location{
		Name: fmt.Sprintf("%s Hub Airport", hubCode),
		Code: hubCode,
		Type: "airport",
	}

	// First leg: Origin to Hub
	firstLeg := TransportOption{
		Mode:      "flight",
		From:      origin,
		To:        hubAirport,
		Duration:  2*time.Hour + 30*time.Minute,
		Price:     fs.estimateFlightPrice(origin.Code, hubCode, "connecting"),
		Currency:  "EUR",
		Departure: date.Add(2 * time.Hour),
		Arrival:   date.Add(4*time.Hour + 30*time.Minute),
		Provider:  "Airlines",
	}

	// Second leg: Hub to Destination (with layover)
	secondLeg := TransportOption{
		Mode:      "flight",
		From:      hubAirport,
		To:        destination,
		Duration:  4 * time.Hour,
		Price:     fs.estimateFlightPrice(hubCode, destination.Code, "connecting"),
		Currency:  "EUR",
		Departure: date.Add(6*time.Hour + 30*time.Minute), // 2-hour layover
		Arrival:   date.Add(10*time.Hour + 30*time.Minute),
		Provider:  "Airlines",
	}

	route := Route{
		Segments: []TransportOption{firstLeg, secondLeg},
		Currency: "EUR",
	}
	route.CalculateTotals()

	return route, nil
}

func (fs *FlightService) isDirectRouteAvailable(fromCode, toCode string) bool {
	// This would be replaced with real route availability checking
	// For now, assume major airports have better connectivity
	majorAirports := map[string]bool{
		"MAD": true, "BCN": true, "LHR": true, "CDG": true,
		"FRA": true, "AMS": true, "FCO": true, "MUC": true,
	}

	return majorAirports[fromCode]
}

func (fs *FlightService) estimateFlightPrice(fromCode, toCode, routeType string) float64 {
	// Base price calculation - replace with real pricing API
	basePrice := 200.0

	if routeType == "connecting" {
		basePrice *= 0.8 // Connecting flights often cheaper
	}

	// Add some variation based on route popularity
	if fromCode == "MAD" || fromCode == "BCN" {
		basePrice += 50 // Premium for major airports
	}

	return basePrice
}
