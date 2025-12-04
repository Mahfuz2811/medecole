package dto

import "quizora-backend/internal/models"

// PackageListRequest represents request parameters for listing packages
type PackageListRequest struct {
	// Removed pagination and filter parameters
	// Sort order will be hardcoded to sort_order in service layer
}

// PackageImageResponse represents responsive image URLs
type PackageImageResponse struct {
	Original  string                `json:"original"`
	Mobile    string                `json:"mobile"`
	Tablet    string                `json:"tablet"`
	Desktop   string                `json:"desktop"`
	Thumbnail string                `json:"thumbnail"`
	AltText   string                `json:"alt_text"`
	Metadata  *models.ImageMetadata `json:"metadata,omitempty"`
}

// PackageResponse represents the API response structure for packages
type PackageResponse struct {
	ID           uint                 `json:"id"`
	Name         string               `json:"name"`
	Slug         string               `json:"slug"`
	Description  *string              `json:"description"`
	PackageType  models.PackageType   `json:"package_type"`
	Price        float64              `json:"price"`
	Images       PackageImageResponse `json:"images"`
	CouponCode   *string              `json:"coupon_code,omitempty"`
	ValidityType models.ValidityType  `json:"validity_type"`
	ValidityDays *int                 `json:"validity_days,omitempty"`
	ValidityDate *string              `json:"validity_date,omitempty"`
	TotalExams   int                  `json:"total_exams"`
	IsActive     bool                 `json:"is_active"`
	SortOrder    int                  `json:"sort_order"`
	CreatedAt    string               `json:"created_at"`
	UpdatedAt    string               `json:"updated_at"`
	// Analytics fields
	EnrollmentCount       int `json:"enrollment_count"`
	ActiveEnrollmentCount int `json:"active_enrollment_count"`
	// Exam schedule data (only included in single package requests)
	Exams []PackageExamScheduleResponse `json:"exams"`
}

// PackageListResponse represents the response for package listing
type PackageListResponse struct {
	Packages []PackageResponse `json:"packages"`
}
