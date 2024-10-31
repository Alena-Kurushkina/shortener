package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type resultResponse struct {
	Result string `json:"result"`
}

func main() {
	endpointApi := "http://localhost:8080/api/shorten"
	endpoint := "http://localhost:8080/"

	// добавляем HTTP-клиент
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	requestText, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(`http://ssite.ru`))
	if err != nil {
		panic(err)
	}
	// в заголовках запроса указываем кодировку
	requestText.Header.Add("Content-Type", "text/plain")

	//отправляем запрос и получаем ответ
	response, err := client.Do(requestText)
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

	//-------------

	request, err := http.NewRequest(http.MethodPost, endpointApi, strings.NewReader(`{"url": "http://some-site.ru"}`))
	if err != nil {
		panic(err)
	}

	request.Header.Add("Content-Type", "application/json")

	// отправляем запрос
	response, err = client.Do(request)
	if err != nil {
		panic(err)
	}

	// ответ
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	body, err = io.ReadAll(response.Body)
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

	getrequest, err := http.NewRequest(http.MethodGet, endpoint+shortening, nil)
	if err != nil {
		panic(err)
	}
	getrequest.Header.Add("Content-Type", "text/plain")

	// отправляем запрос
	origURLResponse, err := client.Do(getrequest)
	if err != nil {
		panic(err)
	}
	fmt.Println("Статус-код ", origURLResponse.Status)
	defer origURLResponse.Body.Close()
	fmt.Println("Header Location ", origURLResponse.Header.Get("Location"))

	// --------

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
	resp, err := client.Do(req)
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

	body, err = io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}
