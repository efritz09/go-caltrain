package caltrain

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestParseDelays(t *testing.T) {
	delay1, _ := time.ParseDuration("12m10s")
	delay2, _ := time.ParseDuration("17m1s")
	tests := []struct {
		name     string
		data     string
		expected []Train
		err      error
	}{
		{
			name: "DelayData1",
			data: "testdata/parseDelayData1.json",
			expected: []Train{
				Train{Number: "258", NextStop: StationSunnyvale, Direction: South, Delay: delay1, Line: Limited},
				Train{Number: "263", NextStop: StationPaloAlto, Direction: North, Delay: delay2, Line: Limited},
			},
			err: nil,
		},
		{
			name:     "DelayData2",
			data:     "testdata/parseDelayData2.json",
			expected: []Train{},
			err:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.data)
			if err != nil {
				t.Fatalf("Could not open test data for %s: %v", tt.name, err)
			}
			data, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatalf("Could not read test data for %s: %v", tt.name, err)
			}

			delays, err := parseDelays(data, defaultDelayThreshold)
			if err != nil && tt.err == nil {
				t.Fatalf("Failed to get trains for %s: %v", tt.name, err)
			} else if err == nil && tt.err != nil {
				t.Fatalf("getTrains improperly succeeded for %s", tt.name)
			}

			if !assertEqual(tt.expected, delays) {
				t.Fatalf("Unexpected delays for %s\nexpected: %v\nreceived: %v", tt.name, tt.expected, delays)
			}
		})
	}
}

func TestGetTrains(t *testing.T) {
	tests := []struct {
		name     string
		data     string // relative file location
		expected []Train
		err      error
	}{
		{
			name: "HillsdaleSouth",
			data: "testdata/parseHillsdaleSouth.json",
			expected: []Train{
				Train{Number: "436", NextStop: StationHillsdale, Direction: South, Delay: 0, Line: Local},
				Train{Number: "804", NextStop: StationHillsdale, Direction: South, Delay: 0, Line: Bullet},
			},
			err: nil,
		},
		{
			name: "HillsdaleNorth",
			data: "testdata/parseHillsdaleNorth.json",
			expected: []Train{
				Train{Number: "437", NextStop: StationHillsdale, Direction: North, Delay: 0, Line: Local},
			},
			err: nil,
		},
		{
			name:     "HillsdaleNorthBad",
			data:     "testdata/parseHillsdaleNorthBad.json",
			expected: []Train{},
			err:      errors.New(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.data)
			if err != nil {
				t.Fatalf("Could not open test data for %s: %v", tt.name, err)
			}
			data, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatalf("Could not read test data for %s: %v", tt.name, err)
			}

			trains, err := getTrains(data)
			if err != nil && tt.err == nil {
				t.Fatalf("Failed to get trains for %s: %v", tt.name, err)
			} else if err == nil && tt.err != nil {
				t.Fatalf("getTrains improperly succeeded for %s", tt.name)
			}

			if !assertEqual(tt.expected, trains) {
				t.Fatalf("Unexpected trains for %s\nexpected: %v\nreceived: %v", tt.name, tt.expected, trains)
			}
		})
	}
}

func TestParseTimetable(t *testing.T) {
	tests := []struct {
		name string
		data string // relative file location
		err  error
	}{
		{
			name: "Bullet",
			data: "testdata/bulletSchedule.json",
			err:  nil,
		},
		{
			name: "Limited",
			data: "testdata/limitedSchedule.json",
			err:  nil,
		},
		{
			name: "Local",
			data: "testdata/localSchedule.json",
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.data)
			if err != nil {
				t.Fatalf("Could not open test data for %s: %v", tt.name, err)
			}
			data, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatalf("Could not read test data for %s: %v", tt.name, err)
			}

			timetable, err := parseTimetable(data)
			if err != nil && tt.err == nil {
				t.Fatalf("Failed to get timetable for %s: %v", tt.name, err)
			} else if err == nil && tt.err != nil {
				t.Fatalf("parseTimetable improperly succeeded for %s", tt.name)
			}
		})
	}
}

// assertEqual compares two Train slices for the same elements
func assertEqual(exp, test []Train) bool {
	if len(exp) != len(test) {
		return false
	}
	// populate a map with number of instances
	m1 := make(map[Train]int)
	m2 := make(map[Train]int)
	for _, k := range exp {
		m1[k]++
	}
	for _, k := range test {
		m2[k]++
	}
	return reflect.DeepEqual(m1, m2)
}
