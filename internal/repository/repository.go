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
	"github.com/Alena-Kurushkina/shortener/internal/sherr"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// A FileRepository represents a file data storage
type FileRepository struct {
	db       map[string]string
	filename string
}

// A MemoryRepository represents a memory data storage
type MemoryRepository struct {
	db map[string]string
}

// A DBRepository store data in database
type DBRepository struct {
	database   *sql.DB
	selectStmt *sql.Stmt
}

// NewRepository defines data storage depending on passed config parameters
func NewRepository(ctx context.Context, config *config.Config) (api.Storager, error) {
	if config.ConnectionStr != "" {
		logger.Log.Info("Database is used as data storage")
		return newDBRepository(ctx, config.ConnectionStr)
	}
	if config.FileStoragePath != "" {
		logger.Log.Info("File is used as data storage")
		return newFileRepository(config.FileStoragePath)
	}
	logger.Log.Info("Memory is used as data storage")
	return newMemoryRepository()
}

// newMemoryRepository initializes data storage in memory
func newMemoryRepository() (api.Storager, error) {
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

// InsertBatch adds array of data to storage
func (r MemoryRepository) InsertBatch(_ context.Context, batch []api.BatchElement) error {
	for _, v := range batch {
		r.db[v.ShortURL] = v.OriginalURL
	}

	return nil
}

// Select returns data from storage
func (r MemoryRepository) Select(_ context.Context, key string) (string, error) {
	if v, ok := r.db[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("can't find value of key")
}

func (rp *MemoryRepository) Close() {}

func (rp *MemoryRepository) Ping(_ context.Context) error { return nil }

// newDBRepository initializes data storage in database
func newDBRepository(ctx context.Context, connectionStr string) (api.Storager, error) {

	db, err := sql.Open("pgx", connectionStr)
	if err != nil {
		return nil, err
	}
	logger.Log.Info("DB connection opened")

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS shortening(
			id SERIAL PRIMARY KEY,
			originalURL varchar(500) NOT NULL,
			shortURL varchar(250) NOT NULL,
			UNIQUE (originalURL)
		);
		CREATE UNIQUE INDEX IF NOT EXISTS short_idx on shortening (shortURL);
	`)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	stmt1, err := db.PrepareContext(ctx, `
		SELECT originalURL 
		FROM shortening 
		WHERE shortURL LIKE $1
	`)
	if err != nil {
		return nil, err
	}

	return &DBRepository{database: db, selectStmt: stmt1}, nil
}

// Insert adds data to storage
func (r DBRepository) Insert(ctx context.Context, insertedShortURL, insertedOriginalURL string) error {
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sqlRow := tx.QueryRowContext(ctx,
		`INSERT INTO shortening (originalURL, shortURL) 
		VALUES ($1, $2) 
		ON CONFLICT (originalurl) 
			DO UPDATE SET originalurl = shortening.originalurl
		RETURNING originalurl, shorturl;`,
		insertedOriginalURL,
		insertedShortURL,
	)

	var (
		dbOriginalURL string
		dbShortURL    string
	)
	err = sqlRow.Scan(&dbOriginalURL, &dbShortURL)
	if err != nil {
		return err
	}

	if dbShortURL != insertedShortURL {
		tx.Rollback()
		return sherr.NewAlreadyExistError(insertedOriginalURL, dbShortURL)
	}

	return tx.Commit()
}

// InsertBatch adds array of data to storage
func (r DBRepository) InsertBatch(ctx context.Context, batch []api.BatchElement) error {
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO shortening (originalURL, shortURL) 
		VALUES ($1, $2) 
		ON CONFLICT (originalurl) 
			DO UPDATE SET originalurl = shortening.originalurl
		RETURNING originalurl, shorturl;
	`)
	if err != nil {
		return err
	}

	for _, v := range batch {
		_, err = stmt.ExecContext(ctx,
			v.OriginalURL,
			v.ShortURL,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (rp *DBRepository) Close() {
	rp.selectStmt.Close()
	rp.database.Close()
}

func (rp *DBRepository) Ping(ctx context.Context) error {
	return rp.database.PingContext(ctx)
}

// Select returns data from storage
func (r DBRepository) Select(ctx context.Context, key string) (string, error) {
	row := r.selectStmt.QueryRowContext(ctx,
		key,
	)
	var longURL string

	err := row.Scan(&longURL)
	if err != nil {
		return "", err
	}

	return longURL, nil
}

// newFileRepository initializes data storage in file
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

func (rp *FileRepository) Close() {}

func (rp *FileRepository) Ping(_ context.Context) error { return nil }

// A record set data representation in file
type record struct {
	UUID        uint   `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// Insert adds array of data to storage
func (r FileRepository) InsertBatch(_ context.Context, batch []api.BatchElement) error {
	// open file
	file, err := os.OpenFile(r.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, v := range batch {
		r.db[v.ShortURL] = v.OriginalURL

		// encode data
		rec := record{UUID: uint(len(r.db)), OriginalURL: v.OriginalURL, ShortURL: v.ShortURL}
		data, err := json.Marshal(&rec)
		if err != nil {
			return err
		}
		data = append(data, '\n')

		// write data to buffer
		if _, err := writer.Write(data); err != nil {
			return err
		}
	}

	// write buffer to file
	if err = writer.Flush(); err != nil {
		return err
	}

	return nil
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
	data = append(data, '\n')

	// write data to buffer
	if _, err := writer.Write(data); err != nil {
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
