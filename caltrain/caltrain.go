// caltrain provides a user API for getting live caltrain updates
package caltrain

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Caltrain interface {
	GetDelays(context.Context) ([]Train, error)
	GetStationStatus(context.Context, string, string) ([]Train, error)
}

type CaltrainClient struct {
	stations  stations // station information struct
	timetable *timeTable

	key string // API key for 511.org

	DelayThreshold time.Duration // delay time to allow before warning user
}

func New(key string) *CaltrainClient {
	return &CaltrainClient{
		stations:       getStations(),
		timetable:      newTimeTable(),
		key:            key,
		DelayThreshold: defaultDelayThreshold,
	}
}

type Train struct {
	number    string        // train number
	nextStop  string        // stop for information
	direction string        // North or South
	delay     time.Duration // time behind schedule
	line      string        // bullet, limited, etc.
}

// GetDelays returns a list of delayed trains and their information
func (c *CaltrainClient) GetDelays(ctx context.Context) ([]Train, error) {
	query := map[string]string{
		"agency":  "CT",
		"api_key": c.key,
	}
	url := baseURL + "StopMonitoring"
	resp, err := c.get(ctx, url, query)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to ready body: %w", err)
	}

	// Now parse the body json string
	return parseDelays(body, c.DelayThreshold)

}

// GetStationStatus returns the status of upcoming trains for a given station
// and direction. Direction should be caltrain.North or caltrain.South
func (c *CaltrainClient) GetStationStatus(ctx context.Context, stationName string, direction string) ([]Train, error) {
	code, err := c.stations.getCode(stationName, direction)
	if err != nil {
		return nil, fmt.Errorf("failed to get station code: %w", err)
	}
	query := map[string]string{
		"agency":   "CT",
		"stopCode": strconv.Itoa(code),
		"api_key":  c.key,
	}
	url := baseURL + "StopMonitoring"
	resp, err := c.get(ctx, url, query)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to ready body: %w", err)
	}

	// Now parse the body json string
	return getTrains(body)
}

// // GetTimeTable returns the time table for the current day for all stations
// func (c *CaltrainClient) GetTimeTable() (*TimeTable, error) {
// 	// TODO: implement
// 	return nil, nil
// }

// GetTrainsBetweenStations returns a list of all trains that go from a to b.
// Trains with statuses available will include the status. This relies on the
// accuracy of the timetable.
func (c *CaltrainClient) GetTrainsBetweenStations(a, b string) ([]*Train, error) {
	// TODO: implement in the future
	return nil, nil
}

// UpdateTimeTable should be called once per day to update the day's timetable
func (c *CaltrainClient) UpdateTimeTable() error {
	// TODO: implement in the future
	return nil
}

func (c *CaltrainClient) get(ctx context.Context, url string, query map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// update the url with the required query parameters
	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
