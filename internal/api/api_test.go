package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Alena-Kurushkina/shortener/internal/config"
)

var cfg *config.Config

type DBMock struct {
	db map[string]string
}

func (mock DBMock) Insert(_ context.Context, key, value string) error {
	mock.db[key] = value

	return nil
}

func (mock DBMock) Select(_ context.Context, key string) (string, error) {
	if v, ok := mock.db[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("can't find value of key")
}

func (mock DBMock) Close() {}
func (mock DBMock) Ping(_ context.Context) error { return nil}

type responseParams struct {
	statusCode  int
	location    string
	respBody    string
	contentType string
}

func testRequest(t *testing.T, ts *httptest.Server, reqMethod, path string, contentType string, body string) responseParams {
	request, err := http.NewRequest(reqMethod, ts.URL+path, strings.NewReader(body))
	request.Header.Add("Content-Type", contentType)
	require.NoError(t, err)

	resp, err := ts.Client().Do(request)
	require.NoError(t, err)
	defer resp.Body.Close()

	rp := responseParams{}

	rp.statusCode = resp.StatusCode

	bodyOutput, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	rp.respBody = strings.TrimSuffix(string(bodyOutput), "\n")

	rp.location = resp.Header.Get("Location")

	rp.contentType = resp.Header.Get("Content-Type")

	return rp
}

func TestRouter(t *testing.T) {
	repo := DBMock{
		db: make(map[string]string),
	}
	cfg = config.InitConfig()
	sh := NewShortener(repo, cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", sh.CreateShortening)
	r.Get("/{id}", sh.GetFullString)

	ts := httptest.NewServer(r)
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer ts.Close()

	type want struct {
		code        int
		response    string
		contentType string
		location    string
	}

	type testData struct {
		method      string
		name        string
		path        string
		contentType string
		body        string
		want        want
	}

	t.Run("positive short and expand test", func(t *testing.T) {
		testDataShort := testData{http.MethodPost, "", "/", "text/plain", "http://site.ru/somelongurl", want{http.StatusCreated, cfg.BaseURL, "text/plain", ""}}
		_ = testDataShort.name
		rp := testRequest(t, ts, testDataShort.method, testDataShort.path, testDataShort.contentType, testDataShort.body)

		splitResult := strings.Split(string(rp.respBody), "/")
		shortening := splitResult[len(splitResult)-1]

		assert.Equal(t, testDataShort.want.code, rp.statusCode, "Short URL: Код статуса ответа не совпадает с ожидаемым")
		assert.Contains(t, rp.respBody, testDataShort.want.response, "Short URL: Тело ответа не совпадает с ожидаемым")
		assert.Contains(t, rp.contentType, testDataShort.want.contentType, "Short URL: Content-Type не совпадает с ожидаемым")
		assert.Equal(t, testDataShort.want.location, rp.location, "Short URL: Location не совпадает с ожидаемым")

		testDataExpand := testData{http.MethodGet, "", "/" + shortening, "text/plain", "", want{http.StatusTemporaryRedirect, "", "text/plain", "http://site.ru/somelongurl"}}
		_ = testDataExpand.name
		rpGet := testRequest(t, ts, testDataExpand.method, testDataExpand.path, testDataExpand.contentType, testDataExpand.body)
		assert.Equal(t, testDataExpand.want.code, rpGet.statusCode, "Expand URL: Код статуса ответа не совпадает с ожидаемым")
		assert.Contains(t, rpGet.respBody, testDataExpand.want.response, "Expand URL: Тело ответа не совпадает с ожидаемым")
		assert.Contains(t, rpGet.contentType, testDataExpand.want.contentType, "Expand URL: Content-Type не совпадает с ожидаемым")
		assert.Equal(t, testDataExpand.want.location, rpGet.location, "Expand URL: Location не совпадает с ожидаемым")
	})

	tests := []testData{
		{http.MethodPost, "negative create shortening test", "/", "text/plain", "", want{http.StatusBadRequest, "Body is empty", "text/plain", ""}},
		{http.MethodGet, "negative get full string test(full string is not found)", "/jfhdgt", "text/plain", "", want{http.StatusBadRequest, "can't find value of key", "text/plain", ""}},
		{http.MethodGet, "negative get full string test(no shortening specified)", "/", "text/plain", "", want{http.StatusMethodNotAllowed, "", "", ""}},
		{http.MethodGet, "negative get full string test(incorrect path)", "/EwddTjks/path", "text/plain", "", want{http.StatusNotFound, "404 page not found", "", ""}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rp := testRequest(t, ts, test.method, test.path, test.contentType, test.body)

			assert.Equal(t, test.want.code, rp.statusCode, "Код статуса ответа не совпадает с ожидаемым")
			assert.Equal(t, test.want.response, rp.respBody, "Тело ответа не совпадает с ожидаемым")
			assert.Contains(t, rp.contentType, test.want.contentType, "Content-Type не совпадает с ожидаемым")
			assert.Equal(t, test.want.location, rp.location, "Location не совпадает с ожидаемым")
		})
	}

}

func TestRouterJSON(t *testing.T) {
	repo := DBMock{
		db: make(map[string]string),
	}
	sh := NewShortener(repo, cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/api/shorten", sh.CreateShorteningJSON)
	r.Get("/{id}", sh.GetFullString)

	ts := httptest.NewServer(r)
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer ts.Close()

	type want struct {
		code        int
		response    string
		contentType string
		location    string
	}

	type testData struct {
		method      string
		name        string
		path        string
		contentType string
		body        string
		want        want
	}

	t.Run("positive short and expand test", func(t *testing.T) {
		testDataShort := testData{http.MethodPost, "", "/api/shorten", "application/json", `{"url": "http://site.ru/somelongurl"}`, want{http.StatusCreated, cfg.BaseURL, "application/json", ""}}
		_ = testDataShort.name
		rp := testRequest(t, ts, testDataShort.method, testDataShort.path, testDataShort.contentType, testDataShort.body)

		rr := struct {
			Result string `json:"result"`
		}{}
		err := json.Unmarshal([]byte(rp.respBody), &rr)
		require.NoError(t, err)

		splitResult := strings.Split(string(rr.Result), "/")
		shortening := splitResult[len(splitResult)-1]

		assert.Equal(t, testDataShort.want.code, rp.statusCode, "Short URL: Код статуса ответа не совпадает с ожидаемым")
		assert.Contains(t, rp.respBody, testDataShort.want.response, "Short URL: Тело ответа не совпадает с ожидаемым")
		assert.Contains(t, rp.contentType, testDataShort.want.contentType, "Short URL: Content-Type не совпадает с ожидаемым")
		assert.Equal(t, testDataShort.want.location, rp.location, "Short URL: Location не совпадает с ожидаемым")

		testDataExpand := testData{http.MethodGet, "", "/" + shortening, "text/plain", "", want{http.StatusTemporaryRedirect, "", "text/plain", "http://site.ru/somelongurl"}}
		_ = testDataExpand.name
		rpGet := testRequest(t, ts, testDataExpand.method, testDataExpand.path, testDataExpand.contentType, testDataExpand.body)
		assert.Equal(t, testDataExpand.want.code, rpGet.statusCode, "Expand URL: Код статуса ответа не совпадает с ожидаемым")
		assert.Contains(t, rpGet.respBody, testDataExpand.want.response, "Expand URL: Тело ответа не совпадает с ожидаемым")
		assert.Contains(t, rpGet.contentType, testDataExpand.want.contentType, "Expand URL: Content-Type не совпадает с ожидаемым")
		assert.Equal(t, testDataExpand.want.location, rpGet.location, "Expand URL: Location не совпадает с ожидаемым")
	})

	tests := []testData{
		{http.MethodPost, "negative create shortening test (can't decode empty body)", "/api/shorten", "application/json", "", want{http.StatusBadRequest, "Can't read body", "text/plain", ""}},
		{http.MethodPost, "negative create shortening test (wrong content type)", "/api/shorten", "text/plain", "", want{http.StatusBadRequest, "Invalid content type", "text/plain", ""}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rp := testRequest(t, ts, test.method, test.path, test.contentType, test.body)

			assert.Equal(t, test.want.code, rp.statusCode, "Код статуса ответа не совпадает с ожидаемым")
			assert.Equal(t, test.want.response, rp.respBody, "Тело ответа не совпадает с ожидаемым")
			assert.Contains(t, rp.contentType, test.want.contentType, "Content-Type не совпадает с ожидаемым")
			assert.Equal(t, test.want.location, rp.location, "Location не совпадает с ожидаемым")
		})
	}

}
