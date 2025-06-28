package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (tf *TravelFinder) handleSearchRoutes(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")
	destination := r.URL.Query().Get("destination")
	dateStr := r.URL.Query().Get("date")

	if origin == "" || destination == "" || dateStr == "" {
		http.Error(w, "Missing required query parameters", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	routes, err := tf.FindRoutes(origin, destination, date)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error finding routes: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routes)
}

// handleNearbyAirports handles GET /airports?location=...&radius=...
func (tf *TravelFinder) handleNearbyAirports(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	radiusStr := r.URL.Query().Get("radius")

	if location == "" {
		http.Error(w, "Missing location parameter", http.StatusBadRequest)
		return
	}

	radius := tf.config.DefaultRadius
	if radiusStr != "" {
		if parsed, err := parseInt(radiusStr); err == nil {
			radius = parsed
		}
	}

	loc, err := tf.airportSvc.GeocodeLocation(location)
	if err != nil {
		http.Error(w, fmt.Sprintf("Geocoding failed: %v", err), http.StatusInternalServerError)
		return
	}

	airports, err := tf.airportSvc.FindNearbyAirports(loc, radius)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error finding airports: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(airports)
}

func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}
