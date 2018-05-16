package main

import (
	"log"

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

	var message raven.Message = raven.Message{Data: "Hello Raven !!"}
	var destination raven.Destination = raven.Destination{
		raven.Q{Name: "Asia"},
	}

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

}
