package dto

import "github.com/Mahfuz2811/medecole/backend/internal/models"

// ExamResponse represents the API response structure for package exams
type ExamResponse struct {
	ID                 int                  `json:"id"`
	Title              string               `json:"title"`
	Slug               string               `json:"slug"`
	Description        string               `json:"description"`
	ExamType           string               `json:"exam_type"`
	TotalQuestions     int                  `json:"total_questions"`
	DurationMinutes    int                  `json:"duration_minutes"`
	PassingScore       float64              `json:"passing_score"`
	TotalMarks         float64              `json:"total_marks"`
	MaxAttempts        int                  `json:"max_attempts"`
	Instructions       *string              `json:"instructions"`
	ScheduledStartDate *string              `json:"scheduled_start_date,omitempty"`
	ScheduledEndDate   *string              `json:"scheduled_end_date,omitempty"`
	AttemptCount       int                  `json:"attempt_count"`
	AverageScore       float64              `json:"average_score"`
	PassRate           float64              `json:"pass_rate"`
	ComputedStatus     string               `json:"computed_status"`
	SortOrder          int                  `json:"sort_order"`
	HasAttempted       bool                 `json:"has_attempted"`
	UserAttempt        *UserAttemptResponse `json:"user_attempt,omitempty"`
}

// UserAttemptResponse represents user attempt data for an exam
type UserAttemptResponse struct {
	ID             uint     `json:"id"`
	Status         string   `json:"status"`
	StartedAt      string   `json:"started_at"`
	CompletedAt    string   `json:"completed_at,omitempty"`
	Score          *float64 `json:"score,omitempty"`
	CorrectAnswers *int     `json:"correct_answers,omitempty"`
	IsPassed       *bool    `json:"is_passed,omitempty"`
	TimeSpent      *int     `json:"time_spent,omitempty"`
	SessionID      *string  `json:"session_id,omitempty"`
}

// PackageInfoResponse represents minimal package data needed for exam listing
type PackageInfoResponse struct {
	ID                    uint                `json:"id"`
	Name                  string              `json:"name"`
	Slug                  string              `json:"slug"`
	Description           *string             `json:"description"`
	PackageType           models.PackageType  `json:"package_type"`
	Price                 float64             `json:"price"`
	ValidityType          models.ValidityType `json:"validity_type"`
	ValidityDays          *int                `json:"validity_days,omitempty"`
	ValidityDate          *string             `json:"validity_date,omitempty"`
	TotalExams            int                 `json:"total_exams"`
	EnrollmentCount       int                 `json:"enrollment_count"`
	ActiveEnrollmentCount int                 `json:"active_enrollment_count"`
}

// ExamListResponse represents the enhanced response that includes both package info and exams
type ExamListResponse struct {
	Package PackageInfoResponse `json:"package"`
	Exams   []ExamResponse      `json:"exams"`
}

// ExamQuestionResponse represents a question in the exam content (includes answers - for admin use)
type ExamQuestionResponse struct {
	ID           uint                   `json:"id"`
	QuestionText string                 `json:"question_text"`
	QuestionType string                 `json:"question_type"`
	Options      map[string]interface{} `json:"options"`
	Points       int                    `json:"points"`
}

// SecureExamQuestionResponse represents a question without answer information (for active exam sessions)
type SecureExamQuestionResponse struct {
	ID           uint                   `json:"id"`
	QuestionText string                 `json:"question_text"`
	QuestionType string                 `json:"question_type"`
	Options      map[string]interface{} `json:"options"` // Only contains 'text' field, no 'is_correct'
	Points       int                    `json:"points"`
}

