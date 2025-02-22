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
			cfg = &Config{}

			// define flags
			flag.StringVar(&cfg.ServerAddress, "a", "", "address of HTTP server")
			flag.StringVar(&cfg.BaseURL, "b", "", "base address of shorten URL")
			flag.StringVar(&cfg.FileStoragePath, "f", "", "path to storage file")
			flag.StringVar(&cfg.ConnectionStr, "d", "", "connection string to database")
			flag.BoolVar(&cfg.EnableHTTPS, "s", false, "enable HTTPS")
			flag.StringVar(&cfg.ConfigPath, "c", "", "path to config file")
			flag.StringVar(&cfg.ConfigPath, "config", "", "path to config file")
			// parse flags
			flag.Parse()

			con, exists := os.LookupEnv("CONFIG")
			if exists {
				cfg.ConfigPath = con
			}
			if cfg.ConfigPath != "" {
				settings := &Settings{}
				readConfigFromFile(cfg.ConfigPath, settings)
				
				if cfg.BaseURL == "" {
					cfg.BaseURL = settings.BaseURL
				}
				if cfg.ConnectionStr == "" {
					cfg.ConnectionStr = settings.ConnectionStr
				}
				if !cfg.EnableHTTPS {
					cfg.EnableHTTPS = settings.EnableHTTPS
				}
				if cfg.FileStoragePath == "" {
					cfg.FileStoragePath = settings.FileStoragePath
				}
				if cfg.ServerAddress == "" {
					cfg.ServerAddress = settings.ServerAddress
				}
			}

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
			_, exists = os.LookupEnv("ENABLE_HTTPS")
			if exists {
				cfg.EnableHTTPS = true
			}

			// form BaseURL variable
			if len(cfg.BaseURL)!=0 && cfg.BaseURL[len(cfg.BaseURL)-1:] != "/" {
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
