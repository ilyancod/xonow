package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"xonow/internal/config"
	"xonow/internal/datastore"
	"xonow/internal/notification"

	"github.com/ilyancod/goqstat"
)

func main() {
	conf := config.GetConfig()
	fileSystem := os.DirFS("../config/")

	err := conf.ReadFromFile(fileSystem, "config.json")
	if err != nil {
		fmt.Println("error opening the config: ", err)
		err = conf.SaveToFile("../config/config.json")
		if err != nil {
			fmt.Println("error creating the config: ", err)
		}
		return
	}

	store := datastore.GetDataStore()
	for serverAddress := range conf.Servers {
		store.AddServer(datastore.ServerAddr(serverAddress), datastore.ServerPayload{})
	}
	notificationSettings := notification.NewNotifierSettings(conf)
	for {

		goqstatData, err := GetGoqstatData(conf)
		if err != nil {
			fmt.Println(err)
			continue
		}

		serverData := datastore.GoqstatToDataServers(&goqstatData)
		dataChanges := store.UpdateServerData(serverData)

		notifyChanges := notification.NewNotifyChanges(dataChanges, notificationSettings)

		notifyDesktop := &notification.NotifyDesktop{
			IconPath: "assets/xonotic.png",
		}
		notifyChanges.Notify(notifyDesktop)
		time.Sleep(time.Second * 5)
	}
}

func GetGoqstatData(conf *config.Store) ([]goqstat.Server, error) {
	servers := make([]string, 0, len(conf.Servers))
	for server := range conf.Servers {
		servers = append(servers, server)
	}
	return goqstat.GetXonotics(servers...)
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
