package caltrain

type stationJson struct {
	Contents struct {
		ResponseTimestamp string `json:"ResponseTimestamp"`
		DataObjects       struct {
			ID                 string               `json:"id"`
			ScheduledStopPoint []scheduledStopPoint `json:"ScheduledStopPoint"`
			StopAreas          interface{}          `json:"stopAreas"`
		} `json:"dataObjects"`
	} `json:"Contents"`
}

type scheduledStopPoint struct {
	ID       string `json:"id"`
	Name     string `json:"Name"`
	Location struct {
		Longitude string `json:"Longitude"`
		Latitude  string `json:"Latitude"`
	} `json:"Location"`
	StopType string `json:"StopType"`
}
