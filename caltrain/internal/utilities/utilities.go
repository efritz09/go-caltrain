package utilities

import (
	"strings"
	"time"
)

// GetWeekday returns the day's weekday given a timezone
func GetWeekday(timezone string) (string, error) {
	// get the current time and day
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}
	tzTime := time.Now().In(loc)
	return strings.ToLower(tzTime.Weekday().String()), nil
}
