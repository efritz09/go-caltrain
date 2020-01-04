package caltrain

import (
	"fmt"
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
func (c *CaltrainClient) getTimetableForStation(stationCode, dir, weekday string) ([]TimetableRouteJourney, error) {
	weekdayRefs := weekdayReferences[weekday]

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
				if isStationInJourney(stationCode, journey) {
					allJourneys = append(allJourneys, journey)
				}
			}
		}
	}
	return allJourneys, nil
}

// getRouteForTrain returns a TimetableRouteJourney and the route's line for
// the given train number
func (c *CaltrainClient) getRouteForTrain(trainNum string) (TimetableRouteJourney, error) {
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
	return TimetableRouteJourney{}, fmt.Errorf("No routes found for train %s", trainNum)
}

// getTrainRoutesBetweenStations returns two maps of line: TimetableRouteJourney
// the first is the routes north, the second is routes south
func (c *CaltrainClient) getTrainRoutesBetweenStations(src, dst, weekday string) ([]TimetableRouteJourney, []TimetableRouteJourney, error) {
	srcN, err := c.getStationCode(src, North)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get station code: %w", err)
	}
	dstN, err := c.getStationCode(dst, North)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get station code: %w", err)
	}
	srcS, err := c.getStationCode(src, South)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get station code: %w", err)
	}
	dstS, err := c.getStationCode(dst, South)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get station code: %w", err)
	}

	weekdayRefs := weekdayReferences[weekday]
	journeyN := []TimetableRouteJourney{}
	journeyS := []TimetableRouteJourney{}
	for line, ttArray := range c.timetable {
		for _, frame := range ttArray {
			// Check the day reference
			if !isInDayRef(weekdayRefs, frame.FrameValidityConditions.AvailabilityCondition.DayTypes.DayTypeRef.Ref) {
				continue
			}
			// convert `Bullet:N :Year Round Weekday (Weekday)` to `North`
			dir := getDirFromChar(strings.Split(frame.Name, ":")[1])
			// if it's north, check that both srcN and dstN are there
			// same for south
			journeys := frame.VehicleJourneys.TimetableRouteJourney
			for _, journey := range journeys {
				if dir == North {
					if areStationsInJourney(srcN, dstN, journey) {
						journey.Line = line
						journeyN = append(journeyN, journey)
					}
				} else if dir == South {
					if areStationsInJourney(srcS, dstS, journey) {
						journey.Line = line
						journeyS = append(journeyS, journey)
					}
				}
			}
		}
	}
	return journeyN, journeyS, nil
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

// TODO: unit test this
func areStationsInJourney(src, dst string, journey TimetableRouteJourney) bool {
	srcT := false
	dstT := false
	for _, call := range journey.Calls.Call {
		if call.ScheduledStopPointRef.Ref == src {
			srcT = true
		} else if call.ScheduledStopPointRef.Ref == dst {
			dstT = true
		}
		if srcT && dstT {
			return true
		}
	}
	return srcT && dstT
}
