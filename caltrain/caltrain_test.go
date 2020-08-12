package caltrain

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

const (
	fakeKey               = "ericisgreat"
	defaultDelayThreshold = 10 * time.Minute
)

func TestGetStations(t *testing.T) {
	exp := map[Station]struct{}{
		Station22ndStreet:   {},
		StationAtherton:     {},
		StationBayshore:     {},
		StationBelmont:      {},
		StationBlossomHill:  {},
		StationBroadway:     {},
		StationBurlingame:   {},
		StationCalAve:       {},
		StationCapitol:      {},
		StationCollegePark:  {},
		StationGilroy:       {},
		StationHaywardPark:  {},
		StationHillsdale:    {},
		StationLawrence:     {},
		StationMenloPark:    {},
		StationMillbrae:     {},
		StationMorganHill:   {},
		StationMountainView: {},
		StationPaloAlto:     {},
		StationRedwoodCity:  {},
		StationSanAntonio:   {},
		StationSanBruno:     {},
		StationSanCarlos:    {},
		StationSanFrancisco: {},
		StationSanJose:      {},
		StationSanMartin:    {},
		StationSanMateo:     {},
		StationSantaClara:   {},
		StationSouthSF:      {},
		StationStanford:     {},
		StationSunnyvale:    {},
		StationTamien:       {},
	}
	stations := GetStations()
	if len(exp) != len(stations) {
		t.Fatalf("incorrect number of stations")
	}
	for _, st := range stations {
		if _, ok := exp[st]; !ok {
			t.Fatalf("unexpected station: %s", st)
		}
	}
}

func TestGetTrainRoute(t *testing.T) {
	ctx := context.Background()
	c := New(fakeKey)
	m := &apiClientMock{}
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
			{Order: 1, Station: StationSanJose, Arrival: time.Date(0, time.January, 1, 9, 51, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 9, 51, 0, 0, time.UTC)},
			{Order: 2, Station: StationSunnyvale, Arrival: time.Date(0, time.January, 1, 10, 1, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 10, 1, 0, 0, time.UTC)},
			{Order: 3, Station: StationMountainView, Arrival: time.Date(0, time.January, 1, 10, 6, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 10, 6, 0, 0, time.UTC)},
			{Order: 4, Station: StationPaloAlto, Arrival: time.Date(0, time.January, 1, 10, 13, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 10, 13, 0, 0, time.UTC)},
			{Order: 5, Station: StationRedwoodCity, Arrival: time.Date(0, time.January, 1, 10, 20, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 10, 20, 0, 0, time.UTC)},
			{Order: 6, Station: StationHillsdale, Arrival: time.Date(0, time.January, 1, 10, 27, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 10, 27, 0, 0, time.UTC)},
			{Order: 7, Station: StationSanMateo, Arrival: time.Date(0, time.January, 1, 10, 32, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 10, 32, 0, 0, time.UTC)},
			{Order: 8, Station: StationMillbrae, Arrival: time.Date(0, time.January, 1, 10, 40, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 10, 40, 0, 0, time.UTC)},
			{Order: 9, Station: StationSanFrancisco, Arrival: time.Date(0, time.January, 1, 11, 0, 0, 0, time.UTC), Departure: time.Date(0, time.January, 1, 11, 0, 0, 0, time.UTC)},
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
func TestGetTrainsBetweenStationsForWeekday(t *testing.T) {
	ctx := context.Background()
	c := New(fakeKey)
	m := &apiClientMock{}
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
		{src: StationSanFrancisco, dst: 9999, numN: 0, numS: 0, day: time.Monday, err: errors.New("")},
	}

	for _, tt := range tests {
		name := tt.src.String() + "_" + tt.dst.String()
		t.Run(name, func(t *testing.T) {
			// verify north
			d1, err := c.GetTrainsBetweenStationsForWeekday(ctx, tt.src, tt.dst, tt.day)
			if err != nil && tt.err == nil {
				t.Fatalf("Failed to get train routes for %s: %v", name, err)
			} else if err == nil && tt.err != nil {
				t.Fatalf("GetTrainsBetweenStationsForWeekday improperly succeeded for %s", name)
			}
			if len(d1) != tt.numN {
				t.Fatalf("Incorrect routes North. Expected %d, received %d", tt.numN, len(d1))
			}

			// verify south
			d2, err := c.GetTrainsBetweenStationsForWeekday(ctx, tt.dst, tt.src, tt.day)
			if err != nil && tt.err == nil {
				t.Fatalf("Failed to get train routes for %s: %v", name, err)
			} else if err == nil && tt.err != nil {
				t.Fatalf("GetTrainsBetweenStationsForWeekday improperly succeeded for %s", name)
			}
			if len(d2) != tt.numS {
				t.Fatalf("Incorrect routes North. Expected %d, received %d", tt.numS, len(d2))
			}
		})
	}
}

