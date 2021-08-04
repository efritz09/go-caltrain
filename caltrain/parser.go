package caltrain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	northVal = 1
)

// parseDelays returns a slice of TrainsStatus for all trains that are delayed
// more than the threshold argument
func parseDelays(raw []byte, threshold time.Duration, lines []Line) ([]TrainStatus, error) {
	trains, err := getTrains(raw, lines)
	if err != nil {
		return nil, err
	}
	delayedTrains := []TrainStatus{}
	for _, t := range trains {
		if t.Delay > threshold {
			delayedTrains = append(delayedTrains, t)
		}
	}
	return delayedTrains, nil
}

// getTrains unmarshals the json blob and returns a slice of trains
func getTrains(raw []byte, lines []Line) ([]TrainStatus, error) {
	data := trainStatusJson{}
	// trim some problematic characters: https://stackoverflow.com/questions/31398044/got-error-invalid-character-%C3%AF-looking-for-beginning-of-value-from-json-unmar
	raw = bytes.TrimPrefix(raw, []byte("\xef\xbb\xbf"))
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	ret := []TrainStatus{}
	trains := data.ServiceDelivery.StopMonitoringDelivery.MonitoredStopVisit
	for _, t := range trains {
		train := t.MonitoredVehicleJourney
		status := train.MonitoredCall
		delay, arrival := getDelay(status)
		if delay < 0 {
			delay = 0
		}
		var next Station
		var dir Direction
		var line Line
		var err error
		if status.StopPointName != "" {
			next, err = ParseStation(strings.Split(status.StopPointName, " Caltrain")[0])
			if err != nil {
				return ret, fmt.Errorf("could not get trains for %s: %w", status.StopPointName, err)
			}
		}
		if train.DirectionRef != "" {
			dir, err = ParseDirection(train.DirectionRef)
			if err != nil {
				return ret, fmt.Errorf("could not get trains: %w", err)
			}
		}
		if train.LineRef != "" {
			line, err = parseLine(train.LineRef, lines)
			if err != nil {
				return ret, fmt.Errorf("could not get trains: %w", err)
			}
		}
		newTrain := TrainStatus{
			TrainNum:  train.FramedVehicleJourneyRef.DatedVehicleJourneyRef,
			NextStop:  next,
			Direction: dir,
			Delay:     delay,
			Arrival:   arrival,
			Line:      line,
		}
		ret = append(ret, newTrain)
	}

	return ret, nil
}

// getDelay returns the time difference between the expected arrival time and
// the aimed arrival time.
func getDelay(status monitoredCall) (time.Duration, time.Time) {
	arrival := status.AimedArrivalTime
	expected := status.ExpectedArrivalTime
	// if arrival is null, the train hasn't left the starting station yet
	if arrival.IsZero() {
		return 0, expected
	}
	if expected.IsZero() {
		// ExpectedArrivalTime can be null if the train is at it's starting station
		expected = status.ExpectedDepartureTime
	}

	// We can use UTC because the API returns UTC time for live updates
	now := time.Now()
	// The API can mess up the aimed arrival time. If the arrival time is
	// earlier than the current time, use the ExpectedDepartureTime
	if arrival.Before(now) {
		arrival = status.AimedDepartureTime
	}

	return expected.Sub(arrival), expected
}

// parseTimetable returns a slice of TimetableFrames from the given raw data
func parseTimetable(raw []byte) ([]timetableFrame, map[string][]string, error) {
	raw = bytes.TrimPrefix(raw, []byte("\xef\xbb\xbf"))
	data := timetableJson{}
	services := make(map[string][]string)
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	frames := data.Content.TimetableFrame
	sframe := data.Content.ServiceCalendarFrame.DayTypes.DayType
	for _, f := range sframe {
		days := strings.Split(strings.TrimSpace(strings.ToLower(f.Properties.PropertyOfDay.DaysOfWeek)), " ")
		services[f.ID] = days
	}
	return frames, services, nil
}

