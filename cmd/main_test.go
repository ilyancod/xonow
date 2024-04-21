package main_test

import (
	"testing"
	"testing/fstest"
	"xonow/internal/config"
)

const (
	validConfig = `{
  {
    "global": {
        "notifications": {
            "maps_appear": [
                "mars",
                "snooker",
                "cofrag",
                "mmkitch_double"
            ],
            "players_appear": [
                "test_user1",
                "test_user2",
                "test_user3"
            ],
            "players_disappear": [
                "test_user4",
                "test_user2",
                "test_user5"
            ],
            "any_player_appear_in_empty_server": false
        }
    },
    "servers": {
        "149.202.87.185:26010": {},
        "168.119.137.110:26000": {}
    }
}`
)

func TestNotification(t *testing.T) {
	conf := config.GetConfig()
	filesystem := fstest.MapFS{"config.json": {Data: []byte(validConfig)}}
	err := conf.ReadFromFile(filesystem, "config.json")
	if err != nil {
		t.Fatal(err)
	}

	//config.Rea
}
