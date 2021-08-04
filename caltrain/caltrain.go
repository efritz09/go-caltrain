package caltrain

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	delayURL         = "http://api.511.org/transit/StopMonitoring"
	holidaysURL      = "http://api.511.org/transit/holidays"
	linesURL         = "http://api.511.org/transit/lines"
	stationsURL      = "http://api.511.org/transit/stops"
	stationStatusURL = "http://api.511.org/transit/StopMonitoring"
	timetableURL     = "http://api.511.org/transit/timetable"
)

// CaltrainClient provides the means for querying information about caltrain
// schedules, getting route information between stations, or getting live train
// status updates
type CaltrainClient struct {
	timetable  map[string][]timetableFrame // map of line name to slice of service journeys
	dayService map[string][]string         // map of id to days of the week that the id corresponds to
	ttLock     sync.RWMutex                // lock in case someone tries to access the timetable during an update
	stations   map[Station]*stationInfo    // station information map
	sLock      sync.RWMutex                // lock in case someone tries to access the stations during an update
	lines      []Line                      // slice of available lines
	lLock      sync.RWMutex                // lock in case someone tries to access the lines during an update
	useCache   bool                        // set by calling the SetupCache method
	holidays   []time.Time                 // slice of days that are on a holiday schedule
	tz         *time.Location              // constant America/LosAngeles time
	key        string                      // API key for 511.org
	cache      cache                       // interface for caching recent request results

	APIClient APIClient // API client for making caltrain queries. Default APIClient511
}

// New returns an instantiated CaltrainClient struct
func New(key string) *CaltrainClient {
	tz, _ := time.LoadLocation("America/Los_Angeles")
	return &CaltrainClient{
		timetable:  make(map[string][]timetableFrame),
		dayService: make(map[string][]string),
		stations:   make(map[Station]*stationInfo),
		lines:      []Line{},
		key:        key,
		tz:         tz,
		APIClient:  NewClient(),
	}
}

// Initialize makes the 511.org API calls to populate the stations and
// timetable. It calls UpdateStations, UpdateTimetable, and UpdateHolidays
func (c *CaltrainClient) Initialize(ctx context.Context) (e error) {
	if err := c.UpdateLines(ctx); err != nil {
		e = fmt.Errorf("failure updating Lines: %w", err)
	}

	if err := c.UpdateStations(ctx); err != nil {
		e = fmt.Errorf("failure updating Stations: %w", err)
	}

	if err := c.UpdateHolidays(ctx); err != nil {
		e = fmt.Errorf("failure updating Holidays: %w", err)
	}

	if err := c.UpdateTimeTable(ctx); err != nil {
		e = fmt.Errorf("failure updating Time Tables: %w", err)
	}
	return
}

// UpdateLines makes an API call to refresh the available lines. This should be
// called before UpdateTimeTable to ensure the time table data is accurate
func (c *CaltrainClient) UpdateLines(ctx context.Context) error {
	logrus.Debug("Updating train lines...")
	c.lLock.Lock()
	defer c.lLock.Unlock()

	query := map[string]string{
		"operator_id": "CT",
		"api_key":     c.key,
	}
	data, err := c.APIClient.Get(ctx, linesURL, query)
	if err != nil {
		return fmt.Errorf("failed to make 'update lines' request: %w", err)
	}

	lines, err := parseLines(data)
	if err != nil {
		return fmt.Errorf("failed to parse lines: %w", err)
	}
	if len(lines) == 0 {
		return errors.New("unable to populate the lines: none found")
	}
	c.lines = lines
	return nil
}

