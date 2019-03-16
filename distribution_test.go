package mcast_test

import (
	"reflect"
	"testing"

	"github.com/jamesbo13/mcast"
)

func TestDeviceDistribution(t *testing.T) {

	d := mcast.Device{
		ControlURL: testServer.URL + "/YamahaExtendedControl/v1/",
	}

	dist, err := d.Distribution()
	if err != nil {
		t.Fatal(err)
	}

	expected := mcast.Distribution{
		GroupID:    "8723d1ceb5cb47c9b4cd4a2a43da23d8",
		GroupName:  "Living Room +3 Rooms",
		Role:       "server",
		Status:     "working",
		ServerZone: "main",
		Clients: []mcast.DistClient{
			{IPAddr: "10.13.1.152", DataType: "base"},
			{IPAddr: "10.13.1.153", DataType: "base"},
			{IPAddr: "10.13.1.154", DataType: "base"},
		},
		AudioDropout: false,
	}

	if !reflect.DeepEqual(dist, expected) {
		t.Errorf("Unexpected values in marshalled response\n\n exp: %+v\n got: %+v\n", expected, dist)
	}
}
