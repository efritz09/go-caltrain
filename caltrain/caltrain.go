// caltrain provides a user API for getting live caltrain updates
package caltrain

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
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

type Train struct {
	Number    string        // train number
	NextStop  string        // stop for information
	Direction string        // North or South
	Delay     time.Duration // time behind schedule
	Line      string        // bullet, limited, etc.
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

// GetTrainsBetweenStations returns a list of all trains that go from a to b.
// Trains with statuses available will include the status. This relies on the
// accuracy of the timetable.
func (c *CaltrainClient) GetTrainsBetweenStations(ctx context.Context, a, b string) ([]*Train, error) {
	c.ttLock.Lock()
	defer c.ttLock.Unlock()
	// TODO: implement in the future
	return nil, nil
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
