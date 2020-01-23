package caltrain

import (
	"fmt"
	"strings"
	"time"
)

// TrainStatus provides the live status of the train
type TrainStatus struct {
	TrainNum  string        // Train reference number
	Direction Direction     // Direction the train is travelling: North or South
	Line      Line          // bullet, limited, etc.
	Delay     time.Duration // Amount of time behind schedule
	Arrival   time.Time     // Expected arrival time at NextStop
	NextStop  Station       // Name of the station that the train will stop at next
}

// Route contains metadata for a given train and the stops that it will make on
// it's route
type Route struct {
	TrainNum  string      // Train reference number
	Direction Direction   // Direction the train is travelling: North or South
	Line      Line        // bullet, limited, etc.
	NumStops  int         // Total number of stops on this route
	Stops     []TrainStop // Slice of stops on this route
}

// TrainStop is a single stop on a route
type TrainStop struct {
	Order     int       // stop number on the route
	Station   Station   // station name of this stop
	Arrival   time.Time // scheduled arrival time at this station
	Departure time.Time // scheduled departure time from this station
}

type stationInfo struct {
	name      Station
	northCode string
	southCode string
	latitude  float64
	longitude float64
}

// A Station specifies a recognized Caltrain station
type Station int

// This ordering is also the station order from north to south along the line
const (
	StationSanFrancisco Station = iota // "San Francisco"
	Station22ndStreet                  // "22nd Street"
	StationBayshore                    // "Bayshore"
	StationSouthSF                     // "South San Francisco"
	StationSanBruno                    // "San Bruno"
	StationMillbrae                    // "Millbrae"
	StationBroadway                    // "Broadway"
	StationBurlingame                  // "Burlingame"
	StationSanMateo                    // "San Mateo"
	StationHaywardPark                 // "Hayward Park"
	StationHillsdale                   // "Hillsdale"
	StationBelmont                     // "Belmont"
	StationSanCarlos                   // "San Carlos"
	StationRedwoodCity                 // "Redwood City"
	StationAtherton                    // "Atherton"
	StationMenloPark                   // "Menlo Park"
	StationPaloAlto                    // "Palo Alto"
	StationStanford                    // "Stanford"
	StationCalAve                      // "California Ave"
	StationSanAntonio                  // "San Antonio"
	StationMountainView                // "Mountain View"
	StationSunnyvale                   // "Sunnyvale"
	StationLawrence                    // "Lawrence"
	StationSantaClara                  // "Santa Clara"
	StationCollegePark                 // "College Park"
	StationSanJose                     // "San Jose Diridon"
	StationTamien                      // "Tamien"
	StationCapitol                     // "Capitol"
	StationBlossomHill                 // "Blossom Hill"
	StationMorganHill                  // "Morgan Hill"
	StationSanMartin                   // "San Martin"
	StationGilroy                      // "Gilroy"
)

var stations = map[Station]string{
	StationSanFrancisco: "San Francisco",
	Station22ndStreet:   "22nd Street",
	StationBayshore:     "Bayshore",
	StationSouthSF:      "South San Francisco",
	StationSanBruno:     "San Bruno",
	StationMillbrae:     "Millbrae",
	StationBroadway:     "Broadway",
	StationBurlingame:   "Burlingame",
	StationSanMateo:     "San Mateo",
	StationHaywardPark:  "Hayward Park",
	StationHillsdale:    "Hillsdale",
	StationBelmont:      "Belmont",
	StationSanCarlos:    "San Carlos",
	StationRedwoodCity:  "Redwood City",
	StationAtherton:     "Atherton",
	StationMenloPark:    "Menlo Park",
	StationPaloAlto:     "Palo Alto",
	StationStanford:     "Stanford",
	StationCalAve:       "California Ave",
	StationSanAntonio:   "San Antonio",
	StationMountainView: "Mountain View",
	StationSunnyvale:    "Sunnyvale",
	StationLawrence:     "Lawrence",
	StationSantaClara:   "Santa Clara",
	StationCollegePark:  "College Park",
	StationSanJose:      "San Jose Diridon",
	StationTamien:       "Tamien",
	StationCapitol:      "Capitol",
	StationBlossomHill:  "Blossom Hill",
	StationMorganHill:   "Morgan Hill",
	StationSanMartin:    "San Martin",
	StationGilroy:       "Gilroy",
}

// String returns the string name of the station
func (s Station) String() string {
	return stations[s]
}

// ParseStation returns a Station from the string passed in. If the string is
// not a recognized station, it will return an error
func ParseStation(s string) (Station, error) {
	l := strings.ToLower(s)
	for k, v := range stations {
		if strings.ToLower(v) == l {
			return k, nil
		}
	}
	return 0, fmt.Errorf("%s is not a recognized station", s)
}

// A Direction specifies a Caltrain route direction (North or South)
type Direction int

const (
	North Direction = iota
	South
)

var directions = [...]string{
	"North",
	"South",
}

// String returns the string name of the direction
func (d Direction) String() string {
	if North <= d && d <= South {
		return directions[d]
	}
	return fmt.Sprintf("unknown direction %d", d)
}

// ParseDirection returns a Direction from the string passed in. If the string
// is not a valid direction it returns an error
func ParseDirection(d string) (Direction, error) {
	l := strings.ToLower(d)
	if l == "north" {
		return North, nil
	} else if l == "south" {
		return South, nil
	} else {
		return 0, fmt.Errorf("%s is not a valid direction. Must be either North or South", d)
	}
}

// A Line specifies a Caltrain route line type (Bullet, Limited, Local)
type Line int

const (
	Bullet Line = iota
	Limited
	Local
)

var lines = [...]string{
	"Bullet",
	"Limited",
	"Local",
}

// String returns the string name of the line
func (l Line) String() string {
	if Bullet <= l && l <= Local {
		return lines[l]
	}
	return fmt.Sprintf("unknown line %d", l)
}

// ParseLine returns a Line from the string passed in. If the string is not a
// valid line it returns an error
func ParseLine(l string) (Line, error) {
	s := strings.ToLower(l)
	if s == "limited" {
		return Limited, nil
	} else if s == "local" {
		return Local, nil
	} else if s == "bullet" {
		return Bullet, nil
	} else {
		return 0, fmt.Errorf("%s is not a valid line. Must be Local, Limited or Bullet", l)
	}
}
