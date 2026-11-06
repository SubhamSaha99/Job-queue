package jobqueue

import (
	"time"

	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/controllers"
	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/repositories"
	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/services"
	"github.com/SubhamSaha99/job-queue/internal/worker"
	"gorm.io/gorm"
)

type Module struct {
	JobRepository *repositories.JobRepository
	JobService    *services.JobService
	JobController *controllers.JobController
	WorkerPool    *worker.WorkerPool
}

func JobModule(db *gorm.DB) *Module {
	jobRepository := repositories.NewJobRepository(db)
	jobService := services.NewJobService(jobRepository)
	jobController := controllers.NewJobController(jobService)
	workerPool := worker.NewWorkerPool(
		jobRepository,
		3,            
		10,             
		10*time.Second, 
	)

	return &Module{
		JobRepository: jobRepository,
		JobService:    jobService,
		JobController: jobController,
		WorkerPool:    workerPool,
	}
}
