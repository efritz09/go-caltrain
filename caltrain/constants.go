package caltrain

import (
	"time"
)

const (
	baseURL      = "http://api.511.org/transit/"
	delayURL     = "http://api.511.org/transit/StopMonitoring"
	stationURL   = "http://api.511.org/transit/StopMonitoring"
	timetableURL = "http://api.511.org/transit/timetable"

	defaultDelayThreshold = 10 * time.Minute

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
	StationSunnyvale    = "Sunnyvale"
	StationTamien       = "Tamien"
)

type station struct {
	name       string
	directions map[string]int
}

// newStation creates a new station struct with the name and direction values
func newStation(name string, n, s int) station {
	return station{
		name:       name,
		directions: map[string]int{North: n, South: s},
	}
}

// getStations returns a Stations struct with a map of station name to station information
func getStations() map[string]station {
	return map[string]station{
		Station22ndStreet:   newStation(Station22ndStreet, 70021, 70022),
		StationAtherton:     newStation(StationAtherton, 70151, 70152),
		StationBayshore:     newStation(StationBayshore, 70031, 70032),
		StationBelmont:      newStation(StationBelmont, 70121, 70122),
		StationBlossomHill:  newStation(StationBlossomHill, 70291, 70292),
		StationBroadway:     newStation(StationBroadway, 70071, 70072),
		StationBurlingame:   newStation(StationBurlingame, 70081, 70082),
		StationCalAve:       newStation(StationCalAve, 70191, 70192),
		StationCapitol:      newStation(StationCapitol, 70281, 70282),
		StationCollegePark:  newStation(StationCollegePark, 70251, 70252),
		StationGilroy:       newStation(StationGilroy, 70321, 70322),
		StationHaywardPark:  newStation(StationHaywardPark, 70101, 70102),
		StationHillsdale:    newStation(StationHillsdale, 70111, 70112),
		StationLawrence:     newStation(StationLawrence, 70231, 70232),
		StationMenloPark:    newStation(StationMenloPark, 70161, 70162),
		StationMillbrae:     newStation(StationMillbrae, 70061, 70062),
		StationMorganHill:   newStation(StationMorganHill, 70301, 70302),
		StationMountainView: newStation(StationMountainView, 70211, 70212),
		StationPaloAlto:     newStation(StationPaloAlto, 70171, 70172),
		StationRedwoodCity:  newStation(StationRedwoodCity, 70141, 70142),
		StationSanAntonio:   newStation(StationSanAntonio, 70201, 70202),
		StationSanBruno:     newStation(StationSanBruno, 70051, 70052),
		StationSanCarlos:    newStation(StationSanCarlos, 70131, 70132),
		StationSanFrancisco: newStation(StationSanFrancisco, 70011, 70012),
		StationSanJose:      newStation(StationSanJose, 70261, 70262),
		StationSanMartin:    newStation(StationSanMartin, 70311, 70312),
		StationSanMateo:     newStation(StationSanMateo, 70091, 70092),
		StationSantaClara:   newStation(StationSantaClara, 70241, 70242),
		StationSouthSF:      newStation(StationSouthSF, 70041, 70042),
		StationSunnyvale:    newStation(StationSunnyvale, 70221, 70222),
		StationTamien:       newStation(StationTamien, 70271, 70272),
	}
}
