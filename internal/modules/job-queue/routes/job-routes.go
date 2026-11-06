package routes

import (
	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterJobRoutes(
	router *gin.RouterGroup,
	jobController *controllers.JobController,
) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	router.POST("/", jobController.CreateJob)
	router.GET("/:id", jobController.FindJobById)
}
