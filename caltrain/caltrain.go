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

type Caltrain struct {
	stations  stations // station information struct
	timetable *timeTable

	key string // API key for 511.org
}

func New(key string) *Caltrain {
	return &Caltrain{
		stations:  getStations(),
		timetable: newTimeTable(),
		key:       key,
	}
}

type TrainStatus struct {
	number    int
	nextStop  station
	direction string
	delay     time.Duration
}

// GetDelays returns a list of delayed trains and their information
func (c *Caltrain) GetDelays(ctx context.Context) ([]*TrainStatus, error) {
	query := map[string]string{
		"agency":  "CT",
		"api_key": c.key,
	}
	url := baseURL + "StopMonitoring"
	resp, err := c.get(ctx, url, query)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	fmt.Printf("resp:\n%+v\n", resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to ready body: %w", err)
	}
	fmt.Println(string(body))
	return nil, nil
}

// GetStationStatus returns the status of upcoming trains for a given station
// and direction. Direction should be caltrain.North or caltrain.South
func (c *Caltrain) GetStationStatus(ctx context.Context, stationName string, direction string) ([]*TrainStatus, error) {
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
	fmt.Printf("resp:\n%+v\n", resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to ready body: %w", err)
	}
	fmt.Println(string(body))

	return nil, nil
}

// // GetTimeTable returns the time table for the current day for all stations
// func (c *Caltrain) GetTimeTable() (*TimeTable, error) {
// 	// TODO: implement
// 	return nil, nil
// }

// GetTrainsBetweenStations returns a list of all trains that go from a to b.
// Trains with statuses available will include the status. This relies on the
// accuracy of the timetable.
func (c *Caltrain) GetTrainsBetweenStations(a, b string) ([]*TrainStatus, error) {
	// TODO: implement in the future
	return nil, nil
}

// UpdateTimeTable should be called once per day to update the day's timetable
func (c *Caltrain) UpdateTimeTable() error {
	// TODO: implement in the future
	return nil
}

func (c *Caltrain) get(ctx context.Context, url string, query map[string]string) (*http.Response, error) {
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
