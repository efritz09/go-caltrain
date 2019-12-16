package caltrain

const (
	baseURL = "http://api.511.org/transit/"
	north   = "North"
	south   = "South"

	// Station constants
	st22ndStreet   = "22nd Street"
	stAtherton     = "Atherton"
	stBayshore     = "Bayshore"
	stBelmont      = "Belmont"
	stBlossomHill  = "Blossom Hill"
	stBroadway     = "Broadway"
	stBurlingame   = "Burlingame"
	stCalAve       = "California Ave"
	stCapitol      = "Capitol"
	stCollegePark  = "College Park"
	stGilroy       = "Gilroy"
	stHaywardPark  = "Hayward Park"
	stHillsdale    = "Hillsdale"
	stLawrence     = "Lawrence"
	stMenloPark    = "Menlo Park"
	stMillbrae     = "Millbrae"
	stMorganHill   = "Morgan Hill"
	stMountainView = "Mountain View"
	stPaloAlto     = "Palo Alto"
	stRedwoodCity  = "Redwood City"
	stSanAntonio   = "San Antonio"
	stSanBruno     = "San Bruno"
	stSanCarlos    = "San Carlos"
	stSanFrancisco = "San Francisco"
	stSanJose      = "San Jose Diridon"
	stSanMartin    = "San Martin"
	stSanMateo     = "San Mateo"
	stSantaClara   = "Santa Clara"
	stSouthSF      = "South San Francisco"
	stSunnyvale    = "Sunnyvale"
	stTamien       = "Tamien"
)

type station struct {
	name  string
	north int
	south int
}

// newStation creates a new station struct with the name and direction values
func newStation(name string, north, south int) station {
	return station{
		name:  name,
		north: north,
		south: south,
	}
}

// getStations returns a map of station namem to station information
func getStations() map[string]station {
	return map[string]station{
		st22ndStreet:   newStation(st22ndStreet, 70021, 70022),
		stAtherton:     newStation(stAtherton, 70151, 70152),
		stBayshore:     newStation(stBayshore, 70031, 70032),
		stBelmont:      newStation(stBelmont, 70121, 70122),
		stBlossomHill:  newStation(stBlossomHill, 70291, 70292),
		stBroadway:     newStation(stBroadway, 70071, 70072),
		stBurlingame:   newStation(stBurlingame, 70081, 70082),
		stCalAve:       newStation(stCalAve, 70191, 70192),
		stCapitol:      newStation(stCapitol, 70281, 70282),
		stCollegePark:  newStation(stCollegePark, 70251, 70252),
		stGilroy:       newStation(stGilroy, 70321, 70322),
		stHaywardPark:  newStation(stHaywardPark, 70101, 70102),
		stHillsdale:    newStation(stHillsdale, 70111, 70112),
		stLawrence:     newStation(stLawrence, 70231, 70232),
		stMenloPark:    newStation(stMenloPark, 70161, 70162),
		stMillbrae:     newStation(stMillbrae, 70061, 70062),
		stMorganHill:   newStation(stMorganHill, 70301, 70302),
		stMountainView: newStation(stMountainView, 70211, 70212),
		stPaloAlto:     newStation(stPaloAlto, 70171, 70172),
		stRedwoodCity:  newStation(stRedwoodCity, 70141, 70142),
		stSanAntonio:   newStation(stSanAntonio, 70201, 70202),
		stSanBruno:     newStation(stSanBruno, 70051, 70052),
		stSanCarlos:    newStation(stSanCarlos, 70131, 70132),
		stSanFrancisco: newStation(stSanFrancisco, 70011, 70012),
		stSanJose:      newStation(stSanJose, 70261, 70262),
		stSanMartin:    newStation(stSanMartin, 70311, 70312),
		stSanMateo:     newStation(stSanMateo, 70091, 70092),
		stSantaClara:   newStation(stSantaClara, 70241, 70242),
		stSouthSF:      newStation(stSouthSF, 70041, 70042),
		stSunnyvale:    newStation(stSunnyvale, 70221, 70222),
		stTamien:       newStation(stTamien, 70271, 70272),
	}
}
