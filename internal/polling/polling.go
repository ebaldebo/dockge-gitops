package polling

import (
	"fmt"
	"strconv"
	"time"
)

const (
	invalidUnitMsg = "invalid unit: %s"
)

func ParsePollingRate(interval string) (time.Duration, error) {
	number := interval[:len(interval)-1]
	unit := interval[len(interval)-1:]

	value, err := strconv.Atoi(number)
	if err != nil {
		return 0, err
	}

	switch unit {
	case "s":
		return time.Duration(value) * time.Second, nil
	case "m":
		return time.Duration(value) * time.Minute, nil
	case "h":
		return time.Duration(value) * time.Hour, nil
	default:
		return 0, fmt.Errorf(invalidUnitMsg, unit)
	}
}
