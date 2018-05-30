package main

import (
	"fmt"
	"log"

	"github.com/sanksons/raven"
)

const SOURCE = "product1"
const BUCKET = "1"

func main() {

	//
	// Initialize raven farm.
	//
	loggerStrict := new(raven.FmtLogger)

	farm, err := raven.InitializeFarm(raven.FARM_TYPE_REDISCLUSTER, raven.RedisClusterConfig{
		Addrs:    []string{"172.17.0.2:30001"},
		PoolSize: 10,
	},
		loggerStrict,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Define a source from which to receive.
	var source raven.Source = raven.CreateSource(SOURCE, BUCKET)

	// Initiate and pick a receiver.
	receiver, err := farm.GetRavenReceiver("one", source)
	if err != nil {
		log.Fatal(err)
	}

	// Mark as Reliable and ordered.
	//receiver.MarkReliable().MarkOrdered()

	//start receiving
	receiver.Start(c)
}

func c(message *raven.Message) error {

	fmt.Printf("Got message: %s\n", message)
	return nil
}
