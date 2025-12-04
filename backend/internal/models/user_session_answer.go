package models

import (
	"time"

	"gorm.io/gorm"
)

// UserSessionAnswer represents a real-time answer stored during an active exam session
// This is used for auto-save functionality and real-time sync
type UserSessionAnswer struct {
	ID        uint `json:"id" gorm:"primarykey"`
	AttemptID uint `json:"attempt_id" gorm:"not null;index:idx_attempt_id;comment:'Reference to user_exam_attempts'"`
	UserID    uint `json:"user_id" gorm:"not null;index:idx_user_id;comment:'For quick user filtering'"`

	// Question and Answer (minimal data for sync)
	QuestionID     uint   `json:"question_id" gorm:"not null;index:idx_question_id;comment:'Question being answered'"`
	SelectedOption string `json:"selected_option" gorm:"not null;comment:'Selected option key (a, b, c, d, e)'"`

	// Timing
	AnsweredAt time.Time `json:"answered_at" gorm:"comment:'When this answer was last updated'"`

	// Metadata
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Attempt UserExamAttempt `json:"attempt,omitempty" gorm:"foreignKey:AttemptID"`
	User    User            `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for UserSessionAnswer
func (UserSessionAnswer) TableName() string {
	return "user_session_answers"
}

// IsValid checks if the selected option is valid (a-e)
func (u *UserSessionAnswer) IsValid() bool {
	validOptions := map[string]bool{
		"a": true, "b": true, "c": true, "d": true, "e": true,
	}
	return validOptions[u.SelectedOption]
}
