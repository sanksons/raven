package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sanksons/raven/childlock"

	"github.com/newrelic/go-agent"

	"github.com/sanksons/raven"
)

const SOURCE = "productQ"
const BUCKETS = 2

func main() {

	//
	// Initialize raven farm.
	//
	loggerStrict := new(raven.DummyLogger)
	//loggerStrict.Level = 10

	farm, err := raven.InitializeFarm(raven.FARM_TYPE_REDIS, raven.RedisSimpleConfig{
		Addr:     "localhost:6379",
		PoolSize: 10,
	},
		loggerStrict,
	)
	if err != nil {
		log.Fatal(err)
	}

	//Make sure lock details are attached.
	farm.AttachLock(childlock.RedisOptions{
		Addres: []string{"localhost:6379"},
	})

	// Define a source from which to receive.
	var source raven.Source = raven.CreateSource(SOURCE, BUCKETS)

	// Initiate and pick a receiver.
	receiver, err := farm.GetRavenReceiver("", source)
	if err != nil {
		log.Fatal(err)
	}

	// Mark as Reliable and Ordered.
	receiver.MarkReliable()
	receiver.SetPort("9001")
	//receiver.SetPort("6379")

	//start receiving

	err1 := receiver.Start(c)
	if err1 != nil {
		log.Fatal(err1)
	}
}

func c(message *raven.Message, txn newrelic.Transaction) error {
	time.Sleep(1 * time.Second)
	fmt.Printf("Got message: %s\n", message)
	return fmt.Errorf("sdsd")
}
