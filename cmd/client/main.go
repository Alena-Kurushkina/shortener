package main

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
	endpointAPI = "http://localhost:8080/api/shorten"
	endpointAPIbatch = "http://localhost:8080/api/shorten/batch"
	endpointAPIselectAll = "http://localhost:8080/api/user/urls"
	endpoint = "http://localhost:8080/"
)

type resultResponse struct {
	Result string `json:"result"`
}

type batchElement struct {
	CorrelarionID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url"`
	ShortURL      string `json:"short_url,omitempty"`
}
type ShClient struct {
	client *http.Client
}

func NewClient() ShClient {
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
	if resp.StatusCode!=http.StatusOK{
		panic("Ping failed")
	}
	return ShClient{client: client}
}

func (cl *ShClient) PostTextPlainRequest(){
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

func (cl ShClient) PostJSONRequest() string {
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

func (cl ShClient) GetTextPlainRequest(shortening string){
	getrequest, err := http.NewRequest(http.MethodGet, endpoint+shortening, nil)
	if err != nil {
		panic(err)
	}
	getrequest.Header.Add("Content-Type", "text/plain")

	// отправляем запрос
	origURLResponse, err := cl.client.Do(getrequest)
	if err != nil {
		panic(err)
	}
	fmt.Println("Статус-код ", origURLResponse.Status)
	defer origURLResponse.Body.Close()
	fmt.Println("Header Location ", origURLResponse.Header.Get("Location"))
}

func (cl ShClient) PostJSONBatchRequest(){
	request, err := http.NewRequest(http.MethodPost, endpointAPIbatch, strings.NewReader(`[{"correlation_id":"dfgh345","original_url": "http://some-site.ru"},{"correlation_id":"kjhg1234","original_url": "http://testsite.ru"}]`))
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

	rr1 := make([]batchElement, 0, 10)
	err = json.Unmarshal(body, &rr1)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response body", rr1)
}

func (cl ShClient) PostGzipRequest(){
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

type Claims struct {
    jwt.RegisteredClaims
    UserID uuid.UUID
}

const TOKEN_EXP = time.Hour * 3
// TODO перенести в env
const SECRET_KEY = "secretkey"

func BuildJWTString(id uuid.UUID) (string, error) {
    // создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims {
        RegisteredClaims: jwt.RegisteredClaims{
            // когда создан токен
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
        },
        // собственное утверждение
        UserID: id,
    })
	_=id
    // создаём строку токена
    tokenString, err := token.SignedString([]byte(SECRET_KEY))
    if err != nil {
        return "", err
    }

    // возвращаем строку токена
    return tokenString, nil
} 

func (cl *ShClient) PostJWTTextPlainRequest(){
	id:=uuid.NewV4()
	token, err:=BuildJWTString(id)
	if err != nil {
		panic(err)
	}

	cookie := &http.Cookie{
        Name:   "token",
        Value:  token,
        MaxAge: 300,
    }

	requestText, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(`http://ssite.ru`))
	if err != nil {
		panic(err)
	}
	// в заголовках запроса указываем кодировку
	requestText.Header.Add("Content-Type", "text/plain")

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

func (cl ShClient) GetJSONBatchRequest(id uuid.UUID){
	request, err := http.NewRequest(http.MethodGet, endpointAPIselectAll, nil)
	if err != nil {
		panic(err)
	}

	request.Header.Add("Content-Type", "application/json")

	token, err:=BuildJWTString(id)
	if err!=nil{
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

	if response.StatusCode!=http.StatusNoContent{
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

func (cl ShClient) DeleteRequest(ids []int){
	var idStr = []string{}
	for _,v:=range ids {
		idStr=append(idStr, string(v))
	}
	param:=`[`+strings.Join(idStr,",")+`]`

	request, err := http.NewRequest(http.MethodDelete, endpointAPIselectAll, strings.NewReader(param))
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
}

func main() {
	cl:=NewClient()
	
	//cl.PostTextPlainRequest()

	//shortening:=cl.PostJSONRequest()
	cl.GetTextPlainRequest("jkafgeyh")

	//cl.PostJSONBatchRequest()

	// cl.PostGzipRequest()

	//cl.PostJWTTextPlainRequest()

	//cl.GetJSONBatchRequest(uuid.FromStringOrNil("2860ca35-5859-4cdd-8662-fe52de9fc4b1"))
}
