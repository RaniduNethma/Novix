package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/streamingplatform/processing-service/internal/config"
	"github.com/streamingplatform/processing-service/pkg/logger"
	"go.uber.org/zap"
)

type MinIOClient struct {
	client          *minio.Client
	rawBucket       string
	processedBucket string
}

func NewMinIOClient(cfg *config.MinIOConfig) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}
	logger.Info("MinIO client connected", zap.String("endpoint", cfg.Endpoint))
	return &MinIOClient{
		client:          client,
		rawBucket:       cfg.RawBucket,
		processedBucket: cfg.ProcessedBucket,
	}, nil
}

// DownloadRawVideo downloads a raw video from MinIO to a local temp path
func (m *MinIOClient) DownloadRawVideo(ctx context.Context, objectKey string) (string, error) {
	// Create temp file
	tmpDir := os.TempDir()
	localPath := filepath.Join(tmpDir, filepath.Base(objectKey))
	err := m.client.FGetObject(
		ctx,
		m.rawBucket,
		objectKey,
		localPath,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return "", fmt.Errorf("failed to download raw video: %w", err)
	}
	logger.Info("Downloaded raw video",
		zap.String("object", objectKey),
		zap.String("localPath", localPath),
	)
	return localPath, nil
}

// UploadProcessedFile uploads a processed file to the processed bucket
func (m *MinIOClient) UploadProcessedFile(ctx context.Context,
	localPath string, objectKey string, contentType string) error {
	_, err := m.client.FPutObject(
		ctx,
		m.processedBucket,
		objectKey,
		localPath,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to upload processed file: %w", err)
	}
	logger.Info("Uploaded processed file",
		zap.String("object", objectKey),
		zap.String("bucket", m.processedBucket),
	)
	return nil
}

// UploadDirectory uploads all files in a directory to MinIO
func (m *MinIOClient) UploadDirectory(ctx context.Context, localDir string, objectPrefix string) error {
	return filepath.Walk(localDir, func(path string, info os.FileInfo,
		err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Determine content type
		contentType := "application/octet-stream"
		ext := filepath.Ext(path)
		switch ext {
		case ".m3u8":
			contentType = "application/vnd.apple.mpegurl"
		case ".ts":
			contentType = "video/MP2T"
		case ".mp4":
			contentType = "video/mp4"
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		}

		// Build object key
		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}
		objectKey := filepath.Join(objectPrefix, relPath)
		return m.UploadProcessedFile(ctx, path, objectKey, contentType)
	})
}

// GetPresignedURL returns a presigned URL for a processed file
func (m *MinIOClient) GetPresignedURL(ctx context.Context, objectKey string) (string, error) {
	url, err := m.client.PresignedGetObject(
		ctx,
		m.processedBucket,
		objectKey,
		0,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return url.String(), nil
}

// DeleteTempFile removes a local temp file after processing
func (m *MinIOClient) DeleteTempFile(path string) {
	if err := os.Remove(path); err != nil {
		logger.Warn("Failed to delete temp file",
			zap.String("path", path),
			zap.Error(err),
		)
	}
}

// DeleteTempDir removes a local temp directory after processing
func (m *MinIOClient) DeleteTempDir(path string) {
	if err := os.RemoveAll(path); err != nil {
		logger.Warn("Failed to delete temp dir",
			zap.String("path", path),
			zap.Error(err),
		)
	}
}

// DownloadToWriter streams an object directly to an io.Writer
func (m *MinIOClient) DownloadToWriter(ctx context.Context, objectKey string, writer io.Writer) error {
	object, err := m.client.GetObject(
		ctx,
		m.rawBucket,
		objectKey,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to get object: %w", err)
	}
	defer object.Close()
	if _, err := io.Copy(writer, object); err != nil {
		return fmt.Errorf("failed to copy object to writer: %w", err)
	}
	return nil
}