// parseStations returns a map of station name to station struct, parsing the
// north and south codes
func parseStations(raw []byte) (map[Station]*stationInfo, error) {
	raw = bytes.TrimPrefix(raw, []byte("\xef\xbb\xbf"))
	data := stationJson{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	ret := make(map[Station]*stationInfo)

	// stops are indexed by id, not by station, so we have to generate a map
	// that gets us halfway there first, then convert to our struct
	stops := data.Contents.DataObjects.ScheduledStopPoint
	for _, stop := range stops {
		if stop.ID == "777403" || stop.ID == "777402" {
			// Skip the Tamien and San Jose stations that are an outlier for some reason
			continue
		}

		name, err := ParseStation(strings.TrimSuffix(stop.Name, " Caltrain Station"))
		if err != nil {
			return ret, fmt.Errorf("failed to parse station: %w", err)
		}
		if st, ok := ret[name]; !ok {
			// create a new station with location
			lat, err := strconv.ParseFloat(stop.Location.Latitude, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse location for %s: %w", name, err)
			}
			lon, err := strconv.ParseFloat(stop.Location.Longitude, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse location for %s: %w", name, err)
			}
			newStation := &stationInfo{
				name:      name,
				latitude:  lat,
				longitude: lon,
			}
			if err := addDirectionToStation(newStation, stop.ID); err != nil {
				return nil, fmt.Errorf("failed to parse stations: %w", err)
			}
			ret[name] = newStation
		} else {
			// the location difference between the north and south side is
			// negligible and we can ignore it
			if err := addDirectionToStation(st, stop.ID); err != nil {
				return nil, fmt.Errorf("failed to parse stations: %w", err)
			}
		}
	}

	return ret, nil
}

// parseLines returns a slice of lines that are available
func parseLines(raw []byte) ([]Line, error) {
	raw = bytes.TrimPrefix(raw, []byte("\xef\xbb\xbf"))
	data := lineJson{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	ret := make([]Line, len(data))
	for i, line := range data {
		ret[i] = Line{
			Id:   line.ID,
			Name: line.Name,
		}
	}

	return ret, nil
}

// parseHolidays returns a slice of dates that are holidays
func parseHolidays(raw []byte) ([]time.Time, error) {
	raw = bytes.TrimPrefix(raw, []byte("\xef\xbb\xbf"))
	data := holidayJson{}
	if err := json.Unmarshal(raw, &data); err != nil {
		// TEMPORARY: try to use the alt struct
		data := holidayJsonAlt{}
		if err := json.Unmarshal(raw, &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %w", err)

		}
		holiday := data.Content.AvailabilityConditions
		id := strings.TrimPrefix(holiday.ID, "CT:")
		date, err := time.Parse("2006-01-02", id)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time value %s: %w", id, err)
		}
		return []time.Time{date}, nil
		// return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	holidays := data.Content.AvailabilityConditions
	ret := make([]time.Time, len(holidays))

	for i, holiday := range holidays {
		id := strings.TrimPrefix(holiday.ID, "CT:")
		date, err := time.Parse("2006-01-02", id)
		if err != nil {
			return ret, fmt.Errorf("failed to parse time value %s: %w", id, err)
		}
		ret[i] = date
	}

	return ret, nil
}

// addDirectionToStation is a helper function to add the code to the proper
// direction in the station struct
func addDirectionToStation(s *stationInfo, id string) error {
	if dir, err := isCodeNorth(id); err != nil {
		return fmt.Errorf("failed to parse stations: %w", err)
	} else if dir {
		s.northCode = id
	} else {
		s.southCode = id
	}
	return nil
}

// isCodeNorth returns true if the code is for a north station
func isCodeNorth(code string) (bool, error) {
	lastChar := code[len(code)-1:]
	i, err := strconv.Atoi(lastChar)
	if err != nil {
		return false, fmt.Errorf("bad station code: %s", code)
	}
	return i <= northVal, nil
}
