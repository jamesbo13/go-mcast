package mcast_test

import (
	"reflect"
	"testing"

	"github.com/jamesbo13/mcast"
)

func TestDeviceZoneStatus(t *testing.T) {

	d := mcast.Device{
		ControlURL: testServer.URL + "/YamahaExtendedControl/v1/",
	}

	status, err := d.ZoneStatus("main")
	if err != nil {
		t.Fatal(err)
	}

	expected := mcast.ZoneStatus{
		Power:           "on",
		Sleep:           0,
		Volume:          84,
		Mute:            false,
		MaxVolume:       161,
		ActualVolume:    mcast.VolumeStatus{"db", -38.5, "dB"},
		LinkAudioDelay:  "balanced",
		LinkControl:     "standard",
		DialogueLevel:   0,
		Input:           "net_radio",
		InputText:       "NET RADIO",
		SoundProgram:    "surr_decoder",
		ContentsDisplay: true,
		Enhancer:        true,
		SurroundDecoder: "auto",
		ToneControl:     mcast.ToneControl{"manual", 0, 0},
	}

	if !reflect.DeepEqual(status, expected) {
		t.Errorf("Unexpected values in marshalled response\n\n exp: %+v\n got: %+v\n", expected, status)
	}
}

func TestDeviceZoneSoundPrograms(t *testing.T) {

	d := mcast.Device{
		ControlURL: testServer.URL + "/YamahaExtendedControl/v1/",
	}

	programs, err := d.ZoneSoundPrograms("main")
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"munich", "vienna", "chamber", "cellar_club", "roxy_theatre", "bottom_line",
		"sports", "action_game", "roleplaying_game", "music_video", "standard",
		"spectacle", "sci-fi", "adventure", "drama", "mono_movie", "2ch_stereo",
		"7ch_stereo", "surr_decoder", "straight",
	}

	if !reflect.DeepEqual(programs, expected) {
		t.Errorf("Unexpected values in marshalled response\n\n exp: %+v\n got: %+v\n", expected, programs)
	}
}
