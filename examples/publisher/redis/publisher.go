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
		raven.FARM_TYPE_REDIS,
		raven.RedisSimpleConfig{
			Addr:     "localhost:6379",
			PoolSize: 10,
		},
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	var destination raven.Destination = raven.CreateDestination(DESTINATION, 2, nil)

	var counter int
	for {
		counter++

		var mess raven.Message = raven.PrepareMessage(
			strconv.Itoa(counter), "", fmt.Sprintf("Hello %d!!", counter))

		fmt.Printf("Publishing message [%+v]\n", mess)

		box, err := destination.GetBox4Msg(mess)
		if err != nil {
			log.Fatal("sdsdsdsd")
		}
		fmt.Println(box.GetName())

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
