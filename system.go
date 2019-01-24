package mcast

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DeviceInfo struct {
	ModelName         string `json:"model_name"`
	Destination       string
	DeviceID          string  `json:"device_id"`
	SystemID          string  `json:"system_id"`
	SystemVersion     float32 `json:"system_version"`
	APIVersion        float32 `json:"api_version"`
	NetModuleGen      int     `json:"netmodule_generation"`
	NetModuleVersion  string  `json:"netmodule_version"`
	NetModuleCheckSum string  `json:"netmodule_checksum"`
	SerialNumber      string  `json:"serial_number"`
	CategoryCode      int     `json:"category_code"`
	OperationMode     string  `json:"operation_mode"`
	UpdateErrorCode   string  `json:"update_error_code"`
}

// Internal type for decoding responses to getDeviceInfo
type jsonGetDeviceInfoResp struct {
	ResponseCode uint `json:"response_code"`
	DeviceInfo
}

type Input struct {
	ID            string
	DistEnable    bool   `json:"distribution_enable"`
	RenameEnable  bool   `json:"rename_enable"`
	AccountEnable bool   `json:"account_enable"`
	PlayInfoType  string `json:"play_info_type"`
}

type RangeStep struct {
	ID   string
	Min  float32
	Max  float32
	Step float32
}

type SurroundMaster struct {
	StereoPair    bool `json:"stereo_pair"`
	SubwooferPair bool `json:"subwoofer_pair"`
	SurroundPair  bool `json:"surround_pair"`
}

type SurroundSlave struct {
	SurroundLR    bool `json:"surround_pair_lr"`
	SurroundLorR  bool `json:"surround_pair_l_or_r"`
	SubwooferPair bool `json:"subwoofer_pair"`
}

type SurroundInfo struct {
	FuncList   []string       `json:"func_list"`
	MasterRole SurroundMaster `json:"master_role"`
	SlaveRole  SurroundSlave  `json:"slave_role"`
	Version    float32        `json:"version"`
}

// TODO: Rename type
type DistributionInfo struct {
	ClientMax         int          `json:"client_max"`
	CompatibleClients []int        `json:"compatible_client"`
	ServerZoneList    []string     `json:"server_zone_list"`
	SurroundInfo      SurroundInfo `json:"mc_surround"`
	Version           float32      `json:"version"` // int does not work?
}

type SystemFeatureInfo struct {
	Functions []string `json:"func_list"`
	NumZones  uint     `json:"zone_num"`
	Inputs    []Input  `json:"input_list"`
}

type Zone struct {
	ID                   string
	ActualVolumeModeList []string    `json:"actual_volume_mode_list"`
	CursorList           []string    `json:"cursor_list"`
	FuncList             []string    `json:"func_list"`
	LinkAudioDelayList   []string    `json:"link_audio_delay_list"`
	LinkControlList      []string    `json:"link_control_list"`
	InputList            []string    `json:"input_list"`
	MenuList             []string    `json:"menu_list"`
	RangeStep            []RangeStep `json:"range_step"`
	SceneNum             int         `json:"scene_num"`
	SoundProgramList     []string    `json:"sound_program_list"`
	SurrDecoderTypeList  []string    `json:"surr_decoder_type_list"`
	ToneControlModeList  []string    `json:"tone_control_mode_list"`
}

type Features struct {
	Distribution DistributionInfo
	System       SystemFeatureInfo
	Zones        []Zone `json:"zone"`
}

// Internal type for decoding responses to getFeatures
type jsonGetFeaturesResp struct {
	ResponseCode uint `json:"response_code"`
	Features
}

// General method for handling JSON responses to HTTP requests
func unmarshalHTTPResp(method, url string, val interface{}) error {

	if url == "" {
		return errors.New("missing url")
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	//fmt.Printf("req: %+v\n", req)

	// TODO: use shared global client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("unexpected HTTP status: " + resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &val)
	if err != nil {
		return err
	}

	return nil

}

// Info returns DeviceInfo for given Device
func (d *Device) Info() (DeviceInfo, error) {

	var resp jsonGetDeviceInfoResp

	err := unmarshalHTTPResp(http.MethodGet, d.ControlURL+"system/getDeviceInfo", &resp)
	if err != nil {
		return DeviceInfo{}, err
	}

	// TODO: Ideally this should be common code in above function call
	if resp.ResponseCode != 0 {
		return DeviceInfo{}, fmt.Errorf("unexpected response code: %d", resp.ResponseCode)
	}

	// TODO: strip spaces from NetModuleVersion ?

	return resp.DeviceInfo, nil
}

func (d *Device) Features() (Features, error) {

	var resp jsonGetFeaturesResp

	err := unmarshalHTTPResp(http.MethodGet, d.ControlURL+"system/getFeatures", &resp)
	if err != nil {
		return Features{}, err
	}

	// TODO: Ideally this should be common code in above function call
	if resp.ResponseCode != 0 {
		return Features{}, fmt.Errorf("unexpected response code: %d", resp.ResponseCode)
	}

	return resp.Features, nil
}
