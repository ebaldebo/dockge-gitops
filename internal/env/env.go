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

func EnvFileExists(envFilePath string) bool {
	if _, err := os.Stat(envFilePath); err == nil {
		return true
	}

	return false
}
