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
		Addrs:    []string{"localhost:30001"},
		PoolSize: 10,
	})
	if err != nil {
		log.Fatal(err)
	}

	var destination raven.Destination = raven.Destination{
		Name: "Asia",
	}

	collector := farm.MessageCollector(destination)

	collector.Start(c)

}

func c(message string) error {

}
