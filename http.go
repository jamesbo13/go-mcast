package mcast

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

// Helper functions for accessing device REST API

// responseError maps errors returned in JSON response_code field to defined messages
type responseError uint8

var respErrorMsg = map[uint8]string{
	0:  "success",
	1:  "initializing",
	2:  "internal error",
	3:  "invalid request",
	4:  "invalid parameter",
	5:  "guarded",
	6:  "time out",
	99: "firmware updating",

	// >= 100 are streaming service errors
	100: "access error",
	101: "other error",
	102: "wrong user name",
	103: "wrong password",
	104: "account expired",
	105: "account disconnected / shut down",
	106: "account number limit",
	107: "server maintenance",
	108: "invalid account",
	109: "license error",
	110: "read-only mode",
	111: "max stations",
	112: "access denied",
}

// fulfill the error interface definition
func (r responseError) Error() string {

	if msg, ok := respErrorMsg[uint8(r)]; ok {
		return msg
	}
	// TODO: map error strings to values
	return fmt.Sprintf("unexpected response: %d", r)
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

	// We expect to be passed a struct with an int ResponseCode field
	// but verify so we don't panic during reflect calls
	v := reflect.ValueOf(val).Elem()
	if v.Kind() == reflect.Struct {
		code := v.FieldByName("ResponseCode")
		if code.IsValid() && code.Kind() == reflect.Int {
			// Invalid value means field not found

			if code.Int() != 0 {
				return responseError(code.Int())
			}

		} else {
			return errors.New("missing response_code field in JSON response")
		}
	} else {
		panic("unexpected unmarshal value")
	}

	return nil

}
