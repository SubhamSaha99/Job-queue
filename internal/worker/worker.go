package worker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/models"
	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/repositories"
)

type Worker struct {
	id            int
	jobRepository *repositories.JobRepository
}

func NewWorker(
	id int,
	jobRepository *repositories.JobRepository,
) *Worker {
	return &Worker{
		id:            id,
		jobRepository: jobRepository,
	}
}

func (w *Worker) ProcessJob(
	ctx context.Context,
	job models.Job,
) {
	// Atomically claim the job.
	if err := w.jobRepository.ClaimJob(ctx, job.ID); err != nil {
		slog.Warn(
			"Job could not be claimed",
			"worker_id", w.id,
			"job_id", job.ID,
			"error", err,
		)
		return
	}

	slog.Info(
		"Worker processing job",
		"worker_id", w.id,
		"job_id", job.ID,
		"type", job.Type,
	)

	// Execute the actual job.
	if err := w.executeJob(&job); err != nil {
		if updateErr := w.jobRepository.MarkFailed(
			ctx,
			job.ID,
			err.Error(),
		); updateErr != nil {
			slog.Error(
				"Failed to mark job as failed",
				"worker_id", w.id,
				"job_id", job.ID,
				"error", updateErr,
			)
		}
		return
	}

	// Mark as completed.
	if err := w.jobRepository.UpdateStatus(
		ctx,
		job.ID,
		models.JobStatusCompleted,
	); err != nil {
		slog.Error(
			"Failed to mark job as completed",
			"worker_id", w.id,
			"job_id", job.ID,
			"error", err,
		)
		return
	}

	slog.Info(
		"Job completed successfully",
		"worker_id", w.id,
		"job_id", job.ID,
	)
}

func (w *Worker) executeJob(
	job *models.Job,
) error {
	slog.Info(
		"Executing job",
		"job_id", job.ID,
		"type", job.Type,
		"payload", job.Payload,
	)

	// Temporary execution logic.
	// Later, you can create handlers for different job types.
	if job.Type == "" {
		return fmt.Errorf("job type is required")
	}

	return nil
}