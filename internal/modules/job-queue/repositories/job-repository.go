package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/models"
	"gorm.io/gorm"
)

type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepository {
	return &JobRepository{db: db}
}

// Create a new job.
func (j *JobRepository) Create(
	ctx context.Context,
	job *models.Job,
) error {
	if err := j.db.WithContext(ctx).Create(job).Error; err != nil {
		return fmt.Errorf("create URL: %w", err)
	}

	return nil
}

// Find job by id
func (j *JobRepository) FindByID(
	ctx context.Context,
	id uint,
) (*models.Job, error) {
	var job models.Job

	if err := j.db.WithContext(ctx).First(&job, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, fmt.Errorf("find Job by ID: %w", err)
	}

	return &job, nil
}

// Find pending jobs
func (j *JobRepository) FindPending(
	ctx context.Context,
	limit int,
) ([]models.Job, error) {
	var jobs []models.Job

	if limit <= 0 {
		limit = 10
	}

	if err := j.db.
		WithContext(ctx).
		Where("status = ?", models.JobStatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&jobs).
		Error; err != nil {
		return nil, fmt.Errorf("find pending jobs: %w", err)
	}

	return jobs, nil
}

// Update Job Status
func (j *JobRepository) UpdateStatus(
	ctx context.Context,
	id uint,
	status models.JobStatus,
) error {
	result := j.db.
		WithContext(ctx).
		Model(&models.Job{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("update job status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Claim job repository
func (j *JobRepository) ClaimJob(
	ctx context.Context,
	id uint,
) error {
	result := j.db.
		WithContext(ctx).
		Model(&models.Job{}).
		Where(
			"id = ? AND status = ?",
			id,
			models.JobStatusPending,
		).
		Updates(map[string]any{
			"status": models.JobStatusProcessing,
			"attempts": gorm.Expr("attempts + 1"),
		})

	if result.Error != nil {
		return fmt.Errorf("claim job: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Mark failed
func (j *JobRepository) MarkFailed(
	ctx context.Context,
	id uint,
	errorMessage string,
) error {
	result := j.db.
		WithContext(ctx).
		Model(&models.Job{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status": models.JobStatusFailed,
			"error":  errorMessage,
		})
	if result.Error != nil {
		return fmt.Errorf("mark job as failed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
