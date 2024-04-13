package datastore

import (
	"reflect"
)

type DataChanges map[ServerAddr]DataProperties
type DataProperties map[PropertyName]interface{}
type PlayersChanges struct {
	Added   Players
	Removed Players
}

func getChanges(first, second ServerData) DataChanges {
	dataChanges := DataChanges{}
	for _, firstData := range first {
		secondData := second[firstData.Address]
		changes := getChangesData(firstData, secondData)
		if len(changes) != 0 {
			dataChanges[firstData.Address] = changes
		}
	}
	return dataChanges
}

func getChangesData(first, second ServerPayload) DataProperties {
	changes := DataProperties{}
	t := reflect.TypeOf(first)
	for i := 0; i < t.NumField(); i++ {
		propertyName := PropertyName(t.Field(i).Name)
		firstValue := reflect.ValueOf(first).Field(i).Interface()
		secondValue := reflect.ValueOf(second).Field(i).Interface()

		changedProperty, found := getChangesProperty(propertyName, firstValue, secondValue)
		if found {
			changes[propertyName] = changedProperty
		}
	}
	return changes
}

func getChangesProperty(propertyName PropertyName, firstValue, secondValue any) (any, bool) {
	if propertyName == "Players" {
		playersChanges := getChangesPlayers(firstValue.(Players), secondValue.(Players))
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

func getChangesPlayers(first, second Players) PlayersChanges {
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

func (pc PlayersChanges) Equal(other PlayersChanges) bool {
	return reflect.DeepEqual(pc.Added, other.Added) &&
		reflect.DeepEqual(pc.Removed, other.Removed)
}

func (pc PlayersChanges) Empty() bool {
	return len(pc.Added) == 0 && len(pc.Removed) == 0
}
