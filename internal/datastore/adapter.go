package datastore

import (
	"strconv"

	"github.com/ilyancod/goqstat"
)

func GoqstatToDataServers(new *[]goqstat.Server) ServerStore {
	newDataMap := make(ServerStore)
	for _, server := range *new {
		if !checkServerPlayersValid(server) {
			continue
		}
		data := serverToData(server)
		newDataMap[data.Address] = data
	}
	return newDataMap
}

func serverToData(server goqstat.Server) ServerPayload {
	numBots, err := getBotsFromRules(server.Rules)
	if err != nil {
		numBots = 0
	}
	dataMap := ServerPayload{
		Address: ServerAddr(server.Address),
		Name:    server.Name,
		Map:     server.Map,
		Ping:    server.Ping,
		Bots:    numBots,
		Players: server.Players,
	}
	return dataMap
}

func checkServerPlayersValid(server goqstat.Server) bool {
	return server.Numplayers == len(server.Players)
}

func getBotsFromRules(rules goqstat.Rules) (int, error) {
	return strconv.Atoi(rules.Bots)
}
