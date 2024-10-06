package config

import "flag"

type Config struct{
	BaseUrl *string
	ServerAddress *string
}

// var Cfg Config

func InitConfig() Config {
	return Config{
		ServerAddress: flag.String("a", "localhost:8080", "address of HTTP server"),
		BaseUrl: flag.String("b", "http://localhost:8080/", "base address of shorten URL"),
	}
}
