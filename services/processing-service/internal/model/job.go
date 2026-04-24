package model

import "time"

type JobStatus string
type JobType string

const(
	JobStatusPending JobStatus = "PENDING"
	JobStatusProcessing JobStatus = "PROCESSING"
	JobStatusCompleted JobStatus = "COMPLETED"
	JobStatusFailed JobStatus = "FAILED"
	JobStatusRetrying JobStatus = "RETRYING"
)

const(
	JobTypeTranscode JobType = "TRANSCODE"
	JobTypeThumbnail JobType = "THUMBNAIL"
	JobTypeHLS JobType = "HLS"
)

type ProcessingJob struct{
	ID uint `gorm:"primaryKey;autoIncrement"`
	VideoID string `gorm:"not null;index"`
	UserID string `gorm:"not null"`
	JobType JobType `gorm:"not null"`
	Status JobStatus `gorm:"not null;default:'PENDING'"`
	RawPath string `gorm:"not null"`
	OutputPath string
	RetryCount int `gorm:"default:0"`
	ErrorMsg string
	StartedAt *time.Time
	CompletedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
