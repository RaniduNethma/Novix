package worker

import (
	"context"

	"github.com/novix/services/processing-service/internal/model"
	"github.com/novix/services/processing-service/pkg/logger"
	"go.uber.org/zap"
)

type Pool struct {
	workers   int
	jobQueue  chan *model.ProcessingJob
	processor *Processor
}

func NewPool(workers int, processor *Processor) *Pool {
	return &Pool{
		workers:   workers,
		jobQueue:  make(chan *model.ProcessingJob, workers*2),
		processor: processor,
	}
}

// Start launches the worker goroutines
func (p *Pool) Start(ctx context.Context) {
	logger.Info("Starting worker pool",
		zap.Int("workers", p.workers),
	)
	for i := 0; i < p.workers; i++ {
		go p.work(ctx, i)
	}
}

// Submit adds a job to the queue
func (p *Pool) Submit(job *model.ProcessingJob) {
	p.jobQueue <- job
}

// work is the goroutine that processes jobs
func (p *Pool) work(ctx context.Context, workerID int) {
	logger.Info("Worker started", zap.Int("workerID", workerID))
	for {
		select {
		case job := <-p.jobQueue:
			logger.Info("Worker picked up job",
				zap.Int("workerID", workerID),
				zap.Uint("jobID", job.ID),
				zap.String("videoID", job.VideoID),
			)
			p.processor.Process(ctx, job)
		case <-ctx.Done():
			logger.Info("Worker shutting down",
				zap.Int("workerID", workerID),
			)
			return
		}
	}
}