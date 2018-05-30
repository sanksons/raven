package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/sanksons/raven"
)

const DESTINATION = "product1"
const BUCKET = "1"

func main() {

	//
	// Initialize raven farm.
	//
	//  loggerStrict := new(raven.FmtLogger)
	//  logger := raven.Logger(loggerStrict)
	farm, err := raven.InitializeFarm(
		raven.FARM_TYPE_REDISCLUSTER,
		raven.RedisClusterConfig{
			Addrs:    []string{"172.17.0.2:30001"},
			PoolSize: 10,
		},
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	var destination raven.Destination = raven.CreateDestination(DESTINATION, BUCKET)

	var counter int
	for {
		counter++

		var mess raven.Message = raven.PrepareMessage(
			strconv.Itoa(counter), "", fmt.Sprintf("Hello %d!!", counter))

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
