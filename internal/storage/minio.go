package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go"
)

type MinioProvider struct {
	client *minio.Client
	bucket string
	region string
}

func NewMinioProvider(config *StorageConfig) (*MinioProvider, error) {
	client, err := minio.New(
		config.Endpoint,
		config.AccessKey,
		config.SecretKey,
		config.UseSSL,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	provider := &MinioProvider{
		client: client,
		bucket: config.Bucket,
		region: config.Region,
	}

	if err := provider.ensureBucketExists(); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	return provider, nil
}

func (p *MinioProvider) ensureBucketExists() error {
	exists, err := p.client.BucketExists(p.bucket)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !exists {
		return p.client.MakeBucket(p.bucket, p.region)
	}

	return nil
}

func (p *MinioProvider) UploadFile(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	_, err := p.client.PutObject(
		p.bucket, key, reader, size, minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil

}

func (p *MinioProvider) DeleteFile(ctx context.Context, key string) error {
	err := p.client.RemoveObject(p.bucket, key)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (p *MinioProvider) ListFiles(ctx context.Context, prefix string) ([]FileInfo, error) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	objects := make([]FileInfo, 0)
	for object := range p.client.ListObjectsV2(p.bucket, prefix, true, doneCh) {
		if object.Err != nil {
			return nil, fmt.Errorf("failed to list files: %w", object.Err)
		}
		objects = append(objects, FileInfo{
			Key:          object.Key,
			Size:         object.Size,
			ContentType:  object.ContentType,
			LastModified: object.LastModified,
			Etag:         object.ETag,
		})
	}
	return objects, nil
}

func (p *MinioProvider) DownloadFile(ctx context.Context, key string) (io.Reader, error) {
	reader, err := p.client.GetObject(p.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return reader, nil
}

func (p *MinioProvider) GetFileURL(ctx context.Context, key string) (string, error) {
	url, err := p.client.PresignedGetObject(p.bucket, key, time.Hour*24, url.Values{})
	if err != nil {
		return "", fmt.Errorf("failed to get file URL: %w", err)
	}

	return url.String(), nil
}

func (p *MinioProvider) FileExists(ctx context.Context, key string) (bool, error) {
	_, err := p.client.StatObject(p.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to check if file exists: %w", err)
	}

	return true, nil
}

func (p *MinioProvider) GetFileInfo(ctx context.Context, key string) (FileInfo, error) {
	object, err := p.client.StatObject(p.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		return FileInfo{}, fmt.Errorf("failed to get file info: %w", err)
	}

	return FileInfo{
		Key:          object.Key,
		Size:         object.Size,
		ContentType:  object.ContentType,
		LastModified: object.LastModified,
		Etag:         object.ETag,
	}, nil
}
