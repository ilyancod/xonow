package main_test

import (
	"github.com/ilyancod/goqstat"
	"sort"
	"strings"
	"testing"
	"testing/fstest"
	"xonow/internal/config"
	"xonow/internal/datastore"
	data "xonow/internal/datastore"
	"xonow/internal/notification"
	"xonow/internal/utils"
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
	Result    string
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

type StubFormatter struct{}

func (sf StubFormatter) FormatTitle(payload data.ServerPayload) string {
	return string(payload.Address)
}

func (sf StubFormatter) FormatMessage(changes notification.NotifyServerChanges) string {
	result := ""
	keys := make([]string, 0, len(changes))
	for key := range changes {
		keys = append(keys, string(key))
	}
	// map keys has a random order
	sort.Strings(keys)
	for _, configName := range keys {
		configValue := changes[notification.ConfigName(configName)]
		result += string(configName) + " " + strings.Join(configValue, " ") + "\n"
	}

	return result
}

func TestNotification(t *testing.T) {
	filesystem := fstest.MapFS{"config.json": {Data: []byte(validConfig)}}
	cases := []struct {
		name          string
		configName    string
		dataStore     *datastore.DataStore
		newData       []goqstat.Server
		want          []NotifyMessageResult
		notifyDesktop bool
	}{
		{
			name:       "valid notification",
			configName: "config.json",
			dataStore:  datastore.GetDataStoreSingleInstance(),
			newData:    []goqstat.Server{goqstat_server1},
			want: []NotifyMessageResult{
				{
					Title: "149.202.87.185:26010",
					Message: `maps_appear snooker
players_appear test_user1 test_user2
`,
				},
			},
			notifyDesktop: true,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			conf := config.GetConfigSingleInstance()
			conf.Clear()
			err := conf.ReadFromFile(filesystem, test.configName)
			if err != nil {
				t.Fatal(err)
			}
			test.dataStore.Clear()
			for serverAddress := range conf.Servers {
				test.dataStore.AddServer(datastore.IpAddr(serverAddress), datastore.ServerPayload{})
			}
			notificationSettings := notification.NewNotifierSettings(conf)

			serverData := datastore.GoqstatToDataServers(&test.newData)
			serverChanges := test.dataStore.UpdateServerData(serverData)
			notifyChanges := notification.NewNotifyChanges(serverChanges, notificationSettings)
			stubNotifier := StubNotifier{}
			stubFormatter := StubFormatter{}
			notifyChanges.Emit(&stubNotifier, stubFormatter)

			assertNotifyMessageResult(t, stubNotifier.Notifiers, test.want)

			if test.notifyDesktop {
				notifyDesktop := &notification.NotifyDesktop{
					IconPath: getIconPath(t),
				}
				formatter := notification.HTMLFormater{}
				notifyChanges.Emit(notifyDesktop, formatter)
			}
		})
	}
}

func assertNotifyMessageResult(t testing.TB, got, want []NotifyMessageResult) {
	t.Helper()
	for _, messageWant := range want {
		found := false
		for _, messageGot := range got {
			if messageGot == messageWant {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("not found %#v\nin %#v", messageWant, got)
		}
	}
}

func getIconPath(t testing.TB) string {
	iconPath, err := utils.GetIconPath()
	if err != nil {
		t.Fatal("fail to get icon path:", err)
	}
	return iconPath
}
