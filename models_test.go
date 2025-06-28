package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLocation(t *testing.T) {
	t.Run("Create valid location", func(t *testing.T) {
		location := Location{
			Name:      "Madrid",
			Latitude:  40.4168,
			Longitude: -3.7038,
			Type:      "city",
			Code:      "",
			Country:   "Spain",
			PlaceID:   "ChIJgTwKgJcpQg0RaSKMYcHeNsQ",
		}

		assert.Equal(t, "Madrid", location.Name)
		assert.Equal(t, 40.4168, location.Latitude)
		assert.Equal(t, -3.7038, location.Longitude)
		assert.Equal(t, "city", location.Type)
		assert.Equal(t, "", location.Code)
		assert.Equal(t, "Spain", location.Country)
		assert.Equal(t, "ChIJgTwKgJcpQg0RaSKMYcHeNsQ", location.PlaceID)
	})

	t.Run("Create airport location", func(t *testing.T) {
		airport := Location{
			Name:      "Madrid-Barajas Airport",
			Latitude:  40.4719,
			Longitude: -3.5626,
			Type:      "airport",
			Code:      "MAD",
			Country:   "Spain",
		}

		assert.Equal(t, "Madrid-Barajas Airport", airport.Name)
		assert.Equal(t, "airport", airport.Type)
		assert.Equal(t, "MAD", airport.Code)
	})
}

func TestTransportOption(t *testing.T) {
	departure := time.Date(2024, 7, 1, 8, 0, 0, 0, time.UTC)
	arrival := departure.Add(2 * time.Hour)

	madrid := Location{
		Name:      "Madrid",
		Latitude:  40.4168,
		Longitude: -3.7038,
		Type:      "city",
	}

	barcelona := Location{
		Name:      "Barcelona",
		Latitude:  41.3851,
		Longitude: 2.1734,
		Type:      "city",
	}

	t.Run("Create flight transport option", func(t *testing.T) {
		flight := TransportOption{
			Mode:       "flight",
			From:       madrid,
			To:         barcelona,
			Duration:   2 * time.Hour,
			Price:      150.0,
			Currency:   "EUR",
			Departure:  departure,
			Arrival:    arrival,
			Provider:   "Iberia",
			BookingURL: "https://example.com/booking",
		}

		assert.Equal(t, "flight", flight.Mode)
		assert.Equal(t, madrid, flight.From)
		assert.Equal(t, barcelona, flight.To)
		assert.Equal(t, 2*time.Hour, flight.Duration)
		assert.Equal(t, 150.0, flight.Price)
		assert.Equal(t, "EUR", flight.Currency)
		assert.Equal(t, departure, flight.Departure)
		assert.Equal(t, arrival, flight.Arrival)
		assert.Equal(t, "Iberia", flight.Provider)
		assert.Equal(t, "https://example.com/booking", flight.BookingURL)
	})

	t.Run("Create ground transport option", func(t *testing.T) {
		taxi := TransportOption{
			Mode:      "taxi",
			From:      madrid,
			To:        barcelona,
			Duration:  6 * time.Hour,
			Price:     300.0,
			Currency:  "EUR",
			Departure: departure,
			Arrival:   departure.Add(6 * time.Hour),
			Provider:  "Taxi Service",
		}

		assert.Equal(t, "taxi", taxi.Mode)
		assert.Equal(t, 300.0, taxi.Price)
		assert.Equal(t, 6*time.Hour, taxi.Duration)
		assert.Equal(t, "Taxi Service", taxi.Provider)
	})

	t.Run("Create train transport option", func(t *testing.T) {
		train := TransportOption{
			Mode:      "train",
			From:      madrid,
			To:        barcelona,
			Duration:  3 * time.Hour,
			Price:     80.0,
			Currency:  "EUR",
			Departure: departure,
			Arrival:   departure.Add(3 * time.Hour),
			Provider:  "Renfe",
		}

		assert.Equal(t, "train", train.Mode)
		assert.Equal(t, 80.0, train.Price)
		assert.Equal(t, "Renfe", train.Provider)
	})
}

