package caltrain

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// timetable.go contains helpers relating to the timetable. All functions must
// have ttLock read locked before calling

// TrainNotFoundError is returned when a provided train number does not exist
// in the current timetable.
type TrainNotFoundError struct {
	Number string
}

func (t *TrainNotFoundError) Error() string {
	return fmt.Sprintf("No routes found for train: %s", t.Number)
}

// getTimetableForStation returns a list of trains that stop at a given station
// code and directions
func (c *CaltrainClient) getTimetableForStation(stationCode string, dir Direction, day time.Weekday) ([]timetableRouteJourney, error) {
	allJourneys := []timetableRouteJourney{}

	weekday := strings.ToLower(day.String())

	for lineId, ttArray := range c.timetable {
		line, err := c.getLine(lineId)
		if err != nil {
			logrus.Errorf("failed to get line!")
			line = Line{Id: "unknown", Name: "Unknown"}
		} else {
			logrus.Debugf("getting line %s-%s for station code %s...", line.Id, line.Name, stationCode)
		}
		for _, frame := range ttArray {
			// Check the day reference
			if !c.isForToday(weekday, frame.FrameValidityConditions.AvailabilityCondition.DayTypes.DayTypeRef.Ref) {
				continue
			}
			// Checkc the direction
			if !isMyDirection(frame.Name, dir) {
				continue
			}
			// loop through all journeys in this frame
			journeys := frame.VehicleJourneys.TimetableRouteJourney
			for _, journey := range journeys {
				if isStationInJourney(stationCode, journey) {
					journey.Line = line.Name
					allJourneys = append(allJourneys, journey)
				}
			}
		}
	}
	return allJourneys, nil
}

// getRouteForTrain returns a TimetableRouteJourney and the route's line for
// the given train number
func (c *CaltrainClient) getRouteForTrain(trainNum string) (timetableRouteJourney, error) {
	// TODO: the train number has metadata on the line type, and the day, it
	// could save time to use that to limit the search
	for line, ttArray := range c.timetable {
		for _, frame := range ttArray {
			journeys := frame.VehicleJourneys.TimetableRouteJourney
			for _, journey := range journeys {
				if journey.ID == trainNum {
					journey.Line = line
					return journey, nil
				}
			}
		}
	}
	return timetableRouteJourney{}, &TrainNotFoundError{Number: trainNum}
}

// getTrainRoutesBetweenStations returns a slice of routes from src to dst on a
// given weekday
func (c *CaltrainClient) getTrainRoutesBetweenStations(src, dst Station, day time.Weekday) ([]timetableRouteJourney, error) {
	sCode, dCode, err := c.getRouteCodes(src, dst)
	if err != nil {
		return nil, fmt.Errorf("failed to get station codes: %w", err)
	}

	weekday := strings.ToLower(day.String())

	routes := []timetableRouteJourney{}
	for line, ttArray := range c.timetable {
		for _, frame := range ttArray {
			// Check the day reference
			if !c.isForToday(weekday, frame.FrameValidityConditions.AvailabilityCondition.DayTypes.DayTypeRef.Ref) {
				continue
			}

			journeys := frame.VehicleJourneys.TimetableRouteJourney
			for _, journey := range journeys {
				if areStationsInJourney([]string{sCode, dCode}, journey) {
					journey.Line = line
					routes = append(routes, journey)
				}
			}
		}
	}
	return routes, nil
}

// getTrainRoutesForAllStops
// TODO: unit test this
func (c *CaltrainClient) getTrainRoutesForAllStops(stops []Station, dir Direction, day time.Weekday) ([]timetableRouteJourney, error) {
	codes := make([]string, len(stops))
	for i, st := range stops {
		c, err := c.getStationCode(st, dir)
		if err != nil {
			return nil, err
		}
		codes[i] = c
	}

	weekday := strings.ToLower(day.String())

	routes := []timetableRouteJourney{}
	for line, ttArray := range c.timetable {
		for _, frame := range ttArray {
			// Check the day reference
			if !c.isForToday(weekday, frame.FrameValidityConditions.AvailabilityCondition.DayTypes.DayTypeRef.Ref) {
				continue
			}

			journeys := frame.VehicleJourneys.TimetableRouteJourney
			for _, journey := range journeys {
				if areStationsInJourney(codes, journey) {
					journey.Line = line
					routes = append(routes, journey)
				}
			}
		}
	}
	return routes, nil
}

// isForToday returns true if the frame is for the day
func (c *CaltrainClient) isForToday(day string, ref string) bool {
	weekdays, ok := c.dayService[ref]
	if !ok {
		return false
	}
	for _, d := range weekdays {
		if d == day {
			return true
		}
	}
	return false
}

// isMyDirection returns true if the frame direction matches dir
func isMyDirection(frame string, dir Direction) bool {
	// convert `Bullet:N :Year Round Weekday (Weekday)` to `N`
	frameDir := strings.Split(strings.Split(frame, ":")[1], "")[0]
	return strings.HasPrefix(dir.String(), frameDir)
}

// TODO: unit test this
func isStationInJourney(st string, journey timetableRouteJourney) bool {
	for _, call := range journey.Calls.Call {
		if call.ScheduledStopPointRef.Ref == st {
			return true
		}
	}
	return false
}

// TODO: unit test this
func areStationsInJourney(stops []string, journey timetableRouteJourney) bool {
	for _, s := range stops {
		if !isStationInJourney(s, journey) {
			return false
		}
	}
	return true
}
