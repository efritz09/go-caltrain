package caltrain

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestGetTimetableForStation(t *testing.T) {
	// Load the timetable for only the bullet schedule
	ctx := context.Background()
	c := New(fakeKey)
	c.lines = allLines
	m := &apiClientMock{}
	m.GetResultFilePath = "testdata/bulletSchedule.json"
	c.APIClient = m
	err := c.UpdateTimeTable(ctx)
	if err != nil {
		t.Fatalf("Unexpected error loading timetable: %v", err)
	}
	m.GetResultFilePath = "testdata/stations.json"
	if err := c.UpdateStations(ctx); err != nil {
		t.Fatalf("Unexpected error loading stations: %v", err)
	}
	// c.UpdateTimeTable currently populates each line with bulletSchedule.
	// remove the other instances
	delete(c.timetable, "Limited")
	delete(c.timetable, "Limited A")
	delete(c.timetable, "Limited B")
	delete(c.timetable, "Local")
	delete(c.timetable, "Special")

	tests := []struct {
		station  Station
		dir      Direction
		day      time.Weekday
		expected int // length of array for now, should be []timetableRouteJourney
	}{
		{station: StationHillsdale, dir: North, day: time.Monday, expected: 5},
		{station: StationHillsdale, dir: North, day: time.Sunday, expected: 2},
		{station: StationHillsdale, dir: South, day: time.Monday, expected: 5},
		{station: StationHillsdale, dir: South, day: time.Sunday, expected: 2},
	}

	for _, tt := range tests {
		name := tt.station.String() + "/" + tt.dir.String() + "/" + tt.day.String()
		t.Run(name, func(t *testing.T) {
			code, err := c.getStationCode(StationHillsdale, tt.dir)
			if err != nil {
				t.Fatalf("failed to get station code: %v", err)
			}

			// Now we know what to expect
			journeys, err := c.getTimetableForStation(code, tt.dir, tt.day)
			if err != nil {
				t.Fatalf("failed to get timetable for station: %v", err)
			}
			// TODO: update with proper checks
			if len(journeys) != tt.expected {
				t.Fatalf("Unexpected journeys\nExpected: %v\nReceived: %v", len(journeys), tt.expected)
			}
		})
	}
}

func TestGetTrainRoutesBetweenStations(t *testing.T) {
	// Load the timetable for only the bullet schedule
	ctx := context.Background()
	c := New(fakeKey)
	c.lines = allLines
	m := &apiClientMock{}
	m.GetResultFilePath = "testdata/bulletSchedule.json"
	c.APIClient = m
	err := c.UpdateTimeTable(ctx)
	if err != nil {
		t.Fatalf("Unexpected error loading timetable: %v", err)
	}
	m.GetResultFilePath = "testdata/stations.json"
	if err := c.UpdateStations(ctx); err != nil {
		t.Fatalf("Unexpected error loading stations: %v", err)
	}
	// c.UpdateTimeTable currently populates each line with bulletSchedule.
	// remove the other instances
	delete(c.timetable, "Limited")
	delete(c.timetable, "Limited A")
	delete(c.timetable, "Limited B")
	delete(c.timetable, "Local")
	delete(c.timetable, "Special")

	tests := []struct {
		src  Station
		dst  Station
		numN int // len of array for now
		numS int
		day  time.Weekday
		err  error
	}{
		{src: StationHillsdale, dst: StationPaloAlto, numN: 5, numS: 5, day: time.Monday, err: nil},
		{src: StationSanJose, dst: StationSanFrancisco, numN: 11, numS: 11, day: time.Monday, err: nil},
		{src: StationSanJose, dst: StationSanFrancisco, numN: 2, numS: 2, day: time.Sunday, err: nil},
		{src: StationHillsdale, dst: StationHaywardPark, numN: 0, numS: 0, day: time.Monday, err: nil},
		{src: StationSanFrancisco, dst: 999, numN: 0, numS: 0, day: time.Monday, err: errors.New("")},
	}

	for _, tt := range tests {
		name := tt.src.String() + "_" + tt.dst.String()
		t.Run(name, func(t *testing.T) {
			// test north
			d1, err := c.getTrainRoutesBetweenStations(tt.src, tt.dst, tt.day)
			if err != nil && tt.err == nil {
				t.Fatalf("Failed to get train routes for %s: %v", name, err)
			} else if err == nil && tt.err != nil {
				t.Fatalf("getTrainRoutesBetweenStations improperly succeeded for %s", name)
			}
			if len(d1) != tt.numN {
				t.Fatalf("Incorrect routes North. Expected %d, received %d", tt.numN, len(d1))
			}

			// test south
			d2, err := c.getTrainRoutesBetweenStations(tt.dst, tt.src, tt.day)
			if err != nil && tt.err == nil {
				t.Fatalf("Failed to get train routes for %s: %v", name, err)
			} else if err == nil && tt.err != nil {
				t.Fatalf("getTrainRoutesBetweenStations improperly succeeded for %s", name)
			}
			if len(d2) != tt.numS {
				t.Fatalf("Incorrect routes North. Expected %d, received %d", tt.numS, len(d2))
			}
		})
	}
}

