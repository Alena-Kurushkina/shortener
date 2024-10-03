package shortener

import (
	"net/http"
	"github.com/Alena-Kurushkina/shortener/internal/api"
	"github.com/Alena-Kurushkina/shortener/internal/repository"	
)

type Server struct {	
	Handler http.Handler
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

func NewRouter(hi api.HandlerInterface) http.Handler {
	mux:=http.NewServeMux()
	mux.HandleFunc(`POST /`, hi.CreateShortening)
	mux.HandleFunc(`GET /{id}`, hi.GetFullString)
	return mux
}

func (s *Server) Run() {	
	err:=http.ListenAndServe(`:8080`, s.Handler)
	if err!=nil{
		panic(err)
	}
}