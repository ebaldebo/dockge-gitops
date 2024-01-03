package env

import (
	"fmt"
	"os"
)

const (
	missingRequiredEnvVar = "missing required environment variable: %s"
)

func GetEnvVar(required bool, key, fallback string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists && required {
		return "", fmt.Errorf(missingRequiredEnvVar, key)
	}

	if !exists {
		return fallback, nil
	}

	return value, nil
}
