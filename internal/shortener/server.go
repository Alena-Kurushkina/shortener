package shortener

import (
	"flag"
	"log"
	"net/http"

	"github.com/Alena-Kurushkina/shortener/internal/api"
	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/repository"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {	
	Handler chi.Router
	Repository *repository.Repository
	Config config.Config
}

func NewServer() *Server{
	repo:=repository.NewRepository()	
	config:=config.InitConfig()
	sh:=api.NewShortener(repo,config)

	return &Server{
		Handler: NewRouter(sh),
		Repository: repo,
		Config: config,
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
	flag.Parse()

	log.Println("I am listening on ", *s.Config.ServerAddress)
	err:=http.ListenAndServe(*s.Config.ServerAddress, s.Handler)
	if err!=nil{
		panic(err)
	}
}