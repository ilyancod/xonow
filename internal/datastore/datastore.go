package datastore

import (
	"sync"

	"github.com/ilyancod/goqstat"
)

type IpAddr string

type DataStore struct {
	serverData ServerStore
}

type Players []goqstat.Player

func (p Players) GetNames() (result []string) {
	for _, player := range p {
		result = append(result, player.Name)
	}
	return result
}

var (
	singleDataStore *DataStore
	once            sync.Once
)

func GetDataStoreSingleInstance() *DataStore {
	once.Do(func() {
		singleDataStore = &DataStore{
			serverData: make(ServerStore),
		}
	})
	return singleDataStore
}

func (ds *DataStore) UpdateServerData(serverData ServerStore) ServerChanges {
	changes := getServerChanges(ds.serverData, serverData)

	for _, data := range serverData {
		ds.serverData[data.Address] = data
	}

	return changes
}

func (ds *DataStore) AddServer(address IpAddr, payload ServerPayload) {
	ds.serverData.Add(address, payload)
}

func (ds *DataStore) RemoveServer(address IpAddr) {
	ds.serverData.Remove(address)
}

func (ds *DataStore) GetServer(address IpAddr) (payload ServerPayload, found bool) {
	payload, found = ds.serverData[address]
	return
}

func (ds *DataStore) Clear() {
	ds.serverData = make(ServerStore)
}

func (p Players) ContainsName(name string) bool {
	for _, player := range p {
		if player.Name == name {
			return true
		}
	}
	return false
}
