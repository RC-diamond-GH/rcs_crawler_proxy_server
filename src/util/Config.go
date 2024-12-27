package util

import (
	"encoding/json"
	"fmt"
	"os"
)

var ExpireTime int = 60

type ConfigStruct struct {
	Redis RedisConfig `json:"Redis"`
	Cache CacheConfig `json:"Cache"`
}

type RedisConfig struct {
	Host     string `json:"Host"`
	Password string `json:"Password"`
	DB       int    `json:"DB"`
}

type CacheConfig struct {
	ExpireTime int `json:"ExpireTime"`
}

var Config ConfigStruct

func ReadConfig(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read config file: %s", err.Error())
		os.Exit(-1)
	}

	if err = json.Unmarshal(data, &Config); err != nil {
		fmt.Printf("Failed to unmarshal JSON: %s", err.Error())
	}
}
