package mcast

import "errors"

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

// GetRequest is a helper function for sending requests to the Zone API.
// See Device.GetRequest() for additional information
func (z Zone) GetRequest(path string, val interface{}, args ...interface{}) error {
	return z.dev.GetRequest(z.name+"/"+path, val, args...)
}

// Status returns the current status of the zone
func (z Zone) Status() (ZoneStatus, error) {

	var st ZoneStatus

	err := z.GetRequest("getStatus", &st)
	return st, err
}

// SoundPrograms returns all supported sound programs for the zone
func (z Zone) SoundPrograms() ([]string, error) {

	var resp struct {
		Programs []string `json:"sound_program_list"`
	}

	err := z.GetRequest("getSoundProgramList", &resp)
	if err != nil {
		return nil, err
	}

	return resp.Programs, nil
}

// Mute enables or disables muting of sound output for the zone.
func (z Zone) Mute(enable bool) error {

	return z.GetRequest("setMute?enable=%t", nil, enable)
}

// SetVolume will set the volume to the provided value
func (z Zone) SetVolume(val int) error {

	// TODO: Check min/max value

	return z.GetRequest("setVolume?volume=%d", nil, val)
}

// SetVolumePercent sets the volume to provided percentage of max volume
func (z Zone) SetVolumePercent(pct float32) error {

	if pct < 0 || pct > 100 {
		return errors.New("out of range value")
	}

	st, err := z.Status()
	if err != nil {
		return err
	}

	vol := int(st.MaxVolume * pct / 100.0)

	return z.SetVolume(vol)
}

// SetVolumeIncr with increment/decrement the volume by the specified step value
func (z Zone) SetVolumeIncr(step int) error {

	var dir string

	switch {
	case step > 0:
		dir = "up"
	case step < 0:
		dir = "down"
		step *= -1 // Need to pass positive step value to API
	default:
		return nil
	}

	return z.GetRequest("setVolume?volume=%s&step=%d", nil, dir, step)
}
}
