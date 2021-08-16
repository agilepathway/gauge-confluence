// Package env provides environment variables functionality.
package env

import (
	"fmt"
	"os"
	"strings"
)

// GetRequired returns an environment variable value or panics if not present.
func GetRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Aborting: %s is not set. Set it and try again. "+
			"See https://github.com/agilepathway/gauge-confluence#plugin-setup", key))
	}

	return value
}

// GetBoolean returns true if the environment variable exists and has a value of "true"
func GetBoolean(key string) bool {
	value := os.Getenv(key)
	return strings.ToLower(value) == "true"
}
