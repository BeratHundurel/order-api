package application

import (
	"os"
	"strconv"
)

type Config struct {
	RedisAddress string
	ServerPort   uint16
}

func LoadConfig() Config {
	cfg := Config{
		RedisAddress: "redis:6379",
		ServerPort:   8080,
	}

	if v := os.Getenv("REDIS_ADDR"); v != "" {
		cfg.RedisAddress = v
	}

	if v := os.Getenv("SERVER_PORT"); v != "" {
		if p, err := strconv.ParseUint(v, 10, 16); err == nil {
			cfg.ServerPort = uint16(p)
		}
	}

	return cfg
}
