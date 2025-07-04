[![Go Coverage](https://img.shields.io/badge/coverage-48.5%25-yellowgreen)](coverage.txt)

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

## 🧪 Test Coverage

To check test coverage:

```sh
go test -coverprofile=coverage.txt ./...
go tool cover -func=coverage.txt
```

Latest coverage: **48.5%** of statements (as of June 28, 2025)

| File                | Coverage |
|---------------------|----------|
| utils.go            | 100%     |
| config.go           | 69.2%    |
| handlers.go         | 62.5%    |
| airport_service.go  | 54.2%    |
| flight_service.go   | 75–100%  |
| transport_service.go| 69–100%  |
| travel_finder.go    | 12.1%    |
| main.go             | 0%       |

> See `coverage.txt` for details. Improve coverage by adding more tests for services and handlers.

✅ TODO

- Add more tests for services and handlers (coverage is improving!)
- Replace mock flight logic with real Amadeus API calls
- Add error logging middleware

Happy hacking! ✈️

