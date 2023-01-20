package config_test

import (
	"reflect"
	"testing"

	"github.com/notEpsilon/ss-cache/config"
)

func TestConfig(t *testing.T) {
	expectedConfig := config.Config{
		Shards: []config.Shard{
			{
				Idx:  0,
				Name: "shard-0",
				ExposedAddress: config.Address{
					Host: "127.0.0.1",
					Port: 8080,
				},
			},
			{
				Idx:  1,
				Name: "shard-1",
				ExposedAddress: config.Address{
					Host: "127.0.0.1",
					Port: 8081,
				},
			},
		},
	}

	testConfigPath := "./config_test.json"

	returnedConfig, err := config.ParseConfig(testConfigPath)
	if err != nil {
		t.Errorf("unexpected error config.ParseConfig(%q): %s\n", testConfigPath, err.Error())
	}

	if !reflect.DeepEqual(expectedConfig, returnedConfig) {
		t.Errorf("incorrect configuration parsing, expected: %#v, got %#v\n", expectedConfig, returnedConfig)
	}
}
