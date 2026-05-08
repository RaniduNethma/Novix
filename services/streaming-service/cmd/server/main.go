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
	"github.com/novix/services/streaming-service/internal/config"
	"github.com/novix/services/streaming-service/internal/handler"
	"github.com/novix/services/streaming-service/internal/middleware"
	"github.com/novix/services/streaming-service/internal/repository"
	"github.com/novix/services/streaming-service/internal/service"
	"github.com/novix/services/streaming-service/pkg/logger"
	redispkg "github.com/novix/services/streaming-service/pkg/redis"
	"github.com/novix/services/streaming-service/pkg/storage"
	"go.uber.org/zap"
)

func main() {
	// Load config
	cfg := config.Load()

	// Init logger
	logger.Init(cfg.App.Env)
	defer logger.Sync()

	logger.Info("Starting streaming service",
		zap.String("env", cfg.App.Env),
		zap.String("port", cfg.Server.Port),
	)

	// Init Redis
	redisClient, err := redispkg.NewClient(&cfg.Redis)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()

	// Init MinIO
	minioClient, err := storage.NewMinIOClient(&cfg.MinIO)
	if err != nil {
		logger.Fatal("Failed to connect to MinIO", zap.Error(err))
	}

	// Init layers
	sessionRepo := repository.NewSessionRepository(redisClient)
	sessionSvc := service.NewSessionService(
		sessionRepo, cfg.Streaming.SessionTTL,
	)
	streamSvc := service.NewStreamService(minioClient)

	// Init handlers
	healthHandler := handler.NewHealthHandler()
	streamHandler := handler.NewStreamHandler(
		streamSvc,
		sessionSvc,
		cfg.Streaming.PresignTTL,
		cfg.Streaming.MaxSessionsPerUser,
	)

	// Set up Gin
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	// Public routes
	router.GET("/actuator/health", healthHandler.HealthCheck)

	// HLS delivery routes — public with session tracking
	hls := router.Group("/hls")
	{
		// GET /hls/{videoId}/master.m3u8
		hls.GET("/:videoId/manifest",
			streamHandler.GetManifest)

		// GET /hls/{videoId}/{quality}/{segment}.ts
		hls.GET("/:videoId/:quality/:segment",
			streamHandler.GetSegment)
	}

	// Protected routes — require JWT
	api := router.Group("/api/v1")
	api.Use(middleware.JWTAuth(cfg.JWT.Secret))
	{
		stream := api.Group("/stream")
		{
			// POST /api/v1/stream/start
			stream.POST("/start", streamHandler.StartStream)

			// DELETE /api/v1/stream/session/{sessionId}
			stream.DELETE("/session/:sessionId",
				streamHandler.EndStream)
		}
	}

	// HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second, // longer for video segments
	}

	// Start server
	go func() {
		logger.Info("Server listening",
			zap.String("port", cfg.Server.Port),
		)
		if err := srv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down streaming service...")

	ctx, cancel := context.WithTimeout(context.Background(),
		10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Streaming service stopped cleanly")
}
