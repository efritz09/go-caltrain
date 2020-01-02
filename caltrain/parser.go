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
		if t.Delay > threshold {
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
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	ret := []Train{}
	trains := data.ServiceDelivery.StopMonitoringDelivery.MonitoredStopVisit
	for _, t := range trains {
		train := t.MonitoredVehicleJourney
		status := train.MonitoredCall
		delay := getDelay(status)
		if delay < 0 {
			delay = 0
		}
		newTrain := Train{
			Number:    train.FramedVehicleJourneyRef.DatedVehicleJourneyRef,
			NextStop:  strings.Split(status.StopPointName, " Caltrain")[0],
			Direction: train.DirectionRef,
			Delay:     delay,
			Line:      train.LineRef,
		}
		ret = append(ret, newTrain)
	}

	return ret, nil
}

// getDelay returns the time difference between the expected arrival time and
// the aimed arrival time.
func getDelay(status monitoredCall) time.Duration {
	arrival := status.AimedArrivalTime
	expected := status.ExpectedArrivalTime
	if expected.IsZero() {
		// ExpectedArrivalTime can be null if the train is at it's starting station
		expected = status.ExpectedDepartureTime
	}

	now := time.Now()
	// The API can mess up the aimed arrival time. If the arrival time is
	// earlier than the current time, use the ExpectedDepartureTime
	if arrival.Before(now) {
		arrival = status.AimedDepartureTime
	}

	return expected.Sub(arrival)
}

func parseTimetable(raw []byte) ([]TimetableFrame, error) {
	raw = bytes.TrimPrefix(raw, []byte("\xef\xbb\xbf"))
	data := timetableJson{}
	if err := json.Unmarshal(raw, &data); err != nil {
		e := fmt.Errorf("failed to unmarshal: %w", err)
		// Try using the alternative struct
		altData := timetableJsonAlternate{}
		if err := json.Unmarshal(raw, &altData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal alternative: %w", e)
		}
		frames := altData.Content.TimetableFrame
		return frames, nil
	} else {
		frames := data.Content.TimetableFrame
		return frames, nil
	}
}
