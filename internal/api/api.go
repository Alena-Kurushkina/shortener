package api

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Alena-Kurushkina/shortener/internal/repository"
)

type HandlerInterface interface {
	createShortening(res http.ResponseWriter, req *http.Request)
	getFullString(res http.ResponseWriter, req *http.Request)
	HandleRequest(res http.ResponseWriter, req *http.Request)
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

func(sh *Shortener) HandleRequest(res http.ResponseWriter, req *http.Request){
	log.Println("Request:", req)
	if req.Method==http.MethodPost {
		sh.createShortening(res,req)
	} else if req.Method==http.MethodGet {
		sh.getFullString(res,req)
	} else {
		http.Error(res, "Request method unsupported", http.StatusBadRequest)
	}
}

func (sh *Shortener) createShortening(res http.ResponseWriter, req *http.Request){
	// if req.Method!=http.MethodPost{
	// 	http.Error(res, "Request method unsupported", http.StatusBadRequest)
	// 	log.Println("Request method unsupported", req.Method, req.URL)
	// 	return
	// }
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
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("EwHXdJfB"))
	log.Println("POST response: ", "EwHXdJfB" )
}

func (sh *Shortener) getFullString(res http.ResponseWriter, req *http.Request){
	// if req.Method!=http.MethodGet{
	// 	http.Error(res, "Request method unsupported", http.StatusBadRequest)
	// 	log.Println("Request method unsupported", req.Method, req.URL)
	// 	return
	// }
	log.Println("GET Request: ", req.Method, req.Host, req.URL )
	param:=strings.TrimPrefix(req.URL.Path, "/")
	param1:=req.PathValue("id")
	_=param1
	if param==""{
		log.Println("Empty parameter")
		http.Error(res, "Bad parameters", http.StatusBadRequest)
		return
	}
	repoOutput:=sh.repository.Select(param)
	if len(repoOutput)==0{
		http.Error(res, "Full string is not found", http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Location", repoOutput)
	res.WriteHeader(http.StatusTemporaryRedirect)

	log.Println("GET Response: ", repoOutput )
}