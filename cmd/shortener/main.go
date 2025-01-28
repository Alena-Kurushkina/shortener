package main

import (
	"context"
	_ "net/http/pprof"

	_ "github.com/golang/mock/mockgen/model"

	"github.com/Alena-Kurushkina/shortener/internal/api"
	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/logger"
	"github.com/Alena-Kurushkina/shortener/internal/repository"
	"github.com/Alena-Kurushkina/shortener/internal/shortener"
)

func main() {
	cfg := config.InitConfig()
	err := logger.Initialize()
	if err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	ctx := context.Background()

	repo, err := repository.NewRepository(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer repo.Close()

	sh := api.NewShortener(repo, cfg)

	server := shortener.NewServer(sh, cfg)

	server.Run()
}