// UpdateTimeTable makes an API call to refresh the timetable data. This
// should be called periodically to ensure correct information.
func (c *CaltrainClient) UpdateTimeTable(ctx context.Context) error {
	logrus.Debug("Updating time tables...")
	c.ttLock.Lock()
	defer c.ttLock.Unlock()
	// request the timetable for each line
	for _, line := range c.lines {
		logrus.Debugf("Fetching time table for %s trains", line.Name)
		query := map[string]string{
			"operator_id": "CT",
			"line_id":     line.Id,
			"api_key":     c.key,
		}
		data, err := c.APIClient.Get(ctx, timetableURL, query)
		if err != nil {
			return fmt.Errorf("failed to make 'update timetable' request: %w", err)
		}

		journeys, services, err := parseTimetable(data)
		if err != nil {
			return fmt.Errorf("failed to parse timetable: %w", err)
		}
		// store the timetable
		c.timetable[line.Name] = journeys

		// overwrite the known data with the timetable's ServiceCalendarFrame
		for key, value := range services {
			c.dayService[key] = value
		}
	}

	if len(c.timetable) == 0 {
		return errors.New("unable to populate the timetables: none found")
	}

	return nil
}

// UpdateStations makes an API call to refresh the station information.
// This should only need to be called during Initialization.
func (c *CaltrainClient) UpdateStations(ctx context.Context) error {
	logrus.Debug("Updating stations...")
	c.sLock.Lock()
	defer c.sLock.Unlock()

	query := map[string]string{
		"operator_id": "CT",
		"api_key":     c.key,
	}
	data, err := c.APIClient.Get(ctx, stationsURL, query)
	if err != nil {
		return fmt.Errorf("failed to make 'update stations' request: %w", err)
	}

	stations, err := parseStations(data)
	if err != nil {
		return fmt.Errorf("failed to parse stations: %w", err)
	}
	if len(stations) == 0 {
		return errors.New("unable to populate the station list: none found")
	}
	c.stations = stations
	return nil
}

// UpdateHolidays makes an API call to refresh the holiday data. This can
// be updated multiple times a year so this should be called periodically.
func (c *CaltrainClient) UpdateHolidays(ctx context.Context) error {
	logrus.Debug("Updating holidays...")
	c.sLock.Lock()
	defer c.sLock.Unlock()

	query := map[string]string{
		"operator_id": "CT",
		"api_key":     c.key,
	}
	data, err := c.APIClient.Get(ctx, holidaysURL, query)
	if err != nil {
		return fmt.Errorf("failed to make 'update holidays' request: %w", err)
	}

	holidays, err := parseHolidays(data)
	if err != nil {
		return fmt.Errorf("failed to parse holidays: %w", err)
	}
	c.holidays = holidays
	return nil
}

// SetupCache enables the use of API caching to prevent going over the API
// limit. Users set the caching expire time.
func (c *CaltrainClient) SetupCache(expire time.Duration) {
	c.cache = newCache(expire)
	c.useCache = true
}

// GetDelays makes an API call and returns a slice of TrainStatus who's
// delay into their next station is greater than the time.Duration argument
func (c *CaltrainClient) GetDelays(ctx context.Context, threshold time.Duration) ([]TrainStatus, time.Time, error) {
	logrus.Debug("Checking for delayed trains...")
	query := map[string]string{
		"agency":  "CT",
		"api_key": c.key,
	}
	t := time.Now()

	var cacheData []TrainStatus
	var cacheTime time.Time
	var cacheError error

	if c.useCache {
		var data []byte
		var ok bool
		data, cacheTime, ok = c.cache.get(delayURL)
		if ok {
			c.lLock.RLock()
			cacheData, cacheError = parseDelays(data, threshold, c.lines)
			c.lLock.RUnlock()
			return cacheData, cacheTime, cacheError
		}
	}

	data, err := c.APIClient.Get(ctx, delayURL, query)
	if err != nil {
		return nil, t, fmt.Errorf("failed to make 'get delays' request: %w", err)
	}

	// Now parse the body json string
	c.lLock.RLock()
	trains, err := parseDelays(data, threshold, c.lines)
	c.lLock.RUnlock()
	if err != nil {
		if c.useCache {
			var limErr *APILimitError
			var apiErr *APIError
			if errors.As(err, &limErr) || errors.As(err, &apiErr) {
				return cacheData, cacheTime, err
			}
			return cacheData, cacheTime, fmt.Errorf("failed to parse delay data: %w", err)
		}
		return nil, t, fmt.Errorf("failed to parse delay data: %w", err)
	}

	if c.useCache {
		c.cache.set(delayURL, data)
	}
	return trains, t, nil
}

