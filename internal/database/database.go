package database

import (
	"fmt"

	"github.com/SubhamSaha99/job-queue/internal/modules/job-queue/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	db, err := gorm.Open(
		sqlite.Open("job-queue.db"),
		&gorm.Config{TranslateError: true},
	)

	if err != nil {
		return nil, fmt.Errorf("connect to SQLite: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get database connection: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping SQLite: %w", err)
	}

	if err := db.AutoMigrate(&models.Job{}); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("migrate URL model: %w", err)
	}

	return db, nil
}
