Travel Route Finder (Go)

A lightweight Go web service that finds travel routes between cities using:

Google Maps APIs (Geocoding, Directions, Places)

Mock flight search logic

Basic public transit or taxi price estimations

🚀 How to Run

1. Clone and Initialize Go Modules

git clone https://github.com/walidBarakeh/travel-routes.git
cd travel-routes
go mod tidy

2. Create a .env file

Refer to the .env.template or see below.

3. Run the Server

go run .

The server will start on port 8080 by default (or PORT env var).

🔧 Environment Variables

Create a .env file in the root directory:

GOOGLE_MAPS_API_KEY=your-google-maps-api-key
AMADEUS_API_KEY=your-amadeus-api-key
AMADEUS_SECRET=your-amadeus-secret
PORT=8080

📡 Available Endpoints

/search

Find travel routes from origin to destination

GET /search?origin=Granada&destination=Tel%20Aviv&date=2024-07-01

/airports

Find nearby airports to a location

GET /airports?location=Granada&radius=30000

/health

Health check endpoint

GET /health

📘 Example Output

Route 1: public_transport (Public Transport) → flight (Airlines)
  Total Price: 215.00 EUR
  Total Time: 8h30m
  Departure: 2024-07-01 08:00
  Arrival:   2024-07-01 16:30

✅ TODO

Add tests for services and handlers

Replace mock flight logic with real Amadeus API calls

Add error logging middleware

Happy hacking! ✈️