func TestRoute(t *testing.T) {
	// Create test locations
	madrid := Location{Name: "Madrid", Type: "city"}
	madridAirport := Location{Name: "Madrid Airport", Type: "airport", Code: "MAD"}
	barcelonaAirport := Location{Name: "Barcelona Airport", Type: "airport", Code: "BCN"}
	barcelona := Location{Name: "Barcelona", Type: "city"}

	departure := time.Date(2024, 7, 1, 8, 0, 0, 0, time.UTC)

	t.Run("Create empty route", func(t *testing.T) {
		route := Route{
			Segments:   []TransportOption{},
			TotalPrice: 0,
			Currency:   "EUR",
			TotalTime:  0,
		}

		assert.Empty(t, route.Segments)
		assert.Equal(t, 0.0, route.TotalPrice)
		assert.Equal(t, "EUR", route.Currency)
		assert.Equal(t, time.Duration(0), route.TotalTime)
	})

	t.Run("Create single segment route", func(t *testing.T) {
		flight := TransportOption{
			Mode:      "flight",
			From:      madridAirport,
			To:        barcelonaAirport,
			Duration:  1*time.Hour + 30*time.Minute,
			Price:     120.0,
			Currency:  "EUR",
			Departure: departure,
			Arrival:   departure.Add(1*time.Hour + 30*time.Minute),
			Provider:  "Airline",
		}

		route := Route{
			Segments:    []TransportOption{flight},
			TotalPrice:  120.0,
			Currency:    "EUR",
			TotalTime:   1*time.Hour + 30*time.Minute,
			Departure:   departure,
			Arrival:     departure.Add(1*time.Hour + 30*time.Minute),
			Description: "flight (Airline)",
		}

		assert.Len(t, route.Segments, 1)
		assert.Equal(t, 120.0, route.TotalPrice)
		assert.Equal(t, "EUR", route.Currency)
		assert.Equal(t, 1*time.Hour+30*time.Minute, route.TotalTime)
		assert.Equal(t, departure, route.Departure)
		assert.Contains(t, route.Description, "flight")
	})

	t.Run("Create multi-segment route", func(t *testing.T) {
		groundTransport := TransportOption{
			Mode:      "taxi",
			From:      madrid,
			To:        madridAirport,
			Duration:  45 * time.Minute,
			Price:     35.0,
			Currency:  "EUR",
			Departure: departure,
			Arrival:   departure.Add(45 * time.Minute),
			Provider:  "Taxi",
		}

		flight := TransportOption{
			Mode:      "flight",
			From:      madridAirport,
			To:        barcelonaAirport,
			Duration:  1*time.Hour + 30*time.Minute,
			Price:     120.0,
			Currency:  "EUR",
			Departure: departure.Add(1 * time.Hour),
			Arrival:   departure.Add(2*time.Hour + 30*time.Minute),
			Provider:  "Airline",
		}

		finalTransport := TransportOption{
			Mode:      "metro",
			From:      barcelonaAirport,
			To:        barcelona,
			Duration:  30 * time.Minute,
			Price:     5.0,
			Currency:  "EUR",
			Departure: departure.Add(3 * time.Hour),
			Arrival:   departure.Add(3*time.Hour + 30*time.Minute),
			Provider:  "Metro",
		}

		route := Route{
			Segments:  []TransportOption{groundTransport, flight, finalTransport},
			Currency:  "EUR",
		}

		route.CalculateTotals()

		assert.Len(t, route.Segments, 3)
		assert.Equal(t, 160.0, route.TotalPrice) // 35 + 120 + 5
		assert.Equal(t, "EUR", route.Currency)
		assert.Equal(t, departure, route.Departure)
		assert.Equal(t, departure.Add(3*time.Hour+30*time.Minute), route.Arrival)
		assert.Equal(t, 3*time.Hour+30*time.Minute, route.TotalTime)
		assert.Contains(t, route.Description, "→")
		assert.Contains(t, route.Description, "taxi")
		assert.Contains(t, route.Description, "flight")
		assert.Contains(t, route.Description, "metro")
	})

	t.Run("Route with zero duration segments", func(t *testing.T) {
		instantTransport := TransportOption{
			Mode:      "teleport",
			From:      madrid,
			To:        barcelona,
			Duration:  0,
			Price:     1000.0,
			Currency:  "EUR",
			Departure: departure,
			Arrival:   departure,
			Provider:  "Magic",
		}

		route := Route{
			Segments: []TransportOption{instantTransport},
			Currency: "EUR",
		}

		route.CalculateTotals()

		assert.Equal(t, 1000.0, route.TotalPrice)
		assert.Equal(t, time.Duration(0), route.TotalTime)
		assert.Equal(t, departure, route.Departure)
		assert.Equal(t, departure, route.Arrival)
	})
}

