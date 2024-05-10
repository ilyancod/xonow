package notification

import (
	"xonow/internal/config"
	. "xonow/internal/datastore"
)

type NotifierSettings struct {
	Global  config.Notifications
	Servers map[IpAddr]config.Notifications
}

func NewNotifierSettings(conf *config.Store) NotifierSettings {
	return NotifierSettings{
		Global: conf.Global.Notifications,
	}
}
