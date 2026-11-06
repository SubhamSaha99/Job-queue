package dtos

type CreateJobDTO struct {
	Type        string `json:"type" binding:"required"`
	Payload     string `json:"payload" binding:"required"`
	MaxAttempts int    `json:"max_attempts" binding:"required,gt=0,lt=10"`
}

type FindJobByIdDTO struct {
	Id uint `json:"id" binding:"required"`
}
