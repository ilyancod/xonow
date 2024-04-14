package datastore

import (
	"testing"

	"github.com/ilyancod/goqstat"
)

var serverPayload1 = ServerPayload{
	Address: ServerAddr(goqstat_server1.Address),
	Name:    goqstat_server1.Name,
	Map:     goqstat_server1.Map,
	Ping:    goqstat_server1.Ping,
}

var serverPayload2 = ServerPayload{
	Address: ServerAddr(goqstat_server2.Address),
	Name:    goqstat_server2.Name,
	Map:     goqstat_server2.Map,
	Ping:    goqstat_server2.Ping,
}

var (
	player1 = goqstat.Player{Name: "player1", Ping: 10}
	player2 = goqstat.Player{Name: "player2", Ping: 20}
	player3 = goqstat.Player{Name: "player3", Ping: 30}
	player4 = goqstat.Player{Name: "player4", Ping: 40}
)

func TestGetChanges(t *testing.T) {
	t.Run("equal data", func(t *testing.T) {
		first := ServerStore{
			serverPayload1.Address: serverPayload1,
			serverPayload2.Address: serverPayload2,
		}
		second := first
		got := getChanges(first, second)
		assertLen(t, len(got), 0)
	})
	t.Run("no equal data", func(t *testing.T) {
		data := serverPayload1
		first := ServerStore{data.Address: data}
		data.Map = "loo"
		second := ServerStore{data.Address: data}
		got := getChanges(first, second)
		assertLen(t, len(got), 1)
	})
	t.Run("empty data", func(t *testing.T) {
		first := ServerStore{serverPayload1.Address: ServerPayload{}}
		second := ServerStore{serverPayload1.Address: serverPayload1}

		got := getChanges(first, second)
		assertLen(t, len(got), 1)
	})
}

func TestGetChangesData(t *testing.T) {
	firstData := ServerPayload{
		Address: ServerAddr(goqstat_server1.Address),
		Name:    goqstat_server1.Name,
		Map:     goqstat_server1.Map,
		Ping:    goqstat_server1.Ping,
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

	t.Run("no equal Data (server)", func(t *testing.T) {
		firstData.Address = ""
		secondData.Address = ServerAddr(goqstat_server2.Address)
		got := getChangesData(firstData, secondData)

		changesAddress, found := got["Address"]

		assertData(t, found, changesAddress, secondData.Address)
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

func assertData(t testing.TB, found bool, got, want any) {
	t.Helper()
	if !found || got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertPlayers(t testing.TB, found bool, got, want PlayersChanges) {
	t.Helper()
	if !found || !want.Equal(got) {
		t.Errorf("got %#v\nwant %#v", got, want)
	}
}
