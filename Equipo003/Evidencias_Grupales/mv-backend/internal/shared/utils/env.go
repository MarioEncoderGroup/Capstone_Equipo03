package utils

import "os"

// GetEnvOrDefault retrieves an environment variable value by key.
// If the environment variable exists and is not empty, returns its value.
// Otherwise, returns the provided defaultValue.
// This utility follows DRY principle to avoid code duplication across the application.
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}