func TestGetRouteForTrain(t *testing.T) {
	// Load the timetable for only the bullet schedule
	ctx := context.Background()
	c := New(fakeKey)
	c.lines = allLines
	m := &apiClientMock{}
	m.GetResultFilePath = "testdata/bulletSchedule.json"
	c.APIClient = m
	err := c.UpdateTimeTable(ctx)
	if err != nil {
		t.Fatalf("Unexpected error loading timetable: %v", err)
	}
	m.GetResultFilePath = "testdata/stations.json"
	if err := c.UpdateStations(ctx); err != nil {
		t.Fatalf("Unexpected error loading stations: %v", err)
	}
	// c.UpdateTimeTable currently populates each line with bulletSchedule.
	// remove the other instances
	delete(c.timetable, "Limited")
	delete(c.timetable, "Limited A")
	delete(c.timetable, "Limited B")
	delete(c.timetable, "Local")
	delete(c.timetable, "Special")

	tests := []struct {
		train string
		line  string
		err   error
	}{
		{train: "801", line: "Bullet", err: nil},
		{train: "324", line: "Bullet", err: nil},
		{train: "101", line: "", err: errors.New("")},
	}

	for _, tt := range tests {
		t.Run(tt.train, func(t *testing.T) {
			r, err := c.getRouteForTrain(tt.train)
			if err != nil && tt.err == nil {
				t.Fatalf("Failed to get train info for %s: %v", tt.train, err)
			} else if err == nil && tt.err != nil {
				t.Fatalf("getRouteForTrain improperly succeeded for %s", tt.train)
			}
			if r.Line != tt.line {
				t.Fatalf("Unexpected train line. Expected %s, received %s", tt.line, r.Line)
			}
		})
	}
}

