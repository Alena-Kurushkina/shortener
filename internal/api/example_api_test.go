package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
)

const (
	endpointAPI = "/api/shorten"

	secretKey = "hgdfsjjhsu7643"
)

type resultResponse struct {
	Result string `json:"result"`
}

type claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

func buildJWTString(id uuid.UUID) (string, error) {	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
		},		
		UserID: id,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func makeCookie(uuid uuid.UUID) *http.Cookie{
	token, err := buildJWTString(uuid)
	if err != nil {
		panic(err)
	}
	cookie := &http.Cookie{
		Name:   "token",
		Value:  token,
		MaxAge: 300,
	}
	return cookie
}

func Example() {
	// make client which will send requests
	client := &http.Client{
		// stop redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
	}}

	// get config parameters
	cfg:=config.InitConfig()

	// make POST request to endpoint /api/shorten with long url in JSON format
	request, err := http.NewRequest(
		http.MethodPost, 
		"http://"+cfg.ServerAddress+endpointAPI, 
		strings.NewReader(`{"url": "http://site.example.ru/long/long/long/long/long/url"}`),
	)
	if err != nil {
		fmt.Println("Error while make request", err.Error())
		return
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept-Encoding", "identity")

	// generate id and place it in cookie 
	id:=uuid.NewV4()
	cookie:=makeCookie(id)
	request.AddCookie(cookie)

	// send request
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error while send request", err.Error())
		return
	}

	// read response
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error while read response", err.Error())
		return
	}
	response.Body.Close()

	// unmarshal response
	rr := resultResponse{}
	err = json.Unmarshal(body, &rr)
	if err != nil {
		fmt.Println("Error while decoding response", err.Error())
		return
	}
	
	// get shortening
	shortening, _ := strings.CutPrefix(rr.Result, cfg.BaseURL)

	// make GET request to endpoint /{shortening}
	getrequest, err := http.NewRequest(http.MethodGet, "http://"+cfg.ServerAddress+"/"+shortening, nil)
	if err != nil {
		fmt.Println("Error while make request", err.Error())
		return
	}
	getrequest.Header.Add("Content-Type", "text/plain")

	//place cookie in request
	getrequest.AddCookie(cookie)

	// make request
	origURLResponse, err := client.Do(getrequest)
	if err != nil {
		fmt.Println("Error while make request", err.Error())
		return
	}

	//read response
	fmt.Println("Status", origURLResponse.Status)
	defer origURLResponse.Body.Close()
	fmt.Println("Long URL", origURLResponse.Header.Get("Location"))
}