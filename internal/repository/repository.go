// Package repository implements routines for manipulating data source
package repository

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	uuid "github.com/satori/go.uuid"

	"github.com/Alena-Kurushkina/shortener/internal/api"
	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/logger"
	"github.com/Alena-Kurushkina/shortener/internal/sherr"
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
	database      *sql.DB
	selectStmt    *sql.Stmt
	selectAllStmt *sql.Stmt
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
func (r MemoryRepository) Insert(_ context.Context, id uuid.UUID, key, value string) error {
	r.db[key] = value

	return nil
}

// InsertBatch adds array of data to storage
func (r MemoryRepository) InsertBatch(_ context.Context, id uuid.UUID, batch []api.BatchElement) error {
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

// Select returns data from storage
func (r MemoryRepository) SelectUserAll(ctx context.Context, id uuid.UUID) ([]api.BatchElement, error) {
	// row := r.selectAllStmt.QueryRowContext(ctx,
	// 	id,
	// )
	// var longURL string

	// err := row.Scan(&longURL)
	// if err != nil {
	// 	return "", err
	// }

	return []api.BatchElement{}, nil
}

func (r MemoryRepository) DeleteRecords(ctx context.Context, deleteItems []api.DeleteItem) error {
	return nil
}

func (r *MemoryRepository) Close() {}

func (r *MemoryRepository) Ping(_ context.Context) error { return nil }

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
			id varchar(50) PRIMARY KEY default substring(md5(random()::text),0,20),
			originalURL varchar(500) NOT NULL,
			shortURL varchar(250) NOT NULL,
			userUUID uuid,
			is_deleted bool DEFAULT(false),
			UNIQUE (originalURL)
		);
		CREATE UNIQUE INDEX IF NOT EXISTS short_idx on shortening (shortURL);
	`)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	stmt1, err := db.PrepareContext(ctx, `
		SELECT originalURL, is_deleted 
		FROM shortening 
		WHERE shortURL LIKE $1
	`)
	if err != nil {
		return nil, err
	}

	stmt2, err := db.PrepareContext(ctx, `
		SELECT originalURL, shortURL
		FROM shortening 
		WHERE shortening.useruuid = $1
	`)
	if err != nil {
		return nil, err
	}

	return &DBRepository{database: db, selectStmt: stmt1, selectAllStmt: stmt2}, nil
}

// Insert adds data to storage
func (r DBRepository) Insert(ctx context.Context, userID uuid.UUID, insertedShortURL, insertedOriginalURL string) error {
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sqlRow := tx.QueryRowContext(ctx,
		`INSERT INTO shortening (userUUID, originalURL, shortURL) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (originalurl) 
			DO UPDATE SET originalurl = shortening.originalurl
		RETURNING originalurl, shorturl;`,
		userID,
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

func (r DBRepository) DeleteRecords(ctx context.Context, deleteItems []api.DeleteItem) error {
	//  UPDATE shortening
	// 	SET is_deleted=true
	// 	FROM (
	// 		VALUES
	// 			(uuid('0de70533-e13c-4269-bba9-e4e7653b4c72'), 'kjhg1234'),
	// 			(uuid('0de70533-e13c-4269-bba9-e4e7653b4c72'), 'dfgh345')
	// 	) AS data(id_user, id_record)
	// 	WHERE shortening.useruuid=data.id_user
	// 		AND shortening.id=data.id_record

	param := ""
	for _, v := range deleteItems {
		for _, i := range v.IDs {
			param = param + `(uuid('` + v.UserID.String() + `'), '` + i + `'),`
		}
	}
	param, _ = strings.CutSuffix(param, ",")

	stmt := `UPDATE shortening
		SET is_deleted=true
		FROM (
			VALUES` + param +
		`) AS data(id_user, shortening)
		WHERE shortening.useruuid=data.id_user
			AND shortening.shorturl=data.shortening`

	logger.Log.Info(stmt)

	sqlRes, err := r.database.ExecContext(ctx, stmt)

	if err != nil {
		logger.Log.Errorf("Error while deletion: ", err.Error())
	}

	sf, _ := sqlRes.RowsAffected()
	logger.Log.Infof("Rows affected while deletion: %s", strconv.FormatInt(sf, 10))

	// _, err := r.database.ExecContext(ctx,
	// 	`UPDATE shortening
	// 	SET is_deleted = true
	// 	WHERE userUUID=$1 AND id IN ('`+strings.Join(recordIDs, "','")+`');`,
	// 	userID,
	// )

	return err
}

// InsertBatch adds array of data to storage
func (r DBRepository) InsertBatch(ctx context.Context, userID uuid.UUID, batch []api.BatchElement) error {
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO shortening (id, userUUID, originalURL, shortURL) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (originalurl) 
			DO UPDATE SET originalurl = shortening.originalurl
		RETURNING originalurl, shorturl;
	`)
	if err != nil {
		return err
	}

	for _, v := range batch {
		_, err = stmt.ExecContext(ctx,
			v.CorrelarionID,
			userID,
			v.OriginalURL,
			v.ShortURL,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *DBRepository) Close() {
	r.selectStmt.Close()
	r.database.Close()
}

func (r *DBRepository) Ping(ctx context.Context) error {
	return r.database.PingContext(ctx)
}

// Select returns data from storage
func (r DBRepository) Select(ctx context.Context, key string) (string, error) {
	row := r.selectStmt.QueryRowContext(ctx,
		key,
	)
	var (
		longURL string
		deleted bool
	)

	err := row.Scan(&longURL, &deleted)
	if err != nil {
		return "", err
	}
	if deleted {
		return "", sherr.ErrDBRecordDeleted
	}

	return longURL, nil
}

// Select returns data from storage
func (r DBRepository) SelectUserAll(ctx context.Context, id uuid.UUID) ([]api.BatchElement, error) {
	rows, err := r.selectAllStmt.QueryContext(ctx,
		id.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]api.BatchElement, 0, 10)

	// пробегаем по всем записям
	for rows.Next() {
		var v api.BatchElement
		err = rows.Scan(&v.OriginalURL, &v.ShortURL)
		if err != nil {
			return nil, err
		}

		records = append(records, v)
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return records, nil
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

func (r *FileRepository) Close() {}

func (r *FileRepository) Ping(_ context.Context) error { return nil }

// A record set data representation in file
type record struct {
	UUID        uuid.UUID `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

// Insert adds array of data to storage
func (r FileRepository) InsertBatch(_ context.Context, userID uuid.UUID, batch []api.BatchElement) error {
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
		rec := record{UUID: userID, OriginalURL: v.OriginalURL, ShortURL: v.ShortURL}
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
func (r FileRepository) Insert(_ context.Context, userID uuid.UUID, key, value string) error {
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
	rec := record{UUID: userID, OriginalURL: value, ShortURL: key}
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

// Select returns data from storage
func (r FileRepository) SelectUserAll(ctx context.Context, id uuid.UUID) ([]api.BatchElement, error) {
	// row := r.selectAllStmt.QueryRowContext(ctx,
	// 	id,
	// )
	// var longURL string

	// err := row.Scan(&longURL)
	// if err != nil {
	// 	return "", err
	// }

	return []api.BatchElement{}, nil
}

func (r FileRepository) DeleteRecords(ctx context.Context, deletedItems []api.DeleteItem) error {
	return nil
}
