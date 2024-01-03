package polling

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePollingRate(t *testing.T) {
	t.Run("should return error if invalid unit", func(t *testing.T) {
		_, err := ParsePollingRate("1x")

		expectedErr := fmt.Sprintf(invalidUnitMsg, "x")

		assert.EqualError(t, err, expectedErr)
	})

	t.Run("should return error if invalid number", func(t *testing.T) {
		_, err := ParsePollingRate("xh")

		assert.Error(t, err)
	})

	t.Run("should return error if invalid number and unit", func(t *testing.T) {
		_, err := ParsePollingRate("::")

		assert.Error(t, err)
	})

	t.Run("should return correct duration for seconds", func(t *testing.T) {
		duration, err := ParsePollingRate("1s")

		assert.NoError(t, err)
		assert.Equal(t, 1, int(duration.Seconds()))
	})

	t.Run("should return correct duration for minutes", func(t *testing.T) {
		duration, err := ParsePollingRate("1m")

		assert.NoError(t, err)
		assert.Equal(t, 1, int(duration.Minutes()))
	})

	t.Run("should return correct duration for hours", func(t *testing.T) {
		duration, err := ParsePollingRate("1h")

		assert.NoError(t, err)
		assert.Equal(t, 1, int(duration.Hours()))
	})
}
