// File: travel_finder.go
package main

import (
	"fmt"
	"net/http"
	"sort"
	"time"
)

// TravelFinder is the main service
type TravelFinder struct {
	config       Config
	client       *http.Client
	airportSvc   *AirportService
	transportSvc *TransportService
	flightSvc    *FlightService
}

// NewTravelFinder creates a new travel finder instance
func NewTravelFinder(config Config) *TravelFinder {
	client := &http.Client{Timeout: 30 * time.Second}

	return &TravelFinder{
		config:       config,
		client:       client,
		airportSvc:   NewAirportService(config, client),
		transportSvc: NewTransportService(config, client),
		flightSvc:    NewFlightService(config, client),
	}
}

// FindRoutes finds all possible routes from origin to destination
func (tf *TravelFinder) FindRoutes(origin, destination string, travelDate time.Time) ([]Route, error) {
	var routes []Route

	// Step 1: Get origin coordinates
	originLocation, err := tf.airportSvc.GeocodeLocation(origin)
	if err != nil {
		return nil, fmt.Errorf("failed to geocode origin %s: %v", origin, err)
	}

	// Step 2: Get destination coordinates and airport info
	destinationLocation, err := tf.airportSvc.GeocodeLocation(destination)
	if err != nil {
		return nil, fmt.Errorf("failed to geocode destination %s: %v", destination, err)
	}

	// Find destination airport
	destAirports, err := tf.airportSvc.FindNearbyAirports(destinationLocation, 50000) // 50km radius for destination
	if err != nil || len(destAirports) == 0 {
		return nil, fmt.Errorf("no airports found near %s", destination)
	}
	destinationAirport := destAirports[0] // Use closest airport

	// Step 3: Find airports reachable from origin
	reachableAirports, err := tf.airportSvc.FindReachableAirports(originLocation)
	if err != nil {
		return nil, fmt.Errorf("error finding reachable airports: %v", err)
	}

	// Step 4: For each reachable airport, find routes
	for _, airport := range reachableAirports {
		// Get ground transport to airport
		groundTransport, err := tf.transportSvc.GetGroundTransport(originLocation, airport, travelDate)
		if err != nil {
			continue // Skip this airport if no ground transport available
		}

		// Check direct flights
		directFlights, err := tf.flightSvc.SearchFlights(airport, destinationAirport, travelDate)
		if err == nil {
			for _, flight := range directFlights {
				route := Route{
					Segments: []TransportOption{groundTransport, flight},
					Currency: "EUR",
				}
				route.CalculateTotals()
				routes = append(routes, route)
			}
		}

		// Check connecting flights
		connectingRoutes, err := tf.flightSvc.FindConnectingFlights(airport, destinationAirport, travelDate)
		if err == nil {
			for _, connectingRoute := range connectingRoutes {
				fullRoute := Route{
					Segments: append([]TransportOption{groundTransport}, connectingRoute.Segments...),
					Currency: "EUR",
				}
				fullRoute.CalculateTotals()
				routes = append(routes, fullRoute)
			}
		}
	}

	// Sort routes by total price
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].TotalPrice < routes[j].TotalPrice
	})

	return routes, nil
}
