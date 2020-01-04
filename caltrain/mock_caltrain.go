package caltrain

import "context"

type MockCaltrain struct {
	GetDelaysFunc        func(context.Context) ([]Train, error)
	GetStationStatusFunc func(context.Context, string, string) ([]Train, error)
	GetStationsFunc      func() []string
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
