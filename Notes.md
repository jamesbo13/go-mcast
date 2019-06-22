# Yamaha MusicCast API Implementation Notes #

## Events ##

Send a GET request to http://<host>/YamahaExtendedControl/v1/main/getStatus
Include headers:
    X-AppName: MusicCast/1   (arbitrary text)
    X-AppPort: 9999          (local port number to listen on)
Will receive UDP JSON responses on dedicated port
Per notes from PHP library it will timeout in 10 minutes so need to resubscribe

For testing:
 curl -H "X-AppName: MusicCast/1" -H "X-AppPort: 9999" http://10.13.1.106/YamahaExtendedControl/v1/main/getStatus

Then in another window:
nc  -u -l 9999

## Device Discovery ##

Uses SSDP.

Search for devices of type 'MediaRenderer' via SSDP.
Decode Location header in response to get URL for XML info
Fetch and inspect XML

    - verify manufacturer
    - decode base URL - yamaha:X_URLBase
    - decode yxc control URL - yamaha:X_yxcControlURL
      (the above is path added to end of base URL)

Icon paths in XML are relative to the XML root path not the base URL

### SSDP ###

Notify responses contain fields:

    - Host (IP:port)
    - Location (URL to desc.xml)
    - Cache-Control (max-age setting ... appears to default to 30s)
    - Server
    - NT
    - NTS
    - USN

USN is same as UDN read from desc.xml with ::upnp:rootdevice appended. Should treat USN/UDN
as the cannonical identifier of the device as other values (eg IP address, name) can change.

## Remote Control API ##

See
https://github.com/PSeitz/yamaha-nodejs/blob/master/simpleCommands.js
http://www.sdu.se/pub/yamaha/yamaha-ynca-receivers-protocol.pdf
https://openremote.github.io/archive-dotorg/forums/Openremote%20controlling%20Yamaha%20RX-V475.html

## FCC Teardown / reports ##

https://fccid.io/A6RNW01A

## Network snooping ##

- Disable network manager

`sudo service NetworkManager stop`

- Setup bridge

```sh  
sudo ip addr del 10.13.1.40/24 dev enp42s0
sudo brctl addbr br0
sudo brctl addif br0 enp42s0
sudo brctl addif br0 enx3c8cf8ff93df
sudo ip addr add 10.13.1.40/24 dev br0
sudo ip link set br0 up
sudo route add default gw 10.13.1.1
```