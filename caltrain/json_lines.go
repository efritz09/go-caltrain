package caltrain

type lineJson []struct {
	ID            string `json:"Id"`
	Name          string `json:"Name"`
	TransportMode string `json:"TransportMode"`
	PublicCode    string `json:"PublicCode"`
	SiriLineRef   string `json:"SiriLineRef"`
	Monitored     bool   `json:"Monitored"`
	OperatorRef   string `json:"OperatorRef"`
}