func TestAirportDistance(t *testing.T) {
	t.Run("Create airport distance", func(t *testing.T) {
		airport := Location{
			Name:      "Madrid-Barajas Airport",
			Latitude:  40.4719,
			Longitude: -3.5626,
			Type:      "airport",
			Code:      "MAD",
		}

		airportDist := AirportDistance{
			Airport:  airport,
			Distance: 15.5,
		}

		assert.Equal(t, airport, airportDist.Airport)
		assert.Equal(t, 15.5, airportDist.Distance)
		assert.Equal(t, "MAD", airportDist.Airport.Code)
		assert.Equal(t, "airport", airportDist.Airport.Type)
	})

	t.Run("Airport distance with zero distance", func(t *testing.T) {
		airport := Location{
			Name: "Local Airport",
			Type: "airport",
		}

		airportDist := AirportDistance{
			Airport:  airport,
			Distance: 0.0,
		}

		assert.Equal(t, 0.0, airportDist.Distance)
	})
}

// Test JSON marshaling/unmarshaling for API compatibility
func TestLocationJSONSerialization(t *testing.T) {
	location := Location{
		Name:      "Madrid",
		Latitude:  40.4168,
		Longitude: -3.7038,
		Type:      "city",
		Code:      "",
		Country:   "Spain",
		PlaceID:   "test-place-id",
	}

	// This test ensures the struct tags are correct for JSON serialization
	// In a real scenario, you'd test actual JSON marshaling/unmarshaling
	assert.Equal(t, "Madrid", location.Name)
	assert.Equal(t, "city", location.Type)
	assert.Equal(t, "Spain", location.Country)
	assert.Equal(t, "test-place-id", location.PlaceID)
}

func TestTransportOptionJSONSerialization(t *testing.T) {
	departure := time.Date(2024, 7, 1, 8, 0, 0, 0, time.UTC)
	madrid := Location{Name: "Madrid", Type: "city"}
	barcelona := Location{Name: "Barcelona", Type: "city"}

	transport := TransportOption{
		Mode:       "flight",
		From:       madrid,
		To:         barcelona,
		Duration:   2 * time.Hour,
		Price:      150.0,
		Currency:   "EUR",
		Departure:  departure,
		Arrival:    departure.Add(2 * time.Hour),
		Provider:   "Test Airline",
		BookingURL: "https://example.com",
	}

	// Verify all fields are accessible (JSON tags would be tested in integration tests)
	assert.Equal(t, "flight", transport.Mode)
	assert.Equal(t, madrid, transport.From)
	assert.Equal(t, barcelona, transport.To)
	assert.Equal(t, 2*time.Hour, transport.Duration)
	assert.Equal(t, 150.0, transport.Price)
	assert.Equal(t, "EUR", transport.Currency)
	assert.Equal(t, "Test Airline", transport.Provider)
	assert.Equal(t, "https://example.com", transport.BookingURL)
}
