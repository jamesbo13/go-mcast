package mcast

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// DistClient is a client within a distribution group
type DistClient struct {
	IPAddr   string `json:"ip_address"`
	DataType string `json:"data_type"`
}

// Distribution is a MusicCast distribution group
type Distribution struct {
	GroupID      string `json:"group_id"`
	GroupName    string `json:"group_name"`
	Role         string
	Status       string
	ServerZone   string       `json:"server_zone"`
	Clients      []DistClient `json:"client_list"`
	AudioDropout bool         `json:"audio_dropout"`
	// Unsupported fields: build_disable, mc_surround
}

// Distribution returns information about any distribution (Link group) the
// device belongs to
func (d Device) Distribution() (Distribution, error) {

	var dist Distribution

	err := d.GetRequest("dist/getDistributionInfo", &dist)
	if err != nil {
		return Distribution{}, err
	}

	return dist, nil
}

type Group struct {
	ID      string
	Server  Zone
	Clients []Zone
}

// TODO - handle special case of all zeros

// Generate a random 32 digit hex string
func newID() (string, error) {
	bytes := make([]byte, 16) // 2 output chars per byte
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func NewGroup(server Zone, clients ...Zone) (Group, error) {

	var g Group
	var err error

	// Generate new ID
	g.ID, err = newID()
	if err != nil {
		return g, err
	}

	g.Clients = make([]Zone, 0, len(clients))
	clientAddrs := make([]string, 0, len(clients))

	// Ensure server is not part of existing Group
	distInfo, err := server.dev.Distribution()
	if err != nil {
		return g, err
	}

	if distInfo.GroupID != "" && distInfo.GroupID != "00000000000000000000000000000000" {
		//	return g, errors.New("Server already part of distribution group")
	}

	// Ensure clients are not part of existing Group, if so leave group
	for _, c := range clients {
		fmt.Printf("Checking status of client %s\n", c.dev.Address)

		distInfo, err = c.dev.Distribution()
		if err != nil {
			return g, err
		}

		if distInfo.GroupID != "" && distInfo.GroupID != "00000000000000000000000000000000" {
			// TODO - leave group
			fmt.Printf("device %s (%s) part of group %s. Leaving group to join new group.\n", c.dev.Name, c.dev.Address, distInfo.GroupID)
			//err = c.LeaveGroup()
			if err != nil {
				return g, err
			}
		}

		// Add client to new group
		err = c.setClientInfo(g.ID)
		if err != nil {
			return g, err
		}

		g.Clients = append(g.Clients, c)
		clientAddrs = append(clientAddrs, c.dev.Address)
	}

	// Add server to group
	err = server.setServerInfo(g.ID, "add", clientAddrs...)
	if err != nil {
		return g, err
	}
	g.Server = server

	// Start distributing
	err = server.dev.GetRequest("dist/startDistribution?num=%d", nil, 1)

	return g, err
}

/*
func (g Group) AddClient(d Device) error {

	err := g.addClient(d)
	if err != nil {
		return err
	}

	return g.addServer(g.Server, d)

}
*/

func (z Zone) setClientInfo(id string) error {

	var jsonStr = `{"group_id":"` + id + `","zone":"` + z.name + `"}`

	return z.dev.PostRequest("dist/setClientInfo", strings.NewReader(jsonStr))
}

func (g Group) addServer(d Device, clients ...Device) error {

	fmt.Println("addServer")

	if g.ID == "" {
		return errors.New("Missing group ID")
	}

	clientIPs := make([]string, 0, len(clients))
	for _, c := range clients {
		clientIPs = append(clientIPs, c.Address)
	}

	req := struct {
		ID      string   `json:"group_id"`
		Zone    string   `json:"zone"`
		Type    string   `json:"type"`
		Clients []string `json:"client_list"`
	}{
		g.ID, "main", "add", clientIPs,
	}

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return d.PostRequest("dist/setServerInfo", bytes.NewReader(jsonBytes))
}

func (g Group) DeleteClient(c Zone) error {

	err := c.setClientInfo("")
	if err != nil {
		return err
	}

	err = g.Server.setServerInfo(g.ID, "remove", c.dev.Address)
	if err == nil {
		return err
	}

	// Remove client from g.Clients
	newSlice := make([]Zone, 0, len(g.Clients))
	for _, c := range g.Clients {
		if c.dev.Address == c.dev.Address && c.name == c.name {
			continue
		}

		newSlice = append(newSlice, c)
	}
	g.Clients = newSlice

	return nil
}

func (z Zone) setServerInfo(id, reqType string, clients ...string) error {

	req := struct {
		ID      string   `json:"group_id"`
		Zone    string   `json:"zone"`
		Type    string   `json:"type"`
		Clients []string `json:"client_list"`
	}{
		id, z.name, reqType, clients,
	}

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return z.dev.PostRequest("dist/setServerInfo", bytes.NewReader(jsonBytes))
}

func (g Group) Delete() error {

	// Remove all clients from the group
	for _, c := range g.Clients {

		err := c.setClientInfo("")
		if err != nil {
			return err
		}
	}

	g.Clients = nil

	return g.Server.setServerInfo("", "remove")
}
