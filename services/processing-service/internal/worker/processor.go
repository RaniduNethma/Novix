package worker

import (
	"context"

	"github.com/novix/services/processing-service/internal/model"
	"github.com/novix/services/processing-service/internal/service"
	"github.com/novix/services/processing-service/pkg/logger"
	"go.uber.org/zap"
)

type Processor struct {
	processingService service.ProcessingService
}

func NewProcessor(processingService service.ProcessingService) *Processor {
	return &Processor{processingService: processingService}
}

// Process routes the job to the correct handler based on job type
func (p *Processor) Process(ctx context.Context, job *model.ProcessingJob) {
	var err error

	switch job.JobType {
	case model.JobTypeHLS, model.JobTypeTranscode:
		err = p.processingService.ProcessVideo(ctx, job)
	case model.JobTypeThumbnail:
		err = p.processingService.GenerateThumbnail(ctx, job)
	default:
		logger.Warn("Unknown job type",
			zap.String("jobType", string(job.JobType)),
		)
		return
	}

	if err != nil {
		logger.Error("Job processing failed",
			zap.Uint("jobID", job.ID),
			zap.String("videoID", job.VideoID),
			zap.Error(err),
		)
		return
	}

	logger.Info("Job completed successfully",
		zap.Uint("jobID", job.ID),
		zap.String("videoID", job.VideoID),
	)
}