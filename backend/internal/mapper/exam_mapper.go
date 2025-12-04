package mapper

import (
	"encoding/json"
	"github.com/Mahfuz2811/medecole/backend/internal/dto"
	"github.com/Mahfuz2811/medecole/backend/internal/models"
	"github.com/Mahfuz2811/medecole/backend/internal/repository"
	"time"
)

// ExamMapper handles conversions between exam models and DTOs
type ExamMapper interface {
	ToExamResponse(exam repository.ExamWithUserData) dto.ExamResponse
	ToExamListResponse(exams []repository.ExamWithUserData) dto.ExamListResponse
	ToExamListResponseWithPackage(pkg models.Package, exams []repository.ExamWithUserData) dto.ExamListResponse
	ToExamContentResponse(exam repository.ExamWithUserData) dto.ExamContentResponse
	ToExamMetaResponse(exam models.Exam) dto.ExamMetaResponse
	ToExamSessionResponse(attempt models.UserExamAttempt, exam models.Exam) dto.ExamSessionResponse
	ToExamSessionResponseWithAnswers(attempt models.UserExamAttempt, exam models.Exam, savedAnswers []dto.UserAnswerResponse) dto.ExamSessionResponse
}

// examMapper implements ExamMapper
type examMapper struct{}

// NewExamMapper creates a new exam mapper
func NewExamMapper() ExamMapper {
	return &examMapper{}
}

// ToExamResponse converts ExamWithUserData to ExamResponse
func (m *examMapper) ToExamResponse(exam repository.ExamWithUserData) dto.ExamResponse {
	// Handle nil values with defaults
	description := ""
	if exam.Description != nil {
		description = *exam.Description
	}

	averageScore := 0.0
	if exam.AverageScore != nil {
		averageScore = *exam.AverageScore
	}

	passRate := 0.0
	if exam.PassRate != nil {
		passRate = *exam.PassRate
	}

	// Format scheduled dates to ISO strings
	var scheduledStartDate *string
	if exam.ScheduledStartDate != nil {
		str := exam.ScheduledStartDate.Format("2006-01-02T15:04:05Z")
		scheduledStartDate = &str
	}

	var scheduledEndDate *string
	if exam.ScheduledEndDate != nil {
		str := exam.ScheduledEndDate.Format("2006-01-02T15:04:05Z")
		scheduledEndDate = &str
	}

	// Prepare user attempt data if available
	var userAttempt *dto.UserAttemptResponse
	if exam.UserAttemptID != nil {
		userAttempt = &dto.UserAttemptResponse{
			ID:             *exam.UserAttemptID,
			Status:         *exam.UserAttemptStatus,
			Score:          exam.AttemptScore,
			CorrectAnswers: exam.AttemptCorrectAnswers,
			IsPassed:       exam.AttemptPassed,
			TimeSpent:      exam.ActualTimeSpent,
			SessionID:      exam.SessionID,
		}

		// Format started_at if available
		if exam.UserAttemptStartedAt != nil {
			userAttempt.StartedAt = exam.UserAttemptStartedAt.Format("2006-01-02T15:04:05Z")
		}

		// Format completed_at if available
		if exam.UserAttemptCompletedAt != nil {
			completedAtStr := exam.UserAttemptCompletedAt.Format("2006-01-02T15:04:05Z")
			userAttempt.CompletedAt = completedAtStr
		}
	}

	return dto.ExamResponse{
		ID:                 int(exam.ID),
		Title:              exam.Title,
		Slug:               exam.Slug,
		Description:        description,
		ExamType:           string(exam.ExamType),
		TotalQuestions:     exam.TotalQuestions,
		DurationMinutes:    exam.DurationMinutes,
		PassingScore:       exam.PassingScore,
		TotalMarks:         exam.TotalMarks,
		MaxAttempts:        exam.MaxAttempts,
		Instructions:       exam.Instructions,
		ScheduledStartDate: scheduledStartDate,
		ScheduledEndDate:   scheduledEndDate,
		AttemptCount:       exam.AttemptCount,
		AverageScore:       averageScore,
		PassRate:           passRate,
		ComputedStatus:     exam.ComputedStatus,
		SortOrder:          exam.SortOrder,
		HasAttempted:       exam.HasAttempted,
		UserAttempt:        userAttempt,
	}
}

// ToExamListResponse converts a slice of ExamWithUserData to ExamListResponse (no pagination)
func (m *examMapper) ToExamListResponse(exams []repository.ExamWithUserData) dto.ExamListResponse {
	examResponses := make([]dto.ExamResponse, len(exams))
	for i, exam := range exams {
		examResponses[i] = m.ToExamResponse(exam)
	}

	return dto.ExamListResponse{
		Exams: examResponses,
	}
}

