package shortener

import (
	"net/http"

	"github.com/Alena-Kurushkina/shortener/internal/api"
	"github.com/Alena-Kurushkina/shortener/internal/repository"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {	
	Handler chi.Router
	Repository *repository.Repository
}

func NewServer() *Server{
	repo:=repository.NewRepository()
	sh:=api.NewShortener(repo)

	return &Server{
		Handler: NewRouter(sh),
		Repository: repo,
	}	
}

func NewRouter(hi api.HandlerInterface) chi.Router {
	r:=chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/",hi.CreateShortening)
	r.Get("/{id}",hi.GetFullString)

	return r
	// mux:=http.NewServeMux()
	// mux.HandleFunc("/", hi.HandleRequest)
	// return mux
}

func (s *Server) Run() {	
	err:=http.ListenAndServe(":8080", s.Handler)
	if err!=nil{
		panic(err)
	}
}