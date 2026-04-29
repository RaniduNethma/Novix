package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/novix/services/processing-service/internal/config"
	"github.com/novix/services/processing-service/internal/handler"
	"github.com/novix/services/processing-service/internal/middleware"
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

	logger.Info("Starting processing service",
		zap.String("env", cfg.App.Env),
		zap.String("port", cfg.Server.Port),
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

	// Auto migrate models
	if err := db.AutoMigrate(
		&model.ProcessingJob{},
		&model.VideoVariant{},
	); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	logger.Info("Database connected and migrated")

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

	// Start Kafka consumer in background
	consumer := kafkapkg.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.TopicUploaded,
		cfg.Kafka.GroupId,
	)

	go func() {
		err := consumer.Consume(ctx,
			func(event kafkapkg.VideoUploadedEvent) error {
				job, err := jobSvc.CreateJob(
					ctx,
					event.VideoID,
					event.UserID,
					model.JobTypeHLS,
					event.RawPath,
				)
				if err != nil {
					return err
				}

				// Submit to worker pool
				pool.Submit(job)

				// Also create thumbnail job
				thumbJob, err := jobSvc.CreateJob(
					ctx,
					event.VideoID,
					event.UserID,
					model.JobTypeThumbnail,
					event.RawPath,
				)
				if err != nil {
					return err
				}
				pool.Submit(thumbJob)

				return nil
			})
		if err != nil {
			logger.Error("Kafka consumer error", zap.Error(err))
		}
	}()

	// Set up Gin router
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.RequestLogger())
	router.Use(gin.Recovery())

	// Init handlers
	jobHandler := handler.NewJobHandler(jobSvc)

	// Routes
	router.GET("/actuator/health", jobHandler.HealthCheck)

	api := router.Group("/api/v1")
	api.Use(middleware.JWTAuth(cfg.JWT.Secret))
	{
		jobs := api.Group("/jobs")
		{
			jobs.GET("/:id", jobHandler.GetJobByID)
			jobs.GET("/video/:videoId", jobHandler.GetJobsByVideoID)
		}
	}

	// HTTP server with graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Server listening",
			zap.String("port", cfg.Server.Port),
		)
		if err := srv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(), 10*time.Second,
	)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited cleanly")
}