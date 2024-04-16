package notification

import (
	"github.com/ilyancod/goqstat"
	"reflect"
	"testing"
	config "xonow/internal/config"
	data "xonow/internal/datastore"
)

var conf = config.Store{
	Global:  config.Global{Notifications: notifConfig1},
	Servers: nil,
}

var notifConfig1 = config.Notifications{
	MapsAppear:                   []string{"mars", "snooker", "cofrag"},
	PlayersAppear:                []string{"user_appear1", "user_appear2"},
	PlayersDisappear:             []string{"user_disappear1", "user_disappear2"},
	AnyPlayerAppearInEmptyServer: true,
}

var (
	user_appear1    = goqstat.Player{Name: "user_appear1", Ping: 10}
	user_appear2    = goqstat.Player{Name: "user_appear2", Ping: 20}
	user_disappear1 = goqstat.Player{Name: "user_disappear1", Ping: 30}
	user_disappear2 = goqstat.Player{Name: "user_disappear2", Ping: 40}
)

func TestConfigToNotifierSettings(t *testing.T) {
	//t.Run("empty conf", func(t *testing.T) {
	//	want := config.Store{}
	//	got := configToNotifierSettings(want)
	//	if !got.Global.Empty() && len(got.Servers) != 0 {
	//		t.Errorf("got no empty NotifierSettints, want empty")
	//	}
	//})
	t.Run("conf with values", func(t *testing.T) {
		want := NotifierSettings{
			Global:  notifConfig1,
			Servers: nil,
		}
		got := configToNotifierSettings(conf.Global)
		assertStruct(t, got, want)
	})
}

func TestGetNotifyResult(t *testing.T) {
	notifierSettings := configToNotifierSettings(conf.Global)
	t.Run("empty changes", func(t *testing.T) {
		got := getNotifyResult(data.ServerChanges{}, notifierSettings)
		assertLen(t, len(got), 0)
		got = getNotifyResult(data.ServerChanges{
			"address": data.ServerProperties{},
		}, notifierSettings)
		assertLen(t, len(got), 0)
	})
}

func TestGetNotifyValue(t *testing.T) {
	changesProperties := data.ServerProperties{
		"Map":  "test",
		"Ping": 51,
	}
	playersChanges := data.PlayersChanges{
		Added:   []goqstat.Player{user_appear1, user_appear2},
		Removed: []goqstat.Player{user_disappear1, user_disappear2},
	}
	t.Run("empty notify value expected", func(t *testing.T) {
		want := NotifyValue{}
		got := getNotifyValue(data.ServerProperties{}, notifConfig1)
		assertStruct(t, got, want)

		got = getNotifyValue(changesProperties, notifConfig1)
		assertStruct(t, got, want)
	})
	t.Run("maps appear", func(t *testing.T) {
		changes := data.ServerProperties{
			"Map":  "mars",
			"Ping": 51,
		}
		want := NotifyValue{
			"maps_appear": []string{"mars"},
		}
		got := getNotifyValue(changes, notifConfig1)
		assertStruct(t, got, want)
	})
	t.Run("players appear and disappear", func(t *testing.T) {
		changes := changesProperties
		changes["Players"] = playersChanges
		want := NotifyValue{
			"players_appear":    []string{"user_appear1", "user_appear2"},
			"players_disappear": []string{"user_disappear1", "user_disappear2"},
		}
		got := getNotifyValue(changes, notifConfig1)
		assertStruct(t, got, want)
	})
	t.Run("maps and players changed", func(t *testing.T) {
		changes := changesProperties
		changes["Map"] = "mars"
		changes["Players"] = playersChanges
		want := NotifyValue{
			"maps_appear":       []string{"mars"},
			"players_appear":    []string{"user_appear1", "user_appear2"},
			"players_disappear": []string{"user_disappear1", "user_disappear2"},
		}
		got := getNotifyValue(changes, notifConfig1)
		assertStruct(t, got, want)
	})

	// t.Run("any player appear in empty server", func(t *testing.T) {
	// 	got := getNotifyValue(data.DataProperties{})
	// 	assertLen(t, len(got), 0)
	// })
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