package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type config struct {
	Port int
}

// parseConfig will read JSON encoded configuration file
func parseConfig(r io.Reader) (config, error) {
	var c config

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(buf, &c)
	return c, err
}
