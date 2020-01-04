package caltrain

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
		TimetableFrame []TimetableFrame `json:"TimetableFrame"`
	} `json:"Content"`
}

type TimetableFrame struct {
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
		TimetableRouteJourney []TimetableRouteJourney `json:"ServiceJourney"`
	} `json:"vehicleJourneys"`
}

type TimetableRouteJourney struct {
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
		Call []TimetableRouteCall `json:"Call"`
	} `json:"calls"`
}

type TimetableRouteCall struct {
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

// timetableJsonAlternate is an alternate struct that the timetable json may
// require if ServiceCalendarFrame.DayTypes.DayType is a map, not an array
type timetableJsonAlternate struct {
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
				DayType struct {
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
		TimetableFrame []TimetableFrame `json:"TimetableFrame"`
	} `json:"Content"`
}
