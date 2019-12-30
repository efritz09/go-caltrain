package caltrain

type trainStatusJson struct {
	ServiceDelivery serviceDelivery
}

type serviceDelivery struct {
	ResponseTimestamp      string
	ProducerRef            string
	Status                 bool
	StopMonitoringDelivery stopMonitoringDelivery
}

type stopMonitoringDelivery struct {
	version            string
	ResponseTimestamp  string
	Status             bool
	MonitoredStopVisit []monitoredStopVisit
}

type monitoredStopVisit struct {
	RecordedAtTime          string
	MonitoringRef           string
	MonitoredVehicleJourney monitoredVehicleJourney
}

type monitoredVehicleJourney struct {
	LineRef                 string
	DirectionRef            string
	FramedVehicleJourneyRef framedVehicleJourneyRef
	PublishedLineName       string
	OperatorRef             string
	OriginRef               string
	OriginName              string
	DestinationRef          string
	DestinationName         string
	Monitored               bool
	InCongestion            bool
	VehicleLocation         vehicleLocation
	Bearing                 int
	Occupancy               string
	VehicleRef              string
	MonitoredCall           monitoredCall
}

type framedVehicleJourneyRef struct {
	DataFrameRef           string
	DatedVehicleJourneyRef string
}

type vehicleLocation struct {
	Longitude string
	Latitude  string
}

type monitoredCall struct {
	StopPointRef          string
	StopPointName         string
	VehicleLocationAtStop string
	VehicleAtStop         string
	AimedArrivalTime      string
	ExpectedArrivalTime   string
	AimedDepartureTime    string
	ExpectedDepartureTime string
	Distances             string
}
