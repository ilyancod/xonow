package main

import (
	"fmt"
	"reflect"
	"strings"
	"xonow/data"
)

// type NotifyResult struct {
// 	ServerAddr string
// 	Type       string
// 	Value      []string
// }

const (
	ErrInterfaceNoString         = NotifierErr("expected String value type")
	ErrInterfaceNoPlayersChanges = NotifierErr("expected PlayersChanges value type")
)

type NotifierErr string

func (e NotifierErr) Error() string {
	return string(e)
}

var settings NotifierSettings

type NotifyResult map[data.ServerAddr]NotifyValue
type NotifyValue map[ConfigName]ConfigValue

type ConfigName string
type ConfigValue []string

type NotifierSettings struct {
	Global  Notifications
	Servers map[data.ServerAddr]Notifications
}

func SetNotifierSettings(config Config) {
	settings = configToNotifierSettings(config)
	// settings.Global = config.Global.Notifications

	// serverMap := map[data.ServerAddr]Notifications{}
	// for serverAddr, server := range config.Servers {
	// 	server[serverAddr] = server.Notifications
	// }
	// notifierSettings = n
}

func configToNotifierSettings(config Config) NotifierSettings {
	return NotifierSettings{
		Global: config.Global.Notifications,
	}
}

func RunNotifier(changes data.DataChanges) {
	notifyResults := getNotifyResult(changes, settings)
	for serverAddr, notifyValue := range notifyResults {
		for configName, configValue := range notifyValue {
			title := "Xonow: changes on the server " + string(serverAddr)
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
	fmt.Println("notifyResult: ", notifyResults)

}

func getNotifyResult(changes data.DataChanges, ns NotifierSettings) NotifyResult {
	notifyResult := NotifyResult{}
	for serverAddr, properties := range changes {
		notifyValue := getNotifyValue(properties, ns.Global)
		if len(notifyValue) != 0 {
			notifyResult[serverAddr] = notifyValue
		}
	}
	return notifyResult
}

func getNotifyValue(properties data.DataProperties, notification Notifications) NotifyValue {
	result := NotifyValue{}
	for name, value := range properties {
		switch name {
		case "Map":
			mapValue, found := getMapsAppear(value, notification.MapsAppear)
			if found {
				result["maps_appear"] = mapValue
			}
		case "Players":
			playersChanges, err := interfaceToPlayersChanges(value)
			if err != nil {
				fmt.Println(err)
			}

			playersAppear, found := getPlayersByNames(playersChanges.Added, notification.PlayersAppear)
			if found {
				result["players_appear"] = playersAppear
			}

			playersDisappear, found := getPlayersByNames(playersChanges.Removed, notification.PlayersDisappear)
			if found {
				result["players_disappear"] = playersDisappear
			}
		}
	}

	return result
}

func getMapsAppear(value any, mapsAppear []string) (result []string, found bool) {
	mapStr, err := interfaceToString(value)
	if err != nil {
		fmt.Println(err)
		return
	}
	if contains(mapsAppear, mapStr) {
		return []string{mapStr}, true
	}
	return
}

func getPlayersByNames(players data.Players, playerNames []string) (result []string, found bool) {
	found = false
	result = []string{}

	for _, playerName := range playerNames {
		if players.ContainsName(playerName) {
			result = append(result, playerName)
		}
	}
	if len(result) > 0 {
		found = true
	}
	return
}

func interfaceToString(value any) (string, error) {
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() != reflect.String {
		return "", ErrInterfaceNoString
	}
	return reflectValue.String(), nil
}

func interfaceToPlayersChanges(value any) (data.PlayersChanges, error) {
	if playersChanges, ok := value.(data.PlayersChanges); ok {
		return playersChanges, nil
	} else {
		return playersChanges, ErrInterfaceNoPlayersChanges
	}
}

func contains(array []string, target string) bool {
	for _, str := range array {
		if str == target {
			return true
		}
	}
	return false
}
