package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var ConfigData Config

type Config struct {
	Global  Global                 `json:"global,omitempty"`
	Servers map[string]interface{} `json:"servers,omitempty"`
}

type Global struct {
	Notifications Notifications `json:"notifications,omitempty"`
}

type Notifications struct {
	MapsAppear                   []string `json:"maps_appear,omitempty"`
	PlayersAppear                []string `json:"players_appear,omitempty"`
	PlayersDisappear             []string `json:"players_disappear,omitempty"`
	AnyPlayerAppearInEmptyServer bool     `json:"any_player_appear_in_empty_server,omitempty"`
}

func (n Notifications) Empty() bool {
	return len(n.MapsAppear) == 0 && len(n.PlayersAppear) == 0 &&
		len(n.PlayersDisappear) == 0 && !n.AnyPlayerAppearInEmptyServer
}

func ReadConfig() {
	file, err := os.Open("/home/ilya/projects/go/src/github.com/ilyancod/xonow/config.json")
	if err != nil {
		fmt.Println("Ошибка чтения конфига:", err)
		return
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()

	if err = dec.Decode(&ConfigData); err != nil {
		fmt.Println("Ошибка декодирования конфига:", err)
		return
	}

	// fmt.Println("config:", Config)
}
