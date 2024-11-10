// Package repository implements routines for manipulating data source
package repository

import (
	"context"
	"database/sql"

	"github.com/Alena-Kurushkina/shortener/internal/api"
	"github.com/Alena-Kurushkina/shortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	database *sql.DB

}

// NewRepository initializes data storage
func NewRepository(connectionStr string) (api.Storager, error) {

	db,err:=sql.Open("pgx", connectionStr)
	if err!=nil{
		return nil, err
	}
	logger.Log.Info("DB connection opened")

	return &Repository{database: db}, nil
}

func (rp *Repository) Ping(ctx context.Context) error{
	return rp.database.PingContext(ctx)
}

func (rp *Repository) Close(){
	rp.database.Close()
}

// A record set data representation in file
// type record struct {
// 	UUID        uint   `json:"uuid"`
// 	ShortURL    string `json:"short_url"`
// 	OriginalURL string `json:"original_url"`
// }

// // NewRepository initializes data storage
// func NewRepository(filename string) (api.Storager, error) {
// 	// open storage file to read
// 	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	// scan all lines from file
// 	scanner := bufio.NewScanner(file)
// 	rmap := make(map[string]string)
// 	record := record{}

// 	for scanner.Scan() {
// 		data := scanner.Bytes()
// 		err := json.Unmarshal(data, &record)
// 		if err != nil {
// 			return nil, err
// 		}
// 		rmap[record.ShortURL] = record.OriginalURL
// 	}
// 	if scanner.Err() != nil {
// 		return nil, scanner.Err()
// 	}

// 	db := Repository{
// 		db:       rmap,
// 		filename: filename,
// 	}

// 	return &db, nil
// }

// Insert adds data to storage
// func (r Repository) Insert(key, value string) error {
// 	// write data to local map
// 	r.db[key] = value

// 	// open file
// 	file, err := os.OpenFile(r.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	writer := bufio.NewWriter(file)

// 	// encode data
// 	rec := record{UUID: uint(len(r.db)), OriginalURL: value, ShortURL: key}
// 	data, err := json.Marshal(&rec)
// 	if err != nil {
// 		return err
// 	}

// 	// write data to buffer
// 	if _, err := writer.Write(data); err != nil {
// 		return err
// 	}

// 	// go to next line
// 	if err := writer.WriteByte('\n'); err != nil {
// 		return err
// 	}

// 	// write buffer to file
// 	if err = writer.Flush(); err != nil {
// 		return err
// 	}

// 	return nil
// }

// Select returns data from storage
// func (r Repository) Select(key string) (string, error) {
// 	if v, ok := r.db[key]; ok {
// 		return v, nil
// 	}
// 	return "", fmt.Errorf("can't find value of key")
// }

func (r Repository) Insert(key, value string) error { return nil}

func (r Repository) Select(key string) (string, error) {return "", nil}