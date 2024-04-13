package main // Xonotic Now

import (
	"fmt"
	"strconv"
	"time"
	"xonow/internal/datastore"

	"github.com/ilyancod/goqstat"
)

func main() {
	ReadConfig()
	store := datastore.GetDataStore()
	for {
		servers := []string{}
		for server := range ConfigData.Servers {
			servers = append(servers, server)
		}
		result, err := goqstat.GetXonotics(servers...)
		if err != nil {
			fmt.Println(err)
		}

		serverData := datastore.GoqstatToDataServers(&result)
		dataChanges := store.UpdateServerData(serverData)

		RunNotifier(dataChanges)
		time.Sleep(time.Second * 5)
	}
}

func PrintCurrentData(data []goqstat.Server) {
	for index, server := range data {
		numBots, _ := strconv.Atoi(server.Rules.Bots)
		fmt.Println("Server name:\t", server.Name)
		fmt.Println("Server map:\t", server.Map)
		fmt.Println("Server players:\t", server.Numplayers-numBots)
		for index, player := range server.Players {
			fmt.Println("\tPlayer name: ", player.Name)
			fmt.Println("\tPlayer team: ", player.Team)
			fmt.Println("\tPlayer score: ", player.Score)
			if index != len(server.Players)-1 {
				fmt.Println("\t- - - - - - - - - - - - - - -")
			}
		}
		if index != len(data)-1 {
			fmt.Println("===========================================")
		}
	}
}
