package config

import (
	"errors"
	"io/fs"
	"os"
	"reflect"
	"testing"
	"testing/fstest"
)

type StubFailingFS struct {
}

func (s StubFailingFS) Open(name string) (fs.File, error) {
	return &os.File{}, errors.New("oh no, i always fail")
}

var validStore = Store{
	Global: Global{
		Notifications: Notifications{
			AnyPlayerAppearInEmptyServer: true,
		},
	},
	Servers: map[string]any{
		"149.202.87.185:26010": make(map[string]any),
	},
}

const (
	invalidConfig = `fail json format`
	validConfig   = `{
  "global": {
    "notifications": {
      "any_player_appear_in_empty_server": true
    }
  },
  "servers": {
    "149.202.87.185:26010": {}
  }
}`
)

func TestGetConfig(t *testing.T) {
	store1 := GetConfigSingleInstance()
	store1.Servers = map[string]any{
		"149.202.87.185:26010": make(map[string]any),
	}
	store2 := GetConfigSingleInstance()
	assertDeepEqual(t, store1, store2)
}

func TestReadFromFile(t *testing.T) {

	cases := []struct {
		name     string
		fs       fs.FS
		fileName string
		error    error
		want     *Store
	}{
		{
			name:     "failing by open",
			fs:       StubFailingFS{},
			fileName: "config.json",
			error:    errors.New("oh no, i always fail"),
			want:     &Store{Servers: make(map[string]any)},
		},
		{
			name:     "invalid json format of config",
			fs:       fstest.MapFS{"config.json": {Data: []byte(invalidConfig)}},
			fileName: "config.json",
			error:    errors.New("oh no, i always fail"),
			want:     &Store{Servers: make(map[string]any)},
		},
		{
			name:     "valid json format config",
			fs:       fstest.MapFS{"config.json": {Data: []byte(validConfig)}},
			fileName: "config.json",
			error:    nil,
			want:     &validStore,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			config := GetConfigSingleInstance()
			config.Clear()
			err := config.ReadFromFile(test.fs, test.fileName)
			assertError(t, err, test.error)
			assertDeepEqual(t, config, test.want)
		})
	}
}

func BenchmarkStore_ReadFromFile(b *testing.B) {
	config := GetConfigSingleInstance()
	config.Clear()
	b.ResetTimer()

	fileSystem := fstest.MapFS{"config.json": {Data: []byte(validConfig)}}
	for i := 0; i < b.N; i++ {
		config.ReadFromFile(fileSystem, "config.json")
	}
}

func assertDeepEqual(t testing.TB, got, want any) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v\nwant %#v", got, want)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if (got == nil) != (want == nil) {
		t.Errorf("got %#v error, want %#v", got, want)
	}
}
