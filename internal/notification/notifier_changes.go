package notification

import (
	"fmt"
	"reflect"
	"xonow/internal/config"
	data "xonow/internal/datastore"
)

const (
	ErrInterfaceNoString         = NotifierErr("expected String value type")
	ErrInterfaceNoPlayersChanges = NotifierErr("expected PlayersChanges value type")
)

type NotifierErr string

func (e NotifierErr) Error() string {
	return string(e)
}

type NotifyChanges map[data.ServerAddr]NotifyServerChanges
type NotifyServerChanges map[ConfigName]ConfigValue

type ConfigName string
type ConfigValue []string

func NewNotifyChanges(serverChanges data.ServerChanges, settings NotifierSettings) NotifyChanges {
	notifyChanges := NotifyChanges{}
	for serverAddr, properties := range serverChanges {
		notifyValue := newNotifyServerChanges(properties, settings.Global)
		if len(notifyValue) != 0 {
			notifyChanges[serverAddr] = notifyValue
		}
	}
	return notifyChanges
}

func newNotifyServerChanges(properties data.ServerProperties, notification config.Notifications) NotifyServerChanges {
	result := NotifyServerChanges{}
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
