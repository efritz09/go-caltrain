package caltrain

type stationJson struct {
	Contents struct {
		ResponseTimestamp string `json:"ResponseTimestamp"`
		DataObjects       struct {
			ID                 string               `json:"id"`
			ScheduledStopPoint []scheduledStopPoint `json:"ScheduledStopPoint"`
		} `json:"dataObjects"`
	} `json:"Contents"`
}

type scheduledStopPoint struct {
	ID         string `json:"id"`
	Extensions struct {
		LocationType  string      `json:"LocationType"`
		PlatformCode  interface{} `json:"PlatformCode"`
		ParentStation interface{} `json:"ParentStation"`
	} `json:"Extensions"`
	Name     string `json:"Name"`
	Location struct {
		Longitude string `json:"Longitude"`
		Latitude  string `json:"Latitude"`
	} `json:"Location"`
	URL      interface{} `json:"Url"`
	StopType string      `json:"StopType"`
}
