// Package api implements handler functions to shorten long URL
// and expand shortenings back to long URL.
package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	uuid "github.com/satori/go.uuid"

	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/generator"
	"github.com/Alena-Kurushkina/shortener/internal/logger"
	"github.com/Alena-Kurushkina/shortener/internal/sherr"
	"github.com/Alena-Kurushkina/shortener/internal/shortener"
)

const timeoutPing time.Duration = 30

// Storager defines operations with data storage.
type Storager interface {
	Insert(ctx context.Context, userID uuid.UUID, key, value string) error
	InsertBatch(_ context.Context, userID uuid.UUID, batch []BatchElement) error
	Select(ctx context.Context, key string) (string, error)
	SelectUserAll(ctx context.Context, userID uuid.UUID) ([]BatchElement, error)
	DeleteRecords(ctx context.Context, deleteItems []DeleteItem) error
	Ping(ctx context.Context) error
	Close()
}

// A Shortener aggregates data storage, configurations and helpful objects.
type Shortener struct {
	repo       Storager
	config     *config.Config
	deleteChan chan DeleteItem
	done       chan struct{}
}

func newShortenerObject(storage Storager, cfg *config.Config) *Shortener {
	return &Shortener{
		repo:       storage,
		config:     cfg,
		deleteChan: make(chan DeleteItem, 1024),
		done:       make(chan struct{}, 1),
	}
}

// NewShortener returns new Shortener pointer initialized by repository and config.
func NewShortener(storage Storager, cfg *config.Config) shortener.Handler {
	shortener := newShortenerObject(storage, cfg)

	go shortener.flushDeleteItems()

	return shortener
}

// DeleteItem represents pair of ids which identify unique record to delete.
type DeleteItem struct {
	IDs    []string
	UserID uuid.UUID
}

func (sh *Shortener) flushDeleteItems() {
	ticker := time.NewTicker(1 * time.Second)

	items := make([]DeleteItem, 0, 1024)

	for {
		select {
		case msg := <-sh.deleteChan:
			items = append(items, msg)
		case <-ticker.C:
			sh.deleteRecords(items)

			items = make([]DeleteItem, 0, 1024)
		case <-sh.done:
			sh.deleteRecords(items)
			return
		}
	}
}

func (sh *Shortener) deleteRecords(items []DeleteItem) {
	if len(items) == 0 {
		return
	}
	err := sh.repo.DeleteRecords(context.TODO(), items)
	if err != nil {
		logger.Log.Infof("Can't delete records", err.Error())
		return
	}
	logger.Log.Info("Patch of shortenings was deleted, patch length: " + strconv.Itoa(len(items)))
}

// Shutdown finishes work gracefully
func (sh *Shortener) Shutdown() {
	logger.Log.Info("Start shortener shutdown")
	close(sh.done)
	sh.repo.Close()
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

	logger.Log.Infof("Handle route /, method POST, body: %s", url)

	// generate shortening
	shortStr := generator.GenerateRandomString(15)

	q := req.URL.Query()
	id, err := uuid.FromString(q.Get("userUUID"))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	insertErr := sh.repo.Insert(req.Context(), id, shortStr, url)

	var existError *sherr.AlreadyExistError
	if errors.As(insertErr, &existError) {
		// make response
		res.WriteHeader(http.StatusConflict)
		res.Write([]byte(sh.config.BaseURL + existError.ExistShortStr))

		return
	} else if insertErr != nil {
		http.Error(res, insertErr.Error(), http.StatusBadRequest)
		return
	}

	// make response
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(sh.config.BaseURL + shortStr))
}

// A URLRequest is for request decoding from json.
type URLRequest struct {
	URL string `json:"url"`
}

// A ResultResponse is for response encoding in json.
type ResultResponse struct {
	Result string `json:"result"`
}