// ExamContentResponse represents the full exam content including questions
type ExamContentResponse struct {
	ID              uint                   `json:"id"`
	Title           string                 `json:"title"`
	Slug            string                 `json:"slug"`
	Description     *string                `json:"description"`
	ExamType        string                 `json:"exam_type"`
	TotalQuestions  int                    `json:"total_questions"`
	DurationMinutes int                    `json:"duration_minutes"`
	PassingScore    float64                `json:"passing_score"`
	TotalMarks      float64                `json:"total_marks"`
	MaxAttempts     int                    `json:"max_attempts"`
	Instructions    *string                `json:"instructions"`
	Questions       []ExamQuestionResponse `json:"questions"`
	ComputedStatus  string                 `json:"computed_status"`
	CanStartExam    bool                   `json:"can_start_exam"`
	HasAttempted    bool                   `json:"has_attempted"`
	RemainingTime   *int                   `json:"remaining_time,omitempty"` // in seconds, for scheduled exams
}

// StartExamRequest represents the request to start an exam
type StartExamRequest struct {
	PackageSlug string            `json:"package_slug" binding:"required"`
	DeviceInfo  map[string]string `json:"device_info,omitempty"`
}

// StartExamResponse represents the response when starting an exam
type StartExamResponse struct {
	SessionID string           `json:"session_id"`
	AttemptID uint             `json:"attempt_id"`
	ExamMeta  ExamMetaResponse `json:"exam_meta"`
}

// ExamMetaResponse represents basic exam metadata for session initialization
type ExamMetaResponse struct {
	ID              uint    `json:"id"`
	Title           string  `json:"title"`
	Slug            string  `json:"slug"`
	DurationMinutes int     `json:"duration_minutes"`
	TotalQuestions  int     `json:"total_questions"`
	PassingScore    float64 `json:"passing_score"`
	TotalMarks      float64 `json:"total_marks"`
	MaxAttempts     int     `json:"max_attempts"`
	Instructions    *string `json:"instructions"`
}

// ExamSessionResponse represents the complete exam session data for active exams
type ExamSessionResponse struct {
	Exam    ExamSessionExamData    `json:"exam"`
	Session ExamSessionSessionData `json:"session"`
}

// ExamSessionExamData represents the exam content and questions for an active session
type ExamSessionExamData struct {
	ID              uint                         `json:"id"`
	Title           string                       `json:"title"`
	Slug            string                       `json:"slug"`
	Description     *string                      `json:"description"`
	ExamType        string                       `json:"exam_type"`
	DurationMinutes int                          `json:"duration_minutes"`
	PassingScore    float64                      `json:"passing_score"`
	TotalMarks      float64                      `json:"total_marks"`
	MaxAttempts     int                          `json:"max_attempts"`
	Instructions    *string                      `json:"instructions"`
	Questions       []SecureExamQuestionResponse `json:"questions"` // Secure questions without answers
}

// ExamSessionSessionData represents the current session state and timing
type ExamSessionSessionData struct {
	SessionID        string               `json:"session_id"`
	AttemptID        uint                 `json:"attempt_id"`
	Status           string               `json:"status"`
	TimeRemaining    int                  `json:"time_remaining"` // in seconds
	TimeLimitSeconds int                  `json:"time_limit_seconds"`
	StartedAt        string               `json:"started_at"`
	CanSubmit        bool                 `json:"can_submit"`
	CanPause         bool                 `json:"can_pause"`
	LastActivity     string               `json:"last_activity"`
	SavedAnswers     []UserAnswerResponse `json:"saved_answers,omitempty"` // Cached answers from Redis
}

// UserAnswerResponse represents a saved user answer from cache
type UserAnswerResponse struct {
	QuestionID     uint   `json:"question_id"`
	SelectedOption string `json:"selected_option"` // Single option or JSON array string for multiple options
}

// SyncSessionRequest represents the request to sync user answers during an exam session
type SyncSessionRequest struct {
	Answers []UserAnswerSync `json:"answers"`
}

// UserAnswerSync represents a single answer to sync (question ID and selected option)
type UserAnswerSync struct {
	QuestionID     uint   `json:"question_id" binding:"required"`
	SelectedOption string `json:"selected_option" binding:"required"` // e.g., "a", "b", "c", "d", "e"
}

// SyncSessionResponse represents the response after syncing answers
type SyncSessionResponse struct {
	Success       bool   `json:"success"`
	SyncedCount   int    `json:"synced_count"`
	LastSyncAt    string `json:"last_sync_at"`
	TimeRemaining int    `json:"time_remaining"` // Current remaining time in seconds
}

