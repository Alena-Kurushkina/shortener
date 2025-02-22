// Package repository implements routines for manipulating data source.
package repository

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Alena-Kurushkina/shortener/internal/api"
	"github.com/Alena-Kurushkina/shortener/internal/config"
	"github.com/Alena-Kurushkina/shortener/internal/logger"
)

// NewRepository defines data storage depending on passed config parameters.
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
