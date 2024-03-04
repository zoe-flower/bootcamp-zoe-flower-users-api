package duration

import (
	"fmt"
	"time"
)

// GreaterThanZero parses a time duration string that must be greater than zero
func GreaterThanZero(in string) (time.Duration, error) {
	d, err := time.ParseDuration(in)
	if err != nil {
		return 0, err
	}
	if d <= 0 {
		return 0, fmt.Errorf("duration must be greater than 0")
	}

	return d, nil
}
