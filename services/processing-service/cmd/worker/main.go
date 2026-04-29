package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/novix/services/processing-service/internal/config"
	"github.com/novix/services/processing-service/internal/model"
	"github.com/novix/services/processing-service/internal/repository"
	"github.com/novix/services/processing-service/internal/service"
	"github.com/novix/services/processing-service/internal/worker"
	ffmpegpkg "github.com/novix/services/processing-service/pkg/ffmpeg"
	kafkapkg "github.com/novix/services/processing-service/pkg/kafka"
	"github.com/novix/services/processing-service/pkg/logger"
	"github.com/novix/services/processing-service/pkg/storage"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load config
	cfg := config.Load()

	// Init logger
	logger.Init(cfg.App.Env)
	defer logger.Sync()

	logger.Info("Starting processing worker",
		zap.String("env", cfg.App.Env),
		zap.Int("poolSize", cfg.Worker.PoolSize),
	)

	// Connect to PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto migrate
	if err := db.AutoMigrate(
		&model.ProcessingJob{},
		&model.VideoVariant{},
	); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	logger.Info("Database connected")

	// Init MinIO
	minioClient, err := storage.NewMinIOClient(&cfg.MinIO)
	if err != nil {
		logger.Fatal("Failed to connect to MinIO", zap.Error(err))
	}

	// Init FFmpeg
	ffmpeg := ffmpegpkg.New(cfg.FFmpeg.Path, cfg.FFmpeg.Threads)

	// Init layers
	jobRepo := repository.NewJobRepository(db)
	jobSvc := service.NewJobService(jobRepo)
	processingSvc := service.NewProcessingService(
		jobRepo, ffmpeg, minioClient, cfg.Worker.MaxRetries,
	)

	// Init worker pool
	processor := worker.NewProcessor(processingSvc)
	pool := worker.NewPool(cfg.Worker.PoolSize, processor)

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker pool
	pool.Start(ctx)

	// Start Kafka consumer — this is the ONLY job of the worker binary
	consumer := kafkapkg.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.TopicUploaded,
		cfg.Kafka.GroupId,
	)
	defer consumer.Close()

	logger.Info("Worker ready — listening for video upload events",
		zap.Strings("brokers", cfg.Kafka.Brokers),
		zap.String("topic", cfg.Kafka.TopicUploaded),
	)

	go func() {
		err := consumer.Consume(ctx,
			func(event kafkapkg.VideoUploadedEvent) error {

				// Create HLS transcoding job
				job, err := jobSvc.CreateJob(
					ctx,
					event.VideoID,
					event.UserID,
					model.JobTypeHLS,
					event.RawPath,
				)
				if err != nil {
					return fmt.Errorf("failed to create HLS job: %w", err)
				}
				pool.Submit(job)

				// Create thumbnail job
				thumbJob, err := jobSvc.CreateJob(
					ctx,
					event.VideoID,
					event.UserID,
					model.JobTypeThumbnail,
					event.RawPath,
				)
				if err != nil {
					return fmt.Errorf("failed to create thumbnail job: %w", err)
				}
				pool.Submit(thumbJob)

				logger.Info("Jobs submitted to worker pool",
					zap.String("videoID", event.VideoID),
				)

				return nil
			})
		if err != nil {
			logger.Error("Kafka consumer error", zap.Error(err))
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Worker shutting down gracefully...")
	cancel()
	logger.Info("Worker stopped")
}