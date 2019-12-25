package main

import (
	"context"
	"fmt"

	"github.com/efritz09/go-caltrain/caltrain"
)

func main() {
	ctx := context.Background()
	c := caltrain.New("844665d9-db36-4209-8a8a-6f49a53f8e6a")

	fmt.Println("Calling GetDelays")
	t, err := c.GetDelays(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(t)
	fmt.Println()

	fmt.Println("Calling GetStationStatus on Hillsdale North")
	s, err := c.GetStationStatus(ctx, caltrain.StationHillsdale, caltrain.North)
	if err != nil {
		panic(err)
	}
	fmt.Println(s)

	fmt.Println("Calling GetStationStatus on Hillsdale South")
	s, err = c.GetStationStatus(ctx, caltrain.StationHillsdale, caltrain.South)
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}