// ExamAlreadySubmittedError represents error details when exam is already submitted
type ExamAlreadySubmittedError struct {
	Error   string               `json:"error"`
	Message string               `json:"message"`
	Data    ExamSubmittedDetails `json:"data"`
}

// ExamSubmittedDetails contains details about the previous submission
type ExamSubmittedDetails struct {
	PreviousAttempt PreviousAttemptInfo `json:"previous_attempt"`
	CanRetry        bool                `json:"can_retry"`
}

// PreviousAttemptInfo contains information about the previous attempt
type PreviousAttemptInfo struct {
	AttemptID   uint     `json:"attempt_id"`
	SubmittedAt string   `json:"submitted_at"`
	Score       *float64 `json:"score"`
	Status      string   `json:"status"`
}

// SubmitExamRequest represents the request to submit an exam
type SubmitExamRequest struct {
	SessionID string `json:"session_id" binding:"required"`
}

// SubmitExamResponse represents the response after submitting an exam
type SubmitExamResponse struct {
	SessionID        string  `json:"session_id"`
	Score            float64 `json:"score"`
	Passed           bool    `json:"passed"`
	TotalQuestions   int     `json:"total_questions"`
	CorrectAnswers   int     `json:"correct_answers"`
	TimeTakenSeconds int     `json:"time_taken_seconds"`
	SubmittedAt      string  `json:"submitted_at"`
}

// ExamResultResponse represents the complete exam result data for the results page
type ExamResultResponse struct {
	Exam      ExamResultExamData    `json:"exam"`
	Attempt   ExamResultAttemptData `json:"attempt"`
	Package   ExamResultPackageData `json:"package"`
	Questions []ExamResultQuestion  `json:"questions"`
}

// ExamResultExamData represents exam information in the result
type ExamResultExamData struct {
	ID              uint    `json:"id"`
	Title           string  `json:"title"`
	Description     *string `json:"description"`
	TotalQuestions  int     `json:"total_questions"`
	DurationMinutes int     `json:"duration_minutes"`
	PassingScore    float64 `json:"passing_score"`
	ExamType        string  `json:"exam_type"`
}

// ExamResultAttemptData represents the attempt information in exam results
type ExamResultAttemptData struct {
	ID                uint    `json:"id"`
	Status            string  `json:"status"`
	StartedAt         string  `json:"started_at"`
	CompletedAt       string  `json:"completed_at"`
	Score             float64 `json:"score"`            // Total earned score (sum of earned points)
	MaxScore          float64 `json:"max_score"`        // Total possible score (sum of max points)
	ScorePercentage   float64 `json:"score_percentage"` // Percentage score (earned/max * 100)
	CorrectAnswers    int     `json:"correct_answers"`
	WrongAnswers      int     `json:"wrong_answers"`
	IsPassed          bool    `json:"is_passed"`
	TimeSpent         int     `json:"time_spent"` // Total time spent in seconds
	Rank              *int    `json:"rank,omitempty"`
	TotalParticipants *int    `json:"total_participants,omitempty"`
}

// ExamResultPackageData represents package information in the result
type ExamResultPackageData struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ExamResultOptionData represents detailed option information in exam results
type ExamResultOptionData struct {
	ID         int    `json:"id"`
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct"`
}

// ExamResultQuestion represents a question in the exam results with enhanced option data
type ExamResultQuestion struct {
	ID            uint                   `json:"id"`
	Question      string                 `json:"question"`
	QuestionType  string                 `json:"question_type"`
	Options       []ExamResultOptionData `json:"options"`
	CorrectAnswer interface{}            `json:"correct_answer"`
	UserAnswer    interface{}            `json:"user_answer"`
	IsCorrect     bool                   `json:"is_correct"`
	Points        int                    `json:"points"`               // Earned points (not max points)
	MaxPoints     *int                   `json:"max_points,omitempty"` // Optional max points
	TimeSpent     *int                   `json:"time_spent,omitempty"` // Optional time spent in seconds
	Explanation   *string                `json:"explanation,omitempty"`
}
