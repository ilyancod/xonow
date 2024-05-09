package datastore

import (
	"encoding/json"
)

type ServerStore map[IpAddr]ServerPayload

func (ss *ServerStore) Add(address IpAddr, payload ServerPayload) {
	(*ss)[address] = payload
}

func (ss *ServerStore) Remove(address IpAddr) {
	delete(*ss, address)
}

type ServerPayload struct {
	Address IpAddr
	Name    string
	Map     string
	Ping    int
	Bots    int
	Players Players
}

func (d ServerPayload) String() string {
	jsonBytes, _ := json.MarshalIndent(d, "", "    ")
	return string(jsonBytes)
}
