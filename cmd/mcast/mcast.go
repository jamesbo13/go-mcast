package main

import (
	"log"

	"github.com/jamesbo13/mcast"
)

func main() {
	err := mcast.Discover()
	if err != nil {
		log.Fatal(err)
	}
}
