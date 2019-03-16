package mcast

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Device is a remotely controllable MusicCast device
// (eg. speaker or A/V receiver)
type Device struct {
	Name         string
	Manufacturer string
	ModelName    string
	ModelDesc    string
	Address      string
	ControlURL   string
	UDN          string
	SerialNum    string // Same as SystemID in DeviceInfo
	YXCVersion   string
}

// TODO - can this be combined with the responseError type?

// Internal type for processing response codes from REST API
type responseCode struct {
	Code int `json:"response_code"`
}

// SendRequest sends request to device using path as trailing part of URL. val
// contains any (optional) returned values from the response (decoded from JSON)
// args are any optional arguments to path if it is a printf style format string.
func (d Device) SendRequest(path string, val interface{}, args ...interface{}) error {

	url := fmt.Sprintf(d.ControlURL+path, args...)
	if url == "" {
		return errors.New("missing url")
	}

	//fmt.Println("URL: " + url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	//fmt.Println(req)
	resp, err := defaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//fmt.Println(resp)

	// TODO: Make an HTTP error type so we can query status code as well as see message
	if resp.StatusCode != 200 {
		return errors.New("unexpected HTTP status: " + resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// First decode only the response code
	var rc responseCode

	//fmt.Println("Start unmarshal")
	err = json.Unmarshal(b, &rc)
	if err != nil {
		return err
	}

	//fmt.Println("End unmarshal")
	if rc.Code != 0 {
		return responseError(rc.Code)
	}

	if val == nil {
		return nil
	}

	// decode other return values (if any)
	return json.Unmarshal(b, val)
}
