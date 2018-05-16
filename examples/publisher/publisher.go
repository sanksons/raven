package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sanksons/raven"
)

func main() {

	//
	// Initialize raven farm.
	//
	farm, err := raven.InitializeFarm(raven.FARM_TYPE_REDISCLUSTER, raven.RedisClusterConfig{
		Addrs:    []string{"172.17.0.2:30001"},
		PoolSize: 10,
	})
	if err != nil {
		log.Fatal(err)
	}

	var message raven.Message = raven.PrepareMessage("", "", "Hello !!")
	var destination raven.Destination = raven.CreateDestination("product", "product")

	for {
		fmt.Println("Publishing message")

		//Pick a Raven from farm
		flyerr := farm.GetRaven().
			// Hand over message to it.
			HandMessage(message).
			// Define Destination.
			SetDestination(destination).
			// make it fly.
			Fly()

		if flyerr != nil {
			log.Fatal(flyerr)
		}
		time.Sleep(2 * time.Second)
	}

}
