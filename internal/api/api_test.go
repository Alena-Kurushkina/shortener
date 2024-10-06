package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/repository"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func TestShortener_CreateShortening(t *testing.T) {
// 	repo := repository.NewRepository()
// 	sh := NewShortener(repo)

// 	type want struct {
// 		code        int
// 		response    string
// 		contentType string
// 	}
// 	tests := []struct {
// 		name    string
// 		request string
// 		contentType string
// 		body    string
// 		want    want
// 	}{
// 		{
// 			name:    "positive create shortening test",
// 			request: "/",
// 			contentType: "text/plain; charset=utf-8",
// 			body:    "http://site.ru",
// 			want: want{
// 				code:        http.StatusCreated,
// 				response:    "EwHXdJfB",
// 				contentType: "text/plain",
// 			},
// 		},
// 		{
// 			name:    "negaitve create shortening test",
// 			request: "/",
// 			contentType: "text/plain",
// 			body:    "",
// 			want: want{
// 				code:        http.StatusBadRequest,
// 				response:    "",
// 				contentType: "text/plain",
// 			},
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			request := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(test.body))
// 			request.Header.Add("Content-Type", test.contentType)
// 			w := httptest.NewRecorder()
// 			sh.CreateShortening(w, request)

// 			res := w.Result()
// 			assert.Equal(t, test.want.code, res.StatusCode, "Код статуса ответа не совпадает с ожидаемым")
// 			if test.want.code == http.StatusBadRequest {
// 				return
// 			}
// 			defer res.Body.Close()
// 			resBody, err := io.ReadAll(res.Body)
// 			require.NoError(t, err)
// 			assert.Equal(t, test.want.response, string(resBody), "Тело ответа не совпадает с ожидаемым")
// 			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"), "Content-Type не совпадает с ожидаемым")
// 		})
// 	}
// }

// func TestShortener_GetFullString(t *testing.T) {
// 	repo:=repository.Repository(map[string]string{"EwddTjks":"http://testsource.ru"})
// 	// repo := repository.NewRepository()
// 	sh := NewShortener(&repo)

// 	type want struct {
// 		code        int
// 		location    string
// 		contentType string
// 	}
// 	tests := []struct {
// 		name    string
// 		request string
// 		body    string
// 		want    want
// 	}{
// 		{
// 			name:    "positive get full string test",
// 			request: "/EwddTjks",
// 			want: want{
// 				code:        http.StatusTemporaryRedirect,
// 				location:    "http://testsource.ru",
// 				contentType: "text/plain",
// 			},
// 		},
// 		{
// 			name:    "negative get full string test(full string is not found)",
// 			request: "/jfhdgt",
// 			want: want{
// 				code:        http.StatusBadRequest,
// 				location:    "",
// 				contentType: "",
// 			},
// 		},
// 		{
// 			name:    "negative get full string test(no shortening specified)",
// 			request: "/",
// 			want: want{
// 				code:        http.StatusBadRequest,
// 				location:    "",
// 				contentType: "",
// 			},
// 		},
// 		{
// 			name:    "negative get full string test(incorrect path)",
// 			request: "/EwddTjks/path",
// 			want: want{
// 				code:        http.StatusBadRequest,
// 				location:    "",
// 				contentType: "",
// 			},
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			request := httptest.NewRequest(http.MethodGet, test.request, nil)
// 			w := httptest.NewRecorder()
// 			sh.GetFullString(w, request)

// 			res := w.Result()
// 			defer res.Body.Close()
// 			assert.Equal(t, test.want.code, res.StatusCode, "Код статуса ответа не совпадает с ожидаемым")
// 			if test.want.code == http.StatusBadRequest{
// 				return
// 			}
// 			assert.Equal(t, test.want.location, res.Header.Get("Location"), "Тело ответа не совпадает с ожидаемым")
// 			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"), "Content-Type не совпадает с ожидаемым")
// 		})
// 	}
// }

type responseParams struct{
	statusCode int 
	location string
	respBody string
	contentType string
}

func testRequest(t *testing.T, ts *httptest.Server, reqMethod,	path string, contentType string, body string) responseParams {
	request, err := http.NewRequest(reqMethod, ts.URL+path, strings.NewReader(body))
	request.Header.Add("Content-Type", contentType)
	require.NoError(t, err)

	resp, err := ts.Client().Do(request)
    require.NoError(t, err)
    defer resp.Body.Close()

	rp:=responseParams{}

	rp.statusCode= resp.StatusCode

	bodyOutput, err := io.ReadAll(resp.Body)	
    require.NoError(t, err)
	rp.respBody=strings.TrimSuffix(string(bodyOutput),"\n")

	rp.location= resp.Header.Get("Location")

	rp.contentType=resp.Header.Get("Content-Type")

	return rp
}


func TestRouter(t *testing.T){
	repo:=repository.NewRepository()
	config:=config.InitConfig()
	sh:=NewShortener(repo,config)

	r:=chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/",sh.CreateShortening)
	r.Get("/{id}",sh.GetFullString)

	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		code        int
		response    string
		contentType string
		location string
	}

	tests := []struct {
		method string
		name    string
		path string
		contentType string
		body    string
		want    want
	}{
		{http.MethodPost, "positive create shortening test", "/", "text/plain", "http://site.ru", want{http.StatusCreated,*config.BaseUrl+"EwHXdJfB","text/plain",""}},
		{http.MethodPost, "negaitve create shortening test", "/", "text/plain", "", want{http.StatusBadRequest, "Body is empty", "text/plain",""}},
		{http.MethodGet, "positive get full string test", "/EwHXdJfB","text/plain","",want{http.StatusTemporaryRedirect, "", "text/plain", "http://site.ru"}},
		{http.MethodGet, "negative get full string test(full string is not found)", "/jfhdgt", "text/plain", "", want{ http.StatusBadRequest,"Full string is not found", "text/plain", ""}},
		{http.MethodGet, "negative get full string test(no shortening specified)", "/", "text/plain", "", want{http.StatusMethodNotAllowed, "", "", ""}},
		{http.MethodGet, "negative get full string test(incorrect path)", "/EwddTjks/path", "text/plain", "", want{http.StatusNotFound,"404 page not found","",""}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			rp:=testRequest(t, ts, test.method, test.path, test.contentType, test.body)

			assert.Equal(t, test.want.code, rp.statusCode, "Код статуса ответа не совпадает с ожидаемым")			
			assert.Equal(t, test.want.response, rp.respBody, "Тело ответа не совпадает с ожидаемым")
			assert.Contains(t, rp.contentType, test.want.contentType, "Content-Type не совпадает с ожидаемым")
			assert.Equal(t, test.want.location, rp.location, "Location не совпадает с ожидаемым")
		})
	}

}
