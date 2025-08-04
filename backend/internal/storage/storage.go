package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/torbenconto/spooler/config"
	"github.com/torbenconto/spooler/internal/types"
)

type StorageClient interface {
	GetFile(ctx context.Context, objectPath string) (io.ReadCloser, error)
	StoreFile(ctx context.Context, objectPath string, file io.Reader) error
	DeleteFile(ctx context.Context, objectPath string) error
}

func NewStorageClient(ctx context.Context, appConfig *config.Config) (StorageClient, error) {
	provider := appConfig.Storage.Provider

	switch provider {
	case types.GoogleCloudStorage:
		return NewGoogleCloudStorageClient(ctx, appConfig.Storage.GoogleCloud.BucketName)
	case types.LocalStorage:
		return NewLocalStorageClient(appConfig.Storage.Local.BasePath)
	default:
		return nil, fmt.Errorf("invalid provider: %s", provider)
	}
}
