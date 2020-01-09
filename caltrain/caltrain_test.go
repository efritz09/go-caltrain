package caltrain

import (
	"context"
	"errors"
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
	if err := c.UpdateTimeTable(ctx); err != nil {
		t.Fatalf("Unexpected error loading timetable: %v", err)
	}
	m.GetResultFilePath = "testdata/stations.json"
	if err := c.UpdateStations(ctx); err != nil {
		t.Fatalf("Unexpected error loading stations: %v", err)
	}
	// c.UpdateTimeTable currently populates each line with bulletSchedule.
	// remove the other instances
	delete(c.timetable, Limited)
	delete(c.timetable, Local)

	exp := &Route{
		TrainNum:  "801",
		Direction: North,
		Line:      Bullet,
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
func TestGetTrainsBetweenStations(t *testing.T) {
	ctx := context.Background()
	c := New(fakeKey)
	m := &MockAPIClient{}
	m.GetResultFilePath = "testdata/bulletSchedule.json"
	c.APIClient = m
	if err := c.UpdateTimeTable(ctx); err != nil {
		t.Fatalf("Unexpected error loading timetable: %v", err)
	}
	m.GetResultFilePath = "testdata/stations.json"
	if err := c.UpdateStations(ctx); err != nil {
		t.Fatalf("Unexpected error loading stations: %v", err)
	}
	// c.UpdateTimeTable currently populates each line with bulletSchedule.
	// remove the other instances
	delete(c.timetable, Limited)
	delete(c.timetable, Local)

	tests := []struct {
		src  string
		dst  string
		numN int // len of array for now
		numS int
		day  string
		err  error
	}{
		{src: StationHillsdale, dst: StationPaloAlto, numN: 5, numS: 5, day: "monday", err: nil},
		{src: StationSanJose, dst: StationSanFrancisco, numN: 11, numS: 11, day: "monday", err: nil},
		{src: StationSanJose, dst: StationSanFrancisco, numN: 2, numS: 2, day: "sunday", err: nil},
		{src: StationHillsdale, dst: StationHaywardPark, numN: 0, numS: 0, day: "monday", err: nil},
		{src: StationSanFrancisco, dst: "BadSation", numN: 0, numS: 0, day: "monday", err: errors.New("")},
	}

	for _, tt := range tests {
		name := tt.src + "_" + tt.dst
		t.Run(name, func(t *testing.T) {
			u := &MockUpdater{}
			u.Weekday = tt.day
			c.Updater = u

			n, s, err := c.GetTrainsBetweenStations(ctx, tt.src, tt.dst)
			if err != nil && tt.err == nil {
				t.Fatalf("Failed to get train routes for %s: %v", name, err)
			} else if err == nil && tt.err != nil {
				t.Fatalf("getTrainRoutesBetweenStations improperly succeeded for %s", name)
			}

			if len(n) != tt.numN {
				t.Fatalf("Incorrect routes North. Expected %d, recieved %d", tt.numN, len(n))
			}
			if len(s) != tt.numS {
				t.Fatalf("Incorrect routes North. Expected %d, recieved %d", tt.numS, len(s))
			}
		})
	}
}

// Simple test to ensure the code runs
func TestGetDelays(t *testing.T) {
	ctx := context.Background()
	c := New(fakeKey)
	tests := []struct {
		name   string
		data   string
		delays int
		err    error
	}{
		{name: "Data1", data: "testdata/parseDelayData1.json", delays: 2, err: nil},
		{name: "Data2", data: "testdata/parseDelayData2.json", delays: 0, err: nil},
		{name: "DataErr", data: "testdata/parseDelayData.json", delays: 0, err: errors.New("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockAPIClient{}
			m.GetResultFilePath = tt.data
			c.APIClient = m
			d, err := c.GetDelays(ctx)
			if err != nil && tt.err == nil {
				t.Fatalf("Failed to get train delays for %s: %v", tt.name, err)
			} else if err == nil && tt.err != nil {
				t.Fatalf("GetDelays improperly succeeded for %s", tt.name)
			}

			if len(d) != tt.delays {
				t.Fatalf("Incorrect number of delays. Expected %d, recieved %d", tt.delays, len(d))
			}
		})
	}
}

func TestGetDelaysCache(t *testing.T) {
	ctx := context.Background()
	c := New(fakeKey)
	m := &MockAPIClient{}
	m.GetResultFilePath = "testdata/parseDelayData1.json"
	c.APIClient = m
	c.SetupCache(DefaultCacheTimeout)

	cache := make(map[string][]byte)
	mock := &MockCache{}
	mock.SetFunc = func(key string, body []byte) { cache[key] = body }
	mock.GetFunc = func(key string) ([]byte, bool) {
		v, ok := cache[key]
		return v, ok
	}
	c.Cache = mock

	// first make a call, the cache should be empty
	if len(cache) != 0 {
		t.Fatalf("Cache is not empty: %v", cache)
	}
	d, err := c.GetDelays(ctx)
	if err != nil {
		t.Fatalf("Failed to get train delays for %v", err)
	}
	if len(d) != 2 {
		t.Fatalf("Incorrect number of delays. Expected %d, recieved %d", 2, len(d))
	}

	// check that the cache was filled
	if len(cache) != 1 {
		t.Fatalf("Cache does not have only 1 key: %v", cache)
	}
	// run it again
	d, err = c.GetDelays(ctx)
	if err != nil {
		t.Fatalf("Failed to get train delays for %v", err)
	}
	if len(d) != 2 {
		t.Fatalf("Incorrect number of delays. Expected %d, recieved %d", 2, len(d))
	}

	// check that the cache was not changed
	if len(cache) != 1 {
		t.Fatalf("Cache does not have only 1 key: %v", cache)
	}
}

// Simple test to ensure the code runs
func TestGetStationStatus(t *testing.T) {
	ctx := context.Background()
	c := New(fakeKey)
	m := &MockAPIClient{}
	m.GetResultFilePath = "testdata/stations.json"
	c.APIClient = m
	if err := c.UpdateStations(ctx); err != nil {
		t.Fatalf("Unexpected error loading stations: %v", err)
	}
	m.GetResultFilePath = "testdata/parseHillsdaleNorth.json"
	_, err := c.GetStationStatus(ctx, StationHillsdale, North)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestGetStationStatusCache(t *testing.T) {
	ctx := context.Background()
	c := New(fakeKey)
	m := &MockAPIClient{}
	m.GetResultFilePath = "testdata/stations.json"
	c.APIClient = m
	if err := c.UpdateStations(ctx); err != nil {
		t.Fatalf("Unexpected error loading stations: %v", err)
	}
	m.GetResultFilePath = "testdata/parseHillsdaleNorth.json"
	c.SetupCache(DefaultCacheTimeout)

	cache := make(map[string][]byte)
	mock := &MockCache{}
	mock.SetFunc = func(key string, body []byte) { cache[key] = body }
	mock.GetFunc = func(key string) ([]byte, bool) {
		v, ok := cache[key]
		return v, ok
	}
	c.Cache = mock

	// first make a call, the cache should be empty
	if len(cache) != 0 {
		t.Fatalf("Cache is not empty: %v", cache)
	}
	d, err := c.GetStationStatus(ctx, StationHillsdale, North)
	if err != nil {
		t.Fatalf("Failed to get train delays for %v", err)
	}
	if len(d) != 1 {
		t.Fatalf("Incorrect number of delays. Expected %d, recieved %d", 1, len(d))
	}
	// check that the cache was filled
	if len(cache) != 1 {
		t.Fatalf("Cache does not have only 1 key: %v", cache)
	}

	// Now replace it with a South call
	m.GetResultFilePath = "testdata/parseHillsdaleSouth.json"
	c.APIClient = m
	// make a south call
	d, err = c.GetStationStatus(ctx, StationHillsdale, South)
	if err != nil {
		t.Fatalf("Failed to get train delays for %v", err)
	}
	if len(d) != 2 {
		t.Fatalf("Incorrect number of delays. Expected %d, recieved %d", 2, len(d))
	}

	// check that the cache was not changed
	if len(cache) != 2 {
		t.Fatalf("Cache does not have only 1 key: %v", cache)
	}
	// make a the same call call
	d, err = c.GetStationStatus(ctx, StationHillsdale, South)
	if err != nil {
		t.Fatalf("Failed to get train delays for %v", err)
	}
	if len(d) != 2 {
		t.Fatalf("Incorrect number of delays. Expected %d, recieved %d", 2, len(d))
	}

	// check that the cache was not changed
	if len(cache) != 2 {
		t.Fatalf("Cache does not have only 1 key: %v", cache)
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
	ctx := context.Background()
	c := New(fakeKey)
	m := &MockAPIClient{}
	m.GetResultFilePath = "testdata/stations.json"
	c.APIClient = m
	if err := c.UpdateStations(ctx); err != nil {
		t.Fatalf("Unexpected error loading stations: %v", err)
	}

	_, err := c.GetStationTimetable(StationHillsdale, North)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
