package caltrain

import (
	"context"
	"testing"
)

const (
	fakeKey = "ericisgreat"
)

// Simple test to ensure the code runs
func TestGetDelays(t *testing.T) {
	ctx := context.Background()
	c := New(fakeKey)
	m := &MockAPIClient{}
	m.GetResultFilePath = "testdata/parseDelayData2.json"
	c.APIClient = m
	_, err := c.GetDelays(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// Simple test to ensure the code runs
func TestGetStationStatus(t *testing.T) {
	ctx := context.Background()
	c := New(fakeKey)
	m := &MockAPIClient{}
	m.GetResultFilePath = "testdata/parseHillsdaleNorth.json"
	c.APIClient = m
	_, err := c.GetStationStatus(ctx, StationHillsdale, North)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// Simple test to ensure the code runs
func TestGetTrainsBetweenStations(t *testing.T) {
	ctx := context.Background()
	c := New(fakeKey)
	m := &MockAPIClient{}
	c.APIClient = m
	_, err := c.GetTrainsBetweenStations(ctx, StationHillsdale, StationPaloAlto)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// Simple test to ensure the code runs
func TestUpdateTimeTable(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
	}{
		{name: Bullet, filepath: "testdata/bulletSchedule.json"},
		{name: Limited, filepath: "testdata/limitedSchedule.json"},
		{name: Local, filepath: "testdata/localSchedule.json"},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(fakeKey)
			m := &MockAPIClient{}
			m.GetResultFilePath = tt.filepath
			c.APIClient = m
			err := c.UpdateTimeTable(ctx)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
		})
	}
}
