package datastore_test

import (
	"reflect"
	"testing"
	. "xonow/internal/datastore"
)

func TestUpdateServerData(t *testing.T) {
	cases := []struct {
		name    string
		servers []ServerPayload
		input   ServerStore
		want    ServerChanges
	}{
		{
			name:    "empty DataStore",
			servers: []ServerPayload{},
			input: ServerStore{
				"149.202.87.185:26010": ServerPayload{
					Address: IpAddr("149.202.87.185:26010"),
					Name:    "test_name",
					Map:     "test_map",
					Ping:    50,
				},
			},
			want: ServerChanges{},
		},
		{
			name: "DataStore with empty payload",
			servers: []ServerPayload{
				{
					Address: "149.202.87.185:26010",
				},
			},
			input: ServerStore{
				"149.202.87.185:26010": ServerPayload{
					Address: IpAddr("149.202.87.185:26010"),
					Name:    "test_name",
					Map:     "test_map",
					Ping:    50,
				},
			},
			want: ServerChanges{
				"149.202.87.185:26010": ServerProperties{
					"Name": "test_name",
					"Map":  "test_map",
					"Ping": 50,
				},
			},
		},
		{
			name: "empty DataStore",
			servers: []ServerPayload{
				{
					Address: IpAddr("149.202.87.185:26010"),
					Name:    "test_name1",
					Map:     "test_map1",
					Ping:    45,
				},
			},
			input: ServerStore{
				"149.202.87.185:26010": ServerPayload{
					Address: IpAddr("149.202.87.185:26010"),
					Name:    "test_name2",
					Map:     "test_map2",
					Ping:    50,
				},
			},
			want: ServerChanges{
				"149.202.87.185:26010": ServerProperties{
					"Name": "test_name2",
					"Map":  "test_map2",
					"Ping": 50,
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			dataStore := GetDataStoreSingleInstance()
			for _, server := range test.servers {
				dataStore.AddServer(server.Address, server)
			}
			got := dataStore.UpdateServerData(test.input)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("got %#v, want %#v", got, test.want)
			}
		})
	}
}
