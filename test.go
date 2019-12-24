package main

import (
	"fmt"

	"github.com/efritz09/go-caltrain/caltrain"
)

func main() {
	c := caltrain.New("844665d9-db36-4209-8a8a-6f49a53f8e6a")

	fmt.Println("Calling GetDelays")
	t, err := c.GetDelays()
	if err != nil {
		panic(err)
	}
	fmt.Println(t)
	fmt.Println()

	fmt.Println("Calling GetStationStatus")
	s, err := c.GetStationStatus()
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}
