package main_test

import (
	"github.com/ilyancod/goqstat"
	"testing"
	"testing/fstest"
	"xonow/internal/config"
	"xonow/internal/datastore"
	"xonow/internal/notification"
)

const (
	validConfig = `{
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

var goqstat_server1 = goqstat.Server{
	Protocol:      "xonotics",
	Address:       "149.202.87.185:26010",
	Status:        "online",
	Hostname:      "149.202.87.185:26010",
	Name:          "[E] TheRegulars B6  Instagib Server [git]",
	Gametype:      "Xonotic",
	Map:           "snooker",
	Numplayers:    2,
	Maxplayers:    48,
	Numspectators: 0,
	Maxspectators: 0,
	Ping:          53,
	Retries:       0,
	Rules:         goqstat.Rules{Bots: "0"},
	Players:       datastore.Players{player1, player2},
}

var (
	player1 = goqstat.Player{Name: "test_user1", Ping: 10}
	player2 = goqstat.Player{Name: "test_user2", Ping: 20}
)

type StubNotifier struct {
	Notifiers []NotifyMessageResult
}

type NotifyMessageResult struct {
	Title   string
	Message string
}

func (sn *StubNotifier) Notify(title, message string) error {
	sn.Notifiers = append(sn.Notifiers, NotifyMessageResult{title, message})
	return nil
}

func TestNotification(t *testing.T) {
	conf := config.GetConfig()
	filesystem := fstest.MapFS{"config.json": {Data: []byte(validConfig)}}
	err := conf.ReadFromFile(filesystem, "config.json")
	if err != nil {
		t.Fatal(err)
	}
	store := datastore.GetDataStore()
	for serverAddress := range conf.Servers {
		store.AddServer(datastore.ServerAddr(serverAddress), datastore.ServerPayload{})
	}
	notificationSettings := notification.NewNotifierSettings(conf)
	goqstatData := []goqstat.Server{goqstat_server1}

	serverData := datastore.GoqstatToDataServers(&goqstatData)
	serverChanges := store.UpdateServerData(serverData)
	notifyChanges := notification.NewNotifyChanges(serverChanges, notificationSettings)
	stubNotifier := StubNotifier{}
	notifyChanges.Emit(&stubNotifier)

	want := []NotifyMessageResult{
		{
			Title:   "Xonow: notifyChanges on the server 149.202.87.185:26010",
			Message: "Map appeared: snooker",
		},
		{
			Title:   "Xonow: notifyChanges on the server 149.202.87.185:26010",
			Message: "Players appeared: test_user1, test_user2",
		},
	}

	assertNotifyMessageResult(t, stubNotifier.Notifiers, want)
}

func assertNotifyMessageResult(t testing.TB, got, want []NotifyMessageResult) {
	for _, messageWant := range want {
		found := false
		for _, messageGot := range got {
			if messageGot == messageWant {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("not found %#v in %#v", messageWant, got)
		}
	}
}
