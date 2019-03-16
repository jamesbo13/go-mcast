package mcast

import "net/http"

type VolumeStatus struct {
	Mode  string
	Value float32
	Unit  string
}

type ToneControl struct {
	Mode   string
	Bass   int
	Treble int
}

type Equalizer struct {
	Mode string
	High int
	Low  int
	Mid  int
}

type ZoneStatus struct {
	Power               string
	Sleep               int
	Volume              float32
	Mute                bool
	MaxVolume           float32      `json:"max_volume"`
	ActualVolume        VolumeStatus `json:"actual_volume"`
	LinkAudioDelay      string       `json:"link_audio_delay"`
	LinkAudioQuality    string       `json:"link_audio_quality"`
	LinkControl         string       `json:"link_control"`
	DialogueLevel       int          `json:"dialogue_level"`
	Input               string
	InputText           string `json:"input_text"`
	DistributionEnabled bool   `json:"distribution_enabled"`
	Direct              bool
	DisableFlags        int `json:"disable_flags"`
	Enhancer            bool
	ContentsDisplay     bool        `json:"contents_display"`
	PartyEnable         bool        `json:"Party_enable"`
	PureDirect          bool        `json:"pure_direct"`
	SoundProgram        string      `json:"sound_program"`
	SurroundDecoder     string      `json:"surr_decoder_type"`
	ToneControl         ToneControl `json:"tone_control"`
	Equalizer           Equalizer

	// TODO: the following fields are listed in docs but not present in sample output:
	//		balance, dialogue_lift, clear_voice, subwooder_volume, bass_extension
}

type Zone struct {
	name string
	dev  *Device
}

func (d Device) Zone(name string) (Zone, error) {

	// TODO: Check if name is valid for given device

	return Zone{name: name, dev: &d}, nil
	}

// SendRequest is a helper function for sending requests to the Zone API.
// See Device.SendRequest() for additional information
func (z Zone) SendRequest(path string, val interface{}, args ...interface{}) error {
	return z.dev.SendRequest(z.name+"/"+path, val, args...)
	}

// Status returns the current status of the zone
func (z Zone) Status() (ZoneStatus, error) {

	var st ZoneStatus

	err := z.SendRequest("getStatus", &st)
	return st, err
}

// SoundPrograms returns all supported sound programs for the zone
func (z Zone) SoundPrograms() ([]string, error) {

	var resp struct {
		Programs []string `json:"sound_program_list"`
	}

	err := z.SendRequest("/getSoundProgramList", &resp)
	if err != nil {
		return nil, err
	}

	return resp.Programs, nil
}
