package caltrain

import (
	"bytes"
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
func parseDelays(raw []byte, threshold time.Duration) ([]Train, error) {
	trains, err := getTrains(raw)
	if err != nil {
		return nil, err
	}
	delayedTrains := []Train{}
	for _, t := range trains {
		if t.delay > threshold {
			delayedTrains = append(delayedTrains, t)
		}
	}
	return delayedTrains, nil
}

// getTrains unmarshals the json blob and returns a slice of trains
func getTrains(raw []byte) ([]Train, error) {
	data := trainStatusJson{}
	// trim some problematic characters: https://stackoverflow.com/questions/31398044/got-error-invalid-character-%C3%AF-looking-for-beginning-of-value-from-json-unmar
	raw = bytes.TrimPrefix(raw, []byte("\xef\xbb\xbf"))
	err := json.Unmarshal(raw, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	ret := []Train{}
	trains := data.ServiceDelivery.StopMonitoringDelivery.MonitoredStopVisit
	for _, t := range trains {
		train := t.MonitoredVehicleJourney
		status := train.MonitoredCall
		delay, err := getDelay(status)
		if delay < 0 {
			delay = 0
		}
		if err != nil {
			fmt.Printf("Error getting the delay: %v\n", err)
		}
		newTrain := Train{
			number:    train.FramedVehicleJourneyRef.DatedVehicleJourneyRef,
			nextStop:  strings.Split(status.StopPointName, " Caltrain")[0],
			direction: train.DirectionRef,
			delay:     delay,
			line:      train.LineRef,
		}
		ret = append(ret, newTrain)
	}

	return ret, nil
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
