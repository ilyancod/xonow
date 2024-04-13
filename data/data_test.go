package data

import (
	"strconv"
	"testing"

	"github.com/ilyancod/goqstat"
)

var server1 = goqstat.Server{
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

var server2 = goqstat.Server{
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

var data1 = Data{
	Address: ServerAddr(server1.Address),
	Name:    server1.Name,
	Map:     server1.Map,
	Ping:    server1.Ping,
}

var data2 = Data{
	Address: ServerAddr(server2.Address),
	Name:    server2.Name,
	Map:     server2.Map,
	Ping:    server2.Ping,
}

var (
	player1 = goqstat.Player{Name: "player1", Ping: 10}
	player2 = goqstat.Player{Name: "player2", Ping: 20}
	player3 = goqstat.Player{Name: "player3", Ping: 30}
	player4 = goqstat.Player{Name: "player4", Ping: 40}
)

func TestCheckServerPlayersValid(t *testing.T) {
	t.Run("empty server", func(t *testing.T) {
		got := checkServerPlayersValid(server1)
		assertBool(t, got, true)
	})
	t.Run("invalid server players", func(t *testing.T) {
		server := server1
		server.Numplayers = 2
		got := checkServerPlayersValid(server)
		assertBool(t, got, false)
	})
	t.Run("valid server players", func(t *testing.T) {
		server := server1
		server.Numplayers = 2
		server.Players = []goqstat.Player{player1, player2}
		got := checkServerPlayersValid(server)
		assertBool(t, got, true)
	})
}

func TestServerToData(t *testing.T) {
	server := server1
	t.Run("invalid Rules", func(t *testing.T) {
		server.Rules.Bots = "1q"
		want := data1
		got := serverToData(server)
		assertStrings(t, got.String(), want.String())
	})
	t.Run("valid server data", func(t *testing.T) {
		server.Rules.Bots = "2"
		want := data1
		want.Bots = 2
		got := serverToData(server)
		assertStrings(t, got.String(), want.String())
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

func TestServersToDataMap(t *testing.T) {
	servers := []goqstat.Server{server1, server2}
	t.Run("valid Servers", func(t *testing.T) {
		dataMap := serversToDataMap(&servers)
		for _, server := range servers {
			data, found := dataMap[ServerAddr(server.Address)]
			if !found {
				t.Errorf("got nothing, but want %v", server.Address)
			}
			if data.Address == data1.Address {
				assertStrings(t, data.String(), data1.String())
			}
			if data.Address == data2.Address {
				assertStrings(t, data.String(), data2.String())
			}
		}
	})
}

func TestGetChanges(t *testing.T) {
	t.Run("equal data", func(t *testing.T) {
		first := DataMap{
			data1.Address: data1,
			data2.Address: data2,
		}
		second := first
		got := getChanges(first, second)
		assertLen(t, len(got), 0)
	})
	t.Run("no equal data", func(t *testing.T) {
		data := data1
		first := DataMap{data.Address: data}
		data.Map = "loo"
		second := DataMap{data.Address: data}
		got := getChanges(first, second)
		assertLen(t, len(got), 1)
	})
}

func TestGetChangesData(t *testing.T) {
	firstData := Data{
		Address: ServerAddr(server1.Address),
		Name:    server1.Name,
		Map:     server1.Map,
		Ping:    server1.Ping,
		Bots:    0,
		Players: Players{},
	}
	secondData := firstData
	t.Run("equal Data", func(t *testing.T) {
		got := getChangesData(firstData, secondData)
		assertLen(t, len(got), 0)
	})
	t.Run("no equal Data (map, ping, bots)", func(t *testing.T) {
		assertData := func(t testing.TB, found bool, got, want any) {
			t.Helper()
			if !found || got != want {
				t.Errorf("got %v want %v", got, want)
			}
		}
		secondData := firstData
		secondData.Map = "snooker"
		secondData.Ping = 73
		secondData.Bots = 3
		got := getChangesData(firstData, secondData)

		assertLen(t, len(got), 3)

		changesMap, found := got["Map"]
		assertData(t, found, changesMap, secondData.Map)

		changesPing, found := got["Ping"]
		assertData(t, found, changesPing, secondData.Ping)

		changesBots, found := got["Bots"]
		assertData(t, found, changesBots, secondData.Bots)
	})

	t.Run("no equal Data (players)", func(t *testing.T) {
		firstData.Players = Players{player1, player2, player3}
		secondData.Players = Players{player4, player2}
		want := PlayersChanges{
			Added:   Players{player4},
			Removed: Players{player1, player3},
		}
		got := getChangesData(firstData, secondData)

		assertLen(t, len(got), 1)

		changesPlayers, found := got["Players"]
		assertPlayers(t, found, changesPlayers.(PlayersChanges), want)
	})
}

func TestGetChangesPlayers(t *testing.T) {
	t.Run("empty players", func(t *testing.T) {
		firstPlayers := Players{}
		secondPlayers := Players{}
		want := PlayersChanges{Players{}, Players{}}
		got := getChangesPlayers(firstPlayers, secondPlayers)
		assertPlayers(t, true, got, want)
	})
	t.Run("equal players", func(t *testing.T) {
		firstPlayers := Players{player1, player2}
		secondPlayers := firstPlayers
		want := PlayersChanges{Players{}, Players{}}
		got := getChangesPlayers(firstPlayers, secondPlayers)
		assertPlayers(t, true, got, want)
	})
	t.Run("changed players", func(t *testing.T) {
		firstPlayers := Players{player1, player2}
		secondPlayers := Players{player2, player3, player4}
		want := PlayersChanges{
			Added:   Players{player3, player4},
			Removed: Players{player1},
		}
		got := getChangesPlayers(firstPlayers, secondPlayers)
		assertPlayers(t, true, got, want)
	})
}

func assertBool(t testing.TB, got, want bool) {
	t.Helper()

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertStrings(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertLen(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %v len, want %v", got, want)
	}
}

func assertPlayers(t testing.TB, found bool, got, want PlayersChanges) {
	t.Helper()
	if !found || !want.Equal(got) {
		t.Errorf("got %#v\nwant %#v", got, want)
	}
}
