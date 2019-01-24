package mcast_test

import (
	"reflect"
	"testing"

	"github.com/jamesbo13/mcast"
)

func TestDeviceInfo(t *testing.T) {

	d := mcast.Device{
		ControlURL: testServer.URL + "/YamahaExtendedControl/v1/",
	}

	info, err := d.Info()
	if err != nil {
		t.Fatal(err)
	}

	expected := mcast.DeviceInfo{
		ModelName:         "TSR-7850",
		Destination:       "U",
		DeviceID:          "946AB03A61CB",
		SystemID:          "0DA87303",
		SystemVersion:     1.47,
		APIVersion:        2.02,
		NetModuleGen:      2,
		NetModuleVersion:  "0502    ",
		NetModuleCheckSum: "BB1E23E6",
		SerialNumber:      "Y022918RT",
		CategoryCode:      1,
		OperationMode:     "normal",
		UpdateErrorCode:   "00000000",
	}

	if !reflect.DeepEqual(info, expected) {
		t.Errorf("Unexpected values in response\n")
	}
}

func TestDeviceFeatures(t *testing.T) {

	d := mcast.Device{
		ControlURL: testServer.URL + "/YamahaExtendedControl/v1/",
	}

	f, err := d.Features()
	if err != nil {
		t.Fatal(err)
	}

	expected := mcast.Features{
		Distribution: mcast.DistributionInfo{
			ClientMax:         19,
			CompatibleClients: []int{2},
			ServerZoneList:    []string{"main", "zone2"},
			SurroundInfo: mcast.SurroundInfo{
				FuncList: []string{},
				MasterRole: mcast.SurroundMaster{
					StereoPair:    false,
					SubwooferPair: true,
					SurroundPair:  true,
				},
				SlaveRole: mcast.SurroundSlave{
					SurroundLR:    false,
					SurroundLorR:  false,
					SubwooferPair: false,
				},
				Version: 1,
			},
			Version: 2,
		},
		System: mcast.SystemFeatureInfo{
			Functions: []string{"wired_lan", "wireless_lan", "network_standby", "network_standby_auto", "bluetooth_standby", "bluetooth_tx_setting", "dfs_option", "hdmi_out_1", "hdmi_out_2", "airplay", "disklavier_settings", "background_download", "remote_info", "network_reboot", "system_reboot", "party_mode"},
			NumZones:  2,
			Inputs: []mcast.Input{
				{"napster", true, false, true, "netusb"},
				{"siriusxm", true, false, true, "netusb"},
				{"pandora", true, false, true, "netusb"},
				{"spotify", true, false, false, "netusb"},
				{"tidal", true, false, true, "netusb"},
				{"deezer", true, false, true, "netusb"},
				{"airplay", false, false, false, "netusb"},
				{"mc_link", false, true, false, "netusb"},
				{"server", true, true, false, "netusb"},
				{"net_radio", true, true, false, "netusb"},
				{"bluetooth", true, false, false, "netusb"},
				{"usb", true, true, false, "netusb"},
				{"tuner", true, true, false, "tuner"},
				{"hdmi1", true, true, false, "none"},
				{"hdmi2", true, true, false, "none"},
				{"hdmi3", true, true, false, "none"},
				{"hdmi4", true, true, false, "none"},
				{"hdmi5", true, true, false, "none"},
				{"av1", true, true, false, "none"},
				{"av2", true, true, false, "none"},
				{"audio1", true, true, false, "none"},
				{"audio2", true, true, false, "none"},
				{"audio3", true, true, false, "none"},
				{"audio4", true, true, false, "none"},
				{"audio5", true, true, false, "none"},
				{"phono", true, true, false, "none"},
				{"aux", true, true, false, "none"},
				{"main_sync", true, true, false, "none"},
			},
			// TODO: bluetooth
			// TODO: web_control_url
		},
		Zones: []mcast.Zone{
			{
				ID:                   "main",
				ActualVolumeModeList: []string{"db", "numeric"},
				CursorList:           []string{"up", "down", "left", "right", "select", "return"},
				FuncList: []string{
					"power", "sleep", "volume", "mute", "sound_program", "pure_direct",
					"enhancer", "tone_control", "dialogue_level", "signal_info",
					"prepare_input_change", "link_control", "link_audio_delay", "scene",
					"contents_display", "cursor", "menu", "actual_volume", "surr_decoder_type",
				},
				InputList: []string{
					"napster", "siriusxm", "pandora", "spotify", "tidal", "deezer",
					"airplay", "mc_link", "server", "net_radio", "bluetooth", "usb",
					"tuner", "hdmi1", "hdmi2", "hdmi3", "hdmi4", "hdmi5", "av1", "av2",
					"audio1", "audio2", "audio3", "audio4", "audio5", "phono", "aux",
				},
				LinkAudioDelayList: []string{"audio_sync", "balanced", "lip_sync"},
				LinkControlList:    []string{"speed", "standard", "stability"},
				MenuList: []string{
					"on_screen", "top_menu", "menu", "option",
					"display", "help", "home", "mode", "red",
					"green", "yellow", "blue",
				},
				RangeStep: []mcast.RangeStep{
					{ID: "volume", Min: 0, Max: 161, Step: 1},
					{ID: "tone_control", Min: -12, Max: 12, Step: 1},
					{ID: "dialogue_level", Min: 0, Max: 3, Step: 1},
					{ID: "actual_volume_db", Min: -80.5, Max: 16.5, Step: 0.5},
					{ID: "actual_volume_numeric", Min: 0, Max: 97, Step: 0.5},
				},
				SceneNum: 8,
				SoundProgramList: []string{
					"munich", "vienna", "chamber", "cellar_club", "roxy_theatre", "bottom_line",
					"sports", "action_game", "roleplaying_game", "music_video", "standard", "spectacle",
					"sci-fi", "adventure", "drama", "mono_movie", "2ch_stereo", "7ch_stereo", "surr_decoder",
					"straight",
				},
				SurrDecoderTypeList: []string{
					"toggle", "auto", "dolby_surround", "dts_neural_x", "dts_neo6_cinema", "dts_neo6_music",
				},
				ToneControlModeList: []string{"manual"},
			},
			{
				ID:                   "zone2",
				ActualVolumeModeList: []string{"db", "numeric"},

				FuncList: []string{
					"power", "sleep", "volume", "mute", "enhancer", "tone_control",
					"prepare_input_change", "link_control", "scene", "actual_volume",
				},
				InputList: []string{
					"napster", "siriusxm", "pandora", "spotify", "tidal", "deezer",
					"airplay", "mc_link", "server", "net_radio", "bluetooth", "usb",
					"tuner", "av1", "av2", "audio1", "audio2", "audio3", "audio4",
					"audio5", "phono", "aux", "main_sync",
				},
				LinkControlList: []string{"speed", "standard", "stability"},
				RangeStep: []mcast.RangeStep{
					{ID: "volume", Min: 0, Max: 161, Step: 1},
					{ID: "tone_control", Min: -12, Max: 12, Step: 1},
					{ID: "actual_volume_db", Min: -80.5, Max: 10, Step: 0.5},
					{ID: "actual_volume_numeric", Min: 0, Max: 90.5, Step: 0.5},
				},
				SceneNum:            8,
				ToneControlModeList: []string{"manual", "auto", "bypass"},
			},
		},
		// TODO: tuner
		// TODO: netusb
		// TODO: ccs
	}

	if !reflect.DeepEqual(f, expected) {
		t.Errorf("Unexpected values in marshalled response\n\n exp: %+v\n got: %+v\n", expected, f)
	}
}

