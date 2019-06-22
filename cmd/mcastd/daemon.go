package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jamesbo13/mcast"
)

// Global
var devices []mcast.Device

func main() {

	var configFile string
	var conf config

	flag.StringVar(&configFile, "config", "", "Name of config file")

	flag.Parse()

	if configFile != "" {
		f, err := os.Open(configFile)
		if err != nil {
			log.Fatal(err)
		}

		conf, err = parseConfig(f)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("%+v\n", conf)

	// Discover devices
	var err error 
	devices, err = mcast.Discover()
	if err != nil {
		log.Fatal(err)
	}

	// start up web server
	s := http.Server{
		Addr:         ":" + strconv.Itoa(conf.Port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	http.HandleFunc("/devices/", listDevices)

	log.Fatal(s.ListenAndServe())
}