// CreateShorteningJSON handle POST HTTP request with long URL in body and retrieves base URL with shortening.
// It handle only requests with content type application/json.
// Response has content type application/json.
// /api/shorten
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

	q := req.URL.Query()
	id, err := uuid.FromString(q.Get("userUUID"))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("Handle route /api/shorten, method POST, body: %s", url.URL)

	// generate shortening
	shortStr := generator.GenerateRandomString(15)
	insertErr := sh.repo.Insert(req.Context(), id, shortStr, url.URL)

	var existError *sherr.AlreadyExistError
	if errors.As(insertErr, &existError) {
		// make response
		responseData, err := json.Marshal(ResultResponse{
			Result: sh.config.BaseURL + existError.ExistShortStr,
		})
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		// make response
		res.WriteHeader(http.StatusConflict)
		res.Write(responseData)
		//res.Write([]byte(sh.config.BaseURL + existError.ExistShortStr))

		return
	} else if insertErr != nil {
		http.Error(res, insertErr.Error(), http.StatusBadRequest)
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

// A BatchElement represent structure to marshal element of request`s json array.
type BatchElement struct {
	CorrelarionID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
	ShortURL      string `json:"short_url,omitempty"`
}

// CreateShorteningJSONBatch handle POST HTTP request with set of long URLs in body and retrieves set of shortenings.
// It handle only requests with content type application/json.
// Response has content type application/json.
// post /api/shorten/batch
func (sh *Shortener) CreateShorteningJSONBatch(res http.ResponseWriter, req *http.Request) {
	// set response content type
	res.Header().Set("Content-Type", "application/json")

	// check content type
	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" && contentType != "application/x-gzip" {
		http.Error(res, "Invalid content type", http.StatusBadRequest)
		return
	}

	// decode request body
	batch := make([]BatchElement, 0, 10)
	if err := json.NewDecoder(req.Body).Decode(&batch); err != nil {
		http.Error(res, "Can't read body", http.StatusBadRequest)
		return
	}
	if len(batch) == 0 {
		http.Error(res, "Body is empty", http.StatusBadRequest)
		return
	}
	// generate shortening
	for k := range batch {
		batch[k].ShortURL = generator.GenerateRandomString(15)
	}

	q := req.URL.Query()
	id, err := uuid.FromString(q.Get("userUUID"))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// write to data storage
	if err = sh.repo.InsertBatch(req.Context(), id, batch); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// make response
	for k, v := range batch {
		batch[k].ShortURL = sh.config.BaseURL + v.ShortURL
	}
	responseData, err := json.Marshal(batch)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write(responseData)
}

// GetFullString handle GET request with shortening in URL parameter named id
// and makes response with long URL in header's location value.
// Response content type is text/plain.
// get /{id}
func (sh *Shortener) GetFullString(res http.ResponseWriter, req *http.Request) {
	// parse parameter id from URL
	param := chi.URLParam(req, "id")
	if param == "" {
		http.Error(res, "Bad parameters", http.StatusBadRequest)
		return
	}

	// get long URL from repository
	repoOutput, err := sh.repo.Select(req.Context(), param)
	if err != nil {
		if errors.Is(err, sherr.ErrDBRecordDeleted) {
			http.Error(res, err.Error(), http.StatusGone)
		}
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// make responce
	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Location", repoOutput)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

// GetUserAllShortenings handle GET request with no parameters and makes response with
// all user's shortenings in body in json format.
// get /api/user/urls
func (sh *Shortener) GetUserAllShortenings(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	q := req.URL.Query()
	id, err := uuid.FromString(q.Get("userUUID"))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	//get all user's long URL from repository
	allRecords, err := sh.repo.SelectUserAll(req.Context(), id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if len(allRecords) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	for k, v := range allRecords {
		allRecords[k].ShortURL = sh.config.BaseURL + v.ShortURL
	}

	responseData, err := json.Marshal(allRecords)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Write(responseData)

	// make responce
	res.WriteHeader(http.StatusOK)
}

// DeleteRecordJSON saves record's id for future deletion. It returns status Accepted on seccuss saving.
// Deletion itself is performed periodically.
// It handle only requests with content type application/json.
// delete /api/user/urls
func (sh *Shortener) DeleteRecordJSON(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	q := req.URL.Query()
	id, err := uuid.FromString(q.Get("userUUID"))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// check content type
	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" && contentType != "application/x-gzip" {
		http.Error(res, "Invalid content type", http.StatusBadRequest)
		return
	}

	// decode request body
	recordIDs := make([]string, 10)
	if err := json.NewDecoder(req.Body).Decode(&recordIDs); err != nil {
		http.Error(res, "Can't read body", http.StatusBadRequest)
		return
	}
	if len(recordIDs) == 0 {
		http.Error(res, "Body is empty", http.StatusBadRequest)
		return
	}

	sh.deleteChan <- DeleteItem{IDs: recordIDs, UserID: id}

	logger.Log.Info("Shortenings' ids were send to chan for deletion")

	// make responce
	res.WriteHeader(http.StatusAccepted)
}

// PingDB check connection to data storage.
func (sh *Shortener) PingDB(res http.ResponseWriter, req *http.Request) {

	ctx, cancel := context.WithTimeout(req.Context(), timeoutPing*time.Second)
	defer cancel()
	if err := sh.repo.Ping(ctx); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// make responce
	res.WriteHeader(http.StatusOK)
}