func TestDeviceNetworkStatus(t *testing.T) {

	d := mcast.Device{
		ControlURL: testServer.URL + "/YamahaExtendedControl/v1/",
	}

	stat, err := d.NetworkStatus()
	if err != nil {
		t.Fatal(err)
	}

	expected := mcast.NetworkStatus{
		Name:           "Living Room",
		Connection:     "wired_lan",
		DHCP:           true,
		IPAddress:      "10.13.1.134",
		SubnetMask:     "255.255.255.0",
		DefaultGateway: "10.13.1.1",
		DNSServer1:     "10.13.1.1",
		DNSServer2:     "0.0.0.0",
		WirelessLAN: mcast.WirelessNet{
			SSID:     "Ferrari-2.4",
			Type:     "wpa2-psk(aes)",
			Key:      "",
			Enable:   true,
			Chan:     0,
			Strength: 0,
		},
		MACAddress: map[string]string{
			"wireless_lan":    "946AB03A61CC",
			"wired_lan":       "946AB03A61CB",
			"wireless_direct": "946AB03A61CC",
		},
		AirplayPIN: "",
		MusicCastNet: mcast.MusicCastNet{
			Ready:      true,
			DeviceType: "standard",
			NumClients: 0,
			Chan:       0,
		},
	}

	if !reflect.DeepEqual(stat, expected) {
		t.Errorf("Unexpected values in marshalled response\n\n exp: %+v\n got: %+v\n", expected, stat)
	}
}
