package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Store struct {
	Global  Global                 `json:"global,omitempty"`
	Servers map[string]interface{} `json:"servers,omitempty"`
}

type Global struct {
	Notifications Notifications `json:"notifications,omitempty"`
}

type Server struct {
	Notifications Notifications `json:"notifications,omitempty"`
}

type Notifications struct {
	MapsAppear                   []string `json:"maps_appear,omitempty"`
	PlayersAppear                []string `json:"players_appear,omitempty"`
	PlayersDisappear             []string `json:"players_disappear,omitempty"`
	AnyPlayerAppearInEmptyServer bool     `json:"any_player_appear_in_empty_server,omitempty"`
}

var singleConfig *Store
var lock = &sync.Mutex{}

func GetConfig() *Store {
	if singleConfig == nil {
		lock.Lock()
		defer lock.Unlock()
		singleConfig = &Store{
			Servers: make(map[string]interface{}),
		}
	}
	return singleConfig
}

func (cs *Store) ReadFromFile(filePath string) error {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return err
	}

	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()

	if err = dec.Decode(&singleConfig); err != nil {
		return err
	}
	return nil
}

func (n Notifications) Empty() bool {
	return len(n.MapsAppear) == 0 && len(n.PlayersAppear) == 0 &&
		len(n.PlayersDisappear) == 0 && !n.AnyPlayerAppearInEmptyServer
}
