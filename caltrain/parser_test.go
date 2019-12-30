package caltrain

import (
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
		t.Fatalf("Could not open test data: %v", err)
	}
	// fmt.Println(string(data))
	delay1, _ := time.ParseDuration("12m10s")
	delay2, _ := time.ParseDuration("17m1s")

	expected := []TrainStatus{
		TrainStatus{number: "258", nextStop: StationSunnyvale, direction: South, delay: delay1, line: Limited},
		TrainStatus{number: "263", nextStop: StationPaloAlto, direction: North, delay: delay2, line: Limited},
	}

	delays, err := parseDelays(data, defaultDelayThreshold)
	if err != nil {
		t.Fatalf("Failed to parse delays: %v", err)
	}

	if !assertEqual(expected, delays) {
		t.Fatalf("Unexpected delays!\nexpected: %v\nreceived %v", expected, delays)
	}
}

// assertEqual compares two TrainStatus slices for the same elements
func assertEqual(exp, test []TrainStatus) bool {
	if len(exp) != len(test) {
		return false
	}
	// populate a map with number of instances
	m1 := make(map[TrainStatus]int)
	m2 := make(map[TrainStatus]int)
	for _, k := range exp {
		m1[k]++
	}
	for _, k := range test {
		m2[k]++
	}
	return reflect.DeepEqual(m1, m2)
}
