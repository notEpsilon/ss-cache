package config

import (
	"encoding/json"
	"os"
)

type Address struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Shard struct {
	Idx            int     `json:"idx"`
	Name           string  `json:"name"`
	ExposedAddress Address `json:"exposedAddress"`
}

type Config struct {
	Shards []Shard `json:"shards"`
}

func ParseConfig(configPath string) (Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
