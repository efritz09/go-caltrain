package caltrain

import (
	"strings"
	"time"
)

type Updater interface {
	GetWeekday(timezone string) (string, error)
}

type CaltrainUpdater struct{}

func NewUpdater() *CaltrainUpdater {
	return &CaltrainUpdater{}
}

// GetWeekday returns the day's weekday given a timezone
func (u *CaltrainUpdater) GetWeekday(timezone string) (string, error) {

	// get the current time and day
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}
	tzTime := time.Now().In(loc)
	return strings.ToLower(tzTime.Weekday().String()), nil
}

type MockUpdater struct {
	Weekday string
}

// GetWeekday returns the value in MockUpdater.Weekday
func (u *MockUpdater) GetWeekday(timezone string) (string, error) {
	return u.Weekday, nil
}
