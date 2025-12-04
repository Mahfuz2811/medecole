package models

import (
	"time"

	"gorm.io/gorm"
)

// AttemptStatus enum for exam attempt states
type AttemptStatus string

const (
	AttemptStatusStarted       AttemptStatus = "STARTED"        // Initial DB record created, exam in progress in Redis
	AttemptStatusCompleted     AttemptStatus = "COMPLETED"      // User submitted voluntarily
	AttemptStatusAutoSubmitted AttemptStatus = "AUTO_SUBMITTED" // Time expired, auto-submitted
	AttemptStatusAbandoned     AttemptStatus = "ABANDONED"      // Never completed (cleanup job)
)

// UserExamAttempt represents user attempts at exams (simplified for Redis workflow)
type UserExamAttempt struct {
	ID        uint `json:"id" gorm:"primarykey"`
	UserID    uint `json:"user_id" gorm:"not null;uniqueIndex:idx_user_exam_package"`
	ExamID    uint `json:"exam_id" gorm:"not null;uniqueIndex:idx_user_exam_package"`
	PackageID uint `json:"package_id" gorm:"not null;uniqueIndex:idx_user_exam_package;index:idx_package_id;comment:'Package context for this attempt - same exam in different packages are separate attempts'"`

	// Single attempt control (ready for future scaling)
	AttemptNumber int `json:"attempt_number" gorm:"default:1;comment:'Always 1 for now, ready for multiple attempts'"`

	// Status and timing (final state only - active state managed in Redis)
	Status      AttemptStatus `json:"status" gorm:"type:enum('STARTED','COMPLETED','AUTO_SUBMITTED','ABANDONED');default:'STARTED';index:idx_status"`
	StartedAt   time.Time     `json:"started_at" gorm:"not null;comment:'When exam was started'"`
	CompletedAt *time.Time    `json:"completed_at" gorm:"comment:'When exam was completed or auto-submitted'"`

	// Session recovery support (minimal data for reconnection)
	SessionID      *string    `json:"session_id" gorm:"type:varchar(64);index:idx_session;comment:'Redis session key for active attempts'"`
	LastActivityAt *time.Time `json:"last_activity_at" gorm:"comment:'Last activity timestamp for session cleanup'"`

	// Time tracking (calculated at completion)
	TimeLimitSeconds int `json:"time_limit_seconds" gorm:"not null;comment:'Snapshot from exam.duration_minutes * 60'"`
	ActualTimeSpent  int `json:"actual_time_spent" gorm:"default:0;comment:'Calculated: completed_at - started_at OR time_limit if auto-submitted'"`

	// Final answers (stored only on completion from Redis)
	AnswersData string `json:"answers_data" gorm:"type:longtext;comment:'Final JSON array of all answers from Redis'"`

	// Exam snapshot (for consistent scoring even if exam changes)
	TotalQuestions int     `json:"total_questions" gorm:"not null;comment:'Snapshot from exam at start time'"`
	PassingScore   float64 `json:"passing_score" gorm:"type:decimal(5,2);not null;comment:'Snapshot from exam at start time'"`

	// Scoring (processed in background after completion)
	IsScored       bool     `json:"is_scored" gorm:"default:false;index:idx_scored;comment:'Whether background scoring is completed'"`
	Score          *float64 `json:"score" gorm:"type:decimal(5,2);comment:'Final score percentage (0-100)'"`
	CorrectAnswers *int     `json:"is_corrects" gorm:"comment:'Number of correct answers'"`
	IsPassed       *bool    `json:"is_passed" gorm:"comment:'Whether attempt passed based on passing_score'"`

	// Metadata
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Exam    Exam    `json:"exam,omitempty" gorm:"foreignKey:ExamID"`
	Package Package `json:"package,omitempty" gorm:"foreignKey:PackageID"`
}

// TableName specifies the table name for UserExamAttempt
func (UserExamAttempt) TableName() string {
	return "user_exam_attempts"
}

// IsCompleted checks if the attempt is in a completed state
func (u *UserExamAttempt) IsCompleted() bool {
	return u.Status == AttemptStatusCompleted || u.Status == AttemptStatusAutoSubmitted
}

// IsInProgress checks if the attempt is currently in progress (should exist in Redis)
func (u *UserExamAttempt) IsInProgress() bool {
	return u.Status == AttemptStatusStarted
}

// GetTimeSpentSeconds returns actual time spent or calculates it if still in progress
func (u *UserExamAttempt) GetTimeSpentSeconds() int {
	if u.IsCompleted() {
		return u.ActualTimeSpent
	}
	// If still in progress, calculate from start time
	return int(time.Since(u.StartedAt).Seconds())
}

// GetRemainingTimeSeconds calculates remaining time for in-progress attempts
func (u *UserExamAttempt) GetRemainingTimeSeconds() int {
	if u.IsCompleted() {
		return 0
	}
	elapsed := int(time.Since(u.StartedAt).Seconds())
	remaining := u.TimeLimitSeconds - elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IsTimeExpired checks if the time limit has been exceeded
func (u *UserExamAttempt) IsTimeExpired() bool {
	return u.GetRemainingTimeSeconds() <= 0
}

// HasActiveSession checks if there's an active Redis session
func (u *UserExamAttempt) HasActiveSession() bool {
	return u.SessionID != nil && u.Status == AttemptStatusStarted
}

// CanResumeSession checks if the attempt can be resumed (not expired and has session)
func (u *UserExamAttempt) CanResumeSession() bool {
	return u.HasActiveSession() && !u.IsTimeExpired()
}

// GetSessionKey returns the Redis session key
func (u *UserExamAttempt) GetSessionKey() string {
	if u.SessionID != nil {
		return *u.SessionID
	}
	return ""
}

// IsSessionStale checks if session hasn't been updated recently (for cleanup)
func (u *UserExamAttempt) IsSessionStale(timeoutMinutes int) bool {
	if u.LastActivityAt == nil {
		return false
	}
	return time.Since(*u.LastActivityAt) > time.Duration(timeoutMinutes)*time.Minute
}

// GetScorePercentage returns the score as a percentage (0-100)
func (u *UserExamAttempt) GetScorePercentage() float64 {
	if u.Score != nil {
		return *u.Score
	}
	return 0.0
}
