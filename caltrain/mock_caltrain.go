package caltrain

import (
	"context"
	"time"
)

type CaltrainMock struct {
	GetDelaysFunc                          func(context.Context, time.Duration) ([]TrainStatus, error)
	GetStationStatusFunc                   func(context.Context, string, string) ([]TrainStatus, error)
	GetTrainsBetweenStationsForWeekdayFunc func(context.Context, string, string, time.Weekday) ([]*Route, error)
	GetTrainsBetweenStationsForDateFunc    func(context.Context, string, string, time.Time) ([]*Route, error)
	GetStationsFunc                        func() []string
	GetDirectionFunc                       func(src, dst string) (Direction, error)
}

func (c *CaltrainMock) GetDelays(ctx context.Context, d time.Duration) ([]TrainStatus, error) {
	if c.GetDelaysFunc != nil {
		return c.GetDelaysFunc(ctx, d)
	}
	return nil, nil
}

func (c *CaltrainMock) GetStationStatus(ctx context.Context, stationName string, direction string) ([]TrainStatus, error) {
	if c.GetStationStatusFunc != nil {
		return c.GetStationStatusFunc(ctx, stationName, direction)
	}
	return nil, nil
}

func (c *CaltrainMock) GetStations() []string {
	if c.GetStationsFunc != nil {
		return c.GetStationsFunc()
	}
	return nil
}

func (c *CaltrainMock) GetTrainsBetweenStationsForWeekday(ctx context.Context, src, dst string, weekday time.Weekday) ([]*Route, error) {
	if c.GetTrainsBetweenStationsForWeekdayFunc != nil {
		return c.GetTrainsBetweenStationsForWeekdayFunc(ctx, src, dst, weekday)
	}
	return nil, nil
}

func (c *CaltrainMock) GetTrainsBetweenStationsForDate(ctx context.Context, src, dst string, date time.Time) ([]*Route, error) {
	if c.GetTrainsBetweenStationsForDateFunc != nil {
		return c.GetTrainsBetweenStationsForDateFunc(ctx, src, dst, date)
	}
	return nil, nil
}

func (c *CaltrainMock) SetupCache(expire time.Duration) {}

func (c *CaltrainMock) UpdateTimeTable(ctx context.Context) error { return nil }

func (c *CaltrainMock) UpdateStations(ctx context.Context) error { return nil }

func (c *CaltrainMock) UpdateHolidays(ctx context.Context) error { return nil }

func (c *CaltrainMock) Initialize(ctx context.Context) error { return nil }

func (c *CaltrainMock) GetDirectionFromSrcToDst(src, dst string) (Direction, error) {
	if c.GetDirectionFunc != nil {
		return c.GetDirectionFunc(src, dst)
	}
	return North, nil
}

func (c *CaltrainMock) IsHoliday(date time.Time) bool { return false }
