package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/novix/services/processing-service/internal/model"
	"github.com/novix/services/processing-service/internal/repository"
	"github.com/novix/services/processing-service/pkg/ffmpeg"
	"github.com/novix/services/processing-service/pkg/logger"
	"github.com/novix/services/processing-service/pkg/storage"
	"go.uber.org/zap"
)

type ProcessingService interface {
	ProcessVideo(ctx context.Context, job *model.ProcessingJob) error
	GenerateThumbnail(ctx context.Context, job *model.ProcessingJob) error
}

type processingService struct {
	jobRepo    repository.JobRepository
	ffmpeg     *ffmpeg.FFmpeg
	storage    *storage.MinIOClient
	maxRetries int
}

func NewProcessingService(
	jobRepo repository.JobRepository,
	ffmpeg *ffmpeg.FFmpeg,
	storage *storage.MinIOClient,
	maxRetries int,
) ProcessingService {
	return &processingService{
		jobRepo:    jobRepo,
		ffmpeg:     ffmpeg,
		storage:    storage,
		maxRetries: maxRetries,
	}
}

// ProcessVideo handles full HLS transcoding pipeline
func (s *processingService) ProcessVideo(ctx context.Context,
	job *model.ProcessingJob) error {

	// Mark job as processing
	now := time.Now()
	job.Status = model.JobStatusProcessing
	job.StartedAt = &now
	if err := s.jobRepo.Update(ctx, job); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	logger.Info("Processing video",
		zap.String("videoID", job.VideoID),
		zap.Uint("jobID", job.ID),
	)

	// Step 1 — Download raw video from MinIO
	localVideoPath, err := s.storage.DownloadRawVideo(ctx, job.RawPath)
	if err != nil {
		return s.failJob(ctx, job, fmt.Sprintf("download failed: %v", err))
	}
	defer s.storage.DeleteTempFile(localVideoPath)

	// Step 2 — Create temp output directory for HLS
	outputDir := filepath.Join(os.TempDir(), "hls", job.VideoID)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return s.failJob(ctx, job,
			fmt.Sprintf("failed to create output dir: %v", err))
	}
	defer s.storage.DeleteTempDir(outputDir)

	// Step 3 — Generate HLS with adaptive bitrate
	hlsOpts := ffmpeg.DefaultHLSOptions(localVideoPath, outputDir)
	masterPath, err := s.ffmpeg.GenerateHLS(ctx, hlsOpts)
	if err != nil {
		return s.failJob(ctx, job,
			fmt.Sprintf("HLS generation failed: %v", err))
	}

	logger.Info("HLS generation complete",
		zap.String("masterPlaylist", masterPath),
	)

	// Step 4 — Upload entire HLS directory to MinIO
	objectPrefix := fmt.Sprintf("videos/%s/hls", job.VideoID)
	if err := s.storage.UploadDirectory(ctx, outputDir, objectPrefix); err != nil {
		return s.failJob(ctx, job,
			fmt.Sprintf("upload failed: %v", err))
	}

	// Step 5 — Mark job complete
	completedAt := time.Now()
	job.Status = model.JobStatusCompleted
	job.CompletedAt = &completedAt
	job.OutputPath = fmt.Sprintf("%s/master.m3u8", objectPrefix)

	if err := s.jobRepo.Update(ctx, job); err != nil {
		return fmt.Errorf("failed to mark job complete: %w", err)
	}
	logger.Info("Video processing complete",
		zap.String("videoID", job.VideoID),
		zap.String("outputPath", job.OutputPath),
	)
	return nil
}

// GenerateThumbnail extracts a thumbnail from the video
func (s *processingService) GenerateThumbnail(ctx context.Context,
	job *model.ProcessingJob) error {

	now := time.Now()
	job.Status = model.JobStatusProcessing
	job.StartedAt = &now
	if err := s.jobRepo.Update(ctx, job); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	// Download raw video
	localVideoPath, err := s.storage.DownloadRawVideo(ctx, job.RawPath)
	if err != nil {
		return s.failJob(ctx, job,
			fmt.Sprintf("download failed: %v", err))
	}
	defer s.storage.DeleteTempFile(localVideoPath)

	// Generate thumbnail
	thumbnailPath := filepath.Join(os.TempDir(),
		fmt.Sprintf("thumb_%s.jpg", job.VideoID))
	defer s.storage.DeleteTempFile(thumbnailPath)

	thumbOpts := ffmpeg.ThumbnailOptions{
		InputPath:  localVideoPath,
		OutputPath: thumbnailPath,
		TimeOffset: "00:00:05",
		Width:      1280,
		Height:     720,
	}

	if err := s.ffmpeg.GenerateThumbnail(ctx, thumbOpts); err != nil {
		return s.failJob(ctx, job,
			fmt.Sprintf("thumbnail generation failed: %v", err))
	}

	// Upload thumbnail to MinIO
	objectKey := fmt.Sprintf("videos/%s/thumbnail.jpg", job.VideoID)
	if err := s.storage.UploadProcessedFile(
		ctx, thumbnailPath, objectKey, "image/jpeg"); err != nil {
		return s.failJob(ctx, job,
			fmt.Sprintf("thumbnail upload failed: %v", err))
	}

	// Mark complete
	completedAt := time.Now()
	job.Status = model.JobStatusCompleted
	job.CompletedAt = &completedAt
	job.OutputPath = objectKey

	if err := s.jobRepo.Update(ctx, job); err != nil {
		return fmt.Errorf("failed to mark job complete: %w", err)
	}
	logger.Info("Thumbnail generation complete",
		zap.String("videoID", job.VideoID),
		zap.String("thumbnailPath", objectKey),
	)
	return nil
}

// failJob marks a job as failed and handles retry logic
func (s *processingService) failJob(ctx context.Context,
	job *model.ProcessingJob, errMsg string) error {

	job.RetryCount++
	job.ErrorMsg = errMsg

	if job.RetryCount >= s.maxRetries {
		job.Status = model.JobStatusFailed
		logger.Error("Job permanently failed",
			zap.String("videoID", job.VideoID),
			zap.String("error", errMsg),
			zap.Int("retries", job.RetryCount),
		)
	} else {
		job.Status = model.JobStatusRetrying
		logger.Warn("Job failed, will retry",
			zap.String("videoID", job.VideoID),
			zap.String("error", errMsg),
			zap.Int("retryCount", job.RetryCount),
		)
	}
	if err := s.jobRepo.Update(ctx, job); err != nil {
		return fmt.Errorf("failed to update failed job: %w", err)
	}
	return fmt.Errorf(errMsg)
}