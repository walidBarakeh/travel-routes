// File: models.go
package main

import (
	"time"
)

// Location represents a geographical location
type Location struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Type      string  `json:"type"` // "city", "airport", "station"
	Code      string  `json:"code"` // IATA code for airports
	Country   string  `json:"country,omitempty"`
	PlaceID   string  `json:"place_id,omitempty"`
}

// TransportOption represents a transportation option
type TransportOption struct {
	Mode        string        `json:"mode"`        // "flight", "train", "bus", "taxi"
	From        Location      `json:"from"`
	To          Location      `json:"to"`
	Duration    time.Duration `json:"duration"`
	Price       float64       `json:"price"`
	Currency    string        `json:"currency"`
	Departure   time.Time     `json:"departure"`
	Arrival     time.Time     `json:"arrival"`
	Provider    string        `json:"provider"`
	BookingURL  string        `json:"booking_url,omitempty"`
}

// Route represents a complete travel route
type Route struct {
	Segments     []TransportOption `json:"segments"`
	TotalPrice   float64           `json:"total_price"`
	Currency     string            `json:"currency"`
	TotalTime    time.Duration     `json:"total_time"`
	Departure    time.Time         `json:"departure"`
	Arrival      time.Time         `json:"arrival"`
	Description  string            `json:"description"`
}

// AirportDistance represents an airport with its distance from origin
type AirportDistance struct {
	Airport  Location `json:"airport"`
	Distance float64  `json:"distance_km"`
}
