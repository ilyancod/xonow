package datastore

import (
	"encoding/json"
)

type ServerStore map[ServerAddr]ServerPayload

func (ss *ServerStore) Add(address ServerAddr, payload ServerPayload) {
	(*ss)[address] = payload
}

func (ss *ServerStore) Remove(address ServerAddr) {
	delete(*ss, address)
}

type ServerPayload struct {
	Address ServerAddr
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
