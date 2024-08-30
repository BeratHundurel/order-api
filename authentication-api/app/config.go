package application

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort uint16
}

func LoadConfig() Config {
	cfg := Config{
		ServerPort: 8081,
	}

	if v := os.Getenv("SERVER_PORT"); v != "" {
		if p, err := strconv.ParseUint(v, 10, 16); err == nil {
			cfg.ServerPort = uint16(p)
		}
	}

	return cfg
}
