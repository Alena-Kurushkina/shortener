package api

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Alena-Kurushkina/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

type HandlerInterface interface {
	CreateShortening(res http.ResponseWriter, req *http.Request)
	GetFullString(res http.ResponseWriter, req *http.Request)
	// HandleRequest(res http.ResponseWriter, req *http.Request)
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

// func(sh *Shortener) HandleRequest(res http.ResponseWriter, req *http.Request){
// 	if req.Method==http.MethodPost {
// 		sh.createShortening(res,req)
// 	} else if req.Method==http.MethodGet {
// 		sh.getFullString(res,req)
// 	} else {
// 		http.Error(res, "Request method unsupported", http.StatusBadRequest)
// 	}
// }

func (sh *Shortener) CreateShortening(res http.ResponseWriter, req *http.Request){
	// if req.Method!=http.MethodPost{
	// 	http.Error(res, "Request method unsupported", http.StatusBadRequest)
	// 	log.Println("Request method unsupported", req.Method, req.URL)
	// 	return
	// }
	contentType := req.Header.Get("Content-Type")
	// log.Println("POST Request: ", req.URL, req.Method, req.Host, contentType )

	res.Header().Set("Content-Type", "text/plain")

	var url = ""
	if contentType == "application/x-www-form-urlencoded" {
		req.ParseForm()
		url = req.FormValue("url")
	} else if strings.Contains(contentType,"text/plain") {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(res, "Can't read body", http.StatusBadRequest)
			return
		}
		url=string(body)
		// url = strings.TrimSuffix(string(urlBytes), "\n")
	} else {
		http.Error(res, "Invalid content type", http.StatusBadRequest)
		return
	}
	log.Println("POST body: ", url )

	if len(url)==0{
		http.Error(res, "Body is empty", http.StatusBadRequest)
		return
	}	
	sh.repository.Insert("EwHXdJfB", url)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("EwHXdJfB"))

	log.Println("POST response: ", "EwHXdJfB" )
}

func (sh *Shortener) GetFullString(res http.ResponseWriter, req *http.Request){
	// if req.Method!=http.MethodGet{
	// 	http.Error(res, "Request method unsupported", http.StatusBadRequest)
	// 	log.Println("Request method unsupported", req.Method, req.URL)
	// 	return
	// }
	// log.Println("GET Request: ", req.Method, req.Host, req.URL )
	// param:=strings.TrimPrefix(req.URL.Path, "/")
	// param1:=req.PathValue("id")
	// _=param1
	param:=chi.URLParam(req,"id")
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