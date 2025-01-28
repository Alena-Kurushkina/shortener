// Package shortener implements HTTP server.
// It tunes requests routing.
package shortener

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"

	"github.com/Alena-Kurushkina/shortener/internal/authenticator"
	"github.com/Alena-Kurushkina/shortener/internal/compress"
	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/logger"
)

// A Handler represent interface for shortening handler
type Handler interface {
	CreateShortening(res http.ResponseWriter, req *http.Request)
	GetFullString(res http.ResponseWriter, req *http.Request)
	CreateShorteningJSON(res http.ResponseWriter, req *http.Request)
	CreateShorteningJSONBatch(res http.ResponseWriter, req *http.Request)
	PingDB(res http.ResponseWriter, req *http.Request)
	GetUserAllShortenings(res http.ResponseWriter, req *http.Request)
	DeleteRecordJSON(res http.ResponseWriter, req *http.Request)
}

// A Server aggregates handler and config
type Server struct {
	Handler chi.Router
	Config  *config.Config
}

// NewRouter creates new routes and middlewares
func newRouter(hi Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/ping", hi.PingDB)
	r.Get("/{id}", hi.GetFullString)

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

// NewServer initializes new server with config and handler
func NewServer(hdl Handler, cfg *config.Config) *Server {
	return &Server{
		Handler: newRouter(hdl),
		Config:  cfg,
	}
}

// Run starts listening to server address and handling requests
func (s *Server) Run() {
	logger.Log.Infof("Server is listening on %s", s.Config.ServerAddress)
	err := http.ListenAndServe(s.Config.ServerAddress, s.Handler)
	if err != nil {
		panic(err)
	}
}
