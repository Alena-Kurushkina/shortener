package api

import (
	"io"
	"net/http"

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
	body, err :=io.ReadAll(req.Body)
	if err!=nil{
		http.Error(res, "Can't read body", http.StatusBadRequest)
		return
	}	
	if req.Host!="localhost:8080"{
		http.Error(res, "Hostname incorrect", http.StatusBadRequest)
		return
	}
	if len(body)==0{
		http.Error(res, "Body is empty", http.StatusBadRequest)
		return
	}	
	sh.repository.Insert("EwHXdJfB", string(body))
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("EwHXdJfB"))
}

func (sh *Shortener) GetFullString(res http.ResponseWriter, req *http.Request){
	param:=req.PathValue("id")
	if param==""{
		http.Error(res, "Bad parameters", http.StatusBadRequest)
		return
	}
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)
	res.Write([]byte(sh.repository.Select(param)))
}