package datastore

import (
	"fmt"
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

func TestGoqstatToDataServers(t *testing.T) {
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
	invalidServer := goqstat_server1
	invalidServer.Numplayers = 2

	validServer := goqstat_server1
	validServer.Numplayers = 2
	validServer.Players = []goqstat.Player{player1, player2}

	cases := []struct {
		name   string
		server goqstat.Server
		want   bool
	}{
		{
			name:   "empty server",
			server: goqstat_server1,
			want:   true,
		},
		{
			name:   "invalid server players",
			server: invalidServer,
			want:   false,
		},
		{
			name:   "valid server players",
			server: validServer,
			want:   true,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := checkServerPlayersValid(test.server)
			assertBool(t, got, test.want)
		})
	}
}

func TestGetBotsFromString(t *testing.T) {
	assertEmptyError := func(t testing.TB, got error) {
		t.Helper()
		if got == nil {
			t.Errorf("expected error, but got no one")
		}
	}

	cases := []struct {
		name     string
		bots     string
		wantErr  error
		wantBots int
	}{
		{
			name:     "empty string",
			bots:     "",
			wantErr:  fmt.Errorf("error"),
			wantBots: 0,
		},
		{
			name:     "invalid bots string 1",
			bots:     "1s",
			wantErr:  fmt.Errorf("error"),
			wantBots: 0,
		},
		{
			name:     "invalid bots string 2",
			bots:     "1.0",
			wantErr:  fmt.Errorf("error"),
			wantBots: 0,
		},
		{
			name:     "valid string",
			bots:     "1",
			wantErr:  nil,
			wantBots: 1,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := getBotsFromString(test.bots)
			if test.wantErr != nil {
				assertEmptyError(t, err)
			}
			assertStrings(t, strconv.Itoa(got), strconv.Itoa(test.wantBots))
		})
	}
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
