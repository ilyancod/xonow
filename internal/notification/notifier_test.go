package notification

import (
	"testing"
	config "xonow/internal/config"
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

func TestNewNotifierSettings(t *testing.T) {
	t.Run("empty conf", func(t *testing.T) {
		got := NewNotifierSettings(config.GetConfig())
		if !got.Global.Empty() && len(got.Servers) != 0 {
			t.Errorf("got no empty NotifierSettints, want empty")
		}
	})
	t.Run("conf with values", func(t *testing.T) {
		want := NotifierSettings{
			Global:  notifConfig1,
			Servers: nil,
		}
		got := NewNotifierSettings(&conf)
		assertStruct(t, got, want)
	})
}
