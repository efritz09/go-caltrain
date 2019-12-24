// caltrain provides a user API for getting live caltrain updates
package caltrain

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Caltrain struct {
	stations  map[string]station // map of station name to station information
	timetable *timeTable
}

func New(key string) *Caltrain {
	return &Caltrain{
		stations:  getStations(),
		timetable: newTimeTable(),

		key: key,
	}
}

type TrainStatus struct {
	number    int
	nextStop  station
	direction string
	delay     time.Duration
}

// GetDelays returns a list of delayed trains and their information
func (c *Caltrain) GetDelays() ([]*TrainStatus, error) {
	// TODO: implement
	return nil, nil
}

// GetStationStatus returns the status of upcoming trains for a given station
func (c *Caltrain) GetStationStatus(stationName string) ([]*TrainStatus, error) {
	// TODO: implement
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

func (c *Caltrain) get(ctx context.Context, url string, query map[string]string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
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
		return err
	}

	fmt.Printf("resp:\n%+v\n", resp)
	return nil
}
