package dto

import "quizora-backend/internal/models"

// ExamScheduleResponse represents minimal exam data for package details page
type ExamScheduleResponse struct {
	ID                 uint              `json:"id"`
	Title              string            `json:"title"`
	Slug               string            `json:"slug"`
	Description        *string           `json:"description"`
	ExamType           models.ExamType   `json:"exam_type"`
	TotalQuestions     int               `json:"total_questions"`
	DurationMinutes    int               `json:"duration_minutes"`
	PassingScore       float64           `json:"passing_score"`
	ScheduledStartDate *string           `json:"scheduled_start_date"`
	ScheduledEndDate   *string           `json:"scheduled_end_date"`
	Status             models.ExamStatus `json:"status"`
	IsActive           bool              `json:"is_active"`
	SortOrder          int               `json:"sort_order"`
	// Basic analytics only
	AttemptCount          int      `json:"attempt_count"`
	CompletedAttemptCount int      `json:"completed_attempt_count"`
	AverageScore          *float64 `json:"average_score"`
	PassRate              *float64 `json:"pass_rate"`
}

// PackageExamScheduleResponse represents the package-exam relationship for scheduling
type PackageExamScheduleResponse struct {
	Exam ExamScheduleResponse `json:"exam"`
}
