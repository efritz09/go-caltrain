package caltrain

type holidayJson struct {
	Content struct {
		ServiceCalendar struct {
			ID       string `json:"id"`
			FromDate string `json:"FromDate"`
			ToDate   string `json:"ToDate"`
		} `json:"ServiceCalendar"`
		AvailabilityConditions []struct {
			Version  string `json:"version"`
			ID       string `json:"id"`
			FromDate string `json:"FromDate"`
			ToDate   string `json:"ToDate"`
		} `json:"AvailabilityConditions"`
	} `json:"Content"`
}

type holidayJsonAlt struct {
	Content struct {
		ServiceCalendar struct {
			ID       string `json:"id"`
			FromDate string `json:"FromDate"`
			ToDate   string `json:"ToDate"`
		} `json:"ServiceCalendar"`
		AvailabilityConditions struct {
			Version  string `json:"version"`
			ID       string `json:"id"`
			FromDate string `json:"FromDate"`
			ToDate   string `json:"ToDate"`
		} `json:"AvailabilityConditions"`
	} `json:"Content"`
}
