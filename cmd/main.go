package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"xonow/internal/config"
	"xonow/internal/datastore"
	"xonow/internal/notification"
	"xonow/internal/utils"

	"github.com/ilyancod/goqstat"
)

func main() {
	conf := config.GetConfigSingleInstance()

	configDir, err := utils.GetConfigDir()
	if err != nil {
		fmt.Println("fail to get config directory:", err)
		return
	}
	fileSystem := os.DirFS(configDir)

	err = conf.ReadFromFile(fileSystem, "config.json")
	if err != nil {
		fmt.Println("error opening the config: ", err)
		err = conf.SaveToFile(filepath.Join(configDir, "config.json"))
		if err != nil {
			fmt.Println("error creating the config: ", err)
		}
		return
	}

	iconPath, err := utils.GetIconPath()
	if err != nil {
		fmt.Println("fail to get icon path:", err)
		return
	}

	store := datastore.GetDataStoreSingleInstance()
	for serverAddress := range conf.Servers {
		store.AddServer(datastore.IpAddr(serverAddress), datastore.ServerPayload{})
	}
	notificationSettings := notification.NewNotifierSettings(conf)
	for {
		goqstatServers, err := GetGoqstatData(conf)
		if err != nil {
			fmt.Println("error getting qstat data:", err)
			time.Sleep(time.Second * 5)
			continue
		}

		serverData := datastore.GoqstatToDataServers(&goqstatServers)
		serverChanges := store.UpdateServerData(serverData)

		notifyChanges := notification.NewNotifyChanges(serverChanges, notificationSettings)

		notifyDesktop := &notification.NotifyDesktop{
			IconPath: iconPath,
		}
		formatter := notification.HTMLFormater{}
		notifyChanges.Emit(notifyDesktop, formatter)
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
