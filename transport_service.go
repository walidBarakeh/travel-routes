// File: transport_service.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type TransportService struct {
	config Config
	client *http.Client
}

func NewTransportService(config Config, client *http.Client) *TransportService {
	return &TransportService{
		config: config,
		client: client,
	}
}

// GetGroundTransport gets ground transportation options
func (ts *TransportService) GetGroundTransport(from, to Location, date time.Time) (TransportOption, error) {
	// Try public transit first
	transitOption, err := ts.getPublicTransit(from, to, date)
	if err == nil {
		return transitOption, nil
	}

	// Fallback to taxi estimate
	return ts.getTaxiEstimate(from, to, date)
}

func (ts *TransportService) getPublicTransit(from, to Location, date time.Time) (TransportOption, error) {
	baseURL := "https://maps.googleapis.com/maps/api/directions/json"
	params := url.Values{}
	params.Add("origin", fmt.Sprintf("%f,%f", from.Latitude, from.Longitude))
	params.Add("destination", fmt.Sprintf("%f,%f", to.Latitude, to.Longitude))
	params.Add("mode", "transit")
	params.Add("departure_time", fmt.Sprintf("%d", date.Unix()))
	params.Add("key", ts.config.GoogleMapsAPIKey)

	resp, err := ts.client.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return TransportOption{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TransportOption{}, err
	}

	var directionsResp GoogleDirectionsResponse
	if err := json.Unmarshal(body, &directionsResp); err != nil {
		return TransportOption{}, err
	}

	if directionsResp.Status != "OK" || len(directionsResp.Routes) == 0 || len(directionsResp.Routes[0].Legs) == 0 {
		return TransportOption{}, fmt.Errorf("no transit routes found")
	}

	leg := directionsResp.Routes[0].Legs[0]
	duration := time.Duration(leg.Duration.Value) * time.Second
	price := ts.estimateTransportPrice(leg.Distance.Value, "transit")

	return TransportOption{
		Mode:      "public_transport",
		From:      from,
		To:        to,
		Duration:  duration,
		Price:     price,
		Currency:  "EUR",
		Departure: date,
		Arrival:   date.Add(duration),
		Provider:  "Public Transport",
	}, nil
}

func (ts *TransportService) getTaxiEstimate(from, to Location, date time.Time) (TransportOption, error) {
	distance := CalculateDistance(from.Latitude, from.Longitude, to.Latitude, to.Longitude)
	duration := time.Duration(distance/60) * time.Hour // Assume 60km/h average
	price := ts.estimateTransportPrice(int(distance*1000), "taxi")

	return TransportOption{
		Mode:      "taxi",
		From:      from,
		To:        to,
		Duration:  duration,
		Price:     price,
		Currency:  "EUR",
		Departure: date,
		Arrival:   date.Add(duration),
		Provider:  "Taxi",
	}, nil
}

func (ts *TransportService) estimateTransportPrice(distanceMeters int, mode string) float64 {
	distanceKm := float64(distanceMeters) / 1000

	switch mode {
	case "transit", "public_transport":
		// Base fare + distance-based pricing
		baseFare := 2.0
		if distanceKm <= 10 {
			return baseFare + (distanceKm * 0.15)
		} else if distanceKm <= 50 {
			return baseFare + (10 * 0.15) + ((distanceKm - 10) * 0.12)
		} else {
			return baseFare + (10 * 0.15) + (40 * 0.12) + ((distanceKm - 50) * 0.10)
		}
	case "taxi":
		// Standard taxi rates: base + per km
		baseFare := 3.0
		return baseFare + (distanceKm * 1.2)
	default:
		return distanceKm * 0.2
	}
}
