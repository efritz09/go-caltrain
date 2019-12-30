package caltrain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	timeLayout = "2006-01-02T15:04:05Z"
)

// parseDelays returns a slice of TrainsStatus for all trains that are delayed
// more than the threshold argument
func parseDelays(raw []byte, threshold time.Duration) ([]TrainStatus, error) {
	data, err := getTrains(raw)
	if err != nil {
		return nil, err
	}

	ret := []TrainStatus{}
	trains := data.ServiceDelivery.StopMonitoringDelivery.MonitoredStopVisit
	for _, t := range trains {
		train := t.MonitoredVehicleJourney
		status := train.MonitoredCall
		delay, err := getDelay(status)
		if err != nil {
			fmt.Printf("Error getting the delay: %v\n", err)
			continue
		}

		if delay > threshold {
			newTrain := TrainStatus{
				number:    train.FramedVehicleJourneyRef.DatedVehicleJourneyRef,
				nextStop:  strings.Split(status.StopPointName, " Caltrain")[0],
				direction: train.DirectionRef,
				delay:     delay,
				line:      train.LineRef,
			}
			ret = append(ret, newTrain)
		}
	}
	return ret, nil
}

// getTrains unmarshals the json blob
func getTrains(raw []byte) (trainStatusJson, error) {
	data := trainStatusJson{}
	err := json.Unmarshal(raw, &data)
	return data, err
}

// getDelay returns the time difference between the expected arrival time and
// the aimed arrival time.
func getDelay(status monitoredCall) (time.Duration, error) {
	arrival, err := time.Parse(timeLayout, status.AimedArrivalTime)
	if err != nil {
		return 0, err
	}
	expected, err := time.Parse(timeLayout, status.ExpectedArrivalTime)
	if err != nil {
		return 0, err
	}
	now := time.Now()

	// The API can mess up the aimed arrival time. If the arrival time is
	// earlier than the current time, use the ExpectedDepartureTime
	if arrival.Before(now) {
		arrival, err = time.Parse(timeLayout, status.AimedDepartureTime)
		if err != nil {
			return 0, err
		}
	}

	return expected.Sub(arrival), nil
}
