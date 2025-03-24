// Package api implements handler functions to shorten long URL
// and expand shortenings back to long URL.
package api

import (
	context "context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	uuid "github.com/satori/go.uuid"

	"github.com/Alena-Kurushkina/shortener/internal/core"
	"github.com/Alena-Kurushkina/shortener/internal/logger"
	"github.com/Alena-Kurushkina/shortener/internal/sherr"
)

const timeoutPing time.Duration = 30

// A Shortener realises handlers.
type Shortener struct {
	Core *core.ShortenerCore
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

	q := req.URL.Query()
	id, err := uuid.FromString(q.Get("userUUID"))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// generate shortening
	shortStr, err := sh.Core.CreateShortening(req.Context(), id, url)

	var existError *sherr.AlreadyExistError
	if errors.As(err, &existError) {
		// make response
		res.WriteHeader(http.StatusConflict)
		res.Write([]byte(shortStr))
		return
	}
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// make response
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortStr))
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
	shortStr, insertErr := sh.Core.CreateShortening(req.Context(), id, url.URL)

	// make response
	responseData, err := json.Marshal(ResultResponse{
		Result: shortStr,
	})
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	var existError *sherr.AlreadyExistError
	if errors.As(insertErr, &existError) {
		// make response
		res.WriteHeader(http.StatusConflict)
		res.Write(responseData)
		return
	}
	if insertErr != nil {
		http.Error(res, insertErr.Error(), http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusCreated)
	res.Write(responseData)
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
	batch := make([]core.BatchElement, 0, 10)
	if err := json.NewDecoder(req.Body).Decode(&batch); err != nil {
		http.Error(res, "Can't read body", http.StatusBadRequest)
		return
	}
	if len(batch) == 0 {
		http.Error(res, "Body is empty", http.StatusBadRequest)
		return
	}

	q := req.URL.Query()
	id, err := uuid.FromString(q.Get("userUUID"))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	batch, err = sh.Core.CreateShorteningBatch(req.Context(), id, batch)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
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
	repoOutput, err := sh.Core.GetFullString(req.Context(), param)
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
	allRecords, err := sh.Core.GetAllUserShortenings(req.Context(), id)
	if err != nil {
		if errors.Is(err, sherr.ErrNoShortenings) {
			http.Error(res, err.Error(), http.StatusNoContent)
		}
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
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

	// sh.deleteChan <- DeleteItem{IDs: recordIDs, UserID: id}
	sh.Core.RegisterToDelete(req.Context(), recordIDs, id)

	logger.Log.Info("Shortenings' ids were send to chan for deletion")

	// make responce
	res.WriteHeader(http.StatusAccepted)
}

// PingDB check connection to data storage.
func (sh *Shortener) PingDB(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), timeoutPing*time.Second)
	defer cancel()
	if err := sh.Core.PingDB(ctx); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// make responce
	res.WriteHeader(http.StatusOK)
}

// GetStats handle GET request with no parameters and makes response with
// number of shorten URLs and users number.
// get /api/internal/stats
func (sh *Shortener) GetStats(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	ipStr := req.Header.Get("X-Real-IP")
	trusted, err := sh.Core.IsTrustedSubnet(ipStr)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if !trusted {
		res.WriteHeader(http.StatusForbidden)
		return
	}

	stats, err := sh.Core.GetStats(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData, err := json.Marshal(stats)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Write(responseData)

	// make responce
	res.WriteHeader(http.StatusOK)
}
