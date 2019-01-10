package mcast

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/koron/go-ssdp"
)

const (
	upnpDevName = "urn:schemas-upnp-org:device:MediaRenderer:1"
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

func deviceInfo(url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP error: %s", resp.Status)
	}

	xmlDoc, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var xmlDev device

	err = xml.Unmarshal(xmlDoc, &xmlDev)
	if err != nil {
		return err
	}

	fmt.Printf("device: %+v\n", xmlDev)

	return nil
}

// Discover returns all MusicCast devices that can be found on local network.
// Uses SSDP (Simple Service Discovery Protocol) to find appropriate hosts
func Discover() error {

	devs, err := ssdp.Search(upnpDevName, 2, "")
	if err != nil {
		return err
	}

	for i, srv := range devs {
		fmt.Printf("%d: %s %s\n", i, srv.Type, srv.Location)
		err = deviceInfo(srv.Location)
		if err != nil {
			return err
		}
	}

	return nil
}
