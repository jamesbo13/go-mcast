package mcast_test

import (
	"net/http"
	"net/http/httptest"
)

var testServer *httptest.Server

func init() {
	// Start a local server using recorded JSON responses in testdata dir
	hdlr := http.StripPrefix("/YamahaExtendedControl/v1/", http.FileServer(http.Dir("testdata/tsr7850")))
	testServer = httptest.NewServer(hdlr)
}
