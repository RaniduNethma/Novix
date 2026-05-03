package service

import (
	"context"
	"fmt"
	"time"

	"github.com/novix/services/streaming-service/pkg/logger"
	"github.com/novix/services/streaming-service/pkg/storage"
	"go.uber.org/zap"
)

type StreamService interface {
	GetManifest(ctx context.Context, videoID string,
		quality string) ([]byte, error)
	GetSegment(ctx context.Context, videoID string,
		quality string, segment string) ([]byte, error)
	GetPresignedManifestURL(ctx context.Context,
		videoID string, ttl time.Duration) (string, error)
	GetThumbnailURL(ctx context.Context,
		videoID string, ttl time.Duration) (string, error)
	VideoExists(ctx context.Context, videoID string) (bool, error)
}

type streamService struct {
	storage *storage.MinIOClient
}

func NewStreamService(storage *storage.MinIOClient) StreamService {
	return &streamService{storage: storage}
}

// GetManifest returns the HLS playlist for a specific quality
func (s *streamService) GetManifest(ctx context.Context,
	videoID string, quality string) ([]byte, error) {

	var objectKey string

	if quality == "" || quality == "master" {
		// Return master playlist
		objectKey = fmt.Sprintf("videos/%s/hls/master.m3u8", videoID)
	} else {
		// Return quality-specific playlist
		objectKey = fmt.Sprintf("videos/%s/hls/%s/playlist.m3u8",
			videoID, quality)
	}

	exists, err := s.storage.ObjectExists(ctx, objectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to check manifest: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("manifest not found for video: %s", videoID)
	}

	data, err := s.storage.GetObjectBytes(ctx, objectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest: %w", err)
	}

	logger.Debug("Served manifest",
		zap.String("videoID", videoID),
		zap.String("quality", quality),
	)

	return data, nil
}

// GetSegment returns a specific HLS .ts segment
func (s *streamService) GetSegment(ctx context.Context,
	videoID string, quality string, segment string) ([]byte, error) {

	objectKey := fmt.Sprintf("videos/%s/hls/%s/%s",
		videoID, quality, segment)

	exists, err := s.storage.ObjectExists(ctx, objectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to check segment: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("segment not found: %s", segment)
	}

	data, err := s.storage.GetObjectBytes(ctx, objectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get segment: %w", err)
	}

	return data, nil
}

// GetPresignedManifestURL returns a presigned URL for the master playlist
func (s *streamService) GetPresignedManifestURL(ctx context.Context,
	videoID string, ttl time.Duration) (string, error) {

	objectKey := fmt.Sprintf("videos/%s/hls/master.m3u8", videoID)

	url, err := s.storage.GeneratePresignedURL(ctx, objectKey, ttl)
	if err != nil {
		return "", fmt.Errorf("failed to generate manifest URL: %w", err)
	}

	return url, nil
}

// GetThumbnailURL returns a presigned URL for the video thumbnail
func (s *streamService) GetThumbnailURL(ctx context.Context,
	videoID string, ttl time.Duration) (string, error) {

	objectKey := fmt.Sprintf("videos/%s/thumbnail.jpg", videoID)

	exists, err := s.storage.ObjectExists(ctx, objectKey)
	if err != nil || !exists {
		return "", nil // thumbnail optional — don't fail
	}

	url, err := s.storage.GeneratePresignedURL(ctx, objectKey, ttl)
	if err != nil {
		return "", fmt.Errorf("failed to generate thumbnail URL: %w", err)
	}

	return url, nil
}

// VideoExists checks if a processed video exists in storage
func (s *streamService) VideoExists(ctx context.Context,
	videoID string) (bool, error) {

	objectKey := fmt.Sprintf("videos/%s/hls/master.m3u8", videoID)
	return s.storage.ObjectExists(ctx, objectKey)
}