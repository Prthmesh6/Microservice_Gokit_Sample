package config

import (
	"encoding/json"
	"log"

	"github.com/hashicorp/consul/api"
)

type Config struct {
	RedisURL string `json:"redisURL"`
	RedisKey string `json:"redisKey"`
}

const (
	defaultRedisKey = "videos"
	defaultRedisURL = "localhost:6379"
)

func SetConfigs(pair *api.KVPair) *Config {

	if pair == nil {
		return &Config{
			RedisURL: defaultRedisURL,
			RedisKey: defaultRedisKey,
		}
	}

	var config Config
	err := json.Unmarshal(pair.Value, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if !isValid(&config) {
		return &Config{
			RedisURL: defaultRedisURL,
			RedisKey: defaultRedisKey,
		}
	}
	return &config
}

func isValid(conf *Config) bool {
	if conf.RedisKey == "" {
		return false
	}
	if conf.RedisURL == "" {
		return false
	}
	return true
}
