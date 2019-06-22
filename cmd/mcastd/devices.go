package main

import (
	"encoding/json"
	"net/http"
)

func listDevices(w http.ResponseWriter, r *http.Request) {

	var buf []byte
	var err error

	path := r.URL.Path

	rootLen := len("/devices/")
	if len(path) <= rootLen {
		buf, err = json.MarshalIndent(devices, "", "  ")
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		id := path[rootLen:]
		for _, d := range devices {
			if d.SerialNum == id {
				buf, err = json.MarshalIndent(d, "", "  ")
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte(err.Error()))
					return
				}
			}
		}
	}
	w.Write(buf)
}
