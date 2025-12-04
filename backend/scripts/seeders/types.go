package seeders

import "quizora-backend/internal/models"

// Sample data structures
type SampleSubject struct {
	Name        string
	Slug        string
	Description string
	Systems     []SampleSystem
}

type SampleSystem struct {
	Name        string
	Slug        string
	Description string
	Questions   []SampleQuestion
}

type SampleQuestion struct {
	QuestionType    models.QuestionType
	QuestionText    string
	Options         map[string]interface{} // JSON structure for options
	Explanation     string
	DifficultyLevel models.DifficultyLevel
	Reference       string
	Tags            []string
}

// Helper functions
func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func BoolPtr(b bool) *bool {
	return &b
}
