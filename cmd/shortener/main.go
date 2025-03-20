// Shortener is a service that accepts and stores long urls and serves requests for corresponding shortenings
package main

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang/mock/mockgen/model"

	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/core"
	"github.com/Alena-Kurushkina/shortener/internal/logger"
	"github.com/Alena-Kurushkina/shortener/internal/repository"
	"github.com/Alena-Kurushkina/shortener/internal/shortener"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

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
	//defer repo.Close() //done in shutdown

	// через него сообщаем основному потоку, что все сетевые соединения обработаны и закрыты
	idleConnsClosed := make(chan struct{})

	core := core.NewShortenerCore(repo, cfg)

	// sh := api.NewShortener(core)
	httpServer := shortener.NewServer(core, cfg, idleConnsClosed)
	rpcServer := shortener.NewServerGRPC(core, cfg, idleConnsClosed)

	// канал для перенаправления прерываний
	sigint := make(chan os.Signal, 1)
	// регистрируем перенаправление прерываний
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	// запускаем горутину обработки пойманных прерываний
	go func() {
		// читаем из канала прерываний
		<-sigint
		// запускаем процедуру graceful shutdown
		if err := httpServer.HTTPServer.Shutdown(context.Background()); err != nil {
			// ошибки закрытия Listener
			logger.Log.Errorf("HTTP server Shutdown: %v", err)
		}
		logger.Log.Info("HTTP server shutdown seccussfully")

		rpcServer.Server.GracefulStop()
		logger.Log.Info("gRPC server was stopped seccussfully")

		// сообщаем основному потоку, что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)

		core.Shutdown()
	}()

	go httpServer.Run()
	rpcServer.Run()
}
