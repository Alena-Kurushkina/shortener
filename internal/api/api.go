// Package api implements handler functions for shorten long URL 
// and expanding shortenings back to long URL
package api

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

// A HandlerInterface represent interface for shortening handler
type HandlerInterface interface {
	CreateShortening(res http.ResponseWriter, req *http.Request)
	GetFullString(res http.ResponseWriter, req *http.Request)
}

// A Shortener aggregates helpfull elements 
type Shortener struct {
	repository *repository.Repository
	config config.Config
}

// NewShortener returns new Shortener pointer initialized by repository and config
func NewShortener(repo *repository.Repository, cfg config.Config) *Shortener {
	shortener:=Shortener{
		repository: repo,
		config: cfg,
	}
	return &shortener
}

// CreateShortening habdle POST HTTP request with long URL in body and retrieves base URL with shortening.
// It handle only requests with content type application/x-www-form-urlencoded or text/plain.
// Response body has content type text/plain.
func (sh *Shortener) CreateShortening(res http.ResponseWriter, req *http.Request){
	// set response content type
	res.Header().Set("Content-Type", "text/plain")
	
	// parse request body
	contentType := req.Header.Get("Content-Type")
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
	if len(url)==0{
		http.Error(res, "Body is empty", http.StatusBadRequest)
		return
	}

	// generate shortening
	shortener:=generateRandomString(15)	
	if err:=sh.repository.Insert(shortener, url); err!=nil{
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// make response
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(sh.config.BaseUrl+shortener))
}

// GetFullString handle GET request with shortening in URL parameter named id 
// and makes response with long URL in header's location value.
// Response content type is text/plain
func (sh *Shortener) GetFullString(res http.ResponseWriter, req *http.Request){
	// parse parameter id from URL
	param:=chi.URLParam(req,"id")
	if param==""{
		http.Error(res, "Bad parameters", http.StatusBadRequest)
		return
	}

	// get long URL from repository
	repoOutput, err:=sh.repository.Select(param)
	if err!=nil{
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// make responce
	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Location", repoOutput)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

// generateRandomString returns string of random characters of passed length
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