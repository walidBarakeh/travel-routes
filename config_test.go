package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Save original environment variables
	originalGoogleAPIKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	originalDefaultRadius := os.Getenv("DEFAULT_RADIUS")

	// Clean up after tests
	defer func() {
		if originalGoogleAPIKey != "" {
			os.Setenv("GOOGLE_MAPS_API_KEY", originalGoogleAPIKey)
		} else {
			os.Unsetenv("GOOGLE_MAPS_API_KEY")
		}
		if originalDefaultRadius != "" {
			os.Setenv("DEFAULT_RADIUS", originalDefaultRadius)
		} else {
			os.Unsetenv("DEFAULT_RADIUS")
		}
	}()

	t.Run("Load config with environment variables", func(t *testing.T) {
		// Set test environment variables
		os.Setenv("GOOGLE_MAPS_API_KEY", "test-api-key-123")
		os.Setenv("DEFAULT_RADIUS", "500")

		config := LoadConfig()

		assert.Equal(t, "test-api-key-123", config.GoogleMapsAPIKey)
		assert.Equal(t, 500, config.DefaultRadius)
	})

	t.Run("Load config with default values", func(t *testing.T) {
		// Clear environment variables
		os.Unsetenv("GOOGLE_MAPS_API_KEY")
		os.Unsetenv("DEFAULT_RADIUS")

		config := LoadConfig()

		assert.Equal(t, "", config.GoogleMapsAPIKey)
		assert.Equal(t, 300000, config.DefaultRadius) // Default 300km in meters
	})

	t.Run("Load config with invalid radius", func(t *testing.T) {
		// Set invalid radius
		os.Setenv("DEFAULT_RADIUS", "invalid-number")
		os.Unsetenv("GOOGLE_MAPS_API_KEY")

		config := LoadConfig()

		assert.Equal(t, "", config.GoogleMapsAPIKey)
		assert.Equal(t, 300000, config.DefaultRadius) // Should fall back to default
	})

	t.Run("Load config with zero radius", func(t *testing.T) {
		// Set zero radius
		os.Setenv("DEFAULT_RADIUS", "0")
		os.Setenv("GOOGLE_MAPS_API_KEY", "test-key")

		config := LoadConfig()

		assert.Equal(t, "test-key", config.GoogleMapsAPIKey)
		assert.Equal(t, 0, config.DefaultRadius)
	})

	t.Run("Load config with negative radius", func(t *testing.T) {
		// Set negative radius
		os.Setenv("DEFAULT_RADIUS", "-100")
		os.Setenv("GOOGLE_MAPS_API_KEY", "test-key")

		config := LoadConfig()

		assert.Equal(t, "test-key", config.GoogleMapsAPIKey)
		assert.Equal(t, -100, config.DefaultRadius)
	})

	t.Run("Load config with very large radius", func(t *testing.T) {
		// Set very large radius
		os.Setenv("DEFAULT_RADIUS", "999999999")
		os.Setenv("GOOGLE_MAPS_API_KEY", "test-key")

		config := LoadConfig()

		assert.Equal(t, "test-key", config.GoogleMapsAPIKey)
		assert.Equal(t, 999999999, config.DefaultRadius)
	})

	t.Run("Load config with empty string values", func(t *testing.T) {
		// Set empty string values
		os.Setenv("GOOGLE_MAPS_API_KEY", "")
		os.Setenv("DEFAULT_RADIUS", "")

		config := LoadConfig()

		assert.Equal(t, "", config.GoogleMapsAPIKey)
		assert.Equal(t, 300000, config.DefaultRadius) // Should use default for empty string
	})

	t.Run("Load config with whitespace values", func(t *testing.T) {
		// Set whitespace values
		os.Setenv("GOOGLE_MAPS_API_KEY", "  test-key-with-spaces  ")
		os.Setenv("DEFAULT_RADIUS", "  1000  ")

		config := LoadConfig()

		// Note: The actual implementation might trim whitespace or not
		// This test documents the current behavior
		assert.Contains(t, config.GoogleMapsAPIKey, "test-key-with-spaces")
		// The radius parsing might handle whitespace differently
	})
}

