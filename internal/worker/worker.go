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
	// Claim the job
	if err := w.jobRepository.ClaimJob(ctx, job.ID); err != nil {
		slog.Warn(
			"Job could not be claimed",
			"worker_id", w.id,
			"job_id", job.ID,
			"error", err,
		)
		return
	}

	// Keep the local job object in sync with the database.
	job.Attempts++

	// Execute the job
	if err := w.executeJob(&job); err != nil {

		// Retry if attempts are remaining
		if job.Attempts < job.MaxAttempts {
			slog.Warn(
				"Job failed, will retry",
				"job_id", job.ID,
				"attempt", job.Attempts,
				"max_attempts", job.MaxAttempts,
				"error", err,
			)

			if err := w.jobRepository.UpdateStatus(
				ctx,
				job.ID,
				models.JobStatusPending,
			); err != nil {
				slog.Error(
					"Failed to reset job status",
					"job_id", job.ID,
					"error", err,
				)
			}

			return
		}

		// No retry attempts remaining
		if err := w.jobRepository.MarkFailed(
			ctx,
			job.ID,
			err.Error(),
		); err != nil {
			slog.Error(
				"Failed to mark job as failed",
				"job_id", job.ID,
				"error", err,
			)
		}

		return
	}

	// Job completed successfully
	if err := w.jobRepository.UpdateStatus(
		ctx,
		job.ID,
		models.JobStatusCompleted,
	); err != nil {
		slog.Error(
			"Failed to mark job as completed",
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

	if job.Type == "" {
		return fmt.Errorf("job type is required")
	}

	return nil
}
