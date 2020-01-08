package caltrain

import (
	"context"
	"errors"
	"testing"
)

func TestGetTimetableForStation(t *testing.T) {
	// Load the timetable for only the bullet schedule
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

	tests := []struct {
		station  string
		dir      string
		day      string
		expected int // length of array for now, should be []TimetableRouteJourney
	}{
		{station: StationHillsdale, dir: North, day: "monday", expected: 5},
		{station: StationHillsdale, dir: North, day: "sunday", expected: 2},
		{station: StationHillsdale, dir: South, day: "monday", expected: 5},
		{station: StationHillsdale, dir: South, day: "sunday", expected: 2},
	}

	for _, tt := range tests {
		name := tt.station + "/" + tt.dir + "/" + tt.day
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
			n, s, err := c.getTrainRoutesBetweenStations(tt.src, tt.dst, tt.day)
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

func TestGetRouteForTrain(t *testing.T) {
	// Load the timetable for only the bullet schedule
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

	tests := []struct {
		train string
		line  string
		err   error
	}{
		{train: "801", line: Bullet, err: nil},
		{train: "324", line: Bullet, err: nil},
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

func TestIsInDayRef(t *testing.T) {
	c := New(fakeKey)
	services := map[string][]string{
		"8005": []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
		"8006": []string{"saturday", "sunday"},
		"8007": []string{"saturday"},
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
		dir string
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
		t.Run(tt.str+"/"+tt.dir, func(t *testing.T) {
			val := isMyDirection(tt.str, tt.dir)
			if val != tt.exp {
				t.Fatalf("isMyDirection unexpectedly returned %t", val)
			}
		})
	}
}

func TestIsStationInJourney(t *testing.T) {}
