// Package shortener implements HTTP server.
// It tunes requests routing.
package shortener

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Alena-Kurushkina/shortener/internal/compress"
	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/logger"
)

// A Handler represent interface for shortening handler
type Handler interface {
	CreateShortening(res http.ResponseWriter, req *http.Request)
	GetFullString(res http.ResponseWriter, req *http.Request)
	CreateShorteningJSON(res http.ResponseWriter, req *http.Request)
	PingDB(res http.ResponseWriter, req *http.Request)
}

// A Server aggregates handler and config
type Server struct {
	Handler chi.Router
	Config  *config.Config
}

// NewRouter creates new routes and middlewares
func newRouter(hi Handler) chi.Router {
	r := chi.NewRouter()
	r.Use(compress.GzipMiddleware)
	r.Use(logger.LogMiddleware)

	r.Post("/", hi.CreateShortening)
	r.Get("/{id}", hi.GetFullString)
	r.Get("/ping", hi.PingDB)
	r.Post("/api/shorten", hi.CreateShorteningJSON)

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
	log.Println("Server is listening on ", s.Config.ServerAddress)
	err := http.ListenAndServe(s.Config.ServerAddress, s.Handler)
	if err != nil {
		panic(err)
	}
}
