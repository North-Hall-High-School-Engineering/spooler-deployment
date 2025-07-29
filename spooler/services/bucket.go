package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/torbenconto/spooler/config"
	"gorm.io/gorm"
)

type BucketService struct {
	db     *gorm.DB
	client *storage.Client
	bucket string
}

func NewBucketService(db *gorm.DB, storageClient *storage.Client) *BucketService {
	return &BucketService{
		db:     db,
		client: storageClient,
		bucket: config.Cfg.GcloudPrintFilesBucket,
	}
}

func (g *BucketService) GetFileBase64(ctx context.Context, objectName string) (string, error) {
	rc, err := g.client.Bucket(g.bucket).Object(objectName).NewReader(ctx)
	if err != nil {
		return "", err
	}
	defer rc.Close()

	var b strings.Builder
	encoder := base64.NewEncoder(base64.StdEncoding, &b)

	if _, err := io.Copy(encoder, rc); err != nil {
		return "", err
	}
	encoder.Close()

	return b.String(), nil
}

func (g *BucketService) UploadFile(ctx context.Context, objectName string, file io.Reader) error {
	wc := g.client.Bucket(g.bucket).Object(objectName).NewWriter(ctx)
	defer wc.Close()

	wc.ContentType = "application/octet-stream"
	wc.ContentDisposition = fmt.Sprintf(`attachment; filename="%s"`, objectName)

	if _, err := io.Copy(wc, file); err != nil {
		return err
	}

	return nil
}

func (s *BucketService) GetFile(ctx context.Context, objectName string) (io.Reader, error) {
	rc, err := s.client.Bucket(s.bucket).Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	return rc, nil
}

func (s *BucketService) DeleteFile(ctx context.Context, objectName string) error {
	err := s.client.Bucket(s.bucket).Object(objectName).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
