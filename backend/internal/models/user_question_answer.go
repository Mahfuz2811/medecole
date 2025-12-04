package models

import (
	"time"

	"gorm.io/gorm"
)

// UserQuestionAnswer represents individual question answers for analytics (populated in background)
type UserQuestionAnswer struct {
	ID        uint `json:"id" gorm:"primarykey"`
	AttemptID uint `json:"attempt_id" gorm:"not null;index:idx_attempt_id;comment:'Reference to user_exam_attempts'"`
	UserID    uint `json:"user_id" gorm:"not null;index:idx_user_id;comment:'For direct user analytics'"`
	ExamID    uint `json:"exam_id" gorm:"not null;index:idx_exam_id;comment:'For exam-level analytics'"`

	// Question details (snapshot for analytics, populated from exam questions)
	QuestionID      uint            `json:"question_id" gorm:"not null;index:idx_question_id;comment:'Original question ID'"`
	QuestionType    QuestionType    `json:"question_type" gorm:"not null;index:idx_question_type;comment:'SBA, TRUE_FALSE, etc.'"`
	QuestionText    string          `json:"question_text" gorm:"type:text;comment:'Snapshot for analytics'"`
	DifficultyLevel DifficultyLevel `json:"difficulty_level" gorm:"index:idx_difficulty;comment:'EASY, MEDIUM, HARD'"`
	QuestionIndex   int             `json:"question_index" gorm:"not null;comment:'Position in exam (0-based)'"`

	// Answer details (extracted from AnswersData JSON during background processing)
	SelectedOptions string `json:"selected_options" gorm:"type:text;comment:'JSON array of selected options'"`
	CorrectOptions  string `json:"correct_options" gorm:"type:text;comment:'JSON array of correct options for comparison'"`

	// Scoring (different mechanisms for different question types)
	IsCorrect    bool    `json:"is_correct" gorm:"index:idx_correct;comment:'Whether answer is completely correct'"`
	PartialScore float64 `json:"partial_score" gorm:"type:decimal(5,2);default:0.00;comment:'Partial credit score (0.00-1.00)'"`
	MaxScore     float64 `json:"max_score" gorm:"type:decimal(5,2);default:1.00;comment:'Maximum possible score for this question'"`

	// Timing analytics (extracted from AnswersData JSON)
	TimeSpent  int       `json:"time_spent" gorm:"default:0;comment:'Seconds spent on this question'"`
	AnsweredAt time.Time `json:"answered_at" gorm:"comment:'When this question was answered'"`

	// Behavioral analytics
	IsSkipped    bool `json:"is_skipped" gorm:"default:false;index:idx_skipped;comment:'Whether question was skipped'"`
	ChangeCount  int  `json:"change_count" gorm:"default:0;comment:'How many times answer was changed'"`
	IsLastAnswer bool `json:"is_last_answer" gorm:"default:true;comment:'Whether this was the final answer or changed later'"`

	// Metadata
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Attempt UserExamAttempt `json:"attempt,omitempty" gorm:"foreignKey:AttemptID"`
	User    User            `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Exam    Exam            `json:"exam,omitempty" gorm:"foreignKey:ExamID"`
}

// TableName specifies the table name for UserQuestionAnswer
func (UserQuestionAnswer) TableName() string {
	return "user_question_answers"
}

// GetSelectedOptionsArray returns the selected options as a string array
func (u *UserQuestionAnswer) GetSelectedOptionsArray() []string {
	var options []string
	if u.SelectedOptions != "" {
		// Parse JSON string to array
		// Implementation would parse the JSON string
	}
	return options
}

// GetCorrectOptionsArray returns the correct options as a string array
func (u *UserQuestionAnswer) GetCorrectOptionsArray() []string {
	var options []string
	if u.CorrectOptions != "" {
		// Parse JSON string to array
		// Implementation would parse the JSON string
	}
	return options
}

// CalculateAccuracy returns the accuracy as a percentage (0-100)
func (u *UserQuestionAnswer) CalculateAccuracy() float64 {
	if u.MaxScore == 0 {
		return 0.0
	}
	return (u.PartialScore / u.MaxScore) * 100.0
}

// IsPartiallyCorrect checks if the answer received partial credit
func (u *UserQuestionAnswer) IsPartiallyCorrect() bool {
	return u.PartialScore > 0 && u.PartialScore < u.MaxScore
}

// GetEfficiencyScore calculates time efficiency (points per second)
func (u *UserQuestionAnswer) GetEfficiencyScore() float64 {
	if u.TimeSpent == 0 {
		return 0.0
	}
	return u.PartialScore / float64(u.TimeSpent)
}
