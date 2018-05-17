package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sanksons/raven"
)

const DESTINATION = "product"
const BUCKET = "product"

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

	var mess raven.Message = raven.PrepareMessage("", "", "Hello !!")
	var destination raven.Destination = raven.CreateDestination(DESTINATION, BUCKET)

	for {
		fmt.Printf("Publishing message [%+v]\n", mess)

		//Pick a Raven from farm
		flyerr := farm.GetRaven().
			// Hand over message to it.
			HandMessage(mess).
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
