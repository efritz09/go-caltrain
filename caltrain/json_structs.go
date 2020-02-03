package caltrain

// json_structs contains all structs for the JSON responses from 511.org. Some
// structs are broken out where appropriate to provide the methods easier
// access to the information

import "time"

type timetableJson struct {
	Content struct {
		ServiceFrame struct {
			ID     string `json:"id"`
			Routes struct {
				Route []struct {
					ID      string `json:"id"`
					Name    string `json:"Name"`
					LineRef struct {
						Ref string `json:"ref"`
					} `json:"LineRef"`
					DirectionRef struct {
						Ref string `json:"ref"`
					} `json:"DirectionRef"`
					PointsInSequence struct {
						PointOnRoute []struct {
							ID       string `json:"id"`
							PointRef struct {
								Ref  string `json:"ref"`
								Type string `json:"type"`
							} `json:"PointRef"`
						} `json:"PointOnRoute"`
					} `json:"pointsInSequence"`
				} `json:"Route"`
			} `json:"routes"`
		} `json:"ServiceFrame"`
		ServiceCalendarFrame struct {
			ID       string `json:"id"`
			DayTypes struct {
				DayType []struct {
					ID         string `json:"id"`
					Name       string `json:"Name"`
					Properties struct {
						PropertyOfDay struct {
							DaysOfWeek string `json:"DaysOfWeek"`
						} `json:"PropertyOfDay"`
					} `json:"properties"`
				} `json:"DayType"`
			} `json:"dayTypes"`
			DayTypeAssignments struct {
				DayTypeAssignment struct {
					DayTypeRef interface{} `json:"DayTypeRef"`
				} `json:"DayTypeAssignment"`
			} `json:"dayTypeAssignments"`
		} `json:"ServiceCalendarFrame"`
		TimetableFrame []timetableFrame `json:"TimetableFrame"`
	} `json:"Content"`
}

type timetableFrame struct {
	ID                      string `json:"id"`
	Name                    string `json:"Name"`
	FrameValidityConditions struct {
		AvailabilityCondition struct {
			ID       string `json:"id"`
			FromDate string `json:"FromDate"`
			ToDate   string `json:"ToDate"`
			DayTypes struct {
				DayTypeRef struct {
					Ref string `json:"ref"`
				} `json:"DayTypeRef"`
			} `json:"dayTypes"`
		} `json:"AvailabilityCondition"`
	} `json:"frameValidityConditions"`
	VehicleJourneys struct {
		TimetableRouteJourney []timetableRouteJourney `json:"ServiceJourney"`
	} `json:"vehicleJourneys"`
}

type timetableRouteJourney struct {
	Line                  string // is not in the json, added for convenience
	ID                    string `json:"id"`
	SiriVehicleJourneyRef string `json:"SiriVehicleJourneyRef"`
	JourneyPatternView    struct {
		RouteRef struct {
			Ref string `json:"ref"`
		} `json:"RouteRef"`
		DirectionRef struct {
			Ref string `json:"ref"`
		} `json:"DirectionRef"`
	} `json:"JourneyPatternView"`
	Calls struct {
		Call []timetableRouteCall `json:"Call"`
	} `json:"calls"`
}

type timetableRouteCall struct {
	Order                 string `json:"order"`
	ScheduledStopPointRef struct {
		Ref string `json:"ref"`
	} `json:"ScheduledStopPointRef"`
	Arrival struct {
		Time       string `json:"Time"`
		DaysOffset string `json:"DaysOffset"`
	} `json:"Arrival"`
	Departure struct {
		Time       string `json:"Time"`
		DaysOffset string `json:"DaysOffset"`
	} `json:"Departure"`
}

type trainStatusJson struct {
	ServiceDelivery struct {
		ResponseTimestamp      time.Time `json:"ResponseTimestamp"`
		ProducerRef            string    `json:"ProducerRef"`
		Status                 bool      `json:"Status"`
		StopMonitoringDelivery struct {
			Version            string    `json:"version"`
			ResponseTimestamp  time.Time `json:"ResponseTimestamp"`
			Status             bool      `json:"Status"`
			MonitoredStopVisit []struct {
				RecordedAtTime          time.Time `json:"RecordedAtTime"`
				MonitoringRef           string    `json:"MonitoringRef"`
				MonitoredVehicleJourney struct {
					LineRef                 string `json:"LineRef"`
					DirectionRef            string `json:"DirectionRef"`
					FramedVehicleJourneyRef struct {
						DataFrameRef           string `json:"DataFrameRef"`
						DatedVehicleJourneyRef string `json:"DatedVehicleJourneyRef"`
					} `json:"FramedVehicleJourneyRef"`
					PublishedLineName string      `json:"PublishedLineName"`
					OperatorRef       string      `json:"OperatorRef"`
					OriginRef         string      `json:"OriginRef"`
					OriginName        string      `json:"OriginName"`
					DestinationRef    string      `json:"DestinationRef"`
					DestinationName   string      `json:"DestinationName"`
					Monitored         bool        `json:"Monitored"`
					InCongestion      interface{} `json:"InCongestion"`
					VehicleLocation   struct {
						Longitude string `json:"Longitude"`
						Latitude  string `json:"Latitude"`
					} `json:"VehicleLocation"`
					Bearing       interface{}   `json:"Bearing"`
					Occupancy     interface{}   `json:"Occupancy"`
					VehicleRef    string        `json:"VehicleRef"`
					MonitoredCall monitoredCall `json:"MonitoredCall"`
				} `json:"MonitoredVehicleJourney"`
			} `json:"MonitoredStopVisit"`
		} `json:"StopMonitoringDelivery"`
	} `json:"ServiceDelivery"`
}

type monitoredCall struct {
	StopPointRef          string    `json:"StopPointRef"`
	StopPointName         string    `json:"StopPointName"`
	VehicleLocationAtStop string    `json:"VehicleLocationAtStop"`
	VehicleAtStop         string    `json:"VehicleAtStop"`
	AimedArrivalTime      time.Time `json:"AimedArrivalTime"`
	ExpectedArrivalTime   time.Time `json:"ExpectedArrivalTime"`
	AimedDepartureTime    time.Time `json:"AimedDepartureTime"`
	ExpectedDepartureTime time.Time `json:"ExpectedDepartureTime"`
	Distances             string    `json:"Distances"`
}

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
