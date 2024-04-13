package datastore

import (
	"encoding/json"
	"sync"

	"github.com/ilyancod/goqstat"
)

type ServerAddr string
type PropertyName string

type DataStore struct {
	serverData ServerData
}

type ServerData map[ServerAddr]ServerPayload

type ServerPayload struct {
	Address ServerAddr
	Name    string
	Map     string
	Ping    int
	Bots    int
	Players Players
}

type Players []goqstat.Player

var singleDataStore *DataStore
var lock = &sync.Mutex{}

func GetDataStore() *DataStore {
	if singleDataStore == nil {
		lock.Lock()
		defer lock.Unlock()
		singleDataStore = &DataStore{
			serverData: make(ServerData),
		}
	}
	return singleDataStore
}

func (ds *DataStore) UpdateServerData(serverData ServerData) DataChanges {
	changes := getChanges(ds.serverData, serverData)

	for _, data := range serverData {
		ds.serverData[data.Address] = data
	}

	return changes
}

func (d ServerPayload) String() string {
	jsonBytes, _ := json.MarshalIndent(d, "", "    ")
	return string(jsonBytes)
}

func (p Players) ContainsName(name string) bool {
	for _, player := range p {
		if player.Name == name {
			return true
		}
	}
	return false
}
