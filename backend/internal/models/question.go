package models

import (
	"time"

	"gorm.io/gorm"
)

// QuestionType enum for different question types
type QuestionType string

const (
	QuestionTypeSBA       QuestionType = "SBA"
	QuestionTypeTrueFalse QuestionType = "TRUE_FALSE"
)

// DifficultyLevel enum for question difficulty
type DifficultyLevel string

const (
	DifficultyEasy   DifficultyLevel = "EASY"
	DifficultyMedium DifficultyLevel = "MEDIUM"
	DifficultyHard   DifficultyLevel = "HARD"
)

// Question represents the questions table
type Question struct {
	ID              uint            `json:"id" gorm:"primarykey"`
	SystemID        uint            `json:"system_id" gorm:"not null;index:idx_system_type;index:idx_system_active"`
	QuestionText    string          `json:"question_text" gorm:"type:text;not null"`
	QuestionType    QuestionType    `json:"question_type" gorm:"type:enum('SBA','TRUE_FALSE');not null;index:idx_system_type;index:idx_type_active"`
	DifficultyLevel DifficultyLevel `json:"difficulty_level" gorm:"type:enum('EASY','MEDIUM','HARD');default:'MEDIUM';index:idx_difficulty_active"`

	// Embedded options as JSON (eliminates options table JOIN)
	Options     string  `json:"options" gorm:"type:json;not null"`
	Explanation *string `json:"explanation" gorm:"type:text"`
	Reference   *string `json:"reference" gorm:"size:255"`
	Tags        string  `json:"tags" gorm:"type:json"`

	// Usage tracking for exam generation
	UsageCount int `json:"usage_count" gorm:"default:0;index:idx_usage_count;comment:'How many times this question has been used in exams'"`

	IsActive  bool  `json:"is_active" gorm:"default:true;index:idx_system_active;index:idx_difficulty_active;index:idx_type_active"`
	CreatedBy *uint `json:"created_by" gorm:"index:idx_created_by"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	System System `json:"system,omitempty" gorm:"foreignKey:SystemID"`
}

// TableName specifies the table name for Question
func (Question) TableName() string {
	return "questions"
}

// IncrementUsage increments the usage count for this question
func (q *Question) IncrementUsage(db *gorm.DB) error {
	return db.Model(q).UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1)).Error
}

// GetUsageLevel returns a string indicating how frequently this question is used
func (q *Question) GetUsageLevel() string {
	switch {
	case q.UsageCount == 0:
		return "UNUSED"
	case q.UsageCount <= 5:
		return "LOW_USAGE"
	case q.UsageCount <= 15:
		return "MEDIUM_USAGE"
	case q.UsageCount <= 30:
		return "HIGH_USAGE"
	default:
		return "OVERUSED"
	}
}

// IsOverused checks if this question has been used too many times
func (q *Question) IsOverused(threshold int) bool {
	if threshold <= 0 {
		threshold = 50 // Default threshold
	}
	return q.UsageCount >= threshold
}
