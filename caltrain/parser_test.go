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
		expected []TrainStatus
		err      error
	}{
		{
			name: "DelayData1",
			data: "testdata/parseDelayData1.json",
			expected: []TrainStatus{
				TrainStatus{TrainNum: "258", NextStop: StationSunnyvale, Direction: South, Delay: delay1, Arrival: time.Date(2019, time.December, 25, 0, 58, 10, 0, time.UTC), Line: Limited},
				TrainStatus{TrainNum: "263", NextStop: StationPaloAlto, Direction: North, Delay: delay2, Arrival: time.Date(2019, time.December, 25, 0, 50, 01, 0, time.UTC), Line: Limited},
			},
			err: nil,
		},
		{
			name:     "DelayData2",
			data:     "testdata/parseDelayData2.json",
			expected: []TrainStatus{},
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

			if !assertTrainStatusEqual(tt.expected, delays) {
				t.Fatalf("Unexpected delays for %s\nexpected: %v\nreceived: %v", tt.name, tt.expected, delays)
			}
		})
	}
}

func TestGetTrains(t *testing.T) {
	tests := []struct {
		name     string
		data     string // relative file location
		expected []TrainStatus
		err      error
	}{
		{
			name: "HillsdaleSouth",
			data: "testdata/parseHillsdaleSouth.json",
			expected: []TrainStatus{
				TrainStatus{TrainNum: "436", NextStop: StationHillsdale, Direction: South, Delay: 0, Arrival: time.Date(2019, time.December, 30, 3, 6, 57, 0, time.UTC), Line: Local},
				TrainStatus{TrainNum: "804", NextStop: StationHillsdale, Direction: South, Delay: 0, Arrival: time.Date(2019, time.December, 30, 3, 59, 45, 0, time.UTC), Line: Bullet},
			},
			err: nil,
		},
		{
			name: "HillsdaleNorth",
			data: "testdata/parseHillsdaleNorth.json",
			expected: []TrainStatus{
				TrainStatus{TrainNum: "437", NextStop: StationHillsdale, Direction: North, Delay: 0, Arrival: time.Date(2019, time.December, 30, 4, 4, 45, 0, time.UTC), Line: Local},
			},
			err: nil,
		},
		{
			name:     "HillsdaleNorthBad",
			data:     "testdata/parseHillsdaleNorthBad.json",
			expected: []TrainStatus{},
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

			if !assertTrainStatusEqual(tt.expected, trains) {
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

	exp := map[Station]*stationInfo{
		Station22ndStreet:   &stationInfo{name: Station22ndStreet, northCode: "70021", southCode: "70022"},
		StationAtherton:     &stationInfo{name: StationAtherton, northCode: "70151", southCode: "70152"},
		StationBayshore:     &stationInfo{name: StationBayshore, northCode: "70031", southCode: "70032"},
		StationBelmont:      &stationInfo{name: StationBelmont, northCode: "70121", southCode: "70122"},
		StationBlossomHill:  &stationInfo{name: StationBlossomHill, northCode: "70291", southCode: "70292"},
		StationBroadway:     &stationInfo{name: StationBroadway, northCode: "70071", southCode: "70072"},
		StationBurlingame:   &stationInfo{name: StationBurlingame, northCode: "70081", southCode: "70082"},
		StationCalAve:       &stationInfo{name: StationCalAve, northCode: "70191", southCode: "70192"},
		StationCapitol:      &stationInfo{name: StationCapitol, northCode: "70281", southCode: "70282"},
		StationCollegePark:  &stationInfo{name: StationCollegePark, northCode: "70251", southCode: "70252"},
		StationGilroy:       &stationInfo{name: StationGilroy, northCode: "70321", southCode: "70322"},
		StationHaywardPark:  &stationInfo{name: StationHaywardPark, northCode: "70101", southCode: "70102"},
		StationHillsdale:    &stationInfo{name: StationHillsdale, northCode: "70111", southCode: "70112"},
		StationLawrence:     &stationInfo{name: StationLawrence, northCode: "70231", southCode: "70232"},
		StationMenloPark:    &stationInfo{name: StationMenloPark, northCode: "70161", southCode: "70162"},
		StationMillbrae:     &stationInfo{name: StationMillbrae, northCode: "70061", southCode: "70062"},
		StationMorganHill:   &stationInfo{name: StationMorganHill, northCode: "70301", southCode: "70302"},
		StationMountainView: &stationInfo{name: StationMountainView, northCode: "70211", southCode: "70212"},
		StationPaloAlto:     &stationInfo{name: StationPaloAlto, northCode: "70171", southCode: "70172"},
		StationRedwoodCity:  &stationInfo{name: StationRedwoodCity, northCode: "70141", southCode: "70142"},
		StationSanAntonio:   &stationInfo{name: StationSanAntonio, northCode: "70201", southCode: "70202"},
		StationSanBruno:     &stationInfo{name: StationSanBruno, northCode: "70051", southCode: "70052"},
		StationSanCarlos:    &stationInfo{name: StationSanCarlos, northCode: "70131", southCode: "70132"},
		StationSanFrancisco: &stationInfo{name: StationSanFrancisco, northCode: "70011", southCode: "70012"},
		StationSanJose:      &stationInfo{name: StationSanJose, northCode: "70261", southCode: "70262"},
		StationSanMartin:    &stationInfo{name: StationSanMartin, northCode: "70311", southCode: "70312"},
		StationSanMateo:     &stationInfo{name: StationSanMateo, northCode: "70091", southCode: "70092"},
		StationSantaClara:   &stationInfo{name: StationSantaClara, northCode: "70241", southCode: "70242"},
		StationSouthSF:      &stationInfo{name: StationSouthSF, northCode: "70041", southCode: "70042"},
		StationSunnyvale:    &stationInfo{name: StationSunnyvale, northCode: "70221", southCode: "70222"},
		StationTamien:       &stationInfo{name: StationTamien, northCode: "70271", southCode: "70272"},
		StationStanford:     &stationInfo{name: StationStanford, northCode: "2537740", southCode: "2537744"},
	}

	s, err := parseStations(data)
	if err != nil {
		t.Fatalf("failed to get stations: %v", err)
	}

	if len(s) != len(exp) {
		t.Fatalf("length mismatch")
	}

	for k, v := range exp {
		st, ok := s[k]
		if !ok {
			t.Fatalf("missing station %s", k)
		}
		if st.northCode != v.northCode {
			t.Fatalf("conflicting northCodes: %s vs %s", v.northCode, st.northCode)
		}
		if st.southCode != v.southCode {
			t.Fatalf("conflicting southCodes: %s vs %s", v.southCode, st.southCode)
		}
	}
}

func TestParseHolidays(t *testing.T) {
	f, err := os.Open("testdata/holiday.json")
	if err != nil {
		t.Fatalf("Could not open test data: %v", err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Could not read test data: %v", err)
	}

	exp := []time.Time{
		time.Date(2019, time.November, 23, 0, 0, 0, 0, time.UTC),
		time.Date(2019, time.November, 24, 0, 0, 0, 0, time.UTC),
		time.Date(2019, time.November, 28, 0, 0, 0, 0, time.UTC),
		time.Date(2019, time.November, 29, 0, 0, 0, 0, time.UTC),
		time.Date(2019, time.December, 25, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.January, 20, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.February, 17, 0, 0, 0, 0, time.UTC),
	}

	holidays, err := parseHolidays(data)
	if err != nil {
		t.Fatalf("Failed to parse holidays: %v", err)
	}

	if len(exp) != len(holidays) {
		t.Fatalf("Unexpected holidays\nexpected: %v\nreceived: %v", exp, holidays)
	}

	m1 := make(map[time.Time]int)
	m2 := make(map[time.Time]int)
	for _, k := range exp {
		m1[k]++
	}
	for _, k := range holidays {
		m2[k]++
	}
	if !reflect.DeepEqual(m1, m2) {
		t.Fatalf("Unexpected holidays\nexpected: %v\nreceived: %v", exp, holidays)
	}
}

// assertTrainStatusEqual compares two TrainStatus slices for the same elements
func assertTrainStatusEqual(exp, test []TrainStatus) bool {
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
