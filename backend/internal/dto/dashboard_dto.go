package dto

// DashboardSummaryResponse represents the complete dashboard data
type DashboardSummaryResponse struct {
	UserStats      UserStatsDTO        `json:"user_stats"`
	RecentActivity []RecentActivityDTO `json:"recent_activity"`
}

// UserStatsDTO represents user performance statistics
type UserStatsDTO struct {
	TotalAttempts  int     `json:"total_attempts"`
	CorrectAnswers int     `json:"is_corrects"`
	AccuracyRate   float64 `json:"accuracy_rate"`
}

// RecentActivityDTO represents a recent exam activity
type RecentActivityDTO struct {
	ID             uint    `json:"id"`
	ExamTitle      string  `json:"exam_title"`
	PackageName    string  `json:"package_name"`
	Date           string  `json:"date"`  // Relative time like "2 hours ago", "Yesterday"
	Score          float64 `json:"score"` // Percentage score
	TotalQuestions int     `json:"total_questions"`
	CorrectAnswers int     `json:"is_corrects"`
	TimeTaken      string  `json:"time_taken"` // Formatted like "12 min"
	Status         string  `json:"status"`     // "completed"
}

// DashboardEnrollmentDTO represents optimized enrollment data for dashboard
type DashboardEnrollmentDTO struct {
	ID             uint    `json:"id"`
	PackageID      uint    `json:"package_id"`
	PackageName    string  `json:"package_name"`
	PackageSlug    string  `json:"package_slug"`
	PackageType    string  `json:"package_type"`
	ExpiryDate     *string `json:"expiry_date"` // ISO string or null
	Status         string  `json:"status"`      // "active" or "enrolled"
	Progress       float64 `json:"progress"`    // Percentage 0-100
	TotalExams     int     `json:"total_exams"`
	CompletedExams int     `json:"completed_exams"`
}

// DashboardEnrollmentsResponse represents dashboard enrollments response
type DashboardEnrollmentsResponse struct {
	Enrollments []DashboardEnrollmentDTO `json:"enrollments"`
	Total       int                      `json:"total"`
	Active      int                      `json:"active"`
}
