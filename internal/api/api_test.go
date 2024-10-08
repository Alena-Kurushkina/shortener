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
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
  		return http.ErrUseLastResponse
	}
	defer ts.Close()

	type want struct {
		code        int
		response    string
		contentType string
		location string
	}

	type testData struct {
		method string
		name    string
		path string
		contentType string
		body    string
		want    want
	}

	t.Run("positive short and expand test", func(t *testing.T) {
		testDataShort:=	testData{http.MethodPost,"", "/", "text/plain", "http://site.ru/somelongurl", want{http.StatusCreated, config.BaseUrl,"text/plain",""}}			
		_=testDataShort.name
		rp:=testRequest(t, ts, testDataShort.method, testDataShort.path, testDataShort.contentType, testDataShort.body)

		splitResult:=strings.Split(string(rp.respBody),"/")
		shortening:=splitResult[len(splitResult)-1]

		assert.Equal(t, testDataShort.want.code, rp.statusCode, "Short URL: Код статуса ответа не совпадает с ожидаемым")			
		assert.Contains(t, rp.respBody, testDataShort.want.response, "Short URL: Тело ответа не совпадает с ожидаемым")
		assert.Contains(t, rp.contentType, testDataShort.want.contentType, "Short URL: Content-Type не совпадает с ожидаемым")
		assert.Equal(t, testDataShort.want.location, rp.location, "Short URL: Location не совпадает с ожидаемым")		
		
		testDataExpand:=testData{http.MethodGet,"", "/" + shortening,"text/plain","",want{http.StatusTemporaryRedirect, "", "text/plain", "http://site.ru/somelongurl"}}
		_=testDataExpand.name
		rpGet:=testRequest(t, ts, testDataExpand.method, testDataExpand.path, testDataExpand.contentType, testDataExpand.body)
		assert.Equal(t, testDataExpand.want.code, rpGet.statusCode, "Expand URL: Код статуса ответа не совпадает с ожидаемым")			
		assert.Contains(t, rpGet.respBody, testDataExpand.want.response, "Expand URL: Тело ответа не совпадает с ожидаемым")
		assert.Contains(t, rpGet.contentType, testDataExpand.want.contentType, "Expand URL: Content-Type не совпадает с ожидаемым")
		assert.Equal(t, testDataExpand.want.location, rpGet.location, "Expand URL: Location не совпадает с ожидаемым")
	})

	tests := []testData {
		{http.MethodPost, "negaitve create shortening test", "/", "text/plain", "", want{http.StatusBadRequest, "Body is empty", "text/plain",""}},
		{http.MethodGet, "negative get full string test(full string is not found)", "/jfhdgt", "text/plain", "", want{ http.StatusBadRequest, "Can't find value of key", "text/plain", ""}},
		{http.MethodGet, "negative get full string test(no shortening specified)", "/", "text/plain", "", want{http.StatusMethodNotAllowed, "", "", ""}},
		{http.MethodGet, "negative get full string test(incorrect path)", "/EwddTjks/path", "text/plain", "", want{http.StatusNotFound, "404 page not found","",""}},
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
