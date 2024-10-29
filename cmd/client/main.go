package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type resultResponse struct {
	Result string `json:"result"`
}

func main() {
	endpointPost := "http://localhost:8080/api/shorten"
	endpointGet := "http://localhost:8080/"
	// контейнер данных для запроса
	// data := url.Values{}
	// приглашение в консоли
	// fmt.Println("Введите длинный URL")
	// // открываем потоковое чтение из консоли
	// reader := bufio.NewReader(os.Stdin)
	// // читаем строку из консоли
	// long, err := reader.ReadString('\n')
	// if err != nil {
	// 	panic(err)
	// }
	// long = strings.TrimSuffix(long, "\n")
	// long = strings.TrimSuffix(long, "\r")

	// заполняем контейнер данными
	// data.Set("url", long)

	// добавляем HTTP-клиент
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	// пишем запрос
	// запрос методом POST должен, помимо заголовков, содержать тело
	// тело должно быть источником потокового чтения io.Reader
	requestText, err := http.NewRequest(http.MethodPost, endpointGet, strings.NewReader(`http://ssite.ru`))
	if err != nil {
		panic(err)
	}
	// в заголовках запроса указываем кодировку
	requestText.Header.Add("Content-Type", "text/plain")
	// отправляем запрос и получаем ответ
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

	request, err := http.NewRequest(http.MethodPost, endpointPost, strings.NewReader(`{"url": "http://some-site.ru"}`)) //strings.NewReader(data.Encode())
	if err != nil {
		panic(err)
	}
	// в заголовках запроса указываем кодировку
	request.Header.Add("Content-Type", "application/json")
	// отправляем запрос и получаем ответ
	response, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	// выводим код ответа
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	// читаем поток из тела ответа
	body, err = io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	// и печатаем его
	fmt.Println(string(body))

	rr := resultResponse{}
	err = json.Unmarshal(body, &rr)
	if err != nil {
		panic(err)
	}

	splitResult := strings.Split(string(rr.Result), "/")
	shortening := splitResult[len(splitResult)-1]
	getrequest, err := http.NewRequest(http.MethodGet, endpointGet+shortening, nil)
	if err != nil {
		panic(err)
	}
	getrequest.Header.Add("Content-Type", "text/plain")
	// отправляем запрос и получаем ответ
	origURLResponse, err := client.Do(getrequest)
	if err != nil {
		panic(err)
	}
	fmt.Println("Статус-код ", origURLResponse.Status)
	defer origURLResponse.Body.Close()
	fmt.Println("Header Location ", origURLResponse.Header.Get("Location"))

	// --------

	var requestBody bytes.Buffer

	// Compress the request body
	gz := gzip.NewWriter(&requestBody)
	gz.Write([]byte("http://some-site-gzip.ru"))
	gz.Close()

	// Create an HTTP request with compressed body
	req, err := http.NewRequest(http.MethodPost, endpointGet, &requestBody)
	if err != nil {
		panic(err)
	}

	// Set the Content-Encoding header to gzip
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Add("Content-Type", "text/plain")

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("Статус-код ", resp.Status)
	defer resp.Body.Close()
	
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// и печатаем его
	fmt.Println(string(body))

}