// ToExamListResponseWithPackage converts package and exams to ExamListResponse with package info
func (m *examMapper) ToExamListResponseWithPackage(pkg models.Package, exams []repository.ExamWithUserData) dto.ExamListResponse {
	// Convert exams
	examResponses := make([]dto.ExamResponse, len(exams))
	for i, exam := range exams {
		examResponses[i] = m.ToExamResponse(exam)
	}

	// Convert package
	packageInfo := dto.PackageInfoResponse{
		ID:                    pkg.ID,
		Name:                  pkg.Name,
		Slug:                  pkg.Slug,
		Description:           pkg.Description,
		PackageType:           pkg.PackageType,
		Price:                 pkg.Price,
		ValidityType:          pkg.ValidityType,
		ValidityDays:          pkg.ValidityDays,
		ValidityDate:          nil, // TODO: Format date if needed
		TotalExams:            pkg.TotalExams,
		EnrollmentCount:       pkg.EnrollmentCount,
		ActiveEnrollmentCount: 0, // TODO: Calculate or fetch from separate query
	}

	// Format validity date if present
	if pkg.ValidityDate != nil {
		validityDateStr := pkg.ValidityDate.Format("2006-01-02")
		packageInfo.ValidityDate = &validityDateStr
	}

	return dto.ExamListResponse{
		Package: packageInfo,
		Exams:   examResponses,
	}
}

// ToExamContentResponse converts ExamWithUserData to ExamContentResponse including questions
func (m *examMapper) ToExamContentResponse(exam repository.ExamWithUserData) dto.ExamContentResponse {
	// Parse questions from JSON
	var questions []dto.ExamQuestionResponse
	if exam.QuestionsData != "" {
		var examQuestions []models.ExamQuestion
		if err := json.Unmarshal([]byte(exam.QuestionsData), &examQuestions); err == nil {
			questions = make([]dto.ExamQuestionResponse, len(examQuestions))
			for i, q := range examQuestions {
				questions[i] = dto.ExamQuestionResponse{
					ID:           q.ID,
					QuestionText: q.QuestionText,
					QuestionType: string(q.QuestionType),
					Options:      q.Options,
					Points:       q.Points,
				}
			}
		}
	}

	// Calculate remaining time for scheduled exams
	var remainingTime *int
	if exam.ScheduledStartDate != nil && exam.ScheduledEndDate != nil {
		now := time.Now()
		if now.Before(*exam.ScheduledEndDate) {
			remaining := int(exam.ScheduledEndDate.Sub(now).Seconds())
			if remaining > 0 {
				remainingTime = &remaining
			}
		}
	}

	// Determine if user can start exam
	canStartExam := exam.ComputedStatus == "AVAILABLE" || exam.ComputedStatus == "LIVE"

	return dto.ExamContentResponse{
		ID:              exam.ID,
		Title:           exam.Title,
		Slug:            exam.Slug,
		Description:     exam.Description,
		ExamType:        string(exam.ExamType),
		TotalQuestions:  exam.TotalQuestions,
		DurationMinutes: exam.DurationMinutes,
		PassingScore:    exam.PassingScore,
		TotalMarks:      exam.TotalMarks,
		MaxAttempts:     exam.MaxAttempts,
		Instructions:    exam.Instructions,
		Questions:       questions,
		ComputedStatus:  exam.ComputedStatus,
		CanStartExam:    canStartExam,
		HasAttempted:    exam.HasAttempted,
		RemainingTime:   remainingTime,
	}
}

// ToExamMetaResponse converts Exam model to ExamMetaResponse for session initialization
func (m *examMapper) ToExamMetaResponse(exam models.Exam) dto.ExamMetaResponse {
	return dto.ExamMetaResponse{
		ID:              exam.ID,
		Title:           exam.Title,
		Slug:            exam.Slug,
		DurationMinutes: exam.DurationMinutes,
		TotalQuestions:  exam.TotalQuestions,
		PassingScore:    exam.PassingScore,
		TotalMarks:      exam.TotalMarks,
		MaxAttempts:     exam.MaxAttempts,
		Instructions:    exam.Instructions,
	}
}

