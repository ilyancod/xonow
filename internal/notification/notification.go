package notification

import (
	"xonow/internal/config"
	data "xonow/internal/datastore"
)

type NotifierSettings struct {
	Global  config.Notifications
	Servers map[data.ServerAddr]config.Notifications
}

func NewNotifierSettings(conf *config.Store) NotifierSettings {
	return NotifierSettings{
		Global: conf.Global.Notifications,
	}
}
