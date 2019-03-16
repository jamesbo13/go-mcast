package mcast

import (
	"fmt"
	"net/http"
	"time"
)

// Helper functions for accessing device REST API

var defaultClient = &http.Client{
	Timeout: time.Second * 5,
}

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

	return fmt.Sprintf("unexpected response: %d", r)
}
