package types

import "time"

// AttemptFilters for filtering exam attempts
type AttemptFilters struct {
	Period    string
	PackageID *uint
	ExamType  string
	StartDate *time.Time
	EndDate   *time.Time
}

// UserStatsData aggregated user statistics
type UserStatsData struct {
	TotalAttempts  int
	CorrectAnswers int
	TotalQuestions int
	TotalTimeSpent int
	AverageScore   float64
}
