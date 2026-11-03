package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/dtos"
	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JobController struct {
	jobService *services.JobService
}

func NewJobController(jobService *services.JobService) *JobController {
	return &JobController{jobService: jobService}
}

// Create job Controller
func (j *JobController) CreateJob(c *gin.Context) {
	var req dtos.CreateJobDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	job, err := j.jobService.CreateJob(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create job",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "job created successfully",
		"data":    job,
	})
}

// Find job by id Controller
func (j *JobController) FindJobById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID!",
		})
		return
	}

	job, err := j.jobService.FindJobById(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "job not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to find job",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Job fetched successfully",
		"data":    job,
	})
}

// Find pending jobs controller
func (j *JobController) FindPendingJobs(c *gin.Context) {
	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil || limit == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please provide a valid limit!",
		})
		return
	}

	jobs, err := j.jobService.FindPendingJobs(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to find pending job",
		})
		return
	}
	if len(jobs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "There is no pending jobs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pending jobs fetched successfully",
		"data":    jobs,
	})
}

