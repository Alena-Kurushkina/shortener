// Package config works with configuration variables
// It parse command flags and read environment variables.
// If environment variable is defined, it has highest priority.
// Otherwise flag values are applied.
package config

import (
	"flag"
	"os"
)

// A Config serves configuration variables
type Config struct {
	BaseURL       string
	ServerAddress string
}

// InitConfig initialize configuration variables from flags values and environment variables
func InitConfig() *Config {
	cfg := Config{}

	// define flags
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "address of HTTP server")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "base address of shorten URL")

	// parse flags
	flag.Parse()

	// read environment variables
	sa, exists := os.LookupEnv("SERVER_ADDRESS")
	if exists {
		cfg.ServerAddress = sa
	}
	bu, exists := os.LookupEnv("BASE_URL")
	if exists {
		cfg.BaseURL = bu
	}

	// form BaseURL variable
	if cfg.BaseURL[len(cfg.BaseURL)-1:] != "/" {
		cfg.BaseURL = cfg.BaseURL + "/"
	}

	return &cfg
}
