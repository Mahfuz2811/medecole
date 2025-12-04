package dto

import (
	"time"
)

// EnrollmentRequest represents the enrollment request payload
type EnrollmentRequest struct {
	PackageID  uint    `json:"package_id" binding:"required"`
	CouponCode *string `json:"coupon_code,omitempty"`
}

// EnrollmentResponse represents the enrollment response
type EnrollmentResponse struct {
	ID                  uint                    `json:"id"`
	UserID              uint                    `json:"user_id"`
	PackageID           uint                    `json:"package_id"`
	EnrollmentType      string                  `json:"enrollment_type"`
	EnrolledAt          time.Time               `json:"enrolled_at"`
	ExpiresAt           *time.Time              `json:"expires_at"`
	EnrolledPackageType string                  `json:"enrolled_package_type"`
	EnrolledPrice       float64                 `json:"enrolled_price"`
	PaymentStatus       string                  `json:"payment_status"`
	PaymentAmount       *float64                `json:"payment_amount"`
	CouponCode          *string                 `json:"coupon_code"`
	OriginalPrice       *float64                `json:"original_price"`
	DiscountPercentage  *float64                `json:"discount_percentage"`
	DiscountAmount      *float64                `json:"discount_amount"`
	FinalPrice          *float64                `json:"final_price"`
	EffectiveStatus     string                  `json:"effective_status"`
	CanAccessContent    bool                    `json:"can_access_content"`
	IsActive            bool                    `json:"is_active"`
	Package             *PackageBasicResponse   `json:"package,omitempty"`
	PriceCalculation    *PriceCalculationResult `json:"price_calculation,omitempty"`
	CreatedAt           time.Time               `json:"created_at"`
	UpdatedAt           time.Time               `json:"updated_at"`
}

// PackageBasicResponse represents basic package info in enrollment response
type PackageBasicResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	PackageType string  `json:"package_type"`
	Price       float64 `json:"price"`
}

// UserEnrollmentsResponse represents user's enrollments list
type UserEnrollmentsResponse struct {
	Enrollments []EnrollmentResponse `json:"enrollments"`
	Total       int                  `json:"total"`
	Active      int                  `json:"active"`
	Expired     int                  `json:"expired"`
}

// CouponValidationRequest represents coupon validation request
type CouponValidationRequest struct {
	CouponCode string `json:"coupon_code" binding:"required"`
	PackageID  uint   `json:"package_id" binding:"required"`
}

// CouponValidationResponse represents coupon validation result
type CouponValidationResponse struct {
	Valid              bool                    `json:"valid"`
	CouponCode         string                  `json:"coupon_code"`
	DiscountPercentage float64                 `json:"discount_percentage"`
	Message            string                  `json:"message"`
	PriceCalculation   *PriceCalculationResult `json:"price_calculation,omitempty"`
}

// PriceCalculationResult represents price calculation details
type PriceCalculationResult struct {
	OriginalPrice      float64 `json:"original_price"`
	DiscountPercentage float64 `json:"discount_percentage"`
	DiscountAmount     float64 `json:"discount_amount"`
	FinalPrice         float64 `json:"final_price"`
	CouponCode         *string `json:"coupon_code"`
}

// EnrollmentStatusResponse represents enrollment status check
type EnrollmentStatusResponse struct {
	HasActiveEnrollment bool                `json:"has_active_enrollment"`
	Enrollment          *EnrollmentResponse `json:"enrollment,omitempty"`
}