// Simple test to ensure the code runs
func TestGetTrainsBetweenStationsForDate(t *testing.T) {
	ctx := context.Background()
	c := New(fakeKey)
	m := &apiClientMock{}
	m.GetResultFilePath = "testdata/bulletSchedule.json"
	c.APIClient = m
	if err := c.UpdateTimeTable(ctx); err != nil {
		t.Fatalf("Unexpected error loading timetable: %v", err)
	}
	m.GetResultFilePath = "testdata/stations.json"
	if err := c.UpdateStations(ctx); err != nil {
		t.Fatalf("Unexpected error loading stations: %v", err)
	}
	m.GetResultFilePath = "testdata/holiday.json"
	if err := c.UpdateHolidays(ctx); err != nil {
		t.Fatalf("Unexpected error loading holidays: %v", err)
	}
	// c.UpdateTimeTable currently populates each line with bulletSchedule.
	// remove the other instances
	delete(c.timetable, Limited)
	delete(c.timetable, Local)

	tests := []struct {
		name string
		src  Station
		dst  Station
		numN int // len of array for now
		numS int
		day  time.Time
		err  error
	}{
		{name: "Weekday-No-Holiday", src: StationSanJose, dst: StationSanFrancisco, numN: 11, numS: 11, day: time.Date(2019, time.November, 22, 0, 0, 0, 0, time.UTC), err: nil},
		{name: "Holiday", src: StationSanJose, dst: StationSanFrancisco, numN: 2, numS: 2, day: time.Date(2019, time.November, 23, 0, 0, 0, 0, time.UTC), err: nil},
		{name: "Error", src: StationSanFrancisco, dst: 999, numN: 0, numS: 0, day: time.Date(2019, time.November, 23, 0, 0, 0, 0, time.UTC), err: errors.New("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// verify north
			d1, err := c.GetTrainsBetweenStationsForDate(ctx, tt.src, tt.dst, tt.day)
			if err != nil && tt.err == nil {
				t.Fatalf("Failed to get train routes for %s: %v", tt.name, err)
			} else if err == nil && tt.err != nil {
				t.Fatalf("GetTrainsBetweenStationsForWeekday improperly succeeded for %s", tt.name)
			}
			if len(d1) != tt.numN {
				t.Fatalf("Incorrect routes North. Expected %d, received %d", tt.numN, len(d1))
			}

			// verify south
			d2, err := c.GetTrainsBetweenStationsForDate(ctx, tt.dst, tt.src, tt.day)
			if err != nil && tt.err == nil {
				t.Fatalf("Failed to get train routes for %s: %v", tt.name, err)
			} else if err == nil && tt.err != nil {
				t.Fatalf("GetTrainsBetweenStationsForWeekday improperly succeeded for %s", tt.name)
			}
			if len(d2) != tt.numS {
				t.Fatalf("Incorrect routes North. Expected %d, received %d", tt.numS, len(d2))
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
			m := &apiClientMock{}
			m.GetResultFilePath = tt.data
			c.APIClient = m
			d, _, err := c.GetDelays(ctx, defaultDelayThreshold)
			if err != nil && tt.err == nil {
				t.Fatalf("Failed to get train delays for %s: %v", tt.name, err)
			} else if err == nil && tt.err != nil {
				t.Fatalf("GetDelays improperly succeeded for %s", tt.name)
			}

			if len(d) != tt.delays {
				t.Fatalf("Incorrect number of delays. Expected %d, received %d", tt.delays, len(d))
			}
		})
	}
}

func TestGetDelaysCache(t *testing.T) {
	ctx := context.Background()
	c := New(fakeKey)
	m := &apiClientMock{}
	m.GetResultFilePath = "testdata/parseDelayData1.json"
	c.APIClient = m
	c.SetupCache(defaultCacheTimeout)

	cache := make(map[string][]byte)
	mock := &mockCache{}
	mock.SetFunc = func(key string, body []byte) { cache[key] = body }
	mock.GetFunc = func(key string) ([]byte, time.Time, bool) {
		v, ok := cache[key]
		return v, time.Now(), ok
	}
	c.cache = mock

	// first make a call, the cache should be empty
	if len(cache) != 0 {
		t.Fatalf("Cache is not empty: %v", cache)
	}
	d, _, err := c.GetDelays(ctx, defaultDelayThreshold)
	if err != nil {
		t.Fatalf("Failed to get train delays for %v", err)
	}
	if len(d) != 2 {
		t.Fatalf("Incorrect number of delays. Expected %d, received %d", 2, len(d))
	}

	// check that the cache was filled
	if len(cache) != 1 {
		t.Fatalf("Cache does not have only 1 key: %v", cache)
	}
	// run it again
	d, _, err = c.GetDelays(ctx, defaultDelayThreshold)
	if err != nil {
		t.Fatalf("Failed to get train delays for %v", err)
	}
	if len(d) != 2 {
		t.Fatalf("Incorrect number of delays. Expected %d, received %d", 2, len(d))
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
	m := &apiClientMock{}
	m.GetResultFilePath = "testdata/stations.json"
	c.APIClient = m
	if err := c.UpdateStations(ctx); err != nil {
		t.Fatalf("Unexpected error loading stations: %v", err)
	}
	m.GetResultFilePath = "testdata/parseHillsdaleNorth.json"
	_, _, err := c.GetStationStatus(ctx, StationHillsdale, North)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestGetStationStatusCache(t *testing.T) {
	ctx := context.Background()
	c := New(fakeKey)
	m := &apiClientMock{}
	m.GetResultFilePath = "testdata/stations.json"
	c.APIClient = m
	if err := c.UpdateStations(ctx); err != nil {
		t.Fatalf("Unexpected error loading stations: %v", err)
	}
	m.GetResultFilePath = "testdata/parseHillsdaleNorth.json"
	c.SetupCache(defaultCacheTimeout)

	cache := make(map[string][]byte)
	mock := &mockCache{}
	mock.SetFunc = func(key string, body []byte) { cache[key] = body }
	mock.GetFunc = func(key string) ([]byte, time.Time, bool) {
		v, ok := cache[key]
		return v, time.Now(), ok
	}
	c.cache = mock

	// first make a call, the cache should be empty
	if len(cache) != 0 {
		t.Fatalf("Cache is not empty: %v", cache)
	}
	d, _, err := c.GetStationStatus(ctx, StationHillsdale, North)
	if err != nil {
		t.Fatalf("Failed to get train delays for %v", err)
	}
	if len(d) != 1 {
		t.Fatalf("Incorrect number of delays. Expected %d, received %d", 1, len(d))
	}
	// check that the cache was filled
	if len(cache) != 1 {
		t.Fatalf("Cache does not have only 1 key: %v", cache)
	}

	// Now replace it with a South call
	m.GetResultFilePath = "testdata/parseHillsdaleSouth.json"
	c.APIClient = m
	// make a south call
	d, _, err = c.GetStationStatus(ctx, StationHillsdale, South)
	if err != nil {
		t.Fatalf("Failed to get train delays for %v", err)
	}
	if len(d) != 2 {
		t.Fatalf("Incorrect number of delays. Expected %d, received %d", 2, len(d))
	}

	// check that the cache was not changed
	if len(cache) != 2 {
		t.Fatalf("Cache does not have only 1 key: %v", cache)
	}
	// make a the same call call
	d, _, err = c.GetStationStatus(ctx, StationHillsdale, South)
	if err != nil {
		t.Fatalf("Failed to get train delays for %v", err)
	}
	if len(d) != 2 {
		t.Fatalf("Incorrect number of delays. Expected %d, received %d", 2, len(d))
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
		{name: "Bullet", filepath: "testdata/bulletSchedule.json"},
		{name: "Limited", filepath: "testdata/limitedSchedule.json"},
		{name: "Local", filepath: "testdata/localSchedule.json"},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(fakeKey)
			m := &apiClientMock{}
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
	m := &apiClientMock{}
	m.GetResultFilePath = "testdata/stations.json"
	c.APIClient = m
	if err := c.UpdateStations(ctx); err != nil {
		t.Fatalf("Unexpected error loading stations: %v", err)
	}

	_, err := c.GetStationTimetable(StationHillsdale, North, time.Date(2019, time.November, 23, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// Simple test to ensure the code runs
func TestUpdateHolidays(t *testing.T) {
	ctx := context.Background()
	c := New(fakeKey)
	m := &apiClientMock{}
	m.GetResultFilePath = "testdata/holiday.json"
	c.APIClient = m
	if err := c.UpdateHolidays(ctx); err != nil {
		t.Fatalf("Unexpected error loading holidays: %v", err)
	}
}

func TestIsHoliday(t *testing.T) {
	ctx := context.Background()
	c := New(fakeKey)
	m := &apiClientMock{}
	m.GetResultFilePath = "testdata/holiday.json"
	c.APIClient = m
	if err := c.UpdateHolidays(ctx); err != nil {
		t.Fatalf("Unexpected error loading holidays: %v", err)
	}

	tests := []struct {
		date time.Time
		ret  bool
	}{
		{date: time.Date(2019, time.November, 23, 0, 0, 0, 0, time.UTC), ret: true},
		{date: time.Date(2020, time.January, 20, 0, 0, 0, 0, time.UTC), ret: true},
		{date: time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC), ret: false},
	}
	for _, tt := range tests {
		t.Run(tt.date.String(), func(t *testing.T) {
			h := c.IsHoliday(tt.date)
			if h != tt.ret {
				t.Fatalf("IsHoliday error. %s improperly returned %t", tt.date, h)
			}
		})
	}

}
