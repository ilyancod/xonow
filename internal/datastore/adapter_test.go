package datastore

import (
	"strconv"
	"testing"

	"github.com/ilyancod/goqstat"
)

var goqstat_server1 = goqstat.Server{
	Protocol:      "xonotics",
	Address:       "149.202.87.185:26010",
	Status:        "online",
	Hostname:      "149.202.87.185:26010",
	Name:          "[E] TheRegulars B6  Instagib Server [git]",
	Gametype:      "Xonotic",
	Map:           "zastavka",
	Numplayers:    0,
	Maxplayers:    48,
	Numspectators: 0,
	Maxspectators: 0,
	Ping:          53,
	Retries:       0,
	Rules:         goqstat.Rules{},
}

var goqstat_server2 = goqstat.Server{
	Protocol:      "xonotics",
	Address:       "168.119.137.110:26000",
	Status:        "online",
	Hostname:      "168.119.137.110:26000",
	Name:          "[E] TheRegulars B6  Instagib Server [git]",
	Gametype:      "Xonotic",
	Map:           "centertest02",
	Numplayers:    0,
	Maxplayers:    26,
	Numspectators: 0,
	Maxspectators: 0,
	Ping:          43,
	Retries:       0,
	Rules:         goqstat.Rules{},
}

func TestServersToDataMap(t *testing.T) {
	servers := []goqstat.Server{goqstat_server1, goqstat_server2}
	t.Run("valid Servers", func(t *testing.T) {
		dataMap := GoqstatToDataServers(&servers)
		for _, server := range servers {
			data, found := dataMap[ServerAddr(server.Address)]
			if !found {
				t.Errorf("got nothing, but want %v", server.Address)
			}
			if data.Address == serverPayload1.Address {
				assertStrings(t, data.String(), serverPayload1.String())
			}
			if data.Address == serverPayload2.Address {
				assertStrings(t, data.String(), serverPayload2.String())
			}
		}
	})
}

func TestServerToData(t *testing.T) {
	server := goqstat_server1
	t.Run("invalid Rules", func(t *testing.T) {
		server.Rules.Bots = "1q"
		want := serverPayload1
		got := serverToData(server)
		assertStrings(t, got.String(), want.String())
	})
	t.Run("valid server data", func(t *testing.T) {
		server.Rules.Bots = "2"
		want := serverPayload1
		want.Bots = 2
		got := serverToData(server)
		assertStrings(t, got.String(), want.String())
	})
}

func TestCheckServerPlayersValid(t *testing.T) {
	t.Run("empty server", func(t *testing.T) {
		got := checkServerPlayersValid(goqstat_server1)
		assertBool(t, got, true)
	})
	t.Run("invalid server players", func(t *testing.T) {
		server := goqstat_server1
		server.Numplayers = 2
		got := checkServerPlayersValid(server)
		assertBool(t, got, false)
	})
	t.Run("valid server players", func(t *testing.T) {
		server := goqstat_server1
		server.Numplayers = 2
		server.Players = []goqstat.Player{player1, player2}
		got := checkServerPlayersValid(server)
		assertBool(t, got, true)
	})
}

func TestGetBotsFromRules(t *testing.T) {
	assertEmptyError := func(t testing.TB, got error) {
		t.Helper()
		if got == nil {
			t.Errorf("expected error, but got no one")
		}
	}
	rules := goqstat.Rules{}
	t.Run("empty Rules", func(t *testing.T) {
		_, err := getBotsFromRules(rules)
		assertEmptyError(t, err)
	})
	t.Run("invalid Rules", func(t *testing.T) {
		rules.Bots = "1s"
		_, err := getBotsFromRules(rules)
		assertEmptyError(t, err)
		rules.Bots = "1.0"
		_, err = getBotsFromRules(rules)
		assertEmptyError(t, err)
	})
	t.Run("valid Rules", func(t *testing.T) {
		rules.Bots = "1"
		want := 1
		got, _ := getBotsFromRules(rules)
		assertStrings(t, strconv.Itoa(got), strconv.Itoa(want))
	})
}
