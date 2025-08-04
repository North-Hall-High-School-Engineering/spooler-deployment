package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorageClient struct {
	basePath string
}

func NewLocalStorageClient(basePath string) (*LocalStorageClient, error) {
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, fmt.Errorf("invalid base path: %w", err)
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorageClient{
		basePath: absPath,
	}, nil
}

func (l *LocalStorageClient) StoreFile(ctx context.Context, objectPath string, file io.Reader) error {
	// Validate and sanitize object path
	if err := l.validatePath(objectPath); err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	fullPath := filepath.Join(l.basePath, objectPath)

	if !strings.HasPrefix(fullPath, l.basePath) {
		return fmt.Errorf("path traversal attempt detected")
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if err := f.Chmod(0644); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	_, err = io.Copy(f, file)
	if err != nil {
		_ = os.Remove(fullPath)
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (l *LocalStorageClient) validatePath(objectPath string) error {
	if objectPath == "" {
		return fmt.Errorf("object path cannot be empty")
	}

	if strings.Contains(objectPath, "..") {
		return fmt.Errorf("path traversal not allowed")
	}

	if filepath.IsAbs(objectPath) {
		return fmt.Errorf("absolute paths not allowed")
	}

	if strings.ContainsAny(objectPath, "<>:\"|?*") {
		return fmt.Errorf("invalid characters in path")
	}

	return nil
}

func (l *LocalStorageClient) GetFile(ctx context.Context, objectPath string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(l.basePath, objectPath))
}

func (l *LocalStorageClient) DeleteFile(ctx context.Context, objectPath string) error {
	if err := l.validatePath(objectPath); err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	fullPath := filepath.Join(l.basePath, objectPath)

	if !strings.HasPrefix(fullPath, l.basePath) {
		return fmt.Errorf("path traversal attempt detected")
	}

	return os.Remove(fullPath)
}
