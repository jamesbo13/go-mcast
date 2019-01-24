package main

import (
	"fmt"
	"log"

	"github.com/jamesbo13/mcast"
)

func main() {
	devs, err := mcast.Discover()
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range devs {
		info, err := d.Info()
		if err != nil {
			log.Printf("error processing device %s: %s", d.Name, err)
			continue
		}

		feats, err := d.Features()
		if err != nil {
			log.Printf("error processing device %s: %s", d.Name, err)
			continue
		}

		fmt.Printf("%s: %+v %+v %+v\n\n", d.Name, d, info, feats)
	}
}