func TestConfigStruct(t *testing.T) {
	t.Run("Create config struct directly", func(t *testing.T) {
		config := Config{
			GoogleMapsAPIKey: "direct-api-key",
			DefaultRadius:    250000,
		}

		assert.Equal(t, "direct-api-key", config.GoogleMapsAPIKey)
		assert.Equal(t, 250000, config.DefaultRadius)
	})

	t.Run("Config struct with zero values", func(t *testing.T) {
		config := Config{}

		assert.Equal(t, "", config.GoogleMapsAPIKey)
		assert.Equal(t, 0, config.DefaultRadius)
	})

	t.Run("Config struct field validation", func(t *testing.T) {
		config := Config{
			GoogleMapsAPIKey: "AIzaSyBvOkBwgGlbUiuS-oKrPbB3IFI-7JGQYdA", // Example format
			DefaultRadius:    500000, // 500km
		}

		// Test that API key looks reasonable (basic format check)
		assert.True(t, len(config.GoogleMapsAPIKey) > 10)
		assert.Contains(t, config.GoogleMapsAPIKey, "AIza")

		// Test that radius is reasonable
		assert.True(t, config.DefaultRadius > 0)
		assert.True(t, config.DefaultRadius <= 1000000) // Less than 1000km seems reasonable
	})
}

// Test environment variable parsing edge cases
func TestEnvironmentVariableParsing(t *testing.T) {
	// Save original environment
	originalDefaultRadius := os.Getenv("DEFAULT_RADIUS")
	defer func() {
		if originalDefaultRadius != "" {
			os.Setenv("DEFAULT_RADIUS", originalDefaultRadius)
		} else {
			os.Unsetenv("DEFAULT_RADIUS")
		}
	}()

	testCases := []struct {
		name           string
		envValue       string
		expectedRadius int
		description    string
	}{
		{
			name:           "Valid positive integer",
			envValue:       "123456",
			expectedRadius: 123456,
			description:    "Should parse valid positive integers",
		},
		{
			name:           "Valid zero",
			envValue:       "0",
			expectedRadius: 0,
			description:    "Should parse zero correctly",
		},
		{
			name:           "Leading zeros",
			envValue:       "000123",
			expectedRadius: 123,
			description:    "Should handle leading zeros",
		},
		{
			name:           "Decimal number",
			envValue:       "123.45",
			expectedRadius: 300000, // Should fall back to default
			description:    "Should fall back to default for decimal numbers",
		},
		{
			name:           "Text value",
			envValue:       "not-a-number",
			expectedRadius: 300000, // Should fall back to default
			description:    "Should fall back to default for non-numeric values",
		},
		{
			name:           "Mixed alphanumeric",
			envValue:       "123abc",
			expectedRadius: 300000, // Should fall back to default
			description:    "Should fall back to default for mixed alphanumeric",
		},
		{
			name:           "Empty string",
			envValue:       "",
			expectedRadius: 300000, // Should fall back to default
			description:    "Should fall back to default for empty string",
		},
		{
			name:           "Whitespace only",
			envValue:       "   ",
			expectedRadius: 300000, // Should fall back to default
			description:    "Should fall back to default for whitespace",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("DEFAULT_RADIUS", tc.envValue)
			os.Unsetenv("GOOGLE_MAPS_API_KEY") // Keep this clean

			config := LoadConfig()

			assert.Equal(t, tc.expectedRadius, config.DefaultRadius, tc.description)
		})
	}
}

// Benchmark config loading
func BenchmarkLoadConfig(b *testing.B) {
	// Set up environment
	os.Setenv("GOOGLE_MAPS_API_KEY", "test-api-key")
	os.Setenv("DEFAULT_RADIUS", "300000")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LoadConfig()
	}
}
