package datastore

import (
	"sync"

	"github.com/ilyancod/goqstat"
)

type ServerAddr string
type PropertyName string

type dataStore struct {
	serverData ServerStore
}

type Players []goqstat.Player

var singleDataStore *dataStore
var lock = &sync.Mutex{}

func GetDataStore() *dataStore {
	if singleDataStore == nil {
		lock.Lock()
		defer lock.Unlock()
		singleDataStore = &dataStore{
			serverData: make(ServerStore),
		}
	}
	return singleDataStore
}

func (ds *dataStore) UpdateServerData(serverData ServerStore) ServerChanges {
	changes := getServerChanges(ds.serverData, serverData)

	for _, data := range serverData {
		ds.serverData[data.Address] = data
	}

	return changes
}

func (ds *dataStore) AddServer(address ServerAddr, payload ServerPayload) {
	ds.serverData.Add(address, payload)
}

func (ds *dataStore) RemoveServer(address ServerAddr) {
	ds.serverData.Remove(address)
}

func (p Players) ContainsName(name string) bool {
	for _, player := range p {
		if player.Name == name {
			return true
		}
	}
	return false
}
