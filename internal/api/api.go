// Package api implements handler functions for shorten long URL
// and expanding shortenings back to long URL
package api

import (
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/shortener"
)

type Storager interface {
	Insert(key, value string) error
	Select(key string) (string, error)
	Ping(ctx context.Context) error
	Close()
}

// A Shortener aggregates helpfull elements
type Shortener struct {
	repo   Storager
	config *config.Config
}

// NewShortener returns new Shortener pointer initialized by repository and config
func NewShortener(storage Storager, cfg *config.Config) shortener.Handler {
	shortener := Shortener{
		repo:   storage,
		config: cfg,
	}
	return &shortener
}



// CreateShortening habdle POST HTTP request with long URL in body and retrieves base URL with shortening.
// It handle only requests with content type application/x-www-form-urlencoded or text/plain.
// Response body has content type text/plain.
func (sh *Shortener) CreateShortening(res http.ResponseWriter, req *http.Request) {
	// set response content type
	res.Header().Set("Content-Type", "text/plain")

	// parse request body
	contentType := req.Header.Get("Content-Type")

	var url string
	if contentType == "application/x-www-form-urlencoded" {
		req.ParseForm()
		url = req.FormValue("url")
	} else if strings.Contains(contentType, "text/plain") || strings.Contains(contentType, "application/x-gzip") {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			http.Error(res, "Can't read body", http.StatusBadRequest)
			return
		}
		url = string(body)
	} else {
		http.Error(res, "Invalid content type", http.StatusBadRequest)
		return
	}
	if len(url) == 0 {
		http.Error(res, "Body is empty", http.StatusBadRequest)
		return
	}

	// generate shortening
	shortStr := generateRandomString(15)
	if err := sh.repo.Insert(shortStr, url); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// make response
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(sh.config.BaseURL + shortStr))
}

// A URLRequest is for request decode from json
type URLRequest struct {
	URL string `json:"url"`
}

// A ResultResponse is for response encode in json
type ResultResponse struct {
	Result string `json:"result"`
}

// CreateShorteningJSON handle POST HTTP request with long URL in body and retrieves base URL with shortening.
// It handle only requests with content type application/json.
// Response body has content type application/json.
func (sh *Shortener) CreateShorteningJSON(res http.ResponseWriter, req *http.Request) {
	// set response content type
	res.Header().Set("Content-Type", "application/json")

	// check content type
	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" && contentType != "application/x-gzip" {
		http.Error(res, "Invalid content type", http.StatusBadRequest)
		return
	}

	// decode request body
	var url URLRequest
	if err := json.NewDecoder(req.Body).Decode(&url); err != nil {
		http.Error(res, "Can't read body", http.StatusBadRequest)
		return
	}
	if len(url.URL) == 0 {
		http.Error(res, "Body is empty", http.StatusBadRequest)
		return
	}

	// generate shortening
	shortStr := generateRandomString(15)
	if err := sh.repo.Insert(shortStr, url.URL); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// make response
	responseData, err := json.Marshal(ResultResponse{
		Result: sh.config.BaseURL + shortStr,
	})
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write(responseData)
}

// GetFullString handle GET request with shortening in URL parameter named id
// and makes response with long URL in header's location value.
// Response content type is text/plain
func (sh *Shortener) GetFullString(res http.ResponseWriter, req *http.Request) {
	// parse parameter id from URL
	param := chi.URLParam(req, "id")
	if param == "" {
		http.Error(res, "Bad parameters", http.StatusBadRequest)
		return
	}

	// get long URL from repository
	repoOutput, err := sh.repo.Select(param)
	if err != nil {
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

func (sh *Shortener) PingDB(res http.ResponseWriter, req *http.Request) {
	
	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
    defer cancel()
    if err := sh.repo.Ping(ctx); err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
		return
    }

	// make responce
	res.WriteHeader(http.StatusOK)
}