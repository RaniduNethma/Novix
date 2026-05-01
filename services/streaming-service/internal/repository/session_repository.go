package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/novix/services/streaming-service/internal/model"
	redispkg "github.com/novix/services/streaming-service/pkg/redis"
)

// Key patterns
const (
	sessionKeyPrefix      = "stream:session:%s"        // stream:session:{sessionId}
	userSessionsKeyPrefix = "stream:user:%s:sessions"  // stream:user:{userId}:sessions
	activeStreamsKey       = "stream:active:count"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *model.StreamSession,
		ttl time.Duration) error
	GetSession(ctx context.Context,
		sessionID string) (*model.StreamSession, error)
	UpdateSession(ctx context.Context,
		session *model.StreamSession, ttl time.Duration) error
	DeleteSession(ctx context.Context, sessionID string,
		userID string) error
	GetUserSessionCount(ctx context.Context, userID string) (int, error)
	GetAllUserSessions(ctx context.Context,
		userID string) ([]string, error)
}

type sessionRepository struct {
	redis *redispkg.Client
}

func NewSessionRepository(redis *redispkg.Client) SessionRepository {
	return &sessionRepository{redis: redis}
}

func (r *sessionRepository) CreateSession(ctx context.Context,
	session *model.StreamSession, ttl time.Duration) error {

	// Store session data
	sessionKey := fmt.Sprintf(sessionKeyPrefix, session.SessionID)
	if err := r.redis.Set(ctx, sessionKey, session, ttl); err != nil {
		return fmt.Errorf("failed to store session: %w", err)
	}

	// Track session under user
	userKey := fmt.Sprintf(userSessionsKeyPrefix, session.UserID)
	userSessions, _ := r.GetAllUserSessions(ctx, session.UserID)
	userSessions = append(userSessions, session.SessionID)

	if err := r.redis.Set(ctx, userKey, userSessions, ttl); err != nil {
		return fmt.Errorf("failed to track user session: %w", err)
	}

	return nil
}

func (r *sessionRepository) GetSession(ctx context.Context,
	sessionID string) (*model.StreamSession, error) {

	sessionKey := fmt.Sprintf(sessionKeyPrefix, sessionID)
	var session model.StreamSession

	if err := r.redis.Get(ctx, sessionKey, &session); err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	return &session, nil
}

func (r *sessionRepository) UpdateSession(ctx context.Context,
	session *model.StreamSession, ttl time.Duration) error {

	sessionKey := fmt.Sprintf(sessionKeyPrefix, session.SessionID)
	return r.redis.Set(ctx, sessionKey, session, ttl)
}

func (r *sessionRepository) DeleteSession(ctx context.Context,
	sessionID string, userID string) error {

	// Delete session
	sessionKey := fmt.Sprintf(sessionKeyPrefix, sessionID)
	if err := r.redis.Delete(ctx, sessionKey); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	// Remove from user sessions
	userKey := fmt.Sprintf(userSessionsKeyPrefix, userID)
	sessions, err := r.GetAllUserSessions(ctx, userID)
	if err != nil {
		return nil // already gone
	}

	filtered := []string{}
	for _, s := range sessions {
		if s != sessionID {
			filtered = append(filtered, s)
		}
	}

	return r.redis.Set(ctx, userKey, filtered,
		time.Duration(3600)*time.Second)
}

func (r *sessionRepository) GetUserSessionCount(ctx context.Context,
	userID string) (int, error) {

	sessions, err := r.GetAllUserSessions(ctx, userID)
	if err != nil {
		return 0, nil
	}
	return len(sessions), nil
}

func (r *sessionRepository) GetAllUserSessions(ctx context.Context,
	userID string) ([]string, error) {

	userKey := fmt.Sprintf(userSessionsKeyPrefix, userID)
	var sessions []string

	if err := r.redis.Get(ctx, userKey, &sessions); err != nil {
		return []string{}, nil
	}

	return sessions, nil
}