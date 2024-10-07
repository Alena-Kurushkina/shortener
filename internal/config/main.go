package config

import (
	"flag"
	"os"
)

type Config struct{
	BaseUrl string
	ServerAddress string
}

// var Cfg Config

func InitConfig() Config {
	// return Config{
	// 	ServerAddress: flag.String("a", "localhost:8080", "address of HTTP server"),
	// 	BaseUrl: flag.String("b", "http://localhost:8080/", "base address of shorten URL"),
	// }

	cfg:=Config{}

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "address of HTTP server")
	flag.StringVar(&cfg.BaseUrl, "b", "http://localhost:8080", "base address of shorten URL")

	flag.Parse()

	sa, exists:= os.LookupEnv("SERVER_ADDRESS"); if exists {
		cfg.ServerAddress=sa
	}
	bu, exists:= os.LookupEnv("BASE_URL"); if exists {
		cfg.BaseUrl=bu
	}

	if cfg.BaseUrl[len(cfg.BaseUrl)-1:]!="/"{
		cfg.BaseUrl=cfg.BaseUrl+"/"
	}

	return cfg
}
