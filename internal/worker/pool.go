package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/models"
	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/repositories"
)

type WorkerPool struct {
	jobRepository *repositories.JobRepository
	workerCount   int
	batchSize     int
	interval      time.Duration
}

// Contructor
func NewWorkerPool(
	jobRepository *repositories.JobRepository,
	workerCount int,
	batchSize int,
	interval time.Duration,
) *WorkerPool {
	return &WorkerPool{
		jobRepository: jobRepository,
		workerCount:   workerCount,
		batchSize:     batchSize,
		interval:      interval,
	}
}

func (p *WorkerPool) Start(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	slog.Info(
		"Worker pool started",
		"workers", p.workerCount,
		"batch_size", p.batchSize,
	)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Worker pool stopped")
			return

		case <-ticker.C:
			p.processBatch(ctx)
		}
	}
}

func (p *WorkerPool) processBatch(ctx context.Context) {
	jobs, err := p.jobRepository.FindPending(ctx, p.batchSize)
	if err != nil {
		slog.Error("Failed to fetch pending jobs", "error", err)
		return
	}

	if len(jobs) == 0 {
		return
	}

	slog.Info("Processing job batch", "count", len(jobs))

	jobQueue := make(chan models.Job, len(jobs))

	var wg sync.WaitGroup

	// Start workers for this batch
	for i := 1; i <= p.workerCount; i++ {
		worker := NewWorker(i, p.jobRepository)

		wg.Add(1)

		go func(w *Worker) {
			defer wg.Done()

			for job := range jobQueue {
				w.ProcessJob(ctx, job)
			}
		}(worker)
	}

	// Send jobs to workers
	for _, job := range jobs {
		jobQueue <- job
	}

	// No more jobs will be added to this batch
	close(jobQueue)

	// Wait until all jobs are processed
	wg.Wait()

	slog.Info("Job batch completed")
}