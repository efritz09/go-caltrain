package caltrain

import (
	"context"
	"reflect"
	"testing"
	"time"
)

const (
	fakeKey = "ericisgreat"
)

func TestGetTrainRoute(t *testing.T) {
	ctx := context.Background()
	c := New(fakeKey)
	m := &MockAPIClient{}
	m.GetResultFilePath = "testdata/bulletSchedule.json"
	c.APIClient = m
	err := c.UpdateTimeTable(ctx)
	if err != nil {
		t.Fatalf("Unexpected error loading timetable: %v", err)
	}
	// c.UpdateTimeTable currently populates each line with bulletSchedule.
	// remove the other instances
	delete(c.timetable, Limited)
	delete(c.timetable, Local)

	exp := Route{
		TrainNum:  "801",
		Direction: North,
		NumStops:  9,
		Stops: []TrainStop{
			TrainStop{Order: 1, Station: StationSanJose, Arrival: time.Date(0, time.January, 1, 9, 51, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 9, 51, 0, 0, time.UTC)},
			TrainStop{Order: 2, Station: StationSunnyvale, Arrival: time.Date(0, time.January, 1, 10, 1, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 10, 1, 0, 0, time.UTC)},
			TrainStop{Order: 3, Station: StationMountainView, Arrival: time.Date(0, time.January, 1, 10, 6, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 10, 6, 0, 0, time.UTC)},
			TrainStop{Order: 4, Station: StationPaloAlto, Arrival: time.Date(0, time.January, 1, 10, 13, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 10, 13, 0, 0, time.UTC)},
			TrainStop{Order: 5, Station: StationRedwoodCity, Arrival: time.Date(0, time.January, 1, 10, 20, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 10, 20, 0, 0, time.UTC)},
			TrainStop{Order: 6, Station: StationHillsdale, Arrival: time.Date(0, time.January, 1, 10, 27, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 10, 27, 0, 0, time.UTC)},
			TrainStop{Order: 7, Station: StationSanMateo, Arrival: time.Date(0, time.January, 1, 10, 32, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 10, 32, 0, 0, time.UTC)},
			TrainStop{Order: 8, Station: StationMillbrae, Arrival: time.Date(0, time.January, 1, 10, 40, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 10, 40, 0, 0, time.UTC)},
			TrainStop{Order: 9, Station: StationSanFrancisco, Arrival: time.Date(0, time.January, 1, 11, 0, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 11, 0, 0, 0, time.UTC)},
		},
	}

	route, err := c.GetTrainRoute("801")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(route, exp) {
		t.Fatalf("Unexpected route\nExpected: %v\nReceived: %v", exp, route)
	}

	noRoute, err := c.GetTrainRoute("101")
	if err == nil {
		t.Fatalf("should not have gotten a route for train 101\n%v", noRoute)
	}

}

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

// Simple test to ensure the code runs
func TestGetStationTimetable(t *testing.T) {
	c := New(fakeKey)
	m := &MockAPIClient{}
	c.APIClient = m
	_, err := c.GetStationTimetable(StationHillsdale, North)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
