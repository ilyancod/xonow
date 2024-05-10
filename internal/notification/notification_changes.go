package notification

import (
	"fmt"
	"reflect"
	"strings"
	"xonow/internal/config"
	. "xonow/internal/datastore"
)

const (
	ErrInterfaceNoString         = NotifierErr("expected String value type")
	ErrInterfaceNoPlayersChanges = NotifierErr("expected PlayersChanges value type")
)

type Notifier interface {
	Notify(title, message string) error
}

type Formatter interface {
	FormatTitle(payload ServerPayload) string
	FormatMessage(changes NotifyServerChanges) string
}

type HTMLFormater struct{}

func (hm HTMLFormater) FormatTitle(payload ServerPayload) string {
	return payload.Name
}

func (hm HTMLFormater) FormatMessage(changes NotifyServerChanges) string {
	result := ""
	for configName, configValue := range changes {
		switch configName {
		case "maps_appear":
			{
				result += "Map appeared: <b>"
			}
		case "players_appear":
			{
				result += "Players appeared: <b>"
			}
		case "players_disappear":
			{
				result += "Players disappeared: <b>"
			}
		case "any_player_appear_in_empty_server":
			{
				result += "Players appeared in empty server: <b>"
			}
		default:
			continue
		}
		result += strings.Join(configValue, ", ") + "</b>\n"
	}
	return result
}

type NotifierErr string

func (e NotifierErr) Error() string {
	return string(e)
}

type NotifyChanges map[IpAddr]NotifyServerChanges
type NotifyServerChanges map[ConfigName]ConfigValue

type ConfigName string
type ConfigValue []string

func NewNotifyChanges(serverChanges ServerChanges, settings NotifierSettings) NotifyChanges {
	notifyChanges := NotifyChanges{}
	for serverAddr, properties := range serverChanges {
		notifyValue := newNotifyServerChanges(properties, settings.Global)
		if len(notifyValue) != 0 {
			notifyChanges[serverAddr] = notifyValue
		}
	}
	return notifyChanges
}

func newNotifyServerChanges(properties ServerProperties, notification config.Notifications) NotifyServerChanges {
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

			found = isAnyPlayersInEmptyServer(playersChanges.Count)
			if found && len(playersAppear) == 0 {
				result["any_player_appear_in_empty_server"] = playersChanges.Added.GetNames()
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

func getPlayersByNames(players Players, playerNames []string) (result []string, found bool) {
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

func isAnyPlayersInEmptyServer(playersCount PlayersCountChanges) bool {
	return playersCount.Was == 0 && playersCount.Become != 0
}

func (nc NotifyChanges) Emit(notifier Notifier, formatter Formatter) {
	dataStore := GetDataStoreSingleInstance()
	for serverAddr, notifyServerChanges := range nc {
		serverPayload, found := dataStore.GetServer(serverAddr)
		if !found {
			continue
		}
		title := formatter.FormatTitle(serverPayload)
		message := formatter.FormatMessage(notifyServerChanges)
		if message == "" {
			continue
		}
		err := notifier.Notify(title, message)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func interfaceToString(value any) (string, error) {
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() != reflect.String {
		return "", ErrInterfaceNoString
	}
	return reflectValue.String(), nil
}

func interfaceToPlayersChanges(value any) (PlayersChanges, error) {
	if playersChanges, ok := value.(PlayersChanges); ok {
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
