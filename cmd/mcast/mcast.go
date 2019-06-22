package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jamesbo13/mcast"
)

func main() {

	// Syntax:  mcast [options] <device> <cmd> [cmd-options] [cmd-args]

	flag.Parse()

	host := flag.Arg(0)
	cmd := flag.Arg(1)

	if host == "" || cmd == "" {
		log.Fatal("Missing required arg(s)")
	}

	/*
		dev, err := mcast.Ping(host)
		if err != nil {
			log.Fatal(err)
		}
	*/

	//	dev := mcast.Device{ControlURL: "http://" + host + ":80/YamahaExtendedControl/v1/", Name: host}
	dev := mcast.Device{Address: host, Name: host}

	/*
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

	*/

	switch cmd {
	case "dev":
		fmt.Printf("%s: %+v\n", dev.Name, dev)

	case "info":
		info, err := dev.Info()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %+v\n", dev.Name, info)

	case "features":
		feat, err := dev.Features()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %+v\n", dev.Name, feat)

	case "dist":
		dist, err := dev.Distribution()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %+v\n", dev.Name, dist)

	case "status":
		z, err := dev.Zone("main")
		if err != nil {
			log.Fatal(err)
		}

		status, err := z.Status()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %+v\n", dev.Name, status)

	case "volume":
		z, err := dev.Zone("main")
		if err != nil {
			log.Fatal(err)
		}

		status, err := z.Status()
		if err != nil {
			log.Fatal(err)
		}

		vol := status.Volume / status.MaxVolume * 100.0

		fmt.Printf("%s: volume: %5.1f %%\n", dev.Name, vol)

	case "mute":
		z, err := dev.Zone("main")
		if err != nil {
			log.Fatal(err)
		}

		err = z.Mute(true)
		if err != nil {
			log.Fatal(err)
		}

	case "unmute":
		z, err := dev.Zone("main")
		if err != nil {
			log.Fatal(err)
		}

		err = z.Mute(false)
		if err != nil {
			log.Fatal(err)
		}

	case "set-volume":
		z, err := dev.Zone("main")
		if err != nil {
			log.Fatal(err)
		}

		if flag.Arg(2) == "" {
			log.Fatalf("Missing volume level")
		}

		vol, err := strconv.Atoi(flag.Arg(2))
		if err != nil {
			log.Fatal(err)
		}

		err = z.SetVolume(vol)
		if err != nil {
			log.Fatal(err)
		}

	case "set-volume-pct":
		z, err := dev.Zone("main")
		if err != nil {
			log.Fatal(err)
		}

		if flag.Arg(2) == "" {
			log.Fatalf("Missing volume level")
		}

		pct, err := strconv.ParseFloat(flag.Arg(2), 32)
		if err != nil {
			log.Fatal(err)
		}

		// TODO: Verify 0 <= pct <= 100

		err = z.SetVolumePercent(float32(pct))
		if err != nil {
			log.Fatal(err)
		}

	case "vol-incr":
		z, err := dev.Zone("main")
		if err != nil {
			log.Fatal(err)
		}

		if flag.Arg(2) == "" {
			log.Fatalf("Missing volume step value")
		}

		vol, err := strconv.Atoi(flag.Arg(2))
		if err != nil {
			log.Fatal(err)
		}

		err = z.SetVolumeIncr(vol)
		if err != nil {
			log.Fatal(err)
		}

	case "power-on":
		z, err := dev.Zone("main")
		if err != nil {
			log.Fatal(err)
		}

		err = z.SetPower(true)
		if err != nil {
			log.Fatal(err)
		}

	case "power-off":
		z, err := dev.Zone("main")
		if err != nil {
			log.Fatal(err)
		}

		err = z.SetPower(false)
		if err != nil {
			log.Fatal(err)
		}

	case "power-toggle":
		z, err := dev.Zone("main")
		if err != nil {
			log.Fatal(err)
		}

		err = z.TogglePower()
		if err != nil {
			log.Fatal(err)
		}

	case "test-dist":
		server := mcast.Device{Address: "10.13.1.153"}
		serverZone, _ := server.Zone("main")

		//		server.LeaveGroup()

		client := mcast.Device{Address: "10.13.1.152"}
		clientZone, _ := client.Zone("main")

		g, err := mcast.NewGroup(serverZone, clientZone)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("group: %+v\n", g)

		distInfo, _ := server.Distribution()
		fmt.Printf("server: %s %+v\n", server.Address, distInfo)

		distInfo, _ = client.Distribution()
		fmt.Printf("client: %s %+v\n", client.Address, distInfo)

		time.Sleep(time.Second * 100)
		fmt.Println("Deleting")

		err = g.Delete()
		if err != nil {
			log.Fatal(err)
		}

		distInfo, _ = server.Distribution()
		fmt.Printf("server: %s %+v\n", server.Address, distInfo)

		distInfo, _ = client.Distribution()
		fmt.Printf("client: %s %+v\n", client.Address, distInfo)

	default:
		fmt.Printf("Unknown command '%s'\n", cmd)
	}
}
