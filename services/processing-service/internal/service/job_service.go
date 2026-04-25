package service

import (
	"context"
	"fmt"

	"github.com/novix/services/processing-service/internal/model"
	"github.com/novix/services/processing-service/internal/repository"
)

type JobService interface {
	CreateJob(ctx context.Context, videoID string, userID string,
		jobType model.JobType, rawPath string) (*model.ProcessingJob, error)
	GetJobByID(ctx context.Context, id uint) (*model.ProcessingJob, error)
	GetJobsByVideoID(ctx context.Context,
		videoID string) ([]*model.ProcessingJob, error)
}

type jobService struct {
	jobRepo repository.JobRepository
}

func NewJobService(jobRepo repository.JobRepository) JobService {
	return &jobService{jobRepo: jobRepo}
}

func (s *jobService) CreateJob(ctx context.Context,
	videoID string, userID string,
	jobType model.JobType, rawPath string) (*model.ProcessingJob, error) {
	job := &model.ProcessingJob{
		VideoID: videoID,
		UserID:  userID,
		JobType: jobType,
		Status:  model.JobStatusPending,
		RawPath: rawPath,
	}
	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}
	return job, nil
}

func (s *jobService) GetJobByID(ctx context.Context,
	id uint) (*model.ProcessingJob, error) {
	return s.jobRepo.GetByID(ctx, id)
}

func (s *jobService) GetJobsByVideoID(ctx context.Context,
	videoID string) ([]*model.ProcessingJob, error) {
	return s.jobRepo.GetByVideoID(ctx, videoID)
}