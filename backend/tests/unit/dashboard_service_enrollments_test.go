package unit

import (
	"errors"
	"github.com/Mahfuz2811/medecole/backend/internal/models"
	"github.com/Mahfuz2811/medecole/backend/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test data generators for enrollment tests
func createTestEnrollments() []models.UserPackageEnrollment {
	now := time.Now()
	expiryFuture := now.Add(30 * 24 * time.Hour) // 30 days from now
	expiryPast := now.Add(-5 * 24 * time.Hour)   // 5 days ago

	return []models.UserPackageEnrollment{
		{
			ID:                  1,
			UserID:              1,
			PackageID:           1,
			EnrollmentType:      models.EnrollmentTypeFull,
			EnrolledAt:          now.Add(-10 * 24 * time.Hour),
			ExpiresAt:           &expiryFuture,
			EnrolledPackageType: models.PackageTypePremium,
			EnrolledPrice:       99.99,
			PaymentStatus:       models.PaymentStatusPaid,
			IsActive:            true,
			Package: models.Package{
				ID:          1,
				Name:        "JavaScript Fundamentals",
				Slug:        "javascript-fundamentals",
				PackageType: models.PackageTypePremium,
			},
		},
		{
			ID:                  2,
			UserID:              1,
			PackageID:           2,
			EnrollmentType:      models.EnrollmentTypeFull,
			EnrolledAt:          now.Add(-20 * 24 * time.Hour),
			ExpiresAt:           &expiryFuture,
			EnrolledPackageType: models.PackageTypePremium,
			EnrolledPrice:       149.99,
			PaymentStatus:       models.PaymentStatusPaid,
			IsActive:            true,
			Package: models.Package{
				ID:          2,
				Name:        "React Advanced",
				Slug:        "react-advanced",
				PackageType: models.PackageTypePremium,
			},
		},
		{
			ID:                  3,
			UserID:              1,
			PackageID:           3,
			EnrollmentType:      models.EnrollmentTypeFull,
			EnrolledAt:          now.Add(-40 * 24 * time.Hour),
			ExpiresAt:           &expiryPast, // Expired enrollment
			EnrolledPackageType: models.PackageTypePremium,
			EnrolledPrice:       89.99,
			PaymentStatus:       models.PaymentStatusExpired,
			IsActive:            false,
			Package: models.Package{
				ID:          3,
				Name:        "Python Basics",
				Slug:        "python-basics",
				PackageType: models.PackageTypePremium,
			},
		},
	}
}

// Test GetDashboardEnrollments - Success
func TestDashboardService_GetDashboardEnrollments_Success(t *testing.T) {
	// Setup
	mockEnrollmentRepo := &MockDashboardEnrollmentRepository{}
	mockAttemptRepo := &MockUserExamAttemptRepository{}

	dashboardService := service.NewDashboardService(nil, mockEnrollmentRepo, mockAttemptRepo)

	userID := uint(1)
	testEnrollments := createTestEnrollments()

	// Only return active, non-expired enrollments
	activeEnrollments := testEnrollments[:2] // First two are active

	// Setup expectations
	mockEnrollmentRepo.On("GetUserEnrollments", userID).Return(activeEnrollments, nil)

	// Execute
	result, err := dashboardService.GetDashboardEnrollments(userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Enrollments, 2)
	assert.Equal(t, 2, result.Total)
	assert.Equal(t, 2, result.Active)

	// Check first enrollment
	enrollment1 := result.Enrollments[0]
	assert.Equal(t, uint(1), enrollment1.ID)
	assert.Equal(t, uint(1), enrollment1.PackageID)
	assert.Equal(t, "JavaScript Fundamentals", enrollment1.PackageName)
	assert.Equal(t, "javascript-fundamentals", enrollment1.PackageSlug)
	assert.Equal(t, "PREMIUM", enrollment1.PackageType)
	assert.Equal(t, "active", enrollment1.Status)
	// Progress will be 50% since we return default test values when db is nil
	assert.Equal(t, 50.0, enrollment1.Progress)
	assert.Equal(t, 10, enrollment1.TotalExams)
	assert.Equal(t, 5, enrollment1.CompletedExams)
	assert.NotNil(t, enrollment1.ExpiryDate)

	// Check second enrollment
	enrollment2 := result.Enrollments[1]
	assert.Equal(t, uint(2), enrollment2.ID)
	assert.Equal(t, "React Advanced", enrollment2.PackageName)
	assert.Equal(t, 50.0, enrollment2.Progress) // 50% from default test values
	assert.Equal(t, 10, enrollment2.TotalExams)
	assert.Equal(t, 5, enrollment2.CompletedExams)

	mockEnrollmentRepo.AssertExpectations(t)
}

// Test GetDashboardEnrollments - Repository Error
func TestDashboardService_GetDashboardEnrollments_RepositoryError(t *testing.T) {
	// Setup
	mockEnrollmentRepo := &MockDashboardEnrollmentRepository{}
	mockAttemptRepo := &MockUserExamAttemptRepository{}
	dashboardService := service.NewDashboardService(nil, mockEnrollmentRepo, mockAttemptRepo)

	userID := uint(1)

	// Setup expectations
	mockEnrollmentRepo.On("GetUserEnrollments", userID).Return([]models.UserPackageEnrollment{}, errors.New("database connection failed"))

	// Execute
	result, err := dashboardService.GetDashboardEnrollments(userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to fetch enrollments")
	assert.Contains(t, err.Error(), "database connection failed")

	mockEnrollmentRepo.AssertExpectations(t)
}

// Test GetDashboardEnrollments - Empty Enrollments
func TestDashboardService_GetDashboardEnrollments_EmptyResult(t *testing.T) {
	// Setup
	mockEnrollmentRepo := &MockDashboardEnrollmentRepository{}
	mockAttemptRepo := &MockUserExamAttemptRepository{}
	dashboardService := service.NewDashboardService(nil, mockEnrollmentRepo, mockAttemptRepo)

	userID := uint(1)

	// Setup expectations
	mockEnrollmentRepo.On("GetUserEnrollments", userID).Return([]models.UserPackageEnrollment{}, nil)

	// Execute
	result, err := dashboardService.GetDashboardEnrollments(userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Enrollments, 0)
	assert.Equal(t, 0, result.Total)
	assert.Equal(t, 0, result.Active)

	mockEnrollmentRepo.AssertExpectations(t)
}

// Test GetDashboardEnrollments - Progress Calculation Error Handling
func TestDashboardService_GetDashboardEnrollments_ProgressCalculationError(t *testing.T) {
	// Setup
	mockEnrollmentRepo := &MockDashboardEnrollmentRepository{}
	mockAttemptRepo := &MockUserExamAttemptRepository{}

	dashboardService := service.NewDashboardService(nil, mockEnrollmentRepo, mockAttemptRepo)

	userID := uint(1)
	testEnrollments := createTestEnrollments()
	activeEnrollments := testEnrollments[:1] // Only first enrollment

	// Setup expectations
	mockEnrollmentRepo.On("GetUserEnrollments", userID).Return(activeEnrollments, nil)

	// Execute (with nil database, progress calculation will fail gracefully)
	result, err := dashboardService.GetDashboardEnrollments(userID)

	// Assert
	assert.NoError(t, err) // Should not fail entire request
	assert.NotNil(t, result)
	assert.Len(t, result.Enrollments, 1)

	// Progress should be 50% due to default test values when db is nil
	enrollment := result.Enrollments[0]
	assert.Equal(t, 50.0, enrollment.Progress)
	assert.Equal(t, 10, enrollment.TotalExams)
	assert.Equal(t, 5, enrollment.CompletedExams)

	mockEnrollmentRepo.AssertExpectations(t)
}

// Test GetDashboardEnrollments - Enrollment Filtering (Expired vs Active)
func TestDashboardService_GetDashboardEnrollments_EnrollmentFiltering(t *testing.T) {
	// Setup
	mockEnrollmentRepo := &MockDashboardEnrollmentRepository{}
	mockAttemptRepo := &MockUserExamAttemptRepository{}

	dashboardService := service.NewDashboardService(nil, mockEnrollmentRepo, mockAttemptRepo)

	userID := uint(1)

	// Create mix of active and expired enrollments
	now := time.Now()
	expiredEnrollment := models.UserPackageEnrollment{
		ID:                  10,
		UserID:              1,
		PackageID:           10,
		EnrollmentType:      models.EnrollmentTypeFull,
		EnrolledAt:          now.Add(-40 * 24 * time.Hour),
		ExpiresAt:           &[]time.Time{now.Add(-5 * 24 * time.Hour)}[0], // Expired 5 days ago
		EnrolledPackageType: models.PackageTypePremium,
		EnrolledPrice:       99.99,
		PaymentStatus:       models.PaymentStatusExpired,
		IsActive:            false,
		Package: models.Package{
			ID:          10,
			Name:        "Expired Course",
			Slug:        "expired-course",
			PackageType: models.PackageTypePremium,
		},
	}

	activeEnrollment := models.UserPackageEnrollment{
		ID:                  11,
		UserID:              1,
		PackageID:           11,
		EnrollmentType:      models.EnrollmentTypeFull,
		EnrolledAt:          now.Add(-10 * 24 * time.Hour),
		ExpiresAt:           &[]time.Time{now.Add(30 * 24 * time.Hour)}[0], // Expires in 30 days
		EnrolledPackageType: models.PackageTypePremium,
		EnrolledPrice:       149.99,
		PaymentStatus:       models.PaymentStatusPaid,
		IsActive:            true,
		Package: models.Package{
			ID:          11,
			Name:        "Active Course",
			Slug:        "active-course",
			PackageType: models.PackageTypePremium,
		},
	}

	mixedEnrollments := []models.UserPackageEnrollment{expiredEnrollment, activeEnrollment}

	// Setup expectations
	mockEnrollmentRepo.On("GetUserEnrollments", userID).Return(mixedEnrollments, nil)

	// Execute
	result, err := dashboardService.GetDashboardEnrollments(userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Should only include active enrollment (expired should be filtered out)
	assert.Len(t, result.Enrollments, 1)
	assert.Equal(t, 1, result.Total)
	assert.Equal(t, 1, result.Active)

	// Verify it's the active enrollment
	enrollment := result.Enrollments[0]
	assert.Equal(t, "Active Course", enrollment.PackageName)
	assert.Equal(t, "active", enrollment.Status)
	assert.Equal(t, 50.0, enrollment.Progress) // 50% from default test values

	mockEnrollmentRepo.AssertExpectations(t)
}

// Test GetDashboardEnrollments - Date Formatting
func TestDashboardService_GetDashboardEnrollments_DateFormatting(t *testing.T) {
	// Setup
	mockEnrollmentRepo := &MockDashboardEnrollmentRepository{}
	mockAttemptRepo := &MockUserExamAttemptRepository{}

	dashboardService := service.NewDashboardService(nil, mockEnrollmentRepo, mockAttemptRepo)

	userID := uint(1)
	now := time.Now()
	expiryDate := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

	enrollmentWithExpiry := models.UserPackageEnrollment{
		ID:                  20,
		UserID:              1,
		PackageID:           20,
		EnrollmentType:      models.EnrollmentTypeFull,
		EnrolledAt:          now.Add(-10 * 24 * time.Hour),
		ExpiresAt:           &expiryDate,
		EnrolledPackageType: models.PackageTypePremium,
		EnrolledPrice:       199.99,
		PaymentStatus:       models.PaymentStatusPaid,
		IsActive:            true,
		Package: models.Package{
			ID:          20,
			Name:        "Course With Expiry",
			Slug:        "course-with-expiry",
			PackageType: models.PackageTypePremium,
		},
	}

	enrollmentWithoutExpiry := models.UserPackageEnrollment{
		ID:                  21,
		UserID:              1,
		PackageID:           21,
		EnrollmentType:      models.EnrollmentTypeFull,
		EnrolledAt:          now.Add(-10 * 24 * time.Hour),
		ExpiresAt:           nil, // No expiry date
		EnrolledPackageType: models.PackageTypeFree,
		EnrolledPrice:       0.00,
		PaymentStatus:       models.PaymentStatusFree,
		IsActive:            true,
		Package: models.Package{
			ID:          21,
			Name:        "Course Without Expiry",
			Slug:        "course-without-expiry",
			PackageType: models.PackageTypeFree,
		},
	}

	enrollments := []models.UserPackageEnrollment{enrollmentWithExpiry, enrollmentWithoutExpiry}

	// Setup expectations
	mockEnrollmentRepo.On("GetUserEnrollments", userID).Return(enrollments, nil)

	// Execute
	result, err := dashboardService.GetDashboardEnrollments(userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Enrollments, 2)

	// Check enrollment with expiry date
	enrollmentWithExpiryResult := result.Enrollments[0]
	assert.NotNil(t, enrollmentWithExpiryResult.ExpiryDate)
	assert.Equal(t, "2025-12-31T23:59:59Z", *enrollmentWithExpiryResult.ExpiryDate)

	// Check enrollment without expiry date
	enrollmentWithoutExpiryResult := result.Enrollments[1]
	assert.Nil(t, enrollmentWithoutExpiryResult.ExpiryDate)

	mockEnrollmentRepo.AssertExpectations(t)
}
