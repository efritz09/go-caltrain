package caltrain

const (
	delayURL         = "http://api.511.org/transit/StopMonitoring"
	stationsURL      = "http://api.511.org/transit/stops"
	stationStatusURL = "http://api.511.org/transit/StopMonitoring"
	timetableURL     = "http://api.511.org/transit/timetable"
	holidaysURL      = "http://api.511.org/transit/holidays"

	// North Direction = "North"
	// South Direction = "South"

	// lines
	// Bullet  Line = "Bullet"
	// Limited Line = "Limited"
	// Local   Line = "Local"

	// public station constants
	Station22ndStreet   Station = "22nd Street"
	StationAtherton     Station = "Atherton"
	StationBayshore     Station = "Bayshore"
	StationBelmont      Station = "Belmont"
	StationBlossomHill  Station = "Blossom Hill"
	StationBroadway     Station = "Broadway"
	StationBurlingame   Station = "Burlingame"
	StationCalAve       Station = "California Ave"
	StationCapitol      Station = "Capitol"
	StationCollegePark  Station = "College Park"
	StationGilroy       Station = "Gilroy"
	StationHaywardPark  Station = "Hayward Park"
	StationHillsdale    Station = "Hillsdale"
	StationLawrence     Station = "Lawrence"
	StationMenloPark    Station = "Menlo Park"
	StationMillbrae     Station = "Millbrae"
	StationMorganHill   Station = "Morgan Hill"
	StationMountainView Station = "Mountain View"
	StationPaloAlto     Station = "Palo Alto"
	StationRedwoodCity  Station = "Redwood City"
	StationSanAntonio   Station = "San Antonio"
	StationSanBruno     Station = "San Bruno"
	StationSanCarlos    Station = "San Carlos"
	StationSanFrancisco Station = "San Francisco"
	StationSanJose      Station = "San Jose Diridon"
	StationSanMartin    Station = "San Martin"
	StationSanMateo     Station = "San Mateo"
	StationSantaClara   Station = "Santa Clara"
	StationSouthSF      Station = "South San Francisco"
	StationStanford     Station = "Stanford"
	StationSunnyvale    Station = "Sunnyvale"
	StationTamien       Station = "Tamien"
)

// stationOrder is a map of the station name to it's order from North to South
// TODO: what if they add or remove a station?
var stationOrder = map[Station]int{
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
