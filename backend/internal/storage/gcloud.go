package storage

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/torbenconto/spooler/internal/util"
)

type GoogleCloudStorageClient struct {
	bucketName string
	client     *storage.Client
}

func NewGoogleCloudStorageClient(ctx context.Context, bucketName string) (*GoogleCloudStorageClient, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GoogleCloudStorageClient{
		bucketName: bucketName,
		client:     client,
	}, nil
}

func (g *GoogleCloudStorageClient) StoreFile(ctx context.Context, objectPath string, file io.Reader) error {
	contentType, reader, err := util.DetectContentType(file)
	if err != nil {
		return err
	}

	w := g.client.Bucket(g.bucketName).Object(objectPath).NewWriter(ctx)
	defer w.Close()

	w.ContentType = contentType
	// No escaping should be neccessary, the only files uploaded using this function will be converted into (uuid).(stl||3mf||gcode.3mf)
	// If sanatization happens to be required, it will be handeled by the caller
	w.ContentDisposition = fmt.Sprintf(`attachment; filename="%s"`, objectPath)

	if _, err := io.Copy(w, reader); err != nil {
		return err
	}

	return nil
}

func (g *GoogleCloudStorageClient) DeleteFile(ctx context.Context, objectPath string) error {
	if err := g.client.Bucket(g.bucketName).Object(objectPath).Delete(ctx); err != nil {
		return err
	}

	return nil
}

func (g *GoogleCloudStorageClient) GetFile(ctx context.Context, objectPath string) (io.ReadCloser, error) {
	r, err := g.client.Bucket(g.bucketName).Object(objectPath).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	return r, nil
}
