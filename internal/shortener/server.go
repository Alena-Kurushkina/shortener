// Package shortener implements HTTP server.
// It tunes requests routing.
package shortener

import (
	"log"
	"net/http"

	"github.com/Alena-Kurushkina/shortener/internal/api"
	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/repository"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

// A Server aggregates used variables
type Server struct {	
	Handler chi.Router
	Repository *repository.Repository
	Config config.Config
}

// NewServer initializes new server with repository, config and handler
func NewServer() *Server{
	repo:=repository.NewRepository()	
	config:=config.InitConfig()
	sh:=api.NewShortener(repo,config)

	return &Server{
		Handler: newRouter(sh),
		Repository: repo,
		Config: config,
	}	
}

// NewRouter creates new routes and middlewares
func newRouter(hi api.HandlerInterface) chi.Router {
	r:=chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/",hi.CreateShortening)
	r.Get("/{id}",hi.GetFullString)

	return r
}

// Run starts listening to server address and handling requests
func (s *Server) Run() {	
	log.Println("Server is listening on ", s.Config.ServerAddress)
	err:=http.ListenAndServe(s.Config.ServerAddress, s.Handler)
	if err!=nil{
		panic(err)
	}
}