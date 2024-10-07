package main

import (
    "bufio"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "strings"
)

func main() {
    endpoint := "http://localhost:8181/"
    // контейнер данных для запроса
    data := url.Values{}
    // приглашение в консоли
    fmt.Println("Введите длинный URL")
    // открываем потоковое чтение из консоли
    reader := bufio.NewReader(os.Stdin)
    // читаем строку из консоли
    long, err := reader.ReadString('\n')
    if err != nil {
        panic(err)
    }
    long = strings.TrimSuffix(long, "\n")
    long = strings.TrimSuffix(long, "\r")
    // заполняем контейнер данными
    data.Set("url", long)
    // добавляем HTTP-клиент
    client := &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    }}
    // пишем запрос
    // запрос методом POST должен, помимо заголовков, содержать тело
    // тело должно быть источником потокового чтения io.Reader
    request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
    if err != nil {
        panic(err)
    }
    // в заголовках запроса указываем кодировку
    request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    // отправляем запрос и получаем ответ
    response, err := client.Do(request)
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

    splitResult:=strings.Split(string(body),"/")
    shortening:=splitResult[len(splitResult)-1]
    getrequest, err := http.NewRequest(http.MethodGet, endpoint+shortening, nil)
    if err != nil {
        panic(err)
    }
    getrequest.Header.Add("Content-Type", "text/plain")
    // отправляем запрос и получаем ответ
    origUrlResponse, err := client.Do(getrequest)
    if err != nil {
        panic(err)
    }
    fmt.Println("Статус-код ", origUrlResponse.Status)
    defer origUrlResponse.Body.Close()
    fmt.Println("Header Location ", origUrlResponse.Header.Get("Location"))
} 