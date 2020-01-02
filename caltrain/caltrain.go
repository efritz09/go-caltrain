// caltrain provides a user API for getting live caltrain updates
package caltrain

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/efritz09/go-caltrain/caltrain/internal/utilities"
)

type Caltrain interface {
	GetDelays(context.Context) ([]Train, error)
	GetStationStatus(context.Context, string, string) ([]Train, error)
	GetStations() []string
}

type CaltrainClient struct {
	timetable map[string][]TimetableFrame // map of line type to slice of service journeys
	ttLock    sync.RWMutex                // lock in case someone tries to access it during and update
	stations  map[string]station          // station information map

	key string // API key for 511.org

	DelayThreshold time.Duration // delay time to allow before warning user
	APIClient      APIClient     // API client for making caltrain queries. Default APIClient511
}

func New(key string) *CaltrainClient {
	return &CaltrainClient{
		timetable:      make(map[string][]TimetableFrame),
		key:            key,
		stations:       getStations(),
		DelayThreshold: defaultDelayThreshold,
		APIClient:      NewClient(),
	}
}

// Information on the current train status
type Train struct {
	Number    string        // train number
	NextStop  string        // stop for information
	Direction string        // North or South
	Delay     time.Duration // time behind schedule
	Line      string        // bullet, limited, etc.
}

// Stops for a given train
type Route struct {
	TrainNum  string // train number
	Direction string
	NumStops  int
	Stops     []TrainStop
	// TODO: define
}

type TrainStop struct {
	Order     int
	Station   string
	Arrival   time.Time
	Departure time.Time
}

// GetStations returns a list of station names
func (c *CaltrainClient) GetStations() []string {
	ret := []string{}
	for k := range c.stations {
		ret = append(ret, k)
	}
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
	data, err := c.APIClient.Get(ctx, delayURL, query)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	// Now parse the body json string
	return parseDelays(data, c.DelayThreshold)

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
		"stopCode": strconv.Itoa(code),
		"api_key":  c.key,
	}
	data, err := c.APIClient.Get(ctx, stationURL, query)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	// Now parse the body json string
	return getTrains(data)
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

		journeys, err := parseTimetable(data)
		if err != nil {
			return fmt.Errorf("failed to parse timetable: %w", err)
		}
		c.timetable[line] = journeys
	}

	return nil
}

// getStationCode returns the code for a given station and direction
func (c *CaltrainClient) getStationCode(st, dir string) (int, error) {
	// first validate the direction
	if dir != North && dir != South {
		return 0, fmt.Errorf("unknown direction %s", dir)
	}

	if station, ok := c.stations[st]; !ok {
		return 0, fmt.Errorf("unknown station %s", st)
	} else {
		return station.directions[dir], nil
	}
}

func (c *CaltrainClient) GetStationTimetable(st, dir string) ([]TimetableRouteJourney, error) {
	c.ttLock.RLock()
	defer c.ttLock.RUnlock()

	code, err := c.getStationCode(st, dir)
	if err != nil {
		return nil, err
	}

	weekday, err := utilities.GetWeekday(timezone)
	if err != nil {
		return nil, err
	}

	return c.getTimetableForStation(code, dir, weekday)

}

// GetTrainRoute returns the Route struct for a given train
func (c *CaltrainClient) GetTrainRoute(trainNum string) (Route, error) {
	c.ttLock.RLock()
	defer c.ttLock.RUnlock()
	route := Route{TrainNum: trainNum}
	r, err := c.getRouteForTrain(trainNum)
	if err != nil {
		return route, fmt.Errorf("failed to get Train Route: %w", err)
	}

	route.Direction = getDirFromChar(r.JourneyPatternView.DirectionRef.Ref)
	route.NumStops = len(r.Calls.Call)
	route.Stops = []TrainStop{}

	for _, s := range r.Calls.Call {
		order, err := strconv.Atoi(s.Order)
		if err != nil {
			return route, fmt.Errorf("could not convert order %s to int: %w", s.Order, err)
		}
		code, err := strconv.Atoi(s.ScheduledStopPointRef.Ref)
		if err != nil {
			return route, fmt.Errorf("could not convert station %s to int: %w", s.ScheduledStopPointRef.Ref, err)
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
			Station:   c.getStationFromCode(code),
			Arrival:   arr,
			Departure: dep,
		}
		route.Stops = append(route.Stops, t)
	}

	return route, nil
}

// GetTrainsBetweenStations returns a list of all trains that go from a to b.
// Trains with statuses available will include the status. This relies on the
// accuracy of the timetable.
func (c *CaltrainClient) GetTrainsBetweenStations(ctx context.Context, a, b string) ([]*Train, error) {
	c.ttLock.RLock()
	defer c.ttLock.RUnlock()
	// TODO: implement in the future
	return nil, nil
}
