package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/SubhamSaha99/job-queue/internal/config"
	"github.com/SubhamSaha99/job-queue/internal/database"
	jobqueue "github.com/SubhamSaha99/job-queue/internal/modules/job-queue"
	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Loading ENVs
	cfg, err := config.LoadEnv()
	if err != nil {
		slog.Error("Load Env", "Error:", err)
		return
	}

	// Connecting Database
	db, err := database.Connect()
	if err != nil {
		slog.Error("DB Startup", "Error:", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Invalid db", "Error:", err)
		return
	}
	defer sqlDB.Close()

	// Create Gin server
	router := gin.Default()

	jobModule := jobqueue.JobModule(db)
	jobAPI := router.Group("/api/v1/jobs")
	
	routes.RegisterJobRoutes(
		jobAPI,
		jobModule.JobController,
	)

	// Starting background worker
	ctx := context.Background()
	go jobModule.WorkerPool.Start(ctx)

	address := fmt.Sprintf(":%s", cfg.Port)
	slog.Info("Server running on", slog.String(":", address))
	if err := router.Run(address); err != nil {
		slog.Error("Server Startup", "Error:", err)
	}

	
}
