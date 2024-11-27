// Package config works with configuration variables
// It parse command flags and read environment variables.
// If environment variable is defined, it has highest priority.
// Otherwise flag values are applied.
package config

import (
	"flag"
	"fmt"
	"os"
)

// A Config serves configuration variables
type Config struct {
	BaseURL         string
	ServerAddress   string
	FileStoragePath string
	ConnectionStr   string
}

// InitConfig initialize configuration variables from flags values and environment variables
func InitConfig() *Config {
	cfg := Config{}

	// define flags
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "address of HTTP server")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "base address of shorten URL")
	//flag.StringVar(&cfg.FileStoragePath, "f", "/Users/shortener_storage.txt", "path to storage file")
	flag.StringVar(&cfg.FileStoragePath, "f", "", "path to storage file")
	flag.StringVar(&cfg.ConnectionStr, "d", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", `127.0.0.1`, `practicum`, `123456`, `practicumdb`), "connection string to database")
	//flag.StringVar(&cfg.ConnectionStr, "d", "", "connection string to database")

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
	fu, exists := os.LookupEnv("FILE_STORAGE_PATH")
	if exists {
		cfg.FileStoragePath = fu
	}
	du, exists := os.LookupEnv("DATABASE_DSN")
	if exists {
		cfg.ConnectionStr = du
	}

	// form BaseURL variable
	if cfg.BaseURL[len(cfg.BaseURL)-1:] != "/" {
		cfg.BaseURL = cfg.BaseURL + "/"
	}

	return &cfg
}
