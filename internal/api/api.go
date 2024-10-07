package api

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

type HandlerInterface interface {
	CreateShortening(res http.ResponseWriter, req *http.Request)
	GetFullString(res http.ResponseWriter, req *http.Request)
}


type Shortener struct {
	repository *repository.Repository
	config config.Config
}

func NewShortener(repo *repository.Repository, cfg config.Config) *Shortener {
	shortener:=Shortener{
		repository: repo,
		config: cfg,
	}
	return &shortener
}

func (sh *Shortener) CreateShortening(res http.ResponseWriter, req *http.Request){
	// if req.Method!=http.MethodPost{
	// 	http.Error(res, "Request method unsupported", http.StatusBadRequest)
	// 	log.Println("Request method unsupported", req.Method, req.URL)
	// 	return
	// }
	contentType := req.Header.Get("Content-Type")	

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
	} else {
		http.Error(res, "Invalid content type", http.StatusBadRequest)
		return
	}
	// log.Println("POST body: ", url )

	if len(url)==0{
		http.Error(res, "Body is empty", http.StatusBadRequest)
		return
	}
	shortener:=generateRandomString(15)	
	sh.repository.Insert(shortener, url)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(sh.config.BaseUrl+shortener))
	// log.Println("POST response: ", shortener )
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
	// log.Println("GET Response: ", repoOutput )
}

func generateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    seed := rand.NewSource(time.Now().UnixNano())
    random := rand.New(seed)

    result := make([]byte, length)
    for i := range result {
        result[i] = charset[random.Intn(len(charset))]
    }
    return string(result)
}