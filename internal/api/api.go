package api

import (
	"io"
	"net/http"
	"log"

	"github.com/Alena-Kurushkina/shortener/internal/repository"
)

type HandlerInterface interface {
	CreateShortening(res http.ResponseWriter, req *http.Request)
	GetFullString(res http.ResponseWriter, req *http.Request)
}


type Shortener struct {
	repository *repository.Repository
}

func NewShortener(repo *repository.Repository) *Shortener {
	shortener:=Shortener{
		repository: repo,
	}
	return &shortener
}

func (sh *Shortener) CreateShortening(res http.ResponseWriter, req *http.Request){
	log.Println("POST Request: ", req.URL, req.Method, req.Host )
	body, err :=io.ReadAll(req.Body)
	if err!=nil{
		http.Error(res, "Can't read body", http.StatusBadRequest)
		return
	}	
	log.Println("POST body: ", string(body) )
	if len(body)==0{
		http.Error(res, "Body is empty", http.StatusBadRequest)
		return
	}	
	sh.repository.Insert("EwHXdJfB", string(body))
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("EwHXdJfB"))
	log.Println("POST response: ", "EwHXdJfB" )
}

func (sh *Shortener) GetFullString(res http.ResponseWriter, req *http.Request){
	log.Println("GET Request: ", req.URL, req.Method, req.Host )
	param:=req.PathValue("id")
	if param==""{
		http.Error(res, "Bad parameters", http.StatusBadRequest)
		return
	}
	repoOutput:=sh.repository.Select(param)
	if len(repoOutput)==0{
		http.Error(res, "Full string is not found", http.StatusBadRequest)
		return
	}
	res.Header().Set("content-type", "text/plain")
	res.Header().Set("Location", repoOutput)
	res.WriteHeader(http.StatusTemporaryRedirect)

	log.Println("GET Response: ", repoOutput )
}