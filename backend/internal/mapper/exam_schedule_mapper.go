package mapper

import (
	"quizora-backend/internal/dto"
	"quizora-backend/internal/models"
	"time"
)

// ExamScheduleMapper handles conversion between exam models and schedule DTOs
type ExamScheduleMapper struct{}

// NewExamScheduleMapper creates a new exam schedule mapper
func NewExamScheduleMapper() *ExamScheduleMapper {
	return &ExamScheduleMapper{}
}

// ToExamScheduleResponse converts an Exam model to ExamScheduleResponse DTO
func (m *ExamScheduleMapper) ToExamScheduleResponse(exam models.Exam) dto.ExamScheduleResponse {
	return dto.ExamScheduleResponse{
		ID:                    exam.ID,
		Title:                 exam.Title,
		Slug:                  exam.Slug,
		Description:           exam.Description,
		ExamType:              exam.ExamType,
		TotalQuestions:        exam.TotalQuestions,
		DurationMinutes:       exam.DurationMinutes,
		PassingScore:          exam.PassingScore,
		ScheduledStartDate:    m.formatTimePtr(exam.ScheduledStartDate),
		ScheduledEndDate:      m.formatTimePtr(exam.ScheduledEndDate),
		Status:                exam.Status,
		IsActive:              exam.IsActive,
		SortOrder:             0, // This will be set from PackageExam
		AttemptCount:          exam.AttemptCount,
		CompletedAttemptCount: exam.CompletedAttemptCount,
		AverageScore:          exam.AverageScore,
		PassRate:              exam.PassRate,
	}
}

// ToPackageExamScheduleResponse converts a PackageExam model to PackageExamScheduleResponse DTO
func (m *ExamScheduleMapper) ToPackageExamScheduleResponse(packageExam models.PackageExam) dto.PackageExamScheduleResponse {
	examResponse := m.ToExamScheduleResponse(packageExam.Exam)
	// Set the sort order from the package-exam relationship
	examResponse.SortOrder = packageExam.SortOrder

	return dto.PackageExamScheduleResponse{
		Exam: examResponse,
	}
}

// formatTimePtr formats a time pointer to ISO string
func (m *ExamScheduleMapper) formatTimePtr(t *time.Time) *string {
	if t != nil {
		str := t.Format("2006-01-02T15:04:05Z")
		return &str
	}
	return nil
}
