// Package repository implements routines for manipulating data source
package repository

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Alena-Kurushkina/shortener/internal/api"
	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// A FileRepository represents a file storage
type FileRepository struct {
	db       map[string]string
	filename string
}

type MemoryRepository struct {
	db       map[string]string
}

type DBRepository struct {
	database *sql.DB
	insertStmt *sql.Stmt
	selectStmt *sql.Stmt

}

func NewRepository(ctx context.Context, config *config.Config) (api.Storager, error){
	if config.ConnectionStr!="" {
		logger.Log.Info("Database is used as data storage")
		return newDBRepository(ctx, config.ConnectionStr)
	}
	if config.FileStoragePath!=""{
		logger.Log.Info("File is used as data storage")
		return newFileRepository(config.FileStoragePath)
	}
	logger.Log.Info("Memory is used as data storage")
	return newMemoryRepository()
}

func newMemoryRepository() (api.Storager, error){
	db := MemoryRepository{
		db: make(map[string]string),
	}
	return &db, nil
}

// Insert adds data to storage
func (r MemoryRepository) Insert(_ context.Context, key, value string) error {
	r.db[key] = value

	return nil
}

// Select returns data from storage
func (r MemoryRepository) Select(_ context.Context, key string) (string, error) {
	if v, ok := r.db[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("can't find value of key")
}

func (rp *MemoryRepository) Close(){

}

func (rp *MemoryRepository) Ping(_ context.Context) error{
	return nil
}

// NewRepository initializes data storage
func newDBRepository(ctx context.Context, connectionStr string) (api.Storager, error) {

	db,err:=sql.Open("pgx", connectionStr)
	if err!=nil{
		return nil, err
	}
	logger.Log.Info("DB connection opened")

	tx, err:=db.BeginTx(ctx,nil)
	if err!=nil{
		return nil, err
	}
	defer tx.Rollback()

	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS shortening(
			id SERIAL PRIMARY KEY,
			originalURL varchar(500) NOT NULL,
			shortURL varchar(250) NOT NULL
		)
	`)
	tx.ExecContext(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS short_idx on shortening (shortURL)`)

	err=tx.Commit()
	if err!=nil{
		return nil, err
	}

	stmt, err:=db.PrepareContext(ctx, "INSERT INTO shortening (originalURL, shortURL) VALUES ($1, $2)")
	if err!=nil{
		return nil, err
	}

	stmt1, err:=db.PrepareContext(ctx, "SELECT originalURL FROM shortening WHERE shortURL LIKE $1")
	if err!=nil{
		return nil, err
	}

	return &DBRepository{database: db, insertStmt: stmt, selectStmt: stmt1}, nil
}

// Insert adds data to storage
func (r DBRepository) Insert(ctx context.Context, key, value string) error {
	tx, err:=r.database.BeginTx(ctx, nil)
	if err!=nil{
		return err
	}
	defer tx.Rollback()

	_, err=r.insertStmt.ExecContext(ctx,
		value,
		key,
	)

	if err!=nil{
		return err
	}

	return tx.Commit()
}

func (rp *DBRepository) Close(){
	rp.insertStmt.Close()
	rp.selectStmt.Close()
	rp.database.Close()
}

func (rp *DBRepository) Ping(ctx context.Context) error{
	return rp.database.PingContext(ctx)
}

func (r DBRepository) Select(ctx context.Context, key string) (string, error) {
	row:=r.selectStmt.QueryRowContext(ctx,
		key,
	)
	var longURL string

	err:=row.Scan(&longURL)
	if err!=nil{
		return "", err
	}

	return longURL, nil
}

func newFileRepository(filename string) (api.Storager, error) {
	// open storage file to read
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// scan all lines from file
	scanner := bufio.NewScanner(file)
	rmap := make(map[string]string)
	record := record{}

	for scanner.Scan() {
		data := scanner.Bytes()
		err := json.Unmarshal(data, &record)
		if err != nil {
			return nil, err
		}
		rmap[record.ShortURL] = record.OriginalURL
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	db := FileRepository{
		db:       rmap,
		filename: filename,
	}

	return &db, nil
}

func (rp *FileRepository) Close(){}

func (rp *FileRepository) Ping(_ context.Context) error{
	return nil
}

// A record set data representation in file
type record struct {
	UUID        uint   `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// Insert adds data to storage
func (r FileRepository) Insert(_ context.Context, key, value string) error {
	// write data to local map
	r.db[key] = value

	// open file
	file, err := os.OpenFile(r.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// encode data
	rec := record{UUID: uint(len(r.db)), OriginalURL: value, ShortURL: key}
	data, err := json.Marshal(&rec)
	if err != nil {
		return err
	}

	// write data to buffer
	if _, err := writer.Write(data); err != nil {
		return err
	}

	// go to next line
	if err := writer.WriteByte('\n'); err != nil {
		return err
	}

	// write buffer to file
	if err = writer.Flush(); err != nil {
		return err
	}

	return nil
}

// Select returns data from storage
func (r FileRepository) Select(_ context.Context, key string) (string, error) {
	if v, ok := r.db[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("can't find value of key")
}
