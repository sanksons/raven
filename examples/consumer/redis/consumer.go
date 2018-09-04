package main

import (
	"fmt"
	"log"

	"github.com/newrelic/go-agent"

	"github.com/sanksons/raven"
)

const SOURCE = "productQ"
const BUCKET = "1"

func main() {

	//
	// Initialize raven farm.
	//
	loggerStrict := new(raven.FmtLogger)
	loggerStrict.Level = 10

	farm, err := raven.InitializeFarm(raven.FARM_TYPE_REDIS, raven.RedisSimpleConfig{
		Addr:     "localhost:6379",
		PoolSize: 10,
	},
		loggerStrict,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Define a source from which to receive.
	var source raven.Source = raven.CreateSource(SOURCE, 2)

	// Initiate and pick a receiver.
	receiver, err := farm.GetRavenReceiver("one", source)
	if err != nil {
		log.Fatal(err)
	}

	// Mark as Reliable and Ordered.
	receiver.MarkReliable()

	//start receiving
	err1 := receiver.Start(c)
	if err1 != nil {
		log.Fatal(err1)
	}
}

func c(message *raven.Message, txn newrelic.Transaction) error {
	//time.Sleep(1 * time.Minute)
	fmt.Printf("Got message: %s\n", message)
	return nil
}
