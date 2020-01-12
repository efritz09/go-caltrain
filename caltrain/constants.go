package caltrain

import (
	"time"
)

const (
	baseURL          = "http://api.511.org/transit/"
	delayURL         = "http://api.511.org/transit/StopMonitoring"
	stationsURL      = "http://api.511.org/transit/stops"
	stationStatusURL = "http://api.511.org/transit/StopMonitoring"
	timetableURL     = "http://api.511.org/transit/timetable"

	defaultDelayThreshold = 10 * time.Minute

	timezone = "America/Los_Angeles"

	North = "North"
	South = "South"

	// lines
	Bullet  = "Bullet"
	Limited = "Limited"
	Local   = "Local"

	// public station constants
	Station22ndStreet   = "22nd Street"
	StationAtherton     = "Atherton"
	StationBayshore     = "Bayshore"
	StationBelmont      = "Belmont"
	StationBlossomHill  = "Blossom Hill"
	StationBroadway     = "Broadway"
	StationBurlingame   = "Burlingame"
	StationCalAve       = "California Ave"
	StationCapitol      = "Capitol"
	StationCollegePark  = "College Park"
	StationGilroy       = "Gilroy"
	StationHaywardPark  = "Hayward Park"
	StationHillsdale    = "Hillsdale"
	StationLawrence     = "Lawrence"
	StationMenloPark    = "Menlo Park"
	StationMillbrae     = "Millbrae"
	StationMorganHill   = "Morgan Hill"
	StationMountainView = "Mountain View"
	StationPaloAlto     = "Palo Alto"
	StationRedwoodCity  = "Redwood City"
	StationSanAntonio   = "San Antonio"
	StationSanBruno     = "San Bruno"
	StationSanCarlos    = "San Carlos"
	StationSanFrancisco = "San Francisco"
	StationSanJose      = "San Jose Diridon"
	StationSanMartin    = "San Martin"
	StationSanMateo     = "San Mateo"
	StationSantaClara   = "Santa Clara"
	StationSouthSF      = "South San Francisco"
	StationStanford     = "Stanford"
	StationSunnyvale    = "Sunnyvale"
	StationTamien       = "Tamien"
)

// stationOrder is a map of the station name to it's order from North to South
var stationOrder = map[string]int{
	StationSanFrancisco: 0,
	Station22ndStreet:   1,
	StationBayshore:     2,
	StationSouthSF:      3,
	StationSanBruno:     4,
	StationMillbrae:     5,
	StationBroadway:     6,
	StationBurlingame:   7,
	StationSanMateo:     8,
	StationHaywardPark:  9,
	StationHillsdale:    10,
	StationBelmont:      11,
	StationSanCarlos:    12,
	StationRedwoodCity:  13,
	StationAtherton:     14,
	StationMenloPark:    15,
	StationPaloAlto:     16,
	StationStanford:     17,
	StationCalAve:       18,
	StationSanAntonio:   19,
	StationMountainView: 20,
	StationSunnyvale:    21,
	StationLawrence:     22,
	StationSantaClara:   23,
	StationCollegePark:  24,
	StationSanJose:      25,
	StationTamien:       26,
	StationCapitol:      27,
	StationBlossomHill:  28,
	StationMorganHill:   29,
	StationSanMartin:    30,
	StationGilroy:       31,
}
