package mcast

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	defaultControlURL string = "YamahaExtendedControl/v1/"
)

// TODO - fields in Device are taken from desc.xml fields ... but are they what we need?

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

// GetRequest sends request to device using path as trailing part of URL. val
// contains any (optional) returned values from the response (decoded from JSON)
// args are any optional arguments to path if it is a printf style format string.
func (d Device) GetRequest(path string, val interface{}, args ...interface{}) error {
	return d.sendRequest(http.MethodGet, path, val, nil, args...)
}

func (d Device) PostRequest(path string, body io.Reader, args ...interface{}) error {

	return d.sendRequest(http.MethodPost, path, nil, body, args...)
}

func (d Device) sendRequest(method, path string, val interface{}, body io.Reader, args ...interface{}) error {

	var baseURL string

	// Use default ControlURL if not defined
	if d.ControlURL == "" {
		if d.Address == "" {
			return errors.New("Device has no defined address")
		}

		baseURL = fmt.Sprintf("http://%s/%s", d.Address, defaultControlURL)
	} else {
		baseURL = d.ControlURL
	}

	url := fmt.Sprintf(baseURL+path, args...)
	if url == "" {
		return errors.New("missing url")
	}

	//fmt.Println("URL: " + url)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-encoded")
	}

	fmt.Println(req)
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
