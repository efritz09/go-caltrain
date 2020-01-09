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
