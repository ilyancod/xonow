package notification

import (
	"strings"
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

func RunNotifier(notifyChanges NotifyChanges) {
	for serverAddr, notifyServerChanges := range notifyChanges {
		for configName, configValue := range notifyServerChanges {
			title := "Xonow: notifyChanges on the server " + string(serverAddr)
			switch configName {
			case "maps_appear":
				{
					notify(title, "Map appeared: "+strings.Join(configValue, " "))
				}
			case "players_appear":
				{
					notify(title, "Players appeared: "+strings.Join(configValue, " "))
				}
			case "players_disappear":
				{
					notify(title, "Players disappeared: "+strings.Join(configValue, " "))
				}
			}
		}
	}
}
