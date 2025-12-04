package mapper

import (
	"quizora-backend/internal/dto"
	"quizora-backend/internal/models"
)

// PackageMapper handles conversion between models and DTOs
type PackageMapper struct {
	examScheduleMapper *ExamScheduleMapper
}

// NewPackageMapper creates a new package mapper
func NewPackageMapper() *PackageMapper {
	return &PackageMapper{
		examScheduleMapper: NewExamScheduleMapper(),
	}
}

// ToPackageResponse converts a Package model to PackageResponse DTO
func (m *PackageMapper) ToPackageResponse(pkg models.Package) dto.PackageResponse {
	return dto.PackageResponse{
		ID:           pkg.ID,
		Name:         pkg.Name,
		Slug:         pkg.Slug,
		Description:  pkg.Description,
		PackageType:  pkg.PackageType,
		Price:        pkg.Price,
		Images:       m.buildImageResponse(pkg),
		CouponCode:   pkg.CouponCode,
		ValidityType: pkg.ValidityType,
		ValidityDays: pkg.ValidityDays,
		ValidityDate: m.formatValidityDate(pkg),
		TotalExams:   pkg.TotalExams,
		IsActive:     pkg.IsActive,
		SortOrder:    pkg.SortOrder,
		CreatedAt:    pkg.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    pkg.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		// Analytics fields
		EnrollmentCount: pkg.EnrollmentCount,
		// Initialize empty exams array - will be filled by ToPackageResponseWithExams if needed
		Exams: []dto.PackageExamScheduleResponse{},
	}
}

// ToPackageListResponse converts packages to list response
func (m *PackageMapper) ToPackageListResponse(packages []models.Package) dto.PackageListResponse {
	packageResponses := make([]dto.PackageResponse, len(packages))
	for i, pkg := range packages {
		packageResponses[i] = m.ToPackageResponse(pkg)
	}

	return dto.PackageListResponse{
		Packages: packageResponses,
	}
}

// ToPackageResponseWithExams converts a Package model with exams to PackageResponse DTO
func (m *PackageMapper) ToPackageResponseWithExams(pkg models.Package) dto.PackageResponse {
	response := m.ToPackageResponse(pkg)

	// Convert package exams
	exams := make([]dto.PackageExamScheduleResponse, len(pkg.PackageExams))
	for i, packageExam := range pkg.PackageExams {
		exams[i] = m.examScheduleMapper.ToPackageExamScheduleResponse(packageExam)
	}

	response.Exams = exams
	return response
}

// buildImageResponse creates the image response for Next.js optimization
func (m *PackageMapper) buildImageResponse(pkg models.Package) dto.PackageImageResponse {
	metadata, _ := pkg.GetImageMetadata()

	// With Next.js Image component, we only need the original URL
	// Next.js automatically generates responsive sizes and optimizations
	originalURL := pkg.GetNextJSImageSrc()

	return dto.PackageImageResponse{
		Original:  originalURL,
		Mobile:    originalURL, // Next.js handles mobile optimization
		Tablet:    originalURL, // Next.js handles tablet optimization
		Desktop:   originalURL, // Next.js handles desktop optimization
		Thumbnail: originalURL, // Next.js handles thumbnail generation
		AltText:   pkg.GetImageAltText(),
		Metadata:  metadata,
	}
}

// formatValidityDate converts validity date to ISO string if exists
func (m *PackageMapper) formatValidityDate(pkg models.Package) *string {
	if pkg.ValidityDate != nil {
		dateStr := pkg.ValidityDate.Format("2006-01-02T15:04:05Z")
		return &dateStr
	}
	return nil
}
