package storage

import (
	"context"
	"io"
	"time"
)

type FileInfo struct {
	Key          string
	Size         int64
	ContentType  string
	LastModified time.Time
	Etag         string
}

type StorageProvider interface {
	UploadFile(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error

	DeleteFile(ctx context.Context, key string) error

	DownloadFile(ctx context.Context, key string) (io.Reader, error)

	ListFiles(ctx context.Context, prefix string) ([]FileInfo, error)

	GetFileURL(ctx context.Context, key string) (string, error)

	FileExists(ctx context.Context, key string) (bool, error)

	GetFileInfo(ctx context.Context, key string) (FileInfo, error)
}

type StorageConfig struct {
	Provider    string `json:"provider"`   // minio, s3, etc.
	Endpoint    string `json:"endpoint"`   // minio endpoint, s3 endpoint, etc.
	AccessKey   string `json:"access_key"` // minio access key, s3 access key, etc.
	SecretKey   string `json:"secret_key"`
	Bucket      string `json:"bucket"`
	Region      string `json:"region"` // minio region, s3 region, etc.
	UseSSL      bool   `json:"use_ssl"`
	Credentials string `json:"credentials"` // for GCS service account
}
