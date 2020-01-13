package caltrain

import (
	"context"
	"time"
)

type MockCaltrain struct {
	GetDelaysFunc                func(context.Context, time.Duration) ([]TrainStatus, error)
	GetStationStatusFunc         func(context.Context, string, string) ([]TrainStatus, error)
	GetTrainsBetweenStationsFunc func(context.Context, string, string, time.Weekday) ([]*Route, error)
	GetStationsFunc              func() []string
	GetDirectionFunc             func(src, dst string) (string, error)
}

func (c *MockCaltrain) GetDelays(ctx context.Context, d time.Duration) ([]TrainStatus, error) {
	if c.GetDelaysFunc != nil {
		return c.GetDelaysFunc(ctx, d)
	}
	return nil, nil
}

func (c *MockCaltrain) GetStationStatus(ctx context.Context, stationName string, direction string) ([]TrainStatus, error) {
	if c.GetStationStatusFunc != nil {
		return c.GetStationStatusFunc(ctx, stationName, direction)
	}
	return nil, nil
}

func (c *MockCaltrain) GetStations() []string {
	if c.GetStationsFunc != nil {
		return c.GetStationsFunc()
	}
	return nil
}

func (c *MockCaltrain) GetTrainsBetweenStations(ctx context.Context, src, dst string, day time.Weekday) ([]*Route, error) {
	if c.GetTrainsBetweenStationsFunc != nil {
		return c.GetTrainsBetweenStationsFunc(ctx, src, dst, day)
	}
	return nil, nil
}

func (c *MockCaltrain) SetupCache(expire time.Duration) {}

func (c *MockCaltrain) UpdateTimeTable(ctx context.Context) error { return nil }

func (c *MockCaltrain) UpdateStations(ctx context.Context) error { return nil }

func (c *MockCaltrain) UpdateHolidays(ctx context.Context) error { return nil }

func (c *MockCaltrain) Initialize(ctx context.Context) error { return nil }

func (c *MockCaltrain) GetDirectionFromSrcToDst(src, dst string) (string, error) {
	if c.GetDirectionFunc != nil {
		return c.GetDirectionFunc(src, dst)
	}
	return North, nil
}
