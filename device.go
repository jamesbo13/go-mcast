package mcast

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
