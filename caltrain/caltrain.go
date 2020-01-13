// caltrain provides a user API for getting live caltrain updates
package caltrain

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Caltrain interface {
	// Initialize makes the 511.org API calls to populate the stations and
	// timetable. It calls UpdateStations, UpdateTimetable, and UpdateHolidays
	Initialize(context.Context) error

	// SetupCache enables the use of API caching to prevent going over the API
	// limit. Users set the caching expire time.
	SetupCache(time.Duration)

	// UpdateStations makes an API call to refresh the station information.
	// This should only need to be called during Initialization.
	UpdateStations(context.Context) error

	// UpdateTimeTable makes an API call to refresh the timetable data. This
	// should be called periodically to ensure correct information.
	UpdateTimeTable(context.Context) error

	// UpdateHolidays makes an API call to refresh the holiday data. This can
	// be updated multiple times a year so this should be called periodically.
	UpdateHolidays(context.Context) error

	// GetDelays makes an API call and returns a slice of TrainStatus who's
	// delay into their next station is greater than the time.Duration argument
	GetDelays(context.Context, time.Duration) ([]TrainStatus, error)

	// GetStationStatus makes an API call and returns a slice of TrainsStatus
	// who have a status reported for the given station and direction.
	GetStationStatus(ctx context.Context, stationName, dir string) ([]TrainStatus, error)

	// GetTrainsBetweenStations returns a slice of Routes that go from src to
	// dst on the given weekday. It uses the cached timetable and does not make
	// an API call
	GetTrainsBetweenStations(ctx context.Context, src, dst string, weekday time.Weekday) ([]*Route, error)

	// GetDirectionFromSrcToDst returns the direction the train would go to get
	// from src to dst. Value is either North or South
	GetDirectionFromSrcToDst(src, dst string) (string, error)

	// GetStations returns a slice of all known stations in alphanumeric order
	GetStations() []string
}

type CaltrainClient struct {
	timetable  map[string][]TimetableFrame // map of line type to slice of service journeys
	dayService map[string][]string         // map of id to days of the week that the id corresponds to
	ttLock     sync.RWMutex                // lock in case someone tries to access the timetable during and update
	stations   map[string]*station         // station information map
	holidays   []time.Time // slice of days that are on a holiday schedule
	sLock      sync.RWMutex                // lock in case someone tries to access the stations during and update
	useCache   bool                        // set by calling the SetupCache method

	key string // API key for 511.org

	APIClient      APIClient     // API client for making caltrain queries. Default APIClient511
	Cache          Cache         // interface for caching recent request results
}

func New(key string) *CaltrainClient {
	return &CaltrainClient{
		timetable:      make(map[string][]TimetableFrame),
		dayService:     make(map[string][]string),
		key:            key,
		APIClient:      NewClient(),
	}
}

// Information on the current train status
type TrainStatus struct {
	TrainNum    string        // train number
	Direction string        // North or South
	Line      string        // bullet, limited, etc.
	Delay     time.Duration // time behind schedule
	Arrival   time.Time     // expected arrival time at NextStop
	NextStop  string        // stop for information
}

// Stops for a given train
type Route struct {
	TrainNum  string // train number
	Direction string
	Line      string // bullet, limited, etc.
	NumStops  int
	Stops     []TrainStop
}

type TrainStop struct {
	Order     int
	Station   string
	Arrival   time.Time
	Departure time.Time
}

type station struct {
	name      string
	northCode string
	southCode string
	latitude  float64
	longitude float64
}

// Initialize calls the methods to get the timetable and stations, and any
// other required info before the package can run properly
func (c *CaltrainClient) Initialize(ctx context.Context) error {
	if err := c.UpdateStations(ctx); err != nil {
		return err
	}
	if err := c.UpdateHolidays(ctx); err != nil {
		return err
	}
	return c.UpdateTimeTable(ctx)
}

// UpdateTimeTable should be called once per day to update the day's timetable
func (c *CaltrainClient) UpdateTimeTable(ctx context.Context) error {
	c.ttLock.Lock()
	defer c.ttLock.Unlock()
	lines := []string{Bullet, Limited, Local}
	// request the timetable for each line
	for _, line := range lines {
		query := map[string]string{
			"operator_id": "CT",
			"line_id":     line,
			"api_key":     c.key,
		}
		data, err := c.APIClient.Get(ctx, timetableURL, query)
		if err != nil {
			return fmt.Errorf("failed to make request: %w", err)
		}

		journeys, services, err := parseTimetable(data)
		if err != nil {
			return fmt.Errorf("failed to parse timetable: %w", err)
		}
		// store the timetable
		c.timetable[line] = journeys

		// overwrite the known data with the timetable's ServiceCalendarFrame
		for key, value := range services {
			c.dayService[key] = value
		}
	}

	return nil
}

// UpdateStations calls the api to get a station list
func (c *CaltrainClient) UpdateStations(ctx context.Context) error {
	c.sLock.Lock()
	defer c.sLock.Unlock()

	query := map[string]string{
		"operator_id": "CT",
		"api_key":     c.key,
	}
	data, err := c.APIClient.Get(ctx, stationsURL, query)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}

	stations, err := parseStations(data)
	if err != nil {
		return fmt.Errorf("failed to parse stations: %w", err)
	}
	c.stations = stations
	return nil
}

// UpdateHolidays calls the api to get upcoming holidays
func (c *CaltrainClient) UpdateHolidays(ctx context.Context) error {
	c.sLock.Lock()
	defer c.sLock.Unlock()

	query := map[string]string{
		"operator_id": "CT",
		"api_key":     c.key,
	}
	data, err := c.APIClient.Get(ctx, stationsURL, query)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}

	holidays, err := parseHolidays(data)
	if err != nil {
		return fmt.Errorf("failed to parse holidays: %w", err)
	}
	c.holidays = holidays
	return nil 
}


// SetupCache defines enables use of endpoint caching
func (c *CaltrainClient) SetupCache(expire time.Duration) {
	c.Cache = NewCache(expire)
	c.useCache = true
}

// GetDelays returns a list of delayed trains and their information
func (c *CaltrainClient) GetDelays(ctx context.Context, threshold time.Duration) ([]TrainStatus, error) {
	query := map[string]string{
		"agency":  "CT",
		"api_key": c.key,
	}

	if c.useCache {
		data, ok := c.Cache.get(delayURL)
		if ok {
			return parseDelays(data, threshold)
		}
	}

	data, err := c.APIClient.Get(ctx, delayURL, query)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	// Now parse the body json string
	trains, err := parseDelays(data, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to parse delay data: %w", err)
	}

	if c.useCache {
		c.Cache.set(delayURL, data)
	}
	return trains, nil
}

// GetStationStatus returns the status of upcoming trains for a given station
// and direction. Direction should be caltrain.North or caltrain.South
func (c *CaltrainClient) GetStationStatus(ctx context.Context, stationName string, direction string) ([]TrainStatus, error) {
	code, err := c.getStationCode(stationName, direction)
	if err != nil {
		return nil, fmt.Errorf("failed to get station code: %w", err)
	}
	query := map[string]string{
		"agency":   "CT",
		"stopCode": code,
		"api_key":  c.key,
	}

	// cache key is stationStatusURL plus the stop code
	if c.useCache {
		data, ok := c.Cache.get(stationStatusURL + code)
		if ok {
			return getTrains(data)
		}
	}

	data, err := c.APIClient.Get(ctx, stationStatusURL, query)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	// Now parse the body json string
	trains, err := getTrains(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trains: %w", err)
	}

	if c.useCache {
		c.Cache.set(stationStatusURL+code, data)
	}
	return trains, nil
}

// GetTrainsBetweenStations a slice of routes from src to dst
func (c *CaltrainClient) GetTrainsBetweenStations(ctx context.Context, src, dst string, weekday time.Weekday) ([]*Route, error) {
	c.ttLock.RLock()
	defer c.ttLock.RUnlock()

	journeys, err := c.getTrainRoutesBetweenStations(src, dst, strings.ToLower(weekday.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to get Train Routes: %w", err)
	}

	routes := make([]*Route, len(journeys))
	for i, journey := range journeys {
		r, err := c.journeyToRoute(journey)
		if err != nil {
			return routes, fmt.Errorf("failed to get Train Routes: %w", err)
		}
		routes[i] = r
	}
	return routes, nil
}

// GetDirectionFromSrcToDst returns North or South given a src and dst station
// TODO: unit test
func (c *CaltrainClient) GetDirectionFromSrcToDst(src, dst string) (string, error) {
	if src == dst {
		return "", fmt.Errorf("The stations are the same: %s to %s", src, dst)
	}
	s, ok := stationOrder[src]
	if !ok {
		return "", fmt.Errorf("Unknown station: %s", src)
	}
	d, ok := stationOrder[dst]
	if !ok {
		return "", fmt.Errorf("Unknown station: %s", dst)
	}

	if s < d {
		return South, nil
	} else if s > d {
		return North, nil
	} else {
		return "", fmt.Errorf("Could not determine direction from %s to %s", src, dst)
	}
}

// GetStations returns a list of station names
func (c *CaltrainClient) GetStations() []string {
	ret := []string{}
	c.sLock.RLock()
	for k := range c.stations {
		ret = append(ret, k)
	}
	c.sLock.RUnlock()
	sort.Strings(ret)
	return ret
}

// GetStationTimetable returns the routes that stop at a given station in the
// given direction
// TODO: export this in the interface???
func (c *CaltrainClient) GetStationTimetable(st, dir string, weekday time.Weekday) ([]TimetableRouteJourney, error) {
	c.ttLock.RLock()
	defer c.ttLock.RUnlock()

	code, err := c.getStationCode(st, dir)
	if err != nil {
		return nil, err
	}

	return c.getTimetableForStation(code, dir, strings.ToLower(weekday.String()))

}

// GetTrainRoute returns the Route struct for a given train
// TODO: export this in the interface???
func (c *CaltrainClient) GetTrainRoute(trainNum string) (*Route, error) {
	c.ttLock.RLock()
	defer c.ttLock.RUnlock()
	journey, err := c.getRouteForTrain(trainNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get Train Route: %w", err)
	}
	return c.journeyToRoute(journey)
}

// getStationCode returns the code for a given station and direction
func (c *CaltrainClient) getStationCode(st, dir string) (string, error) {
	// first validate the direction
	c.sLock.RLock()
	defer c.sLock.RUnlock()
	station, ok := c.stations[st]
	if !ok {
		return "", fmt.Errorf("unknown station %s", st)
	}

	if dir == North {
		return station.northCode, nil
	} else if dir == South {
		return station.southCode, nil
	} else {
		return "", fmt.Errorf("unknown direction %s", dir)
	}
}

// getRouteDirection returns the proper station codes for a route given a
// source and destination station name
func (c *CaltrainClient) getRouteCodes(src, dst string) (string, string, error) {
	c.sLock.RLock()
	defer c.sLock.RUnlock()
	srcSt, ok := c.stations[src]
	if !ok {
		return "", "", fmt.Errorf("unknown station %s", src)
	}
	dstSt, ok := c.stations[dst]
	if !ok {
		return "", "", fmt.Errorf("unknown station %s", dst)
	}

	dir, err := c.GetDirectionFromSrcToDst(src, dst)
	if err != nil {
		return "", "", err
	}

	// if the source is greater than destination, it's moving south
	if dir == South {
		return srcSt.southCode, dstSt.southCode, nil
	} else {
		return srcSt.northCode, dstSt.northCode, nil
	}
}

// journeyToRoute converts a TimetableRouteJourney into a Route
func (c *CaltrainClient) journeyToRoute(r TimetableRouteJourney) (*Route, error) {
	route := &Route{
		TrainNum:  r.ID,
		Direction: getDirFromChar(r.JourneyPatternView.DirectionRef.Ref),
		Line:      r.Line,
		NumStops:  len(r.Calls.Call),
		Stops:     []TrainStop{},
	}

	for _, s := range r.Calls.Call {
		order, err := strconv.Atoi(s.Order)
		if err != nil {
			return route, fmt.Errorf("could not convert order %s to int: %w", s.Order, err)
		}
		arr, err := time.Parse("15:04:05", s.Arrival.Time)
		if err != nil {
			return route, fmt.Errorf("could not parse time from %s: %w", s.Arrival.Time, err)
		}
		if s.Arrival.DaysOffset == "1" {
			arr = arr.Add(24 * time.Hour)
		}
		dep, err := time.Parse("15:04:05", s.Departure.Time)
		if err != nil {
			return route, fmt.Errorf("could not parse time from %s: %w", s.Departure.Time, err)
		}
		if s.Departure.DaysOffset == "1" {
			dep = dep.Add(24 * time.Hour)
		}
		t := TrainStop{
			Order:     order,
			Station:   c.getStationFromCode(s.ScheduledStopPointRef.Ref),
			Arrival:   arr,
			Departure: dep,
		}
		route.Stops = append(route.Stops, t)
	}
	return route, nil
}

// getStationFromCode returns the station name associated with the code
// TODO: unit test this
func (c *CaltrainClient) getStationFromCode(code string) string {
	c.sLock.RLock()
	defer c.sLock.RUnlock()
	for name, st := range c.stations {
		if st.northCode == code || st.southCode == code {
			return name
		}
	}
	return ""
}

// getDirFromChar returns the proper direction string for a given character.
// HasPrefix is used in case the "char" has whitespace
func getDirFromChar(c string) string {
	if strings.HasPrefix(c, "N") {
		return North
	} else if strings.HasPrefix(c, "S") {
		return South
	} else {
		return ""
	}
}
