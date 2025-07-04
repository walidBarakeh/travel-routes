package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (this is okay in production)")
	}

	// Load configuration
	config := LoadConfig()

	// Create travel finder
	tf := NewTravelFinder(config)

	// Set up HTTP routes
	http.HandleFunc("/search", tf.handleSearchRoutes)
	http.HandleFunc("/airports", tf.handleNearbyAirports)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Example usage
	fmt.Println("=== Travel Route Finder ===")
	fmt.Println("Example: Finding routes from Granada to Tel Aviv on July 1st, 2024")

	travelDate := time.Date(2024, 7, 1, 8, 0, 0, 0, time.UTC)
	routes, err := tf.FindRoutes("Granada", "Tel Aviv", travelDate)
	if err != nil {
		log.Printf("Error finding routes: %v", err)
	} else {
		PrintRoutes(routes)
	}

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("Endpoints:\n")
	fmt.Printf("  GET /search?origin=Granada&destination=Tel Aviv&date=2024-07-01\n")
	fmt.Printf("  GET /airports?location=Granada&radius=300\n")
	fmt.Printf("  GET /health\n")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
