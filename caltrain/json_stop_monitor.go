package caltrain

import "time"

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