func TestGetTrainRoutesForAllStops(t *testing.T) {
	// Load the timetable for only the bullet schedule
	ctx := context.Background()
	c := New(fakeKey)
	c.lines = allLines
	m := &apiClientMock{}
	m.GetResultFilePath = "testdata/bulletSchedule.json"
	c.APIClient = m
	err := c.UpdateTimeTable(ctx)
	if err != nil {
		t.Fatalf("Unexpected error loading timetable: %v", err)
	}
	m.GetResultFilePath = "testdata/stations.json"
	if err := c.UpdateStations(ctx); err != nil {
		t.Fatalf("Unexpected error loading stations: %v", err)
	}
	// c.UpdateTimeTable currently populates each line with bulletSchedule.
	// remove the other instances
	delete(c.timetable, "Limited")
	delete(c.timetable, "Limited A")
	delete(c.timetable, "Limited B")
	delete(c.timetable, "Local")
	delete(c.timetable, "Special")

	tests := []struct {
		stops []Station
		numN  int // len of array for now
		numS  int
		day   time.Weekday
		err   error
	}{
		{stops: []Station{StationPaloAlto, StationHillsdale, StationSanFrancisco}, numN: 5, numS: 5, day: time.Monday, err: nil},
		{stops: []Station{StationSanJose, StationSanFrancisco}, numN: 11, numS: 11, day: time.Monday, err: nil},
		{stops: []Station{StationPaloAlto, StationHaywardPark, StationSanFrancisco}, numN: 0, numS: 0, day: time.Monday, err: nil},
		{stops: []Station{StationSanJose, 999}, numN: 0, numS: 0, day: time.Monday, err: errors.New("")},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("test %d", i)
		t.Run(name, func(t *testing.T) {
			// test north
			d1, err := c.getTrainRoutesForAllStops(tt.stops, North, tt.day)
			if err != nil && tt.err == nil {
				t.Fatalf("Failed to get train routes for %s: %v", name, err)
			} else if err == nil && tt.err != nil {
				t.Fatalf("getTrainRoutesForAllStops improperly succeeded for %s", name)
			}
			if len(d1) != tt.numN {
				t.Fatalf("Incorrect routes North. Expected %d, received %d", tt.numN, len(d1))
			}

			// test south
			d2, err := c.getTrainRoutesForAllStops(tt.stops, South, tt.day)
			if err != nil && tt.err == nil {
				t.Fatalf("Failed to get train routes for %s: %v", name, err)
			} else if err == nil && tt.err != nil {
				t.Fatalf("getTrainRoutesForAllStops improperly succeeded for %s", name)
			}
			if len(d2) != tt.numS {
				t.Fatalf("Incorrect routes North. Expected %d, received %d", tt.numS, len(d2))
			}
		})
	}
}

func TestIsInDayRef(t *testing.T) {
	c := New(fakeKey)
	services := map[string][]string{
		"8005": {"monday", "tuesday", "wednesday", "thursday", "friday"},
		"8006": {"saturday", "sunday"},
		"8007": {"saturday"},
	}
	c.dayService = services

	tests := []struct {
		day string
		ref string
		exp bool
	}{
		{day: "monday", ref: "8005", exp: true},
		{day: "tuesday", ref: "8005", exp: true},
		{day: "wednesday", ref: "8005", exp: true},
		{day: "thursday", ref: "8005", exp: true},
		{day: "friday", ref: "8005", exp: true},
		{day: "saturday", ref: "8006", exp: true},
		{day: "saturday", ref: "8007", exp: true},
		{day: "sunday", ref: "8006", exp: true},
		{day: "monday", ref: "8007", exp: false},
		{day: "tuesday", ref: "8006", exp: false},
		{day: "saturday", ref: "8005", exp: false},
		{day: "sunday", ref: "8005", exp: false},
	}
	for _, tt := range tests {
		t.Run(tt.day+"/"+tt.ref, func(t *testing.T) {
			val := c.isForToday(tt.day, tt.ref)
			if val != tt.exp {
				t.Fatalf("isForToday unexpectedly returned %t", val)
			}

		})
	}
}

func TestIsMyDirection(t *testing.T) {
	tests := []struct {
		str string
		dir Direction
		exp bool
	}{
		{str: "Bullet:N :Year Round Weekday (Weekday)", dir: North, exp: true},
		{str: "Bullet:N :Year Round Weekday (Weekday)", dir: South, exp: false},
		{str: "Bullet:S :Year Round Weekday (Weekday)", dir: North, exp: false},
		{str: "Bullet:S :Year Round Weekday (Weekday)", dir: South, exp: true},
		{str: "Local:N :Year Round Weekend (Weekend)", dir: North, exp: true},
		{str: "Local:S :Year Round Weekend (Weekend)", dir: South, exp: true},
		{str: "Limited:S :Year Round Weekday (Weekday)", dir: North, exp: false},
	}
	for _, tt := range tests {
		t.Run(tt.str+"/"+tt.dir.String(), func(t *testing.T) {
			val := isMyDirection(tt.str, tt.dir)
			if val != tt.exp {
				t.Fatalf("isMyDirection unexpectedly returned %t", val)
			}
		})
	}
}

func TestIsStationInJourney(t *testing.T) {}
