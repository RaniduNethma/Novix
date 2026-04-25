package repository

import (
	"context"
	"fmt"

	"github.com/novix/services/processing-service/internal/model"
	"gorm.io/gorm"
)

type JobRepository interface {
	Create(ctx context.Context, job *model.ProcessingJob) error
	GetByID(ctx context.Context, id uint) (*model.ProcessingJob, error)
	GetByVideoID(ctx context.Context, videoID string) ([]*model.ProcessingJob, error)
	UpdateStatus(ctx context.Context, id uint, status model.JobStatus,
		errMsg string) error
	Update(ctx context.Context, job *model.ProcessingJob) error
	GetPendingJobs(ctx context.Context, limit int) ([]*model.ProcessingJob, error)
	GetFailedJobs(ctx context.Context,
		maxRetries int) ([]*model.ProcessingJob, error)
}

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Create(ctx context.Context,
	job *model.ProcessingJob) error {
	if err := r.db.WithContext(ctx).Create(job).Error; err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}
	return nil
}

func (r *jobRepository) GetByID(ctx context.Context,
	id uint) (*model.ProcessingJob, error) {
	var job model.ProcessingJob
	if err := r.db.WithContext(ctx).First(&job, id).Error; err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}
	return &job, nil
}

func (r *jobRepository) GetByVideoID(ctx context.Context,
	videoID string) ([]*model.ProcessingJob, error) {
	var jobs []*model.ProcessingJob
	if err := r.db.WithContext(ctx).
		Where("video_id = ?", videoID).
		Find(&jobs).Error; err != nil {
		return nil, fmt.Errorf("failed to get jobs: %w", err)
	}
	return jobs, nil
}

func (r *jobRepository) UpdateStatus(ctx context.Context, id uint,
	status model.JobStatus, errMsg string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if errMsg != "" {
		updates["error_msg"] = errMsg
	}
	if err := r.db.WithContext(ctx).Model(&model.ProcessingJob{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}
	return nil
}

func (r *jobRepository) Update(ctx context.Context,
	job *model.ProcessingJob) error {
	if err := r.db.WithContext(ctx).Save(job).Error; err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}
	return nil
}

func (r *jobRepository) GetPendingJobs(ctx context.Context,
	limit int) ([]*model.ProcessingJob, error) {
	var jobs []*model.ProcessingJob
	if err := r.db.WithContext(ctx).
		Where("status = ?", model.JobStatusPending).
		Limit(limit).
		Order("created_at ASC").
		Find(&jobs).Error; err != nil {
		return nil, fmt.Errorf("failed to get pending jobs: %w", err)
	}
	return jobs, nil
}

func (r *jobRepository) GetFailedJobs(ctx context.Context,
	maxRetries int) ([]*model.ProcessingJob, error) {
	var jobs []*model.ProcessingJob
	if err := r.db.WithContext(ctx).
		Where("status = ? AND retry_count < ?",
			model.JobStatusFailed, maxRetries).
		Find(&jobs).Error; err != nil {
		return nil, fmt.Errorf("failed to get failed jobs: %w", err)
	}
	return jobs, nil
}