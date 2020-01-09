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
	GetDelays(context.Context) ([]Train, error)
	GetStationStatus(context.Context, string, string) ([]Train, error)
	GetTrainsBetweenStations(ctx context.Context, src, dst string) ([]*Route, []*Route, error)
	GetStations() []string
	SetupCache(time.Duration)
	Initialize(context.Context) error
	UpdateStations(context.Context) error
	UpdateTimeTable(context.Context) error
}

type CaltrainClient struct {
	timetable  map[string][]TimetableFrame // map of line type to slice of service journeys
	dayService map[string][]string         // map of id to days of the week that the id corresponds to
	ttLock     sync.RWMutex                // lock in case someone tries to access the timetable during and update
	stations   map[string]*station         // station information map
	sLock      sync.RWMutex                // lock in case someone tries to access the stations during and update
	useCache   bool                        // set by calling the SetupCache method

	key string // API key for 511.org

	DelayThreshold time.Duration // delay time to allow before warning user
	APIClient      APIClient     // API client for making caltrain queries. Default APIClient511
	Updater        Updater       // interface for applying real world updates, such as the day of the week
	Cache          Cache         // interface for caching recent request results
}

func New(key string) *CaltrainClient {
	return &CaltrainClient{
		timetable:      make(map[string][]TimetableFrame),
		dayService:     make(map[string][]string),
		key:            key,
		DelayThreshold: defaultDelayThreshold,
		APIClient:      NewClient(),
		Updater:        NewUpdater(),
	}
}

// Information on the current train status
type Train struct {
	Number    string        // train number
	NextStop  string        // stop for information
	Direction string        // North or South
	Delay     time.Duration // time behind schedule
	Arrival   time.Time     // expected arrival time at NextStop
	Line      string        // bullet, limited, etc.
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

// SetupCache defines enables use of endpoint caching
func (c *CaltrainClient) SetupCache(expire time.Duration) {
	c.Cache = NewCache(expire)
	c.useCache = true
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

func (c *CaltrainClient) GetTimetable() {
	c.ttLock.Lock()
	defer c.ttLock.Unlock()
	// TODO: return something useful?
	fmt.Printf("%+v\n", c.timetable)
}

// GetDelays returns a list of delayed trains and their information
func (c *CaltrainClient) GetDelays(ctx context.Context) ([]Train, error) {
	query := map[string]string{
		"agency":  "CT",
		"api_key": c.key,
	}

	if c.useCache {
		data, ok := c.Cache.get(delayURL)
		if ok {
			return parseDelays(data, c.DelayThreshold)
		}
	}

	data, err := c.APIClient.Get(ctx, delayURL, query)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	// Now parse the body json string
	trains, err := parseDelays(data, c.DelayThreshold)
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
func (c *CaltrainClient) GetStationStatus(ctx context.Context, stationName string, direction string) ([]Train, error) {
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

// Initialize calls the methods to get the timetable and stations, and any
// other required info before the package can run properly
func (c *CaltrainClient) Initialize(ctx context.Context) error {
	if err := c.UpdateStations(ctx); err != nil {
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

	// if the source is greater than destination, it's moving south
	if srcSt.latitude > dstSt.latitude {
		return srcSt.southCode, dstSt.southCode, nil
	} else {
		return srcSt.northCode, dstSt.northCode, nil
	}
}

// GetStationTimetable returns the routes that stop at a given station in the
// given direction
func (c *CaltrainClient) GetStationTimetable(st, dir string) ([]TimetableRouteJourney, error) {
	c.ttLock.RLock()
	defer c.ttLock.RUnlock()

	code, err := c.getStationCode(st, dir)
	if err != nil {
		return nil, err
	}

	weekday, err := c.Updater.GetWeekday(timezone)
	if err != nil {
		return nil, err
	}

	return c.getTimetableForStation(code, dir, weekday)

}

// GetTrainRoute returns the Route struct for a given train
func (c *CaltrainClient) GetTrainRoute(trainNum string) (*Route, error) {
	c.ttLock.RLock()
	defer c.ttLock.RUnlock()
	journey, err := c.getRouteForTrain(trainNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get Train Route: %w", err)
	}
	return c.journeyToRoute(journey)
}

// GetTrainsBetweenStations a slice of routes from src to dst
func (c *CaltrainClient) GetTrainsBetweenStations(ctx context.Context, src, dst string) ([]*Route, error) {
	c.ttLock.RLock()
	defer c.ttLock.RUnlock()
	weekday, err := c.Updater.GetWeekday(timezone)
	if err != nil {
		return nil, err
	}

	journeys, err := c.getTrainRoutesBetweenStations(src, dst, weekday)
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
			return route, fmt.Errorf("could not parse timem from %s: %w", s.Arrival.Time, err)
		}
		dep, err := time.Parse("15:04:05", s.Departure.Time)
		if err != nil {
			return route, fmt.Errorf("could not parse timem from %s: %w", s.Departure.Time, err)
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
