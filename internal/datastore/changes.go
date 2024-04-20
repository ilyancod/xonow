package datastore

import (
	"reflect"
)

type ServerChanges map[ServerAddr]ServerProperties
type ServerProperties map[PropertyName]interface{}
type PlayersChanges struct {
	Added   Players
	Removed Players
}

func getServerChanges(first, second ServerStore) ServerChanges {
	serverChanges := ServerChanges{}
	for serverAddr, firstPayload := range first {
		secondPayload := second[serverAddr]
		if serverAddr != secondPayload.Address {
			continue
		}
		properties := getServerPropertiesChanges(firstPayload, secondPayload)

		if len(properties) != 0 {
			serverChanges[serverAddr] = properties
		}
	}
	return serverChanges
}

func getServerPropertiesChanges(first, second ServerPayload) ServerProperties {
	changes := ServerProperties{}
	t := reflect.TypeOf(first)
	for i := 0; i < t.NumField(); i++ {
		propertyName := PropertyName(t.Field(i).Name)
		firstValue := reflect.ValueOf(first).Field(i).Interface()
		secondValue := reflect.ValueOf(second).Field(i).Interface()

		changedProperty, found := getAnyPropertyChanges(propertyName, firstValue, secondValue)
		if found {
			changes[propertyName] = changedProperty
		}
	}
	return changes
}

func getAnyPropertyChanges(propertyName PropertyName, firstValue, secondValue any) (any, bool) {
	if propertyName == "Players" {
		playersChanges := getPlayersChanges(firstValue.(Players), secondValue.(Players))
		if !playersChanges.Empty() {
			return playersChanges, true
		}
		return nil, false
	}
	if firstValue != secondValue {
		return secondValue, true
	}
	return nil, false
}

func getPlayersChanges(first, second Players) PlayersChanges {
	changes := PlayersChanges{Players{}, Players{}}
	for _, player := range first {
		if !second.ContainsName(player.Name) {
			changes.Removed = append(changes.Removed, player)
		}
	}
	for _, player := range second {
		if !first.ContainsName(player.Name) {
			changes.Added = append(changes.Added, player)
		}
	}
	return changes
}

func (pc PlayersChanges) Empty() bool {
	return len(pc.Added) == 0 && len(pc.Removed) == 0
}
