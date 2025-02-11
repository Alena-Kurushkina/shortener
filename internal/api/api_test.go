package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Alena-Kurushkina/shortener/internal/authenticator"
	"github.com/Alena-Kurushkina/shortener/internal/config"
)

var cfg *config.Config

type (
	responseParams struct {
		statusCode  int
		location    string
		respBody    string
		contentType string
	}

	want struct {
		code        int
		response    string
		contentType string
		location    string
	}

	testData struct {
		method      string
		name        string
		path        string
		contentType string
		body        string
		want        want
	}

	batchElement struct {
		CorrelarionID string `json:"correlation_id,omitempty"`
		OriginalURL   string `json:"original_url"`
		ShortURL      string `json:"short_url,omitempty"`
	}
)

func testRequest(t *testing.T, ts *httptest.Server, reqMethod, path string, contentType string, body string) responseParams {
	t.Helper()
	request, err := http.NewRequest(reqMethod, ts.URL+path, strings.NewReader(body))
	request.Header.Add("Content-Type", contentType)
	require.NoError(t, err)

	resp, err := ts.Client().Do(request)
	require.NoError(t, err)

	defer func() {
		tempErr:=resp.Body.Close()
		require.NoError(t, tempErr)
	}()

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
	ctrl := gomock.NewController(t)
	m := NewMockStorager(ctrl)

	m.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	cfg = config.InitConfig()
	sh := NewShortener(m, cfg)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/{id}", sh.GetFullString)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(authenticator.AuthMiddleware)
		r.Post("/", sh.CreateShortening)
	})

	ts := httptest.NewServer(r)
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer ts.Close()

	t.Run("positive short and expand test", func(t *testing.T) {
		testDataShort := testData{
			http.MethodPost,
			"",
			"/",
			"text/plain",
			"http://site.ru/somelongurl",
			want{
				http.StatusCreated,
				cfg.BaseURL,
				"text/plain",
				"",
			},
		}
		_ = testDataShort.name

		rp := testRequest(t, ts, testDataShort.method, testDataShort.path, testDataShort.contentType, testDataShort.body)

		splitResult := strings.Split(string(rp.respBody), "/")
		shortening := splitResult[len(splitResult)-1]

		assert.Equal(t, testDataShort.want.code, rp.statusCode, "Short URL: Код статуса ответа не совпадает с ожидаемым")
		assert.Contains(t, rp.respBody, testDataShort.want.response, "Short URL: Тело ответа не совпадает с ожидаемым")
		assert.Contains(t, rp.contentType, testDataShort.want.contentType, "Short URL: Content-Type не совпадает с ожидаемым")
		assert.Equal(t, testDataShort.want.location, rp.location, "Short URL: Location не совпадает с ожидаемым")

		testDataExpand := testData{
			http.MethodGet,
			"",
			"/" + shortening,
			"text/plain",
			"",
			want{
				http.StatusTemporaryRedirect,
				"",
				"text/plain",
				"http://site.ru/somelongurl",
			},
		}
		_ = testDataExpand.name

		m.EXPECT().Select(gomock.Any(), shortening).Return("http://site.ru/somelongurl", nil)

		rpGet := testRequest(t, ts, testDataExpand.method, testDataExpand.path, testDataExpand.contentType, testDataExpand.body)
		assert.Equal(t, testDataExpand.want.code, rpGet.statusCode, "Expand URL: Код статуса ответа не совпадает с ожидаемым")
		assert.Contains(t, rpGet.respBody, testDataExpand.want.response, "Expand URL: Тело ответа не совпадает с ожидаемым")
		assert.Contains(t, rpGet.contentType, testDataExpand.want.contentType, "Expand URL: Content-Type не совпадает с ожидаемым")
		assert.Equal(t, testDataExpand.want.location, rpGet.location, "Expand URL: Location не совпадает с ожидаемым")
	})

	m.EXPECT().Select(gomock.Any(), "jfhdgt").Return("", errors.New("can't find value of key"))

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
	ctrl := gomock.NewController(t)
	m := NewMockStorager(ctrl)

	m.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	//m.EXPECT().InsertBatch(gomock.Any(),gomock.Any(),gomock.Any()).Return(nil)

	sh := NewShortener(m, cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/{id}", sh.GetFullString)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(authenticator.AuthMiddleware)
		r.Post("/api/shorten", sh.CreateShorteningJSON)
	})

	ts := httptest.NewServer(r)
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer ts.Close()

	t.Run("positive short and expand test", func(t *testing.T) {
		testDataShort := testData{
			method:      http.MethodPost,
			name:        "",
			path:        "/api/shorten",
			contentType: "application/json",
			body:        `{"url": "http://site.ru/somelongurl"}`,
			want: want{
				code:        http.StatusCreated,
				response:    cfg.BaseURL,
				contentType: "application/json",
				location:    "",
			},
		}
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

		m.EXPECT().Select(gomock.Any(), shortening).Return("http://site.ru/somelongurl", nil)

		testDataExpand := testData{
			method:      http.MethodGet,
			name:        "",
			path:        "/" + shortening,
			contentType: "text/plain",
			body:        "",
			want: want{
				code:        http.StatusTemporaryRedirect,
				response:    "",
				contentType: "text/plain",
				location:    "http://site.ru/somelongurl",
			},
		}
		_ = testDataExpand.name
		rpGet := testRequest(t, ts, testDataExpand.method, testDataExpand.path, testDataExpand.contentType, testDataExpand.body)
		assert.Equal(t, testDataExpand.want.code, rpGet.statusCode, "Expand URL: Код статуса ответа не совпадает с ожидаемым")
		assert.Contains(t, rpGet.respBody, testDataExpand.want.response, "Expand URL: Тело ответа не совпадает с ожидаемым")
		assert.Contains(t, rpGet.contentType, testDataExpand.want.contentType, "Expand URL: Content-Type не совпадает с ожидаемым")
		assert.Equal(t, testDataExpand.want.location, rpGet.location, "Expand URL: Location не совпадает с ожидаемым")
	})

	tests := []testData{
		{
			http.MethodPost,
			"negative create shortening test (can't decode empty body)",
			"/api/shorten",
			"application/json",
			"",
			want{
				http.StatusBadRequest,
				"Can't read body",
				"text/plain",
				"",
			},
		},
		{
			http.MethodPost,
			"negative create shortening test (wrong content type)",
			"/api/shorten",
			"text/plain",
			"",
			want{
				http.StatusBadRequest,
				"Invalid content type",
				"text/plain",
				"",
			},
		},
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

func TestRouterJSONBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockStorager(ctrl)

	m.EXPECT().InsertBatch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	sh := NewShortener(m, cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/{id}", sh.GetFullString)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(authenticator.AuthMiddleware)
		r.Post("/api/shorten/batch", sh.CreateShorteningJSONBatch)
	})

	ts := httptest.NewServer(r)
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer ts.Close()

	type batchElem struct {
		correlarionID string
		originalURL   string
		baseURL       string
	}
	type wantBatch struct {
		code        int
		contentType string
		batchElems  []batchElem
	}

	t.Run("positive shorten and expand test", func(t *testing.T) {
		testDataShort := testData{
			method:      http.MethodPost,
			name:        "",
			path:        "/api/shorten/batch",
			contentType: "application/json",
			body:        `[{"correlation_id":"dfgh345","original_url": "http://some-site.ru"},{"correlation_id":"kjhg1234","original_url": "http://testsite.ru"}]`,
			want:        want{},
		}
		wantResult := wantBatch{
			code:        http.StatusCreated,
			contentType: "application/json",
			batchElems: []batchElem{
				{
					correlarionID: "dfgh345",
					originalURL:   "http://some-site.ru",
					baseURL:       cfg.BaseURL,
				},
				{
					correlarionID: "kjhg1234",
					originalURL:   "http://testsite.ru",
					baseURL:       cfg.BaseURL,
				},
			},
		}
		_ = testDataShort.name
		rp := testRequest(t, ts, testDataShort.method, testDataShort.path, testDataShort.contentType, testDataShort.body)

		rr1 := make([]batchElement, 0, 10)
		err := json.Unmarshal([]byte(rp.respBody), &rr1)
		require.NoError(t, err)

		assert.Equal(t, wantResult.code, rp.statusCode, "Short URL: Код статуса ответа не совпадает с ожидаемым")
		assert.Contains(t, rp.contentType, wantResult.contentType, "Short URL: Content-Type не совпадает с ожидаемым")

		for k, v := range wantResult.batchElems {
			assert.Equal(t, v.correlarionID, rr1[k].CorrelarionID)
			assert.Equal(t, v.originalURL, rr1[k].OriginalURL)
			assert.Contains(t, rr1[k].ShortURL, v.baseURL)
		}
	})

}
