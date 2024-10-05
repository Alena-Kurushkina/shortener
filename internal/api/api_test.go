package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Alena-Kurushkina/shortener/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortener_CreateShortening(t *testing.T) {
	repo := repository.NewRepository()
	sh := NewShortener(repo)

	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name    string
		request string
		body    string
		want    want
	}{
		{
			name:    "positive create shortening test",
			request: "/",
			body:    "http://site.ru",
			want: want{
				code:        http.StatusCreated,
				response:    `EwHXdJfB`,
				contentType: "text/plain",
			},
		},
		{
			name:    "negaitve create shortening test",
			request: "/",
			body:    "",
			want: want{
				code:        http.StatusBadRequest,
				response:    "",
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			sh.HandleRequest(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode, "Код статуса ответа не совпадает с ожидаемым")
			if test.want.code == http.StatusBadRequest {
				return
			}
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody), "Тело ответа не совпадает с ожидаемым")
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"), "Content-Type не совпадает с ожидаемым")
		})
	}
}

func TestShortener_GetFullString(t *testing.T) {
	repo:=repository.Repository(map[string]string{"EwddTjks":"http://testsource.ru"})
	// repo := repository.NewRepository()
	sh := NewShortener(&repo)

	type want struct {
		code        int
		location    string
		contentType string
	}
	tests := []struct {
		name    string
		request string
		body    string
		want    want
	}{
		{
			name:    "positive get full string test",
			request: "/EwddTjks",			
			want: want{
				code:        http.StatusTemporaryRedirect,
				location:    "http://testsource.ru",
				contentType: "text/plain",
			},
		},
		{
			name:    "negative get full string test(full string is not found)",
			request: "/jfhdgt",			
			want: want{
				code:        http.StatusBadRequest,
				location:    "",
				contentType: "",
			},
		},
		{
			name:    "negative get full string test(no shortening specified)",
			request: "/",			
			want: want{
				code:        http.StatusBadRequest,
				location:    "",
				contentType: "",
			},
		},
		{
			name:    "negative get full string test(incorrect path)",
			request: "/EwddTjks/path",			
			want: want{
				code:        http.StatusBadRequest,
				location:    "",
				contentType: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.request, nil)
			w := httptest.NewRecorder()
			sh.HandleRequest(w, request)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, test.want.code, res.StatusCode, "Код статуса ответа не совпадает с ожидаемым")
			if test.want.code == http.StatusBadRequest{
				return
			}
			assert.Equal(t, test.want.location, res.Header.Get("Location"), "Тело ответа не совпадает с ожидаемым")
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"), "Content-Type не совпадает с ожидаемым")
		})
	}
}