// ToExamSessionResponse converts attempt and exam models to ExamSessionResponse
func (m *examMapper) ToExamSessionResponse(attempt models.UserExamAttempt, exam models.Exam) dto.ExamSessionResponse {
	// Parse questions from JSON and create secure versions (without answers)
	var questions []dto.SecureExamQuestionResponse
	if exam.QuestionsData != "" {
		var examQuestions []models.ExamQuestion
		if err := json.Unmarshal([]byte(exam.QuestionsData), &examQuestions); err == nil {
			questions = make([]dto.SecureExamQuestionResponse, len(examQuestions))
			for i, q := range examQuestions {
				// Create secure options without 'is_correct' field
				secureOptions := make(map[string]interface{})
				if q.Options != nil {
					for key, value := range q.Options {
						if optionMap, ok := value.(map[string]interface{}); ok {
							// Only include 'text' field, exclude 'is_correct'
							if text, exists := optionMap["text"]; exists {
								secureOptions[key] = map[string]interface{}{
									"text": text,
								}
							}
						}
					}
				}

				questions[i] = dto.SecureExamQuestionResponse{
					ID:           q.ID,
					QuestionText: q.QuestionText,
					QuestionType: string(q.QuestionType),
					Options:      secureOptions,
					Points:       q.Points,
				}
			}
		}
	}

	// Format timestamps
	startedAt := attempt.StartedAt.Format("2006-01-02T15:04:05Z")
	lastActivity := ""
	if attempt.LastActivityAt != nil {
		lastActivity = attempt.LastActivityAt.Format("2006-01-02T15:04:05Z")
	}

	return dto.ExamSessionResponse{
		Exam: dto.ExamSessionExamData{
			ID:              exam.ID,
			Title:           exam.Title,
			Slug:            exam.Slug,
			Description:     exam.Description,
			ExamType:        string(exam.ExamType),
			DurationMinutes: exam.DurationMinutes,
			PassingScore:    exam.PassingScore,
			TotalMarks:      exam.TotalMarks,
			MaxAttempts:     exam.MaxAttempts,
			Instructions:    exam.Instructions,
			Questions:       questions,
		},
		Session: dto.ExamSessionSessionData{
			SessionID:        attempt.GetSessionKey(),
			AttemptID:        attempt.ID,
			Status:           string(attempt.Status),
			TimeRemaining:    attempt.GetRemainingTimeSeconds(),
			TimeLimitSeconds: attempt.TimeLimitSeconds,
			StartedAt:        startedAt,
			CanSubmit:        attempt.IsInProgress() && !attempt.IsTimeExpired(),
			CanPause:         attempt.IsInProgress(),
			LastActivity:     lastActivity,
		},
	}
}

// ToExamSessionResponseWithAnswers converts attempt and exam models to ExamSessionResponse with saved answers
func (m *examMapper) ToExamSessionResponseWithAnswers(attempt models.UserExamAttempt, exam models.Exam, savedAnswers []dto.UserAnswerResponse) dto.ExamSessionResponse {
	// Parse questions from JSON and create secure versions (without answers)
	var questions []dto.SecureExamQuestionResponse
	if exam.QuestionsData != "" {
		var examQuestions []models.ExamQuestion
		if err := json.Unmarshal([]byte(exam.QuestionsData), &examQuestions); err == nil {
			questions = make([]dto.SecureExamQuestionResponse, len(examQuestions))
			for i, q := range examQuestions {
				// Create secure options without 'is_correct' field
				secureOptions := make(map[string]interface{})
				if q.Options != nil {
					for key, value := range q.Options {
						if optionMap, ok := value.(map[string]interface{}); ok {
							// Only include 'text' field, exclude 'is_correct'
							if text, exists := optionMap["text"]; exists {
								secureOptions[key] = map[string]interface{}{
									"text": text,
								}
							}
						}
					}
				}

				questions[i] = dto.SecureExamQuestionResponse{
					ID:           q.ID,
					QuestionText: q.QuestionText,
					QuestionType: string(q.QuestionType),
					Options:      secureOptions,
					Points:       q.Points,
				}
			}
		}
	}

	// Format timestamps
	startedAt := attempt.StartedAt.Format("2006-01-02T15:04:05Z")
	lastActivity := ""
	if attempt.LastActivityAt != nil {
		lastActivity = attempt.LastActivityAt.Format("2006-01-02T15:04:05Z")
	}

	return dto.ExamSessionResponse{
		Exam: dto.ExamSessionExamData{
			ID:              exam.ID,
			Title:           exam.Title,
			Slug:            exam.Slug,
			Description:     exam.Description,
			ExamType:        string(exam.ExamType),
			DurationMinutes: exam.DurationMinutes,
			PassingScore:    exam.PassingScore,
			TotalMarks:      exam.TotalMarks,
			MaxAttempts:     exam.MaxAttempts,
			Instructions:    exam.Instructions,
			Questions:       questions,
		},
		Session: dto.ExamSessionSessionData{
			SessionID:        attempt.GetSessionKey(),
			AttemptID:        attempt.ID,
			Status:           string(attempt.Status),
			TimeRemaining:    attempt.GetRemainingTimeSeconds(),
			TimeLimitSeconds: attempt.TimeLimitSeconds,
			StartedAt:        startedAt,
			CanSubmit:        attempt.IsInProgress() && !attempt.IsTimeExpired(),
			CanPause:         attempt.IsInProgress(),
			LastActivity:     lastActivity,
			SavedAnswers:     savedAnswers, // Include saved answers from cache
		},
	}
}