// GetStationStatus makes an API call and returns a slice of TrainsStatus
// who have a status reported for the given station and direction.
func (c *CaltrainClient) GetStationStatus(ctx context.Context, stationName Station, direction Direction) ([]TrainStatus, time.Time, error) {
	logrus.Debugf("Getting station status for %s...", stationName.String())
	t := time.Now()
	code, err := c.getStationCode(stationName, direction)
	if err != nil {
		return nil, t, fmt.Errorf("failed to get station code: %w", err)
	}
	query := map[string]string{
		"agency":   "CT",
		"stopCode": code,
		"api_key":  c.key,
	}

	var cacheData []TrainStatus
	var cacheTime time.Time
	var cacheError error

	// cache key is stationStatusURL plus the stop code
	if c.useCache {
		var data []byte
		var ok bool
		data, cacheTime, ok = c.cache.get(stationStatusURL + code)
		if ok {
			c.lLock.RLock()
			cacheData, cacheError = getTrains(data, c.lines)
			c.lLock.RUnlock()
			return cacheData, cacheTime, cacheError
		}
	}

	data, err := c.APIClient.Get(ctx, stationStatusURL, query)
	if err != nil {
		return nil, t, fmt.Errorf("failed to make 'get station status' request: %w", err)
	}

	// Now parse the body json string
	c.lLock.RLock()
	trains, err := getTrains(data, c.lines)
	c.lLock.RUnlock()
	if err != nil {
		if c.useCache {
			var limErr *APILimitError
			var apiErr *APIError
			if errors.As(err, &limErr) || errors.As(err, &apiErr) {
				return cacheData, cacheTime, err
			}
			return cacheData, cacheTime, fmt.Errorf("failed to parse trains: %w", err)
		}
		return nil, t, fmt.Errorf("failed to parse trains: %w", err)
	}

	if c.useCache {
		c.cache.set(stationStatusURL+code, data)
	}
	return trains, t, nil
}

