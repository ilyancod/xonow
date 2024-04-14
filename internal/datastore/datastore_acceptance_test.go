package datastore_test

import (
	"reflect"
	"testing"
	data "xonow/internal/datastore"
)

func TestUpdateServerData(t *testing.T) {
	cases := []struct {
		name    string
		servers []data.ServerPayload
		input   data.ServerStore
		want    data.ServerChanges
	}{
		{
			name:    "empty DataStore",
			servers: []data.ServerPayload{},
			input: data.ServerStore{
				"149.202.87.185:26010": data.ServerPayload{
					Address: data.ServerAddr("149.202.87.185:26010"),
					Name:    "test_name",
					Map:     "test_map",
					Ping:    50,
				},
			},
			want: data.ServerChanges{},
		},
		{
			name: "DataStore with empty payload",
			servers: []data.ServerPayload{
				{
					Address: "149.202.87.185:26010",
				},
			},
			input: data.ServerStore{
				"149.202.87.185:26010": data.ServerPayload{
					Address: data.ServerAddr("149.202.87.185:26010"),
					Name:    "test_name",
					Map:     "test_map",
					Ping:    50,
				},
			},
			want: data.ServerChanges{
				"149.202.87.185:26010": data.ServerProperties{
					"Name": "test_name",
					"Map":  "test_map",
					"Ping": 50,
				},
			},
		},
		{
			name: "empty DataStore",
			servers: []data.ServerPayload{
				{
					Address: data.ServerAddr("149.202.87.185:26010"),
					Name:    "test_name1",
					Map:     "test_map1",
					Ping:    45,
				},
			},
			input: data.ServerStore{
				"149.202.87.185:26010": data.ServerPayload{
					Address: data.ServerAddr("149.202.87.185:26010"),
					Name:    "test_name2",
					Map:     "test_map2",
					Ping:    50,
				},
			},
			want: data.ServerChanges{
				"149.202.87.185:26010": data.ServerProperties{
					"Name": "test_name2",
					"Map":  "test_map2",
					"Ping": 50,
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			dataStore := data.GetDataStore()
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
