package caltrain

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestParseDelays(t *testing.T) {
	f, err := os.Open("testdata/parseDelayData.json")
	if err != nil {
		t.Fatalf("Could not open test data: %v", err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Could not read test data: %v", err)
	}
	delay1, _ := time.ParseDuration("12m10s")
	delay2, _ := time.ParseDuration("17m1s")

	expected := []Train{
		Train{number: "258", nextStop: StationSunnyvale, direction: South, delay: delay1, line: Limited},
		Train{number: "263", nextStop: StationPaloAlto, direction: North, delay: delay2, line: Limited},
	}

	delays, err := parseDelays(data, defaultDelayThreshold)
	if err != nil {
		t.Fatalf("Failed to parse delays: %v", err)
	}

	if !assertEqual(expected, delays) {
		t.Fatalf("Unexpected delays!\nexpected: %v\nreceived %v", expected, delays)
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
				Train{number: "436", nextStop: StationHillsdale, direction: South, delay: 0, line: Local},
				Train{number: "804", nextStop: StationHillsdale, direction: South, delay: 0, line: Bullet},
			},
			err: nil,
		},
		{
			name: "HillsdaleNorth",
			data: "testdata/parseHillsdaleNorth.json",
			expected: []Train{
				Train{number: "437", nextStop: StationHillsdale, direction: North, delay: 0, line: Local},
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
			fmt.Println(err)

			if !assertEqual(tt.expected, trains) {
				t.Fatalf("Unexpected trains for %s\nexpected: %v\nreceived: %v", tt.name, tt.expected, trains)
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
