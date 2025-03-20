// Package shortener implements HTTP server.
// It tunes requests routing.
package shortener

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/acme/autocert"

	"github.com/Alena-Kurushkina/shortener/internal/api"
	"github.com/Alena-Kurushkina/shortener/internal/authenticator"
	"github.com/Alena-Kurushkina/shortener/internal/compress"
	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/core"
	"github.com/Alena-Kurushkina/shortener/internal/logger"
)

// A Handler represent interface for shortening handler.
type Handler interface {
	CreateShortening(res http.ResponseWriter, req *http.Request)
	GetFullString(res http.ResponseWriter, req *http.Request)
	CreateShorteningJSON(res http.ResponseWriter, req *http.Request)
	CreateShorteningJSONBatch(res http.ResponseWriter, req *http.Request)
	PingDB(res http.ResponseWriter, req *http.Request)
	GetUserAllShortenings(res http.ResponseWriter, req *http.Request)
	DeleteRecordJSON(res http.ResponseWriter, req *http.Request)
	GetStats(res http.ResponseWriter, req *http.Request)
	//Shutdown()
}

// A Server aggregates handler and config.
type Server struct {
	HTTPServer http.Server
	Shortener  Handler
	//Handler chi.Router
	Config          *config.Config
	IdleConnsClosed chan struct{}
}

// NewRouter creates new routes and middlewares.
func newRouter(hi Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/ping", hi.PingDB)
	r.Get("/{id}", hi.GetFullString)

	r.Get("/api/internal/stats", hi.GetStats)

	r.Get("/debug/pprof/", pprof.Index)
	r.Get("/debug/pprof/profile", pprof.Profile)
	r.Get("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)

	r.Group(func(r chi.Router) {
		r.Use(compress.GzipMiddleware, logger.LogMiddleware, authenticator.AuthMiddleware)

		r.Post("/", hi.CreateShortening)
		// r.Get("/{id}", hi.GetFullString)
		r.Get("/api/user/urls", hi.GetUserAllShortenings)
		r.Post("/api/shorten", hi.CreateShorteningJSON)
		r.Post("/api/shorten/batch", hi.CreateShorteningJSONBatch)
		r.Delete("/api/user/urls", hi.DeleteRecordJSON)
	})

	return r
}

// NewServer initializes new server with config and handler.
func NewServer(core *core.ShortenerCore, cfg *config.Config, idleConnsChan chan struct{}) *Server {
	hdl := api.Shortener{
		Core: core,
	}

	srv := &Server{
		//Handler: newRouter(hdl),
		HTTPServer: http.Server{
			Handler: newRouter(&hdl),
			Addr:    cfg.ServerAddress,
		},
		Config: cfg,
		//Shortener: hdl,
		// через этот канал сообщим основному потоку, что соединения закрыты
		IdleConnsClosed: idleConnsChan,
	}
	// // канал для перенаправления прерываний
	// sigint := make(chan os.Signal, 1)
	// // регистрируем перенаправление прерываний
	// signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	// // запускаем горутину обработки пойманных прерываний
	// go func() {
	// 	// читаем из канала прерываний
	// 	<-sigint
	// 	// запускаем процедуру graceful shutdown
	// 	if err := srv.HTTPServer.Shutdown(context.Background()); err != nil {
	// 		// ошибки закрытия Listener
	// 		logger.Log.Errorf("HTTP server Shutdown: %v", err)
	// 	}
	// 	logger.Log.Info("HTTP server shutdown seccussfully")
	// 	// сообщаем основному потоку, что все сетевые соединения обработаны и закрыты
	// 	close(srv.IdleConnsClosed)
	// }()

	return srv
}

// Run starts listening to server address and handling requests.
func (s *Server) Run() {
	logger.Log.Infof("Server is listening on %s", s.Config.ServerAddress)
	logger.Log.Infof("Base URL: %s", s.Config.BaseURL)

	if !s.Config.EnableHTTPS {
		logger.Log.Infof("HTTPS disabled")

		if err := s.HTTPServer.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	} else {
		logger.Log.Infof("HTTPS enabled")

		manager := &autocert.Manager{
			// директория для хранения сертификатов
			Cache: autocert.DirCache("cache-dir"),
			// функция, принимающая Terms of Service издателя сертификатов
			Prompt: autocert.AcceptTOS,
		}

		s.HTTPServer.TLSConfig = manager.TLSConfig()
		if err := s.HTTPServer.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
			logger.Log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}
	// ждём завершения процедуры graceful shutdown
	<-s.IdleConnsClosed
	// получили оповещение о завершении
	// здесь можно освобождать ресурсы перед выходом,
	// например закрыть соединение с базой данных
	//s.Shortener.Shutdown()
	logger.Log.Info("Server Shutdown gracefully")
}
