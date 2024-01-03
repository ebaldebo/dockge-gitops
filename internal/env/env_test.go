package env

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetEnvVar(t *testing.T) {
	t.Run("should return error if requires env var is not set", func(t *testing.T) {
		os.Clearenv()

		key := "TEST"
		expectedErr := fmt.Sprintf(missingRequiredEnvVar, key)

		_, err := GetEnvVar(true, key, "fallback")

		assert.EqualError(t, err, expectedErr)
	})

	t.Run("should return fallback if env var is not set", func(t *testing.T) {
		os.Clearenv()

		key := "TEST"
		fallback := "fallback"

		value, err := GetEnvVar(false, key, fallback)

		assert.NoError(t, err)
		assert.Equal(t, fallback, value)
	})

	t.Run("should return env var if set", func(t *testing.T) {
		os.Clearenv()

		key := "TEST"
		value := "value"

		os.Setenv(key, value)

		result, err := GetEnvVar(false, key, "fallback")

		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})
}
