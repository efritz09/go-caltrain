package caltrain

import (
	"context"
	"time"
)

type MockCaltrain struct {
	GetDelaysFunc                func(context.Context) ([]Train, error)
	GetStationStatusFunc         func(context.Context, string, string) ([]Train, error)
	GetTrainsBetweenStationsFunc func(context.Context, string, string) ([]*Route, error)
	GetStationsFunc              func() []string
	GetDirectionFunc             func(src, dst string) (error, string)
}

func (c *MockCaltrain) GetDelays(ctx context.Context) ([]Train, error) {
	if c.GetDelaysFunc != nil {
		return c.GetDelaysFunc(ctx)
	}
	return nil, nil
}

func (c *MockCaltrain) GetStationStatus(ctx context.Context, stationName string, direction string) ([]Train, error) {
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

func (c *MockCaltrain) GetTrainsBetweenStations(ctx context.Context, src, dst string) ([]*Route, error) {
	if c.GetTrainsBetweenStationsFunc != nil {
		return c.GetTrainsBetweenStationsFunc(ctx, src, dst)
	}
	return nil, nil
}

func (c *MockCaltrain) SetupCache(expire time.Duration) {}

func (c *MockCaltrain) UpdateTimeTable(ctx context.Context) error { return nil }

func (c *MockCaltrain) UpdateStations(ctx context.Context) error { return nil }

func (c *MockCaltrain) Initialize(ctx context.Context) error { return nil }

func (c *MockCaltrain) GetDirectionFromSrcToDst(src, dst string) (string, error) {
	if c.GetDirectionFunc != nil {
		return c.GetDirectionFunc(src, dst)
	}
	return North, nil
}
