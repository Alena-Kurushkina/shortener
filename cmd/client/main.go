package main

//lint:file-ignore U1000 игнорируем неиспользуемый код, так как он нужен только при разработке

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
)

const (
	endpointAPI          = "http://localhost:8080/api/shorten"
	endpointAPIbatch     = "http://localhost:8080/api/shorten/batch"
	endpointAPIselectAll = "http://localhost:8080/api/user/urls"
	endpoint             = "http://localhost:8080/"
)

type resultResponse struct {
	Result string `json:"result"`
}

type batchElement struct {
	CorrelarionID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url"`
	ShortURL      string `json:"short_url,omitempty"`
}

// ShortenerClient represents http client to shortener service
type ShortenerClient struct {
	client *http.Client
}

func newClient() ShortenerClient {
	// добавляем HTTP-клиент
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	reqPing, err := http.NewRequest(http.MethodGet, endpoint+"ping", nil)
	if err != nil {
		panic(err)
	}
	//отправляем запрос и получаем ответ
	resp, err := client.Do(reqPing)
	if err != nil {
		panic(err)
	}

	// выводим код ответа
	fmt.Println("Статус-код ", resp.Status)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic("Ping failed")
	}
	return ShortenerClient{client: client}
}

func (cl *ShortenerClient) postTextPlainRequest() {
	requestText, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(`http://ssite.ru`))
	if err != nil {
		panic(err)
	}
	// в заголовках запроса указываем кодировку
	requestText.Header.Add("Content-Type", "text/plain")

	//отправляем запрос и получаем ответ
	response, err := cl.client.Do(requestText)
	if err != nil {
		panic(err)
	}

	// выводим код ответа
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()

	// читаем поток из тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	// и печатаем его
	fmt.Println(string(body))
}

func (cl ShortenerClient) postJSONRequest() string {
	request, err := http.NewRequest(http.MethodPost, endpointAPI, strings.NewReader(`{"url": "http://ssgreh.ru"}`))
	if err != nil {
		panic(err)
	}

	request.Header.Add("Content-Type", "application/json")

	// отправляем запрос
	response, err := cl.client.Do(request)
	if err != nil {
		panic(err)
	}

	// ответ
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))

	rr := resultResponse{}
	err = json.Unmarshal(body, &rr)
	if err != nil {
		panic(err)
	}

	splitResult := strings.Split(string(rr.Result), "/")
	shortening := splitResult[len(splitResult)-1]

	return shortening
}

func (cl ShortenerClient) getTextPlainRequest(id uuid.UUID, shortening string) {
	getrequest, err := http.NewRequest(http.MethodGet, endpoint+shortening, nil)
	if err != nil {
		panic(err)
	}
	getrequest.Header.Add("Content-Type", "text/plain")

	token, err := buildJWTString(id)
	if err != nil {
		panic(err)
	}
	cookie := &http.Cookie{
		Name:   "token",
		Value:  token,
		MaxAge: 300,
	}
	getrequest.AddCookie(cookie)

	// отправляем запрос
	origURLResponse, err := cl.client.Do(getrequest)
	if err != nil {
		panic(err)
	}
	fmt.Println("Статус-код ", origURLResponse.Status)
	defer origURLResponse.Body.Close()
	fmt.Println("Header Location ", origURLResponse.Header.Get("Location"))
}

func (cl ShortenerClient) postJSONBatchRequest(id uuid.UUID, param string) {
	request, err := http.NewRequest(http.MethodPost, endpointAPIbatch, strings.NewReader(param))
	if err != nil {
		panic(err)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept-Encoding", "identity")

	token, err := buildJWTString(id)
	if err != nil {
		panic(err)
	}
	cookie := &http.Cookie{
		Name:   "token",
		Value:  token,
		MaxAge: 300,
	}
	request.AddCookie(cookie)

	// отправляем запрос
	response, err := cl.client.Do(request)
	if err != nil {
		panic(err)
	}

	// ответ
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))

	rr1 := make([]batchElement, 0, 10)
	err = json.Unmarshal(body, &rr1)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response body", rr1)
}

func (cl ShortenerClient) postGzipRequest() {
	var requestBody bytes.Buffer

	// запрос с компрессией
	gz := gzip.NewWriter(&requestBody)
	gz.Write([]byte("http://some-site-gzip.ru"))
	gz.Close()

	req, err := http.NewRequest(http.MethodPost, endpoint, &requestBody)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	// отправить запрос
	resp, err := cl.client.Do(req)
	if err != nil {
		panic(err)
	}

	// ответ
	fmt.Println("Статус-код ", resp.Status)
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}

type claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

const tokenExp = time.Hour * 3

// TODO перенести в env
const secretKey = "secretkey"

func buildJWTString(id uuid.UUID) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		// собственное утверждение
		UserID: id,
	})
	_ = id
	// создаём строку токена
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func (cl *ShortenerClient) postJWTTextPlainRequest() {
	requestText, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(`http://ssite.ru`))
	if err != nil {
		panic(err)
	}
	// в заголовках запроса указываем кодировку
	requestText.Header.Add("Content-Type", "text/plain")

	id := uuid.NewV4()
	token, err := buildJWTString(id)
	if err != nil {
		panic(err)
	}

	cookie := &http.Cookie{
		Name:   "token",
		Value:  token,
		MaxAge: 300,
	}

	requestText.AddCookie(cookie)

	//отправляем запрос и получаем ответ
	response, err := cl.client.Do(requestText)
	if err != nil {
		panic(err)
	}

	// выводим код ответа
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()

	// читаем поток из тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	// и печатаем его
	fmt.Println(string(body))
}

func (cl ShortenerClient) getJSONBatchRequest(id uuid.UUID) {
	request, err := http.NewRequest(http.MethodGet, endpointAPIselectAll, nil)
	if err != nil {
		panic(err)
	}

	request.Header.Add("Content-Type", "application/json")

	token, err := buildJWTString(id)
	if err != nil {
		panic(err)
	}
	request.AddCookie(&http.Cookie{Name: "token", Value: token, MaxAge: 0})

	// отправляем запрос
	response, err := cl.client.Do(request)
	if err != nil {
		panic(err)
	}

	// ответ
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			panic(err)
		}
		defer reader.Close()

		body, err := io.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(body))

		rr1 := make([]batchElement, 0, 10)
		err = json.Unmarshal(body, &rr1)
		if err != nil {
			panic(err)
		}

		fmt.Println("Response body", rr1)
	}
}

func (cl ShortenerClient) deleteRequest(ids []string, id uuid.UUID) {
	// var idStr = []string{}
	// for _,v:=range ids {
	// 	idStr=append(idStr, strconv.Itoa(v))
	// }
	param, err := json.Marshal(ids)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest(http.MethodDelete, endpointAPIselectAll, bytes.NewBuffer(param))
	if err != nil {
		panic(err)
	}

	request.Header.Add("Content-Type", "application/json")

	token, err := buildJWTString(id)
	if err != nil {
		panic(err)
	}
	request.AddCookie(&http.Cookie{Name: "token", Value: token, MaxAge: 0})

	// отправляем запрос
	response, err := cl.client.Do(request)
	if err != nil {
		panic(err)
	}

	// ответ
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

func main() {
	cl := newClient()

	//cl.postTextPlainRequest()

	//shortening:=cl.postJSONRequest()
	//cl.getTextPlainRequest("0yofl1hsoCo3nlK")

	id1 := uuid.NewV4()
	cl.postJSONBatchRequest(id1, `[{"correlation_id":"8f4f4159-85d2-4aa6-bce8-4d9eb249c01b","original_url":"http://uk8d4ovutebb2.ru"},{"correlation_id":"450cffae-147a-4653-8b91-4b3c2e06df30","original_url":"http://yq1xxhwihp4l1.net/jelbsck49bdkp"}]`)
	// id2:=uuid.NewV4()
	// cl.postJSONBatchRequest(id2,`[{"correlation_id":"23456tg","original_url": "http://so-site.ru"},{"correlation_id":"sghgrh4","original_url": "http://tmdssujh.ru"}]`)

	// cl.postGzipRequest()

	//cl.postJWTTextPlainRequest()

	//cl.getJSONBatchRequest(uuid.FromStringOrNil("2240318c-b936-4795-b8e5-82d421142fc4"))

	//cl.deleteRequest([]string{"d70561a2addfe213ca3"}, uuid.FromStringOrNil("56b4fc0f-406b-48f7-9026-aa8b685762d6"))
	//cl.getTextPlainRequest(uuid.FromStringOrNil("56b4fc0f-406b-48f7-9026-aa8b685762d6"), "tUfUTzrkrFIyZoI")
}
