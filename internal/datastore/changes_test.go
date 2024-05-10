package datastore

import (
	"reflect"
	"testing"

	. "github.com/ilyancod/goqstat"
)

var serverPayload1 = ServerPayload{
	Address: IpAddr(goqstat_server1.Address),
	Name:    goqstat_server1.Name,
	Map:     goqstat_server1.Map,
	Ping:    goqstat_server1.Ping,
}

var serverPayload2 = ServerPayload{
	Address: IpAddr(goqstat_server2.Address),
	Name:    goqstat_server2.Name,
	Map:     goqstat_server2.Map,
	Ping:    goqstat_server2.Ping,
}

var (
	player1 = Player{Name: "player1", Ping: 10}
	player2 = Player{Name: "player2", Ping: 20}
	player3 = Player{Name: "player3", Ping: 30}
	player4 = Player{Name: "player4", Ping: 40}
)

func TestGetServerChanges(t *testing.T) {
	serverPayloadChanged := serverPayload1
	serverPayloadChanged.Map = "loo"

	serverStore1 := ServerStore{
		serverPayload1.Address: serverPayload1,
		serverPayload2.Address: serverPayload2,
	}
	serverStore2 := ServerStore{
		serverPayload1.Address: serverPayloadChanged,
		serverPayload2.Address: serverPayload2,
	}
	serverStoreEmptyPayload := ServerStore{
		serverPayload1.Address: ServerPayload{},
	}

	cases := []struct {
		name         string
		firstServer  ServerStore
		secondServer ServerStore
		wantLength   int
	}{
		{
			name:         "empty ServerStore",
			firstServer:  ServerStore{},
			secondServer: ServerStore{},
			wantLength:   0,
		},
		{
			name:         "ServerStore with empty payload",
			firstServer:  serverStoreEmptyPayload,
			secondServer: serverStore1,
			wantLength:   1,
		},
		{
			name:         "equal ServerStore",
			firstServer:  serverStore1,
			secondServer: serverStore1,
			wantLength:   0,
		},
		{
			name:         "no equal ServerStore",
			firstServer:  serverStore1,
			secondServer: serverStore2,
			wantLength:   1,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := getServerChanges(test.firstServer, test.secondServer)
			assertLen(t, len(got), test.wantLength)
		})
	}
}

func TestGetServerPropertiesChanges(t *testing.T) {
	firstData := ServerPayload{
		Address: IpAddr(goqstat_server1.Address),
		Name:    goqstat_server1.Name,
		Map:     goqstat_server1.Map,
		Ping:    goqstat_server1.Ping,
		Bots:    0,
		Players: Players{},
	}

	cases := []struct {
		name   string
		first  ServerPayload
		second ServerPayload
		want   ServerProperties
	}{
		{
			name:   "equal ServerPayload",
			first:  firstData,
			second: firstData,
			want:   ServerProperties{},
		},
		{
			name:  "no equal ServerPayload (map, ping, bots)",
			first: firstData,
			second: ServerPayload{
				Address: IpAddr(goqstat_server1.Address),
				Name:    goqstat_server1.Name,
				Map:     "snooker",
				Ping:    73,
				Bots:    3,
			},
			want: ServerProperties{
				"Map":  "snooker",
				"Ping": 73,
				"Bots": 3,
			},
		},
		{
			name: "no equal ServerPayload (players)",
			first: ServerPayload{
				Address: IpAddr(goqstat_server1.Address),
				Name:    goqstat_server1.Name,
				Players: Players{player1, player2, player3},
			},
			second: ServerPayload{
				Address: IpAddr(goqstat_server1.Address),
				Name:    goqstat_server1.Name,
				Players: Players{player4, player2},
			},
			want: ServerProperties{
				"Players": PlayersChanges{
					Added:   Players{player4},
					Removed: Players{player1, player3},
					Count:   PlayersCountChanges{3, 2},
				},
			},
		},
		{
			name: "no equal ServerPayload (server address)",
			first: ServerPayload{
				Address: "",
			},
			second: ServerPayload{
				Address: IpAddr(goqstat_server2.Address),
			},
			want: ServerProperties{
				"Address": IpAddr(goqstat_server2.Address),
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := getServerPropertiesChanges(test.first, test.second)
			assertDeepEqual(t, got, test.want)
		})
	}
}

func TestGetChangesPlayers(t *testing.T) {
	cases := []struct {
		name   string
		first  Players
		second Players
		want   PlayersChanges
	}{
		{
			name:   "empty Players",
			first:  Players{},
			second: Players{},
			want:   PlayersChanges{Players{}, Players{}, PlayersCountChanges{}},
		},
		{
			name:   "equal Players",
			first:  Players{player1, player2},
			second: Players{player1, player2},
			want:   PlayersChanges{Players{}, Players{}, PlayersCountChanges{2, 2}},
		},
		{
			name:   "changed Players",
			first:  Players{player1, player2},
			second: Players{player2, player3, player4},
			want: PlayersChanges{
				Added:   Players{player3, player4},
				Removed: Players{player1},
				Count:   PlayersCountChanges{2, 3},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := getPlayersChanges(test.first, test.second)
			assertDeepEqual(t, got, test.want)
		})
	}
}

func assertDeepEqual(t testing.TB, got, want any) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v\nwant %#v", got, want)
	}
}
