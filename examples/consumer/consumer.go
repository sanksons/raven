package main

import (
	"fmt"
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

	var source raven.Source = raven.CreateSource("product", "product")

	receiver, err := farm.GetRavenReceiver("one", source)
	if err != nil {
		log.Fatal(err)
	}
	receiver.MarkReliable().MarkOrdered()
	receiver.Start(c)
}

func c(message string) error {

	fmt.Printf("Got message: %s\n", message)
	return nil
}
