// Package repository implements routines for manipulating data source
package repository

import (
	"fmt"

	"github.com/Alena-Kurushkina/shortener/internal/api"
)

// A Repository represents a data storage
type Repository struct {
	db map[string]string
}

// NewRepository initializes data storage
func NewRepository() api.Storager {
	db := Repository{
		db: make(map[string]string),
	}
	return &db
}

// Insert adds data to storage
func (r Repository) Insert(key, value string) error {
	r.db[key] = value

	return nil
}

// Select returns data from storage
func (r Repository) Select(key string) (string, error) {
	if v, ok := r.db[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("can't find value of key")
}
