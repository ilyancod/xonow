package notification

import (
	"github.com/ilyancod/goqstat"
	"reflect"
	"testing"
	"xonow/internal/config"
	data "xonow/internal/datastore"
)

var (
	userAppear1    = goqstat.Player{Name: "user_appear1", Ping: 10}
	userAppear2    = goqstat.Player{Name: "user_appear2", Ping: 20}
	userDisappear1 = goqstat.Player{Name: "user_disappear1", Ping: 30}
	userDisappear2 = goqstat.Player{Name: "user_disappear2", Ping: 40}
)

func TestNewNotifyChanges(t *testing.T) {
	notifierSettings := NewNotifierSettings(config.GetConfig())
	t.Run("empty changes", func(t *testing.T) {
		got := NewNotifyChanges(data.ServerChanges{}, notifierSettings)
		assertLen(t, len(got), 0)
		got = NewNotifyChanges(data.ServerChanges{
			"address": data.ServerProperties{},
		}, notifierSettings)
		assertLen(t, len(got), 0)
	})
}

func TestNewNotifyServerChanges(t *testing.T) {
	changesProperties := data.ServerProperties{
		"Map":  "test",
		"Ping": 51,
	}
	playersChanges := data.PlayersChanges{
		Added:   []goqstat.Player{userAppear1, userAppear2},
		Removed: []goqstat.Player{userDisappear1, userDisappear2},
	}
	t.Run("empty notify value expected", func(t *testing.T) {
		want := NotifyServerChanges{}
		got := newNotifyServerChanges(data.ServerProperties{}, notifConfig1)
		assertStruct(t, got, want)

		got = newNotifyServerChanges(changesProperties, notifConfig1)
		assertStruct(t, got, want)
	})
	t.Run("maps appear", func(t *testing.T) {
		changes := data.ServerProperties{
			"Map":  "mars",
			"Ping": 51,
		}
		want := NotifyServerChanges{
			"maps_appear": []string{"mars"},
		}
		got := newNotifyServerChanges(changes, notifConfig1)
		assertStruct(t, got, want)
	})
	t.Run("players appear and disappear", func(t *testing.T) {
		changes := changesProperties
		changes["Players"] = playersChanges
		want := NotifyServerChanges{
			"players_appear":    []string{"user_appear1", "user_appear2"},
			"players_disappear": []string{"user_disappear1", "user_disappear2"},
		}
		got := newNotifyServerChanges(changes, notifConfig1)
		assertStruct(t, got, want)
	})
	t.Run("maps and players changed", func(t *testing.T) {
		changes := changesProperties
		changes["Map"] = "mars"
		changes["Players"] = playersChanges
		want := NotifyServerChanges{
			"maps_appear":       []string{"mars"},
			"players_appear":    []string{"user_appear1", "user_appear2"},
			"players_disappear": []string{"user_disappear1", "user_disappear2"},
		}
		got := newNotifyServerChanges(changes, notifConfig1)
		assertStruct(t, got, want)
	})

	// t.Run("any player appear in empty server", func(t *testing.T) {
	// 	got := newNotifyServerChanges(data.DataProperties{})
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
