package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/novix/services/streaming-service/internal/model"
	"github.com/novix/services/streaming-service/internal/repository"
	"github.com/novix/services/streaming-service/pkg/logger"
	"go.uber.org/zap"
)

type SessionService interface {
	CreateSession(ctx context.Context, userID string, videoID string,
		quality string, userAgent string,
		ipAddress string) (*model.StreamSession, error)
	GetSession(ctx context.Context,
		sessionID string) (*model.StreamSession, error)
	EndSession(ctx context.Context, sessionID string, userID string) error
	ValidateSessionLimit(ctx context.Context, userID string,
		maxSessions int) error
	RefreshSession(ctx context.Context, sessionID string) error
}

type sessionService struct {
	sessionRepo repository.SessionRepository
	sessionTTL  time.Duration
}

func NewSessionService(sessionRepo repository.SessionRepository,
	sessionTTL int) SessionService {
	return &sessionService{
		sessionRepo: sessionRepo,
		sessionTTL:  time.Duration(sessionTTL) * time.Second,
	}
}

func (s *sessionService) CreateSession(ctx context.Context,
	userID string, videoID string, quality string,
	userAgent string, ipAddress string) (*model.StreamSession, error) {

	session := &model.StreamSession{
		SessionID:  uuid.New().String(),
		UserID:     userID,
		VideoID:    videoID,
		Quality:    quality,
		StartedAt:  time.Now(),
		LastSeenAt: time.Now(),
		UserAgent:  userAgent,
		IPAddress:  ipAddress,
	}

	if err := s.sessionRepo.CreateSession(ctx, session,
		s.sessionTTL); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	logger.Info("Stream session created",
		zap.String("sessionID", session.SessionID),
		zap.String("userID", userID),
		zap.String("videoID", videoID),
	)

	return session, nil
}

func (s *sessionService) GetSession(ctx context.Context,
	sessionID string) (*model.StreamSession, error) {
	return s.sessionRepo.GetSession(ctx, sessionID)
}

func (s *sessionService) EndSession(ctx context.Context,
	sessionID string, userID string) error {

	if err := s.sessionRepo.DeleteSession(ctx, sessionID,
		userID); err != nil {
		return fmt.Errorf("failed to end session: %w", err)
	}

	logger.Info("Stream session ended",
		zap.String("sessionID", sessionID),
		zap.String("userID", userID),
	)

	return nil
}

func (s *sessionService) ValidateSessionLimit(ctx context.Context,
	userID string, maxSessions int) error {

	count, err := s.sessionRepo.GetUserSessionCount(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get session count: %w", err)
	}

	if count >= maxSessions {
		return fmt.Errorf(
			"maximum concurrent streams reached (%d/%d)",
			count, maxSessions,
		)
	}

	return nil
}

func (s *sessionService) RefreshSession(ctx context.Context,
	sessionID string) error {

	session, err := s.sessionRepo.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	session.LastSeenAt = time.Now()

	return s.sessionRepo.UpdateSession(ctx, session, s.sessionTTL)
}