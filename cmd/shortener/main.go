package main

import (
	"github.com/Alena-Kurushkina/shortener/internal/shortener"
)

func main() {
	server := shortener.NewServer()
	server.Run()
}
