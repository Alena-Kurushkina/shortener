// Package repository implements routines for manipulating data source
package repository

import "fmt"

// A Repository represents a data storage
type Repository map[string]string

// NewRepository initializes data storage
func NewRepository() *Repository {
	db := make(Repository)
	return &db
}

// Insert adds data to storage
func (r Repository) Insert(key, value string) error {
	r[key] = value

	return nil
}

// Select returns data from storage
func (r Repository) Select(key string) (string, error) {
	if v, ok := r[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("can't find value of key")
}
