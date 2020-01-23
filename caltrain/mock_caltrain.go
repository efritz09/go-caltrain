package caltrain

import (
	"context"
	"time"
)

// CaltrainMock implements Caltrain and can be used for unit testing with the
// required methods already mocked. Override the struct functions to make the
// mocked methods run your implementation.
type CaltrainMock struct {
	GetDelaysFunc                          func(context.Context, time.Duration) ([]TrainStatus, error)
	GetStationStatusFunc                   func(context.Context, string, string) ([]TrainStatus, error)
	GetTrainsBetweenStationsForWeekdayFunc func(context.Context, string, string, time.Weekday) ([]*Route, error)
	GetTrainsBetweenStationsForDateFunc    func(context.Context, string, string, time.Time) ([]*Route, error)
}

// GetDelays returns the GetDelaysFunc if it exists. Default nil, nil
func (c *CaltrainMock) GetDelays(ctx context.Context, d time.Duration) ([]TrainStatus, error) {
	if c.GetDelaysFunc != nil {
		return c.GetDelaysFunc(ctx, d)
	}
	return nil, nil
}

// GetDelays returns the GetStationStatusFunc if it exists. Default nil, nil
func (c *CaltrainMock) GetStationStatus(ctx context.Context, stationName string, direction string) ([]TrainStatus, error) {
	if c.GetStationStatusFunc != nil {
		return c.GetStationStatusFunc(ctx, stationName, direction)
	}
	return nil, nil
}

// GetDelays returns the GetTrainsBetweenStationsForWeekdayFunc if it exists.
// Default nil, nil
func (c *CaltrainMock) GetTrainsBetweenStationsForWeekday(ctx context.Context, src, dst string, weekday time.Weekday) ([]*Route, error) {
	if c.GetTrainsBetweenStationsForWeekdayFunc != nil {
		return c.GetTrainsBetweenStationsForWeekdayFunc(ctx, src, dst, weekday)
	}
	return nil, nil
}

// GetDelays returns the GetTrainsBetweenStationsForDateFunc if it exists.
// Default nil, nil
func (c *CaltrainMock) GetTrainsBetweenStationsForDate(ctx context.Context, src, dst string, date time.Time) ([]*Route, error) {
	if c.GetTrainsBetweenStationsForDateFunc != nil {
		return c.GetTrainsBetweenStationsForDateFunc(ctx, src, dst, date)
	}
	return nil, nil
}

// SetupCache returns without doing anything
func (c *CaltrainMock) SetupCache(expire time.Duration) {}

// UpdateTimeTable returns nil
func (c *CaltrainMock) UpdateTimeTable(ctx context.Context) error { return nil }

// UpdateStations returns nil
func (c *CaltrainMock) UpdateStations(ctx context.Context) error { return nil }

// UpdateHolidays returns nil
func (c *CaltrainMock) UpdateHolidays(ctx context.Context) error { return nil }

// Initialize returns nil
func (c *CaltrainMock) Initialize(ctx context.Context) error { return nil }

// IsHoliday returns false
func (c *CaltrainMock) IsHoliday(date time.Time) bool { return false }
