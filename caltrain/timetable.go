package caltrain

import (
	"fmt"
	"strconv"
	"strings"
)

// timetable.go contains helpers relating to the timetable. All functions must
// have ttLock read locked before calling

// weekdayReferences are hard coded DayType values. These are from ServiceCalendarFrame
// but that is currently inconsistent. Future work would be to make this adaptive
var weekdayReferences = map[string][]string{
	"monday":    []string{"8005"},
	"tuesday":   []string{"8005"},
	"wednesday": []string{"8005"},
	"thursday":  []string{"8005"},
	"friday":    []string{"8005"},
	"saturday":  []string{"8006", "8007"},
	"sunday":    []string{"8006"},
}

// getTimetableForStation returns a list of trains that stop at a given station
// code and directions
func (c *CaltrainClient) getTimetableForStation(stationCode int, dir, weekday string) ([]TimetableRouteJourney, error) {
	weekdayRefs := weekdayReferences[weekday]
	st := strconv.Itoa(stationCode)

	allJourneys := []TimetableRouteJourney{}

	for _, ttArray := range c.timetable {
		for _, frame := range ttArray {
			// Check the day reference
			if !isInDayRef(weekdayRefs, frame.FrameValidityConditions.AvailabilityCondition.DayTypes.DayTypeRef.Ref) {
				continue
			}
			// Checkc the direction
			if !isMyDirection(frame.Name, dir) {
				continue
			}
			// loop through all journeys in this frame
			journeys := frame.VehicleJourneys.TimetableRouteJourney
			for _, journey := range journeys {
				if isStationInJourney(st, journey) {
					allJourneys = append(allJourneys, journey)
				}
			}
		}
	}
	return allJourneys, nil
}

// getRouteForTrain returns a TimetableRouteJourney for the given train number
func (c *CaltrainClient) getRouteForTrain(trainNum string) (TimetableRouteJourney, error) {
	// TODO: the train number has metadata on the line type, and the day, it
	// could save time to use that to limit the search
	for _, ttArray := range c.timetable {
		for _, frame := range ttArray {
			journeys := frame.VehicleJourneys.TimetableRouteJourney
			for _, journey := range journeys {
				if journey.ID == trainNum {
					return journey, nil
				}
			}
		}
	}
	return TimetableRouteJourney{}, fmt.Errorf("No routes found for train %s", trainNum)
}

// isInDayRef returns true if the value is in the slice day
func isInDayRef(day []string, val string) bool {
	for _, d := range day {
		if d == val {
			return true
		}
	}
	return false
}

// isMyDirection returns true if the frame direction matches dir
func isMyDirection(frame, dir string) bool {
	// convert `Bullet:N :Year Round Weekday (Weekday)` to `N`
	frameDir := strings.Split(strings.Split(frame, ":")[1], "")[0]
	return strings.HasPrefix(dir, frameDir)
}

// TODO: unit test this
func isStationInJourney(st string, journey TimetableRouteJourney) bool {
	for _, call := range journey.Calls.Call {
		if call.ScheduledStopPointRef.Ref == st {
			return true
		}
	}
	return false
}
