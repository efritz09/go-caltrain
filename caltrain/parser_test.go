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
				Train{Number: "258", NextStop: StationSunnyvale, Direction: South, Delay: delay1, Arrival: time.Date(2019, time.December, 25, 0, 58, 10, 0, time.UTC), Line: Limited},
				Train{Number: "263", NextStop: StationPaloAlto, Direction: North, Delay: delay2, Arrival: time.Date(2019, time.December, 25, 0, 50, 01, 0, time.UTC), Line: Limited},
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
				Train{Number: "436", NextStop: StationHillsdale, Direction: South, Delay: 0, Arrival: time.Date(2019, time.December, 30, 3, 6, 57, 0, time.UTC), Line: Local},
				Train{Number: "804", NextStop: StationHillsdale, Direction: South, Delay: 0, Arrival: time.Date(2019, time.December, 30, 3, 59, 45, 0, time.UTC), Line: Bullet},
			},
			err: nil,
		},
		{
			name: "HillsdaleNorth",
			data: "testdata/parseHillsdaleNorth.json",
			expected: []Train{
				Train{Number: "437", NextStop: StationHillsdale, Direction: North, Delay: 0, Arrival: time.Date(2019, time.December, 30, 4, 4, 45, 0, time.UTC), Line: Local},
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

			_, _, err = parseTimetable(data)
			if err != nil && tt.err == nil {
				t.Fatalf("Failed to get timetable for %s: %v", tt.name, err)
			} else if err == nil && tt.err != nil {
				t.Fatalf("parseTimetable improperly succeeded for %s", tt.name)
			}
		})
	}
}

func TestParseStations(t *testing.T) {
	f, err := os.Open("testdata/stations.json")
	if err != nil {
		t.Fatalf("Could not open test data: %v", err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Could not read test data: %v", err)
	}

	exp := map[string]*station{
		Station22ndStreet:   &station{name: Station22ndStreet, northCode: "70021", southCode: "70022"},
		StationAtherton:     &station{name: StationAtherton, northCode: "70151", southCode: "70152"},
		StationBayshore:     &station{name: StationBayshore, northCode: "70031", southCode: "70032"},
		StationBelmont:      &station{name: StationBelmont, northCode: "70121", southCode: "70122"},
		StationBlossomHill:  &station{name: StationBlossomHill, northCode: "70291", southCode: "70292"},
		StationBroadway:     &station{name: StationBroadway, northCode: "70071", southCode: "70072"},
		StationBurlingame:   &station{name: StationBurlingame, northCode: "70081", southCode: "70082"},
		StationCalAve:       &station{name: StationCalAve, northCode: "70191", southCode: "70192"},
		StationCapitol:      &station{name: StationCapitol, northCode: "70281", southCode: "70282"},
		StationCollegePark:  &station{name: StationCollegePark, northCode: "70251", southCode: "70252"},
		StationGilroy:       &station{name: StationGilroy, northCode: "70321", southCode: "70322"},
		StationHaywardPark:  &station{name: StationHaywardPark, northCode: "70101", southCode: "70102"},
		StationHillsdale:    &station{name: StationHillsdale, northCode: "70111", southCode: "70112"},
		StationLawrence:     &station{name: StationLawrence, northCode: "70231", southCode: "70232"},
		StationMenloPark:    &station{name: StationMenloPark, northCode: "70161", southCode: "70162"},
		StationMillbrae:     &station{name: StationMillbrae, northCode: "70061", southCode: "70062"},
		StationMorganHill:   &station{name: StationMorganHill, northCode: "70301", southCode: "70302"},
		StationMountainView: &station{name: StationMountainView, northCode: "70211", southCode: "70212"},
		StationPaloAlto:     &station{name: StationPaloAlto, northCode: "70171", southCode: "70172"},
		StationRedwoodCity:  &station{name: StationRedwoodCity, northCode: "70141", southCode: "70142"},
		StationSanAntonio:   &station{name: StationSanAntonio, northCode: "70201", southCode: "70202"},
		StationSanBruno:     &station{name: StationSanBruno, northCode: "70051", southCode: "70052"},
		StationSanCarlos:    &station{name: StationSanCarlos, northCode: "70131", southCode: "70132"},
		StationSanFrancisco: &station{name: StationSanFrancisco, northCode: "70011", southCode: "70012"},
		StationSanJose:      &station{name: StationSanJose, northCode: "70261", southCode: "70262"},
		StationSanMartin:    &station{name: StationSanMartin, northCode: "70311", southCode: "70312"},
		StationSanMateo:     &station{name: StationSanMateo, northCode: "70091", southCode: "70092"},
		StationSantaClara:   &station{name: StationSantaClara, northCode: "70241", southCode: "70242"},
		StationSouthSF:      &station{name: StationSouthSF, northCode: "70041", southCode: "70042"},
		StationSunnyvale:    &station{name: StationSunnyvale, northCode: "70221", southCode: "70222"},
		StationTamien:       &station{name: StationTamien, northCode: "70271", southCode: "70272"},
		StationStanford:     &station{name: StationStanford, northCode: "2537740", southCode: "2537744"},
	}

	s, err := parseStations(data)
	if err != nil {
		t.Fatalf("failed to get stations: %v", err)
	}

	if !reflect.DeepEqual(s, exp) {
		t.Fatalf("station parsing failed.\nexpected %v\nreceived %v", exp, s)
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
