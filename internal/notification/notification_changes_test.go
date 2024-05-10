package notification

import (
	"github.com/ilyancod/goqstat"
	"reflect"
	"testing"
	"xonow/internal/config"
	. "xonow/internal/datastore"
)

var (
	userAppear1    = goqstat.Player{Name: "user_appear1", Ping: 10}
	userAppear2    = goqstat.Player{Name: "user_appear2", Ping: 20}
	userDisappear1 = goqstat.Player{Name: "user_disappear1", Ping: 30}
	userDisappear2 = goqstat.Player{Name: "user_disappear2", Ping: 40}
	userMismatch1  = goqstat.Player{Name: "user_mismatch1", Ping: 40}
	userMismatch2  = goqstat.Player{Name: "user_mismatch2", Ping: 40}
)

func TestNewNotifyChanges(t *testing.T) {
	notifierSettings := NewNotifierSettings(config.GetConfigSingleInstance())
	t.Run("empty changes", func(t *testing.T) {
		got := NewNotifyChanges(ServerChanges{}, notifierSettings)
		assertLen(t, len(got), 0)
		got = NewNotifyChanges(ServerChanges{
			"address": ServerProperties{},
		}, notifierSettings)
		assertLen(t, len(got), 0)
	})
}

func TestNewNotifyServerChanges(t *testing.T) {
	var (
		playersChanges = PlayersChanges{
			Added:   []goqstat.Player{userAppear1, userAppear2},
			Removed: []goqstat.Player{userDisappear1, userDisappear2},
			Count:   PlayersCountChanges{Was: 2, Become: 2},
		}
		playersChangesMismatch = PlayersChanges{
			Added:   []goqstat.Player{userMismatch1},
			Removed: []goqstat.Player{userMismatch2},
			Count:   PlayersCountChanges{Was: 2, Become: 2},
		}
		playersAppearInEmptyServer = PlayersChanges{
			Added:   []goqstat.Player{userMismatch1},
			Removed: []goqstat.Player{},
			Count:   PlayersCountChanges{Was: 0, Become: 1},
		}
	)
	var (
		propertiesMismatch = ServerProperties{
			"Map":     "test_map",
			"Players": playersChangesMismatch,
			"Ping":    50,
			"Bots":    3,
			"Name":    "test_name",
		}
		propertiesMapAndPing = ServerProperties{
			"Map":  "mars",
			"Ping": 50,
		}
		propertiesPlayers = ServerProperties{
			"Players": playersChanges,
			"Ping":    50,
		}
		propertiesPlayersAndMap = ServerProperties{
			"Map":     "mars",
			"Players": playersChanges,
			"Ping":    50,
		}
		propertiesPlayersInEmptyServer = ServerProperties{
			"Players": playersAppearInEmptyServer,
			"Ping":    50,
		}
	)
	cases := []struct {
		name       string
		properties ServerProperties
		config     config.Notifications
		want       NotifyServerChanges
	}{
		{
			name:       "empty ServerProperties and Notify config",
			properties: ServerProperties{},
			config:     config.Notifications{},
			want:       NotifyServerChanges{},
		},
		{
			name:       "empty Notify config",
			properties: propertiesMapAndPing,
			config:     config.Notifications{},
			want:       NotifyServerChanges{},
		},
		{
			name:       "mismatch map",
			properties: propertiesMismatch,
			config:     config.Notifications{},
			want:       NotifyServerChanges{},
		},
		{
			name:       "maps appear",
			properties: propertiesMapAndPing,
			config:     notifConfig,
			want: NotifyServerChanges{
				"maps_appear": []string{"mars"},
			},
		},
		{
			name:       "players appear and disappear",
			properties: propertiesPlayers,
			config:     notifConfig,
			want: NotifyServerChanges{
				"players_appear":    []string{"user_appear1", "user_appear2"},
				"players_disappear": []string{"user_disappear1", "user_disappear2"},
			},
		},
		{
			name:       "maps and players changed",
			properties: propertiesPlayersAndMap,
			config:     notifConfig,
			want: NotifyServerChanges{
				"maps_appear":       []string{"mars"},
				"players_appear":    []string{"user_appear1", "user_appear2"},
				"players_disappear": []string{"user_disappear1", "user_disappear2"},
			},
		},
		{
			name:       "any player appear in empty server",
			properties: propertiesPlayersInEmptyServer,
			config:     notifConfig,
			want: NotifyServerChanges{
				"any_player_appear_in_empty_server": []string{"user_mismatch1"},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := newNotifyServerChanges(test.properties, test.config)
			assertStruct(t, got, test.want)
		})
	}
}

func assertLen(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %v len, want %v", got, want)
	}
}

func assertStruct(t testing.TB, got, want any) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v\n got %#v", want, got)
	}
}
