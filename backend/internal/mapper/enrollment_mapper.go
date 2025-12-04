package mapper

import (
	"quizora-backend/internal/dto"
	"quizora-backend/internal/models"
)

// EnrollmentMapper handles mapping between enrollment models and DTOs
type EnrollmentMapper struct{}

// NewEnrollmentMapper creates a new enrollment mapper
func NewEnrollmentMapper() *EnrollmentMapper {
	return &EnrollmentMapper{}
}

// ToEnrollmentResponse converts UserPackageEnrollment model to EnrollmentResponse DTO
func (m *EnrollmentMapper) ToEnrollmentResponse(enrollment *models.UserPackageEnrollment) *dto.EnrollmentResponse {
	if enrollment == nil {
		return nil
	}

	response := &dto.EnrollmentResponse{
		ID:                  enrollment.ID,
		UserID:              enrollment.UserID,
		PackageID:           enrollment.PackageID,
		EnrollmentType:      string(enrollment.EnrollmentType),
		PaymentStatus:       string(enrollment.PaymentStatus),
		EnrolledAt:          enrollment.EnrolledAt,
		ExpiresAt:           enrollment.ExpiresAt,
		EnrolledPackageType: string(enrollment.EnrolledPackageType),
		EnrolledPrice:       enrollment.EnrolledPrice,
		PaymentAmount:       enrollment.PaymentAmount,
		IsActive:            enrollment.IsActive,
		CanAccessContent:    enrollment.CanAccessContent(),
		EffectiveStatus:     enrollment.GetEffectiveStatus(),
		CreatedAt:           enrollment.CreatedAt,
		UpdatedAt:           enrollment.UpdatedAt,
	}

	// Add coupon information if used
	if enrollment.CouponID != nil {
		response.CouponCode = enrollment.CouponCode
		response.OriginalPrice = enrollment.OriginalPrice
		response.DiscountPercentage = enrollment.DiscountPercentage
		response.DiscountAmount = enrollment.DiscountAmount
		response.FinalPrice = enrollment.FinalPrice
	}

	// Add package information if loaded
	if enrollment.Package.ID != 0 {
		response.Package = m.toPackageBasicResponse(&enrollment.Package)
	}

	return response
}

// toPackageBasicResponse converts Package model to PackageBasicResponse DTO
func (m *EnrollmentMapper) toPackageBasicResponse(pkg *models.Package) *dto.PackageBasicResponse {
	if pkg == nil {
		return nil
	}

	return &dto.PackageBasicResponse{
		ID:          pkg.ID,
		Name:        pkg.Name,
		Slug:        pkg.Slug,
		PackageType: string(pkg.PackageType),
		Price:       pkg.Price,
	}
}

// ToEnrollmentListResponse converts slice of enrollments to list response
func (m *EnrollmentMapper) ToEnrollmentListResponse(enrollments []models.UserPackageEnrollment) []dto.EnrollmentResponse {
	responses := make([]dto.EnrollmentResponse, len(enrollments))
	for i, enrollment := range enrollments {
		responses[i] = *m.ToEnrollmentResponse(&enrollment)
	}
	return responses
}
