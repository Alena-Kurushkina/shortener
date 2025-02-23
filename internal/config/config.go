// Package config works with configuration variables
// It parse command flags and read environment variables.
// If environment variable is defined, it has highest priority.
// Otherwise flag values are applied.
package config

import (
	"encoding/json"
	"flag"
	"os"
	"sync"
)

// A Config serves configuration variables.
type Config struct {
	ConfigPath string
	Settings
}

type Settings struct {
	BaseURL         string `json:"base_url"`
	ServerAddress   string `json:"server_address"`
	FileStoragePath string `json:"file_storage_path"`
	ConnectionStr   string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}

var (
	cfg  *Config
	once sync.Once
)

// InitConfig initialize configuration variables from flags values and environment variables.
func InitConfig() *Config {
	once.Do(
		func() {
			// default values
			cfg = &Config{}
			cfg.ServerAddress = "localhost:8080"
			cfg.BaseURL = "http://localhost:8080"

			// define flags
			flagValues := &Config{}
			flag.StringVar(&flagValues.ServerAddress, "a", "", "address of HTTP server")
			flag.StringVar(&flagValues.BaseURL, "b", "", "base address of shorten URL")
			flag.StringVar(&flagValues.FileStoragePath, "f", "", "path to storage file")
			flag.StringVar(&flagValues.ConnectionStr, "d", "", "connection string to database")
			flag.BoolVar(&flagValues.EnableHTTPS, "s", false, "enable HTTPS")

			flag.StringVar(&cfg.ConfigPath, "c", "", "path to config file")
			flag.StringVar(&cfg.ConfigPath, "config", "", "path to config file")
			// parse flags
			flag.Parse()

			// read configs from file
			con, exists := os.LookupEnv("CONFIG")
			if exists {
				cfg.ConfigPath = con
			}
			if cfg.ConfigPath != "" {
				settings := &Settings{}
				readConfigFromFile(cfg.ConfigPath, settings)

				if settings.BaseURL != "" {
					cfg.BaseURL = settings.BaseURL
				}
				if settings.ConnectionStr != "" {
					cfg.ConnectionStr = settings.ConnectionStr
				}
				if !settings.EnableHTTPS {
					cfg.EnableHTTPS = settings.EnableHTTPS
				}
				if settings.FileStoragePath != "" {
					cfg.FileStoragePath = settings.FileStoragePath
				}
				if settings.ServerAddress != "" {
					cfg.ServerAddress = settings.ServerAddress
				}
			}

			// read environment variables
			sa, exists := os.LookupEnv("SERVER_ADDRESS")
			if exists {
				cfg.ServerAddress = sa
			} else if flagValues.ServerAddress != "" {
				cfg.ServerAddress = flagValues.ServerAddress
			}
			bu, exists := os.LookupEnv("BASE_URL")
			if exists {
				cfg.BaseURL = bu
			} else if flagValues.BaseURL != "" {
				cfg.BaseURL = flagValues.BaseURL
			}
			fu, exists := os.LookupEnv("FILE_STORAGE_PATH")
			if exists {
				cfg.FileStoragePath = fu
			} else if flagValues.FileStoragePath != "" {
				cfg.FileStoragePath = flagValues.FileStoragePath
			}
			du, exists := os.LookupEnv("DATABASE_DSN")
			if exists {
				cfg.ConnectionStr = du
			} else if flagValues.ConnectionStr != "" {
				cfg.ConnectionStr = flagValues.ConnectionStr
			}
			_, exists = os.LookupEnv("ENABLE_HTTPS")
			if exists {
				cfg.EnableHTTPS = true
			} else if flagValues.EnableHTTPS {
				cfg.EnableHTTPS = flagValues.EnableHTTPS
			}

			// form BaseURL variable
			if len(cfg.BaseURL) != 0 && cfg.BaseURL[len(cfg.BaseURL)-1:] != "/" {
				cfg.BaseURL = cfg.BaseURL + "/"
			}
		})
	return cfg
}

func readConfigFromFile(pathToConfig string, settings *Settings) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(dir + "/" + pathToConfig)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, settings); err != nil {
		return err
	}
	return nil
}
