package main

import (
	"fmt"
	"log"

	"github.com/sanksons/raven/childlock"
)

func main() {

	lock := childlock.New("foo.lock", 300, childlock.RedisOptions{
		Addres: []string{"localhost:6379"},
	})
	err := lock.Acquire("vvvv")
	if err != nil {
		log.Fatal(err)
	}
	rerr := lock.Release()
	if rerr != nil {
		log.Fatal(rerr)
	}
	fmt.Println("fine")

}
