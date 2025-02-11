package repository

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/Alena-Kurushkina/shortener/internal/api"
)

// A MemoryRepository represents a memory data storage.
type MemoryRepository struct {
	db map[string]string
}

// newMemoryRepository initializes data storage in memory.
func newMemoryRepository() (api.Storager, error) {
	db := MemoryRepository{
		db: make(map[string]string),
	}
	return &db, nil
}

// Insert adds data to storage.
func (r MemoryRepository) Insert(_ context.Context, id uuid.UUID, key, value string) error {
	r.db[key] = value

	return nil
}

// InsertBatch adds array of data to storage.
func (r MemoryRepository) InsertBatch(_ context.Context, id uuid.UUID, batch []api.BatchElement) error {
	for _, v := range batch {
		r.db[v.ShortURL] = v.OriginalURL
	}

	return nil
}

// Select returns data from storage.
func (r MemoryRepository) Select(_ context.Context, key string) (string, error) {
	if v, ok := r.db[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("can't find value of key")
}

// SelectUserAll returns data from storage.
func (r MemoryRepository) SelectUserAll(ctx context.Context, id uuid.UUID) ([]api.BatchElement, error) {
	return []api.BatchElement{}, nil
}

// DeleteRecords delete data from storage.
func (r MemoryRepository) DeleteRecords(ctx context.Context, deleteItems []api.DeleteItem) error {
	return nil
}

// Close satisfies the interface.
func (r *MemoryRepository) Close() {}

// Ping satisfies the interface.
func (r *MemoryRepository) Ping(_ context.Context) error { return nil }
