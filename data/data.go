package data

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/ilyancod/goqstat"
)

type ServerAddr string
type PropertyName string

type Data struct {
	Address ServerAddr
	Name    string
	Map     string
	Ping    int
	Bots    int
	Players Players
}

func (d Data) String() string {
	jsonBytes, _ := json.MarshalIndent(d, "", "    ")
	return string(jsonBytes)
}

type Players []goqstat.Player

func (p Players) ContainsName(name string) bool {
	for _, player := range p {
		if player.Name == name {
			return true
		}
	}
	return false
}

type DataMap map[ServerAddr]Data

type DataChanges map[ServerAddr]DataProperties
type DataProperties map[PropertyName]interface{}
type PlayersChanges struct {
	Added   Players
	Removed Players
}

func (pc PlayersChanges) Equal(other PlayersChanges) bool {
	return reflect.DeepEqual(pc.Added, other.Added) &&
		reflect.DeepEqual(pc.Removed, other.Removed)
}

func (pc PlayersChanges) Empty() bool {
	return len(pc.Added) == 0 && len(pc.Removed) == 0
}

var currentData DataMap

func SetData(new *[]goqstat.Server) (DataChanges, error) {
	if currentData == nil {
		currentData = make(DataMap)
	}
	validNew := []goqstat.Server{}
	for _, server := range *new {
		if !checkServerPlayersValid(server) {
			continue
		}
		validNew = append(validNew, server)
	}

	newData := serversToDataMap(&validNew)
	changes := getChanges(currentData, newData)

	for _, data := range newData {
		currentData[data.Address] = data
	}

	return changes, nil
}

func checkServerPlayersValid(server goqstat.Server) bool {
	return server.Numplayers == len(server.Players)
}

func serversToDataMap(new *[]goqstat.Server) DataMap {
	newDataMap := make(DataMap)
	for _, server := range *new {
		data := serverToData(server)
		newDataMap[data.Address] = data
	}
	return newDataMap
}

func serverToData(server goqstat.Server) Data {
	numBots, err := getBotsFromRules(server.Rules)
	if err != nil {
		numBots = 0
	}
	dataMap := Data{
		Address: ServerAddr(server.Address),
		Name:    server.Name,
		Map:     server.Map,
		Ping:    server.Ping,
		Bots:    numBots,
		Players: server.Players,
	}
	return dataMap
}

func getChanges(first, second DataMap) DataChanges {
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

func getChangesData(first, second Data) DataProperties {
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

func getBotsFromRules(rules goqstat.Rules) (int, error) {
	return strconv.Atoi(rules.Bots)
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
