package main

import (
	"fmt"
	"math"
	"strings"
)

// CalculateDistance calculates distance between two points using Haversine formula
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Earth radius in kilometers

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Haversine formula
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// ExtractIATACode attempts to extract IATA code from airport name
func ExtractIATACode(airportName string) string {
	// Try to extract 3-letter code from parentheses
	if strings.Contains(airportName, "(") && strings.Contains(airportName, ")") {
		start := strings.Index(airportName, "(")
		end := strings.Index(airportName, ")")
		if end > start+1 {
			potential := strings.TrimSpace(airportName[start+1 : end])
			if len(potential) == 3 && strings.ToUpper(potential) == potential {
				return potential
			}
		}
	}

	// Try to match with known patterns
	name := strings.ToLower(airportName)

	// Common patterns for extracting codes
	if strings.Contains(name, "madrid") || strings.Contains(name, "barajas") {
		return "MAD"
	}
	if strings.Contains(name, "barcelona") || strings.Contains(name, "el prat") {
		return "BCN"
	}
	if strings.Contains(name, "málaga") || strings.Contains(name, "malaga") {
		return "AGP"
	}
	if strings.Contains(name, "sevilla") || strings.Contains(name, "seville") {
		return "SVQ"
	}
	if strings.Contains(name, "valencia") {
		return "VLC"
	}
	if strings.Contains(name, "bilbao") {
		return "BIO"
	}
	if strings.Contains(name, "granada") {
		return "GRX"
	}
	if strings.Contains(name, "tel aviv") || strings.Contains(name, "ben gurion") {
		return "TLV"
	}

	return ""
}

// PrintRoutes prints routes in a formatted way
func PrintRoutes(routes []Route) {
	fmt.Printf("\nFound %d routes:\n\n", len(routes))
	for i, route := range routes {
		fmt.Printf("Route %d: %s\n", i+1, route.Description)
		fmt.Printf("  Total Price: %.2f %s\n", route.TotalPrice, route.Currency)
		fmt.Printf("  Total Time: %v\n", route.TotalTime)
		fmt.Printf("  Departure: %s\n", route.Departure.Format("2006-01-02 15:04"))
		fmt.Printf("  Arrival: %s\n", route.Arrival.Format("2006-01-02 15:04"))

		for j, segment := range route.Segments {
			fmt.Printf("    Segment %d: %s from %s to %s\n", j+1, segment.Mode, segment.From.Name, segment.To.Name)
			fmt.Printf("      Duration: %v, Price: %.2f %s\n", segment.Duration, segment.Price, segment.Currency)
		}
		fmt.Println()
	}
}

// CalculateTotals calculates total price and time for a route
func (r *Route) CalculateTotals() {
	r.TotalPrice = 0
	r.TotalTime = 0

	if len(r.Segments) == 0 {
		return
	}

	r.Departure = r.Segments[0].Departure
	r.Arrival = r.Segments[len(r.Segments)-1].Arrival
	r.TotalTime = r.Arrival.Sub(r.Departure)

	var descriptions []string
	for _, segment := range r.Segments {
		r.TotalPrice += segment.Price
		descriptions = append(descriptions, fmt.Sprintf("%s (%s)", segment.Mode, segment.Provider))
	}

	r.Description = strings.Join(descriptions, " → ")
}
