// File: config.go
package main

import (
	"os"
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
	return Config{
		GoogleMapsAPIKey: os.Getenv("GOOGLE_MAPS_API_KEY"),
		AmadeusAPIKey:    os.Getenv("AMADEUS_API_KEY"),
		AmadeusSecret:    os.Getenv("AMADEUS_SECRET"),
		DefaultRadius:    300000, // 300km in meters
		MaxAirports:      10,
		MaxDistance:      500.0, // 500km max for ground transport
	}
}