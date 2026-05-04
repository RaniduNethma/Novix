package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/novix/services/streaming-service/internal/config"
	"github.com/novix/services/streaming-service/pkg/logger"
	"go.uber.org/zap"
)

type MinIOClient struct {
	client          *minio.Client
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

	logger.Info("MinIO client connected",
		zap.String("endpoint", cfg.Endpoint),
	)

	return &MinIOClient{
		client:          client,
		processedBucket: cfg.ProcessedBucket,
	}, nil
}

// GetObject streams an object directly to a writer
// Used for proxying HLS segments to the client
func (m *MinIOClient) GetObject(ctx context.Context,
	objectKey string, writer io.Writer) error {

	object, err := m.client.GetObject(
		ctx,
		m.processedBucket,
		objectKey,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to get object: %w", err)
	}
	defer object.Close()

	if _, err := io.Copy(writer, object); err != nil {
		return fmt.Errorf("failed to stream object: %w", err)
	}

	return nil
}

// GetObjectBytes returns object content as bytes
// Used for serving HLS manifests
func (m *MinIOClient) GetObjectBytes(ctx context.Context,
	objectKey string) ([]byte, error) {

	object, err := m.client.GetObject(
		ctx,
		m.processedBucket,
		objectKey,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	return data, nil
}

// GeneratePresignedURL creates a time-limited URL for direct client access
func (m *MinIOClient) GeneratePresignedURL(ctx context.Context,
	objectKey string, ttl time.Duration) (string, error) {

	reqParams := make(url.Values)

	presignedURL, err := m.client.PresignedGetObject(
		ctx,
		m.processedBucket,
		objectKey,
		ttl,
		reqParams,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}

// ObjectExists checks if an object exists in the bucket
func (m *MinIOClient) ObjectExists(ctx context.Context,
	objectKey string) (bool, error) {

	_, err := m.client.StatObject(
		ctx,
		m.processedBucket,
		objectKey,
		minio.StatObjectOptions{},
	)
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to stat object: %w", err)
	}

	return true, nil
}