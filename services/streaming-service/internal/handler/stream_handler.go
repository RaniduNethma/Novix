package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/novix/services/streaming-service/internal/model"
	"github.com/novix/services/streaming-service/internal/service"
	"github.com/novix/services/streaming-service/pkg/logger"
	"go.uber.org/zap"
)

type StreamHandler struct {
	streamService  service.StreamService
	sessionService service.SessionService
	presignTTL     time.Duration
	maxSessions    int
}

func NewStreamHandler(
	streamService service.StreamService,
	sessionService service.SessionService,
	presignTTL int,
	maxSessions int,
) *StreamHandler {
	return &StreamHandler{
		streamService:  streamService,
		sessionService: sessionService,
		presignTTL:     time.Duration(presignTTL) * time.Second,
		maxSessions:    maxSessions,
	}
}

// StartStream creates a session and returns the manifest URL
func (h *StreamHandler) StartStream(c *gin.Context) {
	var req model.StreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	userID := c.GetString("userID")

	// Check video exists
	exists, err := h.streamService.VideoExists(c.Request.Context(),
		req.VideoID)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Video not found or not ready",
		})
		return
	}

	// Validate concurrent session limit
	if err := h.sessionService.ValidateSessionLimit(
		c.Request.Context(), userID, h.maxSessions); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create stream session
	session, err := h.sessionService.CreateSession(
		c.Request.Context(),
		userID,
		req.VideoID,
		req.Quality,
		c.GetHeader("User-Agent"),
		c.ClientIP(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create stream session",
		})
		return
	}

	// Generate presigned manifest URL
	manifestURL, err := h.streamService.GetPresignedManifestURL(
		c.Request.Context(), req.VideoID, h.presignTTL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate stream URL",
		})
		return
	}

	// Get thumbnail URL
	thumbnailURL, _ := h.streamService.GetThumbnailURL(
		c.Request.Context(), req.VideoID, h.presignTTL,
	)

	logger.Info("Stream started",
		zap.String("sessionID", session.SessionID),
		zap.String("videoID", req.VideoID),
		zap.String("userID", userID),
	)

	c.JSON(http.StatusOK, model.StreamResponse{
		SessionID:    session.SessionID,
		ManifestURL:  manifestURL,
		ThumbnailURL: thumbnailURL,
		ExpiresIn:    int(h.presignTTL.Seconds()),
	})
}

// GetManifest serves the HLS master or quality playlist directly
func (h *StreamHandler) GetManifest(c *gin.Context) {
	videoID := c.Param("videoId")
	quality := c.DefaultQuery("quality", "master")

	// Refresh session if provided
	sessionID := c.GetHeader("Session-Id")
	if sessionID != "" {
		_ = h.sessionService.RefreshSession(c.Request.Context(), sessionID)
	}

	data, err := h.streamService.GetManifest(
		c.Request.Context(), videoID, quality,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Manifest not found",
		})
		return
	}

	c.Header("Content-Type", "application/vnd.apple.mpegurl")
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "application/vnd.apple.mpegurl", data)
}

// GetSegment serves a specific HLS .ts segment
func (h *StreamHandler) GetSegment(c *gin.Context) {
	videoID := c.Param("videoId")
	quality := c.Param("quality")
	segment := c.Param("segment")

	data, err := h.streamService.GetSegment(
		c.Request.Context(), videoID, quality, segment,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Segment not found",
		})
		return
	}

	c.Header("Content-Type", "video/MP2T")
	c.Header("Cache-Control", "max-age=3600")
	c.Data(http.StatusOK, "video/MP2T", data)
}

// EndStream terminates a streaming session
func (h *StreamHandler) EndStream(c *gin.Context) {
	sessionID := c.Param("sessionId")
	userID := c.GetString("userID")

	if err := h.sessionService.EndSession(
		c.Request.Context(), sessionID, userID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to end session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stream session ended",
	})
}