// GetTrainsBetweenStationsForWeekday returns a slice of Routes that travel
// from src to dst on the given weekday. It uses the cached timetable and
// does not make an API call
func (c *CaltrainClient) GetTrainsBetweenStationsForWeekday(ctx context.Context, src, dst Station, weekday time.Weekday) ([]*Route, error) {
	logrus.Debugf("Getting trains between stations '%s' and '%s' for a '%s'", src.String(), dst.String(), weekday.String())
	c.ttLock.RLock()
	defer c.ttLock.RUnlock()

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

// GetTrainsBetweenStationsForDate returns a slice of Routes that travel
// from src to dst for a given date. It uses the cached timetable and does
// not make an API call. It checks against the known holidays. Date must be
// in the correct time zone
func (c *CaltrainClient) GetTrainsBetweenStationsForDate(ctx context.Context, src, dst Station, date time.Time) ([]*Route, error) {
	if c.IsHoliday(date) {
		return c.GetTrainsBetweenStationsForWeekday(ctx, src, dst, time.Sunday)
	}
	return c.GetTrainsBetweenStationsForWeekday(ctx, src, dst, date.Weekday())
}

// IsHoliday returns true if the date passed in is a holiday
func (c *CaltrainClient) IsHoliday(date time.Time) bool {
	d := date.Truncate(24 * time.Hour)
	c.sLock.RLock()
	defer c.sLock.RUnlock()
	for _, h := range c.holidays {
		if d.Equal(h.Truncate(24 * time.Hour)) {
			return true
		}
	}
	return false
}

// GetRoutesForAllStops works the same as GetTrainsBetweenStationsForDate
// except many stations will be checked instead of just two
func (c *CaltrainClient) GetRoutesForAllStops(ctx context.Context, stops []Station, dir Direction, date time.Time) ([]*Route, error) {
	c.ttLock.RLock()
	defer c.ttLock.RUnlock()

	var d time.Weekday
	if c.IsHoliday(date) {
		d = time.Sunday
	} else {
		d = date.Weekday()
	}

	journeys, err := c.getTrainRoutesForAllStops(stops, dir, d)
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

// GetStationTimetable returns the routes that stop at a given station in the
// given direction
func (c *CaltrainClient) GetStationTimetable(st Station, dir Direction, date time.Time) ([]*Route, error) {
	c.ttLock.RLock()
	defer c.ttLock.RUnlock()

	code, err := c.getStationCode(st, dir)
	if err != nil {
		return nil, err
	}
	weekday := date.Weekday()
	if c.IsHoliday(date) {
		weekday = time.Sunday
	}
	journeys, err := c.getTimetableForStation(code, dir, weekday)
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

// GetTrainRoute returns the Route for a given train
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
func (c *CaltrainClient) getStationCode(st Station, dir Direction) (string, error) {
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

// getRouteCodes returns the proper station codes for a route given a
// source and destination station name
func (c *CaltrainClient) getRouteCodes(src, dst Station) (string, string, error) {
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

	dir, err := GetDirectionFromSrcToDst(src, dst)
	if err != nil {
		return "", "", err
	}

	// if the source is greater than destination, it's moving south
	if dir == South {
		return srcSt.southCode, dstSt.southCode, nil
	}
	return srcSt.northCode, dstSt.northCode, nil
}

// journeyToRoute converts a timetableRouteJourney into a Route
func (c *CaltrainClient) journeyToRoute(r timetableRouteJourney) (*Route, error) {
	c.lLock.RLock()
	line, err := parseLine(r.Line, c.lines)
	c.lLock.RUnlock()
	if err != nil {
		return nil, err
	}

	route := &Route{
		TrainNum:  r.ID,
		Direction: getDirFromChar(r.JourneyPatternView.DirectionRef.Ref),
		Line:      line,
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
func (c *CaltrainClient) getStationFromCode(code string) Station {
	c.sLock.RLock()
	defer c.sLock.RUnlock()
	for name, st := range c.stations {
		if st.northCode == code || st.southCode == code {
			return name
		}
	}
	return 0
}

// AllLines returns a slice of all available train lines
func (c *CaltrainClient) AllLines() []Line {
	c.lLock.RLock()
	defer c.lLock.RUnlock()
	return c.lines
}

// GetDirectionFromSrcToDst returns the direction the train would go to get
// from src to dst. Value is either North or South
func GetDirectionFromSrcToDst(src, dst Station) (Direction, error) {
	var dir Direction
	if src == dst {
		return dir, fmt.Errorf("The stations are the same: %s to %s", src, dst)
	}
	// Station is an int, with the southern station having a larger value than
	// the northern station
	if src > dst {
		return North, nil
	} else if dst > src {
		return South, nil
	} else {
		return dir, fmt.Errorf("could not determine direction from %s to %s", src, dst)
	}
}

// GetStations returns a slice of all recognized stations in order from North
// to South
func GetStations() []Station {
	return stationSlice
}

// getDirFromChar returns the proper direction string for a given character.
// HasPrefix is used in case the "char" has whitespace
func getDirFromChar(c string) Direction {
	if strings.HasPrefix(c, "N") {
		return North
	} else if strings.HasPrefix(c, "S") {
		return South
	} else {
		return 0
	}
}

// parseLine returns the Line corresponding with the string argument. It
// accepts names or line IDs
func parseLine(line string, lines []Line) (Line, error) {
	for _, l := range lines {
		if strings.EqualFold(l.Id, line) || strings.EqualFold(l.Name, line) {
			return l, nil
		}
	}

	return Line{}, fmt.Errorf("%s is not a valid line", line)
}
