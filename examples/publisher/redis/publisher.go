package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/sanksons/raven"
)

const DESTINATION = "productQ"
const BUCKET = "1"

type m struct {
	Name  string
	Id    int
	Class string
}

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

		d := m{
			Name:  fmt.Sprintf("I am message number: %d", counter),
			Id:    counter,
			Class: "My class",
		}
		dataBytes, _ := json.Marshal(d)
		var mess raven.Message = raven.PrepareMessage(
			strconv.Itoa(counter), "", string(dataBytes), strconv.Itoa(counter))

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
