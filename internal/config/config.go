// Package config works with configuration variables
// It parse command flags and read environment variables.
// If environment variable is defined, it has highest priority.
// Otherwise flag values are applied.
package config

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"sync"
)

// A Config serves all configuration variables.
type Config struct {
	ConfigPath string
	Settings
}

// A Settings keeps service main configurations.
type Settings struct {
	BaseURL           string `json:"base_url"`
	ServerAddress     string `json:"server_address"`
	GRPCServerAddress string `json:"grpc_server_address"`
	FileStoragePath   string `json:"file_storage_path"`
	ConnectionStr     string `json:"database_dsn"`
	EnableHTTPS       bool   `json:"enable_https"`
	TrustedSubnet     string `json:"trusted_subnet"`
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
			cfg.GRPCServerAddress = ":3200"

			// define flags
			flagValues := &Config{}
			flag.StringVar(&flagValues.ServerAddress, "a", "", "address of HTTP server")
			flag.StringVar(&flagValues.GRPCServerAddress, "g", "", "address of gRPC server")
			flag.StringVar(&flagValues.BaseURL, "b", "", "base address of shorten URL")
			flag.StringVar(&flagValues.FileStoragePath, "f", "", "path to storage file")
			flag.StringVar(&flagValues.ConnectionStr, "d", "", "connection string to database")
			flag.BoolVar(&flagValues.EnableHTTPS, "s", false, "enable HTTPS")
			flag.StringVar(&flagValues.TrustedSubnet, "t", "", "trusted subnet")

			flag.StringVar(&cfg.ConfigPath, "c", "", "path to config file")
			flag.StringVar(&cfg.ConfigPath, "config", "", "path to config file")
			flag.Parse()

			// check if config file was set
			con, exists := os.LookupEnv("CONFIG")
			if exists {
				cfg.ConfigPath = con
			}
			if cfg.ConfigPath != "" {
				// read configs from file to settings var
				settings := &Settings{}
				readConfigFromFile(cfg.ConfigPath, settings)
				// if config param exists, save it to main config var
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
				if settings.GRPCServerAddress != "" {
					cfg.GRPCServerAddress = settings.GRPCServerAddress
				}
				if settings.TrustedSubnet != "" {
					cfg.TrustedSubnet = settings.TrustedSubnet
				}
			}

			// read environment variables
			// if it exists, save value to main config var
			// else save flag value
			sa, exists := os.LookupEnv("SERVER_ADDRESS")
			if exists {
				cfg.ServerAddress = sa
			} else if flagValues.ServerAddress != "" {
				cfg.ServerAddress = flagValues.ServerAddress
			}
			ga, exists := os.LookupEnv("GRPC_SERVER_ADDRESS")
			if exists {
				cfg.GRPCServerAddress = ga
			} else if flagValues.GRPCServerAddress != "" {
				cfg.GRPCServerAddress = flagValues.GRPCServerAddress
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

			tu, exists := os.LookupEnv("TRUSTED_SUBNET")
			if exists {
				cfg.TrustedSubnet = tu
			} else if flagValues.TrustedSubnet != "" {
				cfg.TrustedSubnet = flagValues.TrustedSubnet
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
	data, err := os.ReadFile(filepath.Join(dir,pathToConfig))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, settings); err != nil {
		return err
	}
	return nil
}
