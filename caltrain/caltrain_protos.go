package caltrain

import (
	"context"
	"fmt"

	gtfs "github.com/efritz09/go-caltrain/caltrain/transit_realtime"
	"github.com/golang/protobuf/proto"
)

const (
	tripUpdatesURL      = "http://api.511.org/Transit/TripUpdates"
	vehiclePositionsURL = "http://api.511.org/Transit/VehiclePositions"
)

func (c *CaltrainClient) TripUpdates(ctx context.Context) error {
	query := map[string]string{
		"agency":  "CT",
		"api_key": c.key,
	}
	raw, err := c.APIClient.Get(ctx, tripUpdatesURL, query)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	// TODO: unmarshal
	data := &gtfs.FeedMessage{}
	err = proto.Unmarshal(raw, data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}
	fmt.Printf("%+v\n\n", data)
	return nil

}

func (c *CaltrainClient) VehiclePositions(ctx context.Context) error {
	query := map[string]string{
		"agency":  "CT",
		"api_key": c.key,
	}
	raw, err := c.APIClient.Get(ctx, vehiclePositionsURL, query)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	// TODO: unmarshal
	data := &gtfs.FeedMessage{}
	err = proto.Unmarshal(raw, data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}
	fmt.Printf("%+v\n\n", data)
	return nil

}
