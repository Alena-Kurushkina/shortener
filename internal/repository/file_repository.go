package repository

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	uuid "github.com/satori/go.uuid"

	"github.com/Alena-Kurushkina/shortener/internal/api"
)

// A FileRepository represents a file data storage.
type FileRepository struct {
	db       map[string]string
	filename string
}

// newFileRepository initializes data storage in file.
func newFileRepository(filename string) (db api.Storager, err error) {
	// open storage file to read
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer func() {
		if tErr := file.Close(); tErr != nil {
			err = tErr
		}
	}()

	// scan all lines from file
	scanner := bufio.NewScanner(file)
	rmap := make(map[string]string)
	record := record{}

	for scanner.Scan() {
		data := scanner.Bytes()
		err = json.Unmarshal(data, &record)
		if err != nil {
			return nil, err
		}
		rmap[record.ShortURL] = record.OriginalURL
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	db = &FileRepository{
		db:       rmap,
		filename: filename,
	}

	return db, err
}

// Close satisfies interface.
func (r *FileRepository) Close() {}

// Ping satisfies interface.
func (r *FileRepository) Ping(_ context.Context) error { return nil }

// A record sets data representation in file.
type record struct {
	UUID        uuid.UUID `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

// InsertBatch adds array of data to storage.
func (r FileRepository) InsertBatch(_ context.Context, userID uuid.UUID, batch []api.BatchElement) (err error) {
	// open file
	file, err := os.OpenFile(r.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if tErr := file.Close(); tErr != nil {
			err = tErr
		}
	}()

	writer := bufio.NewWriter(file)

	for _, v := range batch {
		r.db[v.ShortURL] = v.OriginalURL

		// encode data
		rec := record{UUID: userID, OriginalURL: v.OriginalURL, ShortURL: v.ShortURL}
		data, errm := json.Marshal(&rec)
		if errm != nil {
			return errm
		}
		data = append(data, '\n')

		// write data to buffer
		if _, err = writer.Write(data); err != nil {
			return err
		}
	}

	// write buffer to file
	if err = writer.Flush(); err != nil {
		return err
	}

	return err
}

// Insert adds data to storage.
func (r FileRepository) Insert(_ context.Context, userID uuid.UUID, key, value string) (err error) {
	// write data to local map
	r.db[key] = value

	// open file
	file, err := os.OpenFile(r.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if tErr := file.Close(); tErr != nil {
			err = tErr
		}
	}()

	writer := bufio.NewWriter(file)

	// encode data
	rec := record{UUID: userID, OriginalURL: value, ShortURL: key}
	data, err := json.Marshal(&rec)
	if err != nil {
		return err
	}
	data = append(data, '\n')

	// write data to buffer
	if _, err = writer.Write(data); err != nil {
		return err
	}

	// write buffer to file
	if err = writer.Flush(); err != nil {
		return err
	}

	return err
}

// Select returns data from storage.
func (r FileRepository) Select(_ context.Context, key string) (string, error) {
	if v, ok := r.db[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("can't find value of key")
}

// SelectUserAll returns data from storage.
func (r FileRepository) SelectUserAll(ctx context.Context, id uuid.UUID) ([]api.BatchElement, error) {
	return []api.BatchElement{}, nil
}

// DeleteRecords deletes data from storage.
func (r FileRepository) DeleteRecords(ctx context.Context, deletedItems []api.DeleteItem) error {
	return nil
}

func (r FileRepository) SelectStats(ctx context.Context) (stats *api.Stats, err error){
	return nil, nil
}
