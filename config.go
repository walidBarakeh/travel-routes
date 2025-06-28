// File: config.go
package main

import (
	"os"
	"strconv"
)

type Config struct {
	GoogleMapsAPIKey string
	AmadeusAPIKey    string
	AmadeusSecret    string
	DefaultRadius    int
	MaxAirports      int
	MaxDistance      float64
}

func LoadConfig() Config {
	defaultRadius := 300000
	if val := os.Getenv("DEFAULT_RADIUS"); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			defaultRadius = v
		}
	}

	maxAirports := 10
	if val := os.Getenv("MAX_AIRPORTS"); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			maxAirports = v
		}
	}

	maxDistance := 500.0
	if val := os.Getenv("MAX_DISTANCE"); val != "" {
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			maxDistance = v
		}
	}

	return Config{
		GoogleMapsAPIKey: os.Getenv("GOOGLE_MAPS_API_KEY"),
		AmadeusAPIKey:    os.Getenv("AMADEUS_API_KEY"),
		AmadeusSecret:    os.Getenv("AMADEUS_SECRET"),
		DefaultRadius:    defaultRadius,
		MaxAirports:      maxAirports,
		MaxDistance:      maxDistance,
	}
}