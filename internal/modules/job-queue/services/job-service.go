package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/dtos"
	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/models"
	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/repositories"
	"gorm.io/gorm"
)

type JobService struct {
	jobRepository *repositories.JobRepository
}

func NewJobService(jobRepository *repositories.JobRepository) *JobService {
	return &JobService{jobRepository: jobRepository}
}

// Create new job service
func (j *JobService) CreateJob(
	ctx context.Context,
	jobDTO *dtos.CreateJobDTO,
) (*models.Job, error) {

	job := &models.Job{
		Type:        jobDTO.Type,
		Payload:     jobDTO.Payload,
		Status:      models.JobStatusPending,
		Attempts:    0,
		MaxAttempts: jobDTO.MaxAttempts,
	}

	if err := j.jobRepository.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("create job: %w", err)
	}

	return job, nil
}

// Find job by id
func (j *JobService) FindJobById(
	ctx context.Context,
	id uint,
) (*models.Job, error) {
	job, err := j.jobRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("Find job: %w", err)
	}
	return job, nil
}

