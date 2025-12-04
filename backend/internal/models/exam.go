package models

import (
	"time"

	"gorm.io/gorm"
)

// ExamStatus enum for different exam states
type ExamStatus string

const (
	ExamStatusDraft     ExamStatus = "DRAFT"
	ExamStatusScheduled ExamStatus = "SCHEDULED"
	ExamStatusActive    ExamStatus = "ACTIVE"
	ExamStatusCompleted ExamStatus = "COMPLETED"
)

// ExamType enum for different exam categories
type ExamType string

const (
	ExamTypeDaily  ExamType = "DAILY"
	ExamTypeMock   ExamType = "MOCK"
	ExamTypeReview ExamType = "REVIEW"
	ExamTypeFinal  ExamType = "FINAL"
)

// Exam represents the exams table
type Exam struct {
	ID          uint    `json:"id" gorm:"primarykey"`
	Title       string  `json:"title" gorm:"size:200;not null"`
	Slug        string  `json:"slug" gorm:"size:200;uniqueIndex;not null"`
	Description *string `json:"description" gorm:"type:text"`

	// Exam Type & Configuration
	ExamType        ExamType `json:"exam_type" gorm:"type:enum('DAILY','MOCK','REVIEW','FINAL');not null;default:'DAILY';index:idx_exam_type"`
	TotalQuestions  int      `json:"total_questions" gorm:"not null"`
	DurationMinutes int      `json:"duration_minutes" gorm:"not null;default:60"`
	TotalMarks      float64  `json:"total_marks" gorm:"type:decimal(8,2);not null;default:0.00;comment:'Total marks/points for this exam'"`
	PassingScore    float64  `json:"passing_score" gorm:"type:decimal(5,2);default:60.00"`
	MaxAttempts     int      `json:"max_attempts" gorm:"default:1"`

	// Embedded Questions (JSON) - No JOIN needed!
	QuestionsData string `json:"questions_data" gorm:"type:longtext;not null;comment:'JSON array of complete question objects with options, answers, explanations'"`

	// Scheduling
	ScheduledStartDate *time.Time `json:"scheduled_start_date"`
	ScheduledEndDate   *time.Time `json:"scheduled_end_date"`

	// Settings
	Instructions *string `json:"instructions" gorm:"type:text"`

	// Analytics & Statistics (denormalized for performance)
	AttemptCount          int        `json:"attempt_count" gorm:"default:0;index:idx_attempt_count;comment:'Total number of exam attempts by all users'"`
	CompletedAttemptCount int        `json:"completed_attempt_count" gorm:"default:0;comment:'Number of completed attempts (excludes abandoned)'"`
	AverageScore          *float64   `json:"average_score" gorm:"type:decimal(5,2);comment:'Average score of all completed attempts'"`
	PassRate              *float64   `json:"pass_rate" gorm:"type:decimal(5,2);comment:'Percentage of attempts that passed'"`
	LastAttemptAt         *time.Time `json:"last_attempt_at" gorm:"comment:'When the last attempt was made'"`

	// Status
	Status    ExamStatus `json:"status" gorm:"type:enum('DRAFT','SCHEDULED','ACTIVE','COMPLETED');default:'DRAFT';index:idx_status"`
	IsActive  bool       `json:"is_active" gorm:"default:true;index:idx_active"`
	CreatedBy *uint      `json:"created_by" gorm:"index:idx_created_by"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships (will add Package relationship via junction table)
}

// TableName specifies the table name for Exam
func (Exam) TableName() string {
	return "exams"
}

// ExamQuestion represents the structure of questions stored in QuestionsData JSON
type ExamQuestion struct {
	ID           uint                   `json:"id"`
	QuestionText string                 `json:"question_text"`
	QuestionType QuestionType           `json:"question_type"`
	Options      map[string]interface{} `json:"options"`
	Points       int                    `json:"points"` // Points for this question
}

// UpdateAttemptCount increments the attempt counter
func (e *Exam) UpdateAttemptCount(db *gorm.DB) error {
	now := time.Now()
	return db.Model(e).Updates(map[string]interface{}{
		"attempt_count":   gorm.Expr("attempt_count + 1"),
		"last_attempt_at": now,
	}).Error
}

// UpdateCompletedAttemptStats updates statistics when an attempt is completed
func (e *Exam) UpdateCompletedAttemptStats(db *gorm.DB, score float64, passed bool) error {
	return db.Model(e).Update("completed_attempt_count", gorm.Expr("completed_attempt_count + 1")).Error
}

// RecalculateExamStats recalculates exam statistics from actual attempt data
func (e *Exam) RecalculateExamStats(db *gorm.DB) error {
	var totalAttempts int64
	var completedAttempts int64
	var averageScore float64
	var passCount int64
	var lastAttempt time.Time

	// Count total attempts
	if err := db.Model(&UserExamAttempt{}).Where("exam_id = ?", e.ID).Count(&totalAttempts).Error; err != nil {
		return err
	}

	// Count completed attempts and calculate average score
	var result struct {
		Count    int64
		AvgScore *float64
	}

	if err := db.Model(&UserExamAttempt{}).
		Where("exam_id = ? AND status IN ? AND is_scored = true", e.ID, []AttemptStatus{AttemptStatusCompleted, AttemptStatusAutoSubmitted}).
		Select("COUNT(*) as count, AVG(score) as avg_score").
		Scan(&result).Error; err != nil {
		return err
	}

	completedAttempts = result.Count
	if result.AvgScore != nil {
		averageScore = *result.AvgScore
	}

	// Count passed attempts
	if err := db.Model(&UserExamAttempt{}).
		Where("exam_id = ? AND is_passed = true", e.ID).
		Count(&passCount).Error; err != nil {
		return err
	}

	// Get last attempt date
	var attempt UserExamAttempt
	if err := db.Where("exam_id = ?", e.ID).
		Order("started_at DESC").
		First(&attempt).Error; err == nil {
		lastAttempt = attempt.StartedAt
	}

	// Calculate pass rate
	var passRate float64
	if completedAttempts > 0 {
		passRate = (float64(passCount) / float64(completedAttempts)) * 100
	}

	// Update the exam
	updates := map[string]interface{}{
		"attempt_count":           totalAttempts,
		"completed_attempt_count": completedAttempts,
	}

	if completedAttempts > 0 {
		updates["average_score"] = averageScore
		updates["pass_rate"] = passRate
	}

	if !lastAttempt.IsZero() {
		updates["last_attempt_at"] = lastAttempt
	}

	return db.Model(e).Updates(updates).Error
}
