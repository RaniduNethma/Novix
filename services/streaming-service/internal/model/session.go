package model

import "time"

// StreamSession tracks an active viewing session in Redis
type StreamSession struct {
	SessionID  string    `json:"sessionId"`
	UserID     string    `json:"userId"`
	VideoID    string    `json:"videoId"`
	Quality    string    `json:"quality"`
	StartedAt  time.Time `json:"startedAt"`
	LastSeenAt time.Time `json:"lastSeenAt"`
	BytesSent  int64     `json:"bytesSent"`
	UserAgent  string    `json:"userAgent"`
	IPAddress  string    `json:"ipAddress"`
}