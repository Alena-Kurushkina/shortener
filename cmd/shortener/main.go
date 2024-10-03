package main

import (
	"github.com/Alena-Kurushkina/shortener/internal/shortener"
)

// type Storage map[string]string

// func (s Storage) Insert(key, value string) {
// 	s[key]=value
// }

// func (s Storage) Select(key string) (value string) {
// 	return s[key]
// }

// func shorten(res http.ResponseWriter, req *http.Request){
// 	// if req.Method!=http.MethodPost{
// 	// 	http.Error(res, "Only POST requests are allowed", http.StatusBadRequest)
// 	// 	return
// 	// }
// 	body, err :=io.ReadAll(req.Body)
// 	if err!=nil{
// 		http.Error(res, "Can't read body", http.StatusBadRequest)
// 		return
// 	}	
// 	if len(body)==0{
// 		http.Error(res, "Body is empty", http.StatusBadRequest)
// 		return
// 	}	
// 	st.Insert("EwHXdJfB")
// 	res.Header().Set("content-type", "text/plain")
// 	res.WriteHeader(http.StatusCreated)
// 	res.Write([]byte("EwHXdJfB"))
// }

// func origin(res http.ResponseWriter, req *http.Request){
// 	// if req.Method!=http.MethodGet{
// 	// 	http.Error(res, "Only GET requests are allowed", http.StatusBadRequest)
// 	// 	return
// 	// }
// 	param:=req.PathValue("id")
// 	if param==""{
// 		http.Error(res, "Bad parameters", http.StatusBadRequest)
// 		return
// 	}
// 	res.Header().Set("content-type", "text/plain")
// 	res.WriteHeader(http.StatusTemporaryRedirect)
// 	res.Write([]byte(st.Select(param)))
// }

func main() {

	// st:= make(Storage)

	// mux:=http.NewServeMux()
	// mux.HandleFunc(`POST /`, shorten)
	// mux.HandleFunc(`GET /{id}`, origin)

	// err:=http.ListenAndServe(`:8080`, mux)
	// if err!=nil{
	// 	panic(err)
	// }
	

	server:=shortener.NewServer()
	server.Run()
}
