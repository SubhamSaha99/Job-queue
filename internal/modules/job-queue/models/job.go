package models

import "time"

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type Job struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Type        string    `gorm:"not null" json:"type"`
	Payload     string    `gorm:"type:text;not null" json:"payload"`
	Status      JobStatus `gorm:"not null;index" json:"status"`
	Attempts    int       `gorm:"default:0" json:"attempts"`
	MaxAttempts int       `gorm:"default:3" json:"max_attempts"`
	Error       string    `gorm:"type:text" json:"error,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
