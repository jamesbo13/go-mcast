package mcast

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	ssdp "github.com/koron/go-ssdp"
)

const (
	yamahaMfr   = "Yamaha Corporation"
	upnpDevName = "urn:schemas-upnp-org:device:MediaRenderer:1"

	yamahaExtendedControlService = "urn:schemas-yamaha-com:service:X_YamahaExtendedControl:1"
)

// Internal types for decoding XML description file

type yamahaService struct {
	SpecType      string `xml:"X_specType"`
	ControlURL    string `xml:"X_controlURL"`
	UnitDescURL   string `xml:"X_unitDescURL"`
	YXCControlURL string `xml:"X_yxcControlURL"`
	YXCVersion    string `xml:"X_yxcVersion"`
}

type device struct {
	//XMLName  xml.Name `xml:"root"`
	//YamahaNS string   `xml:"yamaha,attr,omitempty"`

	// SpecVersion
	MajorVers uint `xml:"specVersion>major"`
	MinorVers uint `xml:"specVersion>minor"`

	// Device
	Type            string `xml:"device>deviceType"`
	Name            string `xml:"device>friendlyName"`
	Manufacturer    string `xml:"device>manufacturer"`
	ManufacturerURL string `xml:"device>manufacturerURL"`
	ModelDesc       string `xml:"device>modelDescription"`
	ModelName       string `xml:"device>modelName"`
	ModelNumber     string `xml:"device>modelNumber"`
	ModelURL        string `xml:"device>modelURL"`
	UDN             string `xml:"device>UDN"`
	SerialNum       string `xml:"device>serialNumber"`

	// TODO iconList>icon
	// TODO serviceList>service

	// yamaha:X_device
	YamahaBaseURL  string          `xml:"urn:schemas-yamaha-com:device-1-0 X_device>X_URLBase"`
	YamahaServices []yamahaService `xml:"urn:schemas-yamaha-com:device-1-0 X_device>X_serviceList>X_service"`
}

func deviceInfo(urlStr string) (Device, error) {

	var d Device
	var xmlDevice device

	// TODO: use client with timeout
	resp, err := http.Get(urlStr)
	if err != nil {
		return d, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return d, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	xmlDoc, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return d, err
	}

	err = xml.Unmarshal(xmlDoc, &xmlDevice)
	if err != nil {
		return d, err
	}

	// Per YXC documentation, after searching for MediaRenderer
	// verify the following from devices description:
	//	1. manufacturer = "Yamaha Corporation"
	//	2. yamaha:X_device tag exists

	if xmlDevice.Manufacturer != yamahaMfr || xmlDevice.YamahaBaseURL == "" {
		return d, fmt.Errorf("unsupported device")
	}

	// Generate an exportable Device struct from internal device type
	d = Device{
		Name:         xmlDevice.Name,
		Manufacturer: xmlDevice.Manufacturer,
		ModelName:    xmlDevice.ModelName,
		ModelDesc:    xmlDevice.ModelDesc,
		UDN:          xmlDevice.UDN,
		SerialNum:    xmlDevice.SerialNum,
	}

	// Extract host address from base URL
	uri, err := url.Parse(xmlDevice.YamahaBaseURL)
	if err != nil {
		return d, fmt.Errorf("could not parse URL")
	}

	d.Address = uri.Hostname()

	// Find extended control service definition
	for _, s := range xmlDevice.YamahaServices {
		if s.SpecType != yamahaExtendedControlService {
			continue
		}

		if s.YXCControlURL != "" {
			// FIXME - this adds a duplicate / in the middle
			if s.YXCControlURL[0] == '/' {
				d.ControlURL = xmlDevice.YamahaBaseURL + s.YXCControlURL[1:]
			} else {
				d.ControlURL = xmlDevice.YamahaBaseURL + s.YXCControlURL
			}
		}

		d.YXCVersion = strings.TrimSpace(s.YXCVersion)
		break
	}

	return d, nil
}

// Hack to get local IP address.
func getLocalAddr() (string, error) {
	conn, err := net.Dial("udp", "224.0.0.251:5353")
	if err != nil {
		return "", err
	}

	s := conn.LocalAddr().String()

	// strip off :port from address
	return s[:strings.Index(s, ":")], nil
}

// Discover returns all MusicCast devices that can be found on local network.
// Uses SSDP (Simple Service Discovery Protocol) to find appropriate hosts
func Discover() ([]Device, error) {

	var ret []Device

	// TODO: If we pass no address this can fail because other programs are listening
	//       for mcast traffic on same addr/port. Set to local IP Address to get to work.
	localAddr, err := getLocalAddr()
	if err != nil {
		return nil, err
	}

	fmt.Println(localAddr)
	devs, err := ssdp.Search(upnpDevName, 2, localAddr+":0")
	if err != nil {
		return nil, err
	}

	for _, srv := range devs {
		// Fetch desc.xml from srv.Location and parse into device
		d, err := deviceInfo(srv.Location)
		if err != nil {
			continue
		}

		ret = append(ret, d)
	}

	return ret, nil
}

func Ping(addr string) (Device, error) {

	URL := "http://" + addr + ":49154/MediaRenderer/desc.xml"

	return deviceInfo(URL)
}
