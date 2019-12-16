// caltrain provides a user API for getting live caltrain updates
package caltrain

import "time"

type Caltrain struct {
	stations  map[string]station // map of station name to station information
	timetable *timeTable
}

func New() *Caltrain {
	return &Caltrain{
		stations:  getStations(),
		timetable: newTimeTable(),
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
