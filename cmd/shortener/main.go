package main

import (
	"github.com/Alena-Kurushkina/shortener/internal/api"
	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/repository"
	"github.com/Alena-Kurushkina/shortener/internal/shortener"
)

func main() {
	cfg := config.InitConfig()
	repo := repository.NewRepository()

	sh := api.NewShortener(repo, cfg)

	server := shortener.NewServer(sh, cfg)

	server.Run()
}
