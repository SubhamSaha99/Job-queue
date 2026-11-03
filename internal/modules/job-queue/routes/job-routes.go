package routes

import (
	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/controllers"
	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/repositories"
	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterJobRoutes(router *gin.RouterGroup, db *gorm.DB) {
	jobRepository := repositories.NewJobRepository(db)
	jobService := services.NewJobService(jobRepository)
	jobController := controllers.NewJobController(jobService)

	router.POST("/", jobController.CreateJob)
	router.GET("/:id", jobController.FindJobById)
	router.GET("/pending/:limit", jobController.FindPendingJobs)
}
