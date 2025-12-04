package unit

import (
	"context"
	"errors"
	"github.com/Mahfuz2811/medecole/backend/internal/dto"
	"github.com/Mahfuz2811/medecole/backend/internal/mapper"
	"github.com/Mahfuz2811/medecole/backend/internal/models"
	"github.com/Mahfuz2811/medecole/backend/internal/repository"
	"github.com/Mahfuz2811/medecole/backend/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockEnrollmentRepository mocks the EnrollmentRepository interface
type MockEnrollmentRepository struct {
	mock.Mock
}

func (m *MockEnrollmentRepository) CreateEnrollment(enrollment *models.UserPackageEnrollment) error {
	args := m.Called(enrollment)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) GetEnrollmentByID(id uint) (*models.UserPackageEnrollment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserPackageEnrollment), args.Error(1)
}

func (m *MockEnrollmentRepository) GetUserEnrollments(userID uint) ([]models.UserPackageEnrollment, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.UserPackageEnrollment), args.Error(1)
}

func (m *MockEnrollmentRepository) GetActiveEnrollment(userID, packageID uint) (*models.UserPackageEnrollment, error) {
	args := m.Called(userID, packageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserPackageEnrollment), args.Error(1)
}

func (m *MockEnrollmentRepository) GetPackageByID(packageID uint) (*models.Package, error) {
	args := m.Called(packageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Package), args.Error(1)
}

func (m *MockEnrollmentRepository) GetCouponByCode(code string) (*models.Coupon, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Coupon), args.Error(1)
}

func (m *MockEnrollmentRepository) ValidateCoupon(coupon *models.Coupon, packageID uint) error {
	args := m.Called(coupon, packageID)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) IncrementCouponUsage(couponID uint) error {
	args := m.Called(couponID)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) CreateCouponUsage(usage *models.CouponUsage) error {
	args := m.Called(usage)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) WithTransaction(tx *gorm.DB) repository.EnrollmentRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.EnrollmentRepository)
}

// IsUserEnrolledInPackage mock implementation for Phase 2 optimization
func (m *MockEnrollmentRepository) IsUserEnrolledInPackage(userID, packageID uint) (bool, error) {
	args := m.Called(userID, packageID)
	return args.Bool(0), args.Error(1)
}

// Test data generators
func createEnrollmentTestPackage() *models.Package {
	return &models.Package{
		ID:          1,
		Name:        "Test Package",
		Slug:        "test-package",
		PackageType: models.PackageTypePremium,
		Price:       99.99,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func createEnrollmentTestCoupon() *models.Coupon {
	validUntil := time.Now().Add(24 * time.Hour)
	usageLimit := 100 // Use int instead of uint for UsageLimit

	return &models.Coupon{
		ID:                 1,
		Code:               "SAVE20",
		DiscountPercentage: 20.0,
		ValidFrom:          time.Now().Add(-24 * time.Hour),
		ValidUntil:         &validUntil,
		UsageLimit:         &usageLimit,
		UsageCount:         0,
		IsActive:           true,
		Status:             models.CouponStatusActive,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

func createEnrollmentTestEnrollment() *models.UserPackageEnrollment {
	return &models.UserPackageEnrollment{
		ID:                  1,
		UserID:              1,
		PackageID:           1,
		EnrollmentType:      models.EnrollmentTypeFull, // Use correct enum
		EnrolledAt:          time.Now(),
		EnrolledPackageType: models.PackageTypePremium,
		EnrolledPrice:       79.99,
		PaymentStatus:       models.PaymentStatusPaid, // Use correct enum
		PaymentAmount:       ptrFloat64(79.99),
		IsActive:            true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}

// Helper functions
func ptrFloat64(f float64) *float64 {
	return &f
}

// Since the service uses direct GORM database transactions, we'll test the non-transaction methods
// and create separate integration tests for transaction-based methods

// Test CheckEnrollmentStatus - User has active enrollment
func TestEnrollmentService_CheckEnrollmentStatus_ActiveEnrollment(t *testing.T) {
	// Setup
	mockRepo := &MockEnrollmentRepository{}
	realMapper := mapper.NewEnrollmentMapper()
	enrollmentService := service.NewEnrollmentService(mockRepo, realMapper, &gorm.DB{})

	userID := uint(1)
	packageID := uint(1)

	testEnrollment := createEnrollmentTestEnrollment()

	// Setup expectations
	mockRepo.On("GetActiveEnrollment", userID, packageID).Return(testEnrollment, nil)

	// Execute
	result, err := enrollmentService.CheckEnrollmentStatus(context.Background(), userID, packageID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.HasActiveEnrollment)
	assert.NotNil(t, result.Enrollment)
	assert.Equal(t, testEnrollment.ID, result.Enrollment.ID)

	mockRepo.AssertExpectations(t)
}

// Test CheckEnrollmentStatus - No active enrollment
func TestEnrollmentService_CheckEnrollmentStatus_NoActiveEnrollment(t *testing.T) {
	// Setup
	mockRepo := &MockEnrollmentRepository{}
	realMapper := mapper.NewEnrollmentMapper()
	enrollmentService := service.NewEnrollmentService(mockRepo, realMapper, &gorm.DB{})

	userID := uint(1)
	packageID := uint(1)

	// Setup expectations
	mockRepo.On("GetActiveEnrollment", userID, packageID).Return(nil, nil)

	// Execute
	result, err := enrollmentService.CheckEnrollmentStatus(context.Background(), userID, packageID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.HasActiveEnrollment)
	assert.Nil(t, result.Enrollment)

	mockRepo.AssertExpectations(t)
}

// Test CheckEnrollmentStatus - Repository error
func TestEnrollmentService_CheckEnrollmentStatus_RepositoryError(t *testing.T) {
	// Setup
	mockRepo := &MockEnrollmentRepository{}
	realMapper := mapper.NewEnrollmentMapper()
	enrollmentService := service.NewEnrollmentService(mockRepo, realMapper, &gorm.DB{})

	userID := uint(1)
	packageID := uint(1)

	// Setup expectations
	mockRepo.On("GetActiveEnrollment", userID, packageID).Return(nil, errors.New("database error"))

	// Execute
	result, err := enrollmentService.CheckEnrollmentStatus(context.Background(), userID, packageID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to check enrollment status")

	mockRepo.AssertExpectations(t)
}

// Test ValidateCoupon - Valid coupon
func TestEnrollmentService_ValidateCoupon_ValidCoupon(t *testing.T) {
	// Setup
	mockRepo := &MockEnrollmentRepository{}
	realMapper := mapper.NewEnrollmentMapper()
	enrollmentService := service.NewEnrollmentService(mockRepo, realMapper, &gorm.DB{})

	req := dto.CouponValidationRequest{
		CouponCode: "SAVE20",
		PackageID:  1,
	}

	testPackage := createEnrollmentTestPackage()
	testCoupon := createEnrollmentTestCoupon()

	// Setup expectations
	mockRepo.On("GetPackageByID", uint(1)).Return(testPackage, nil)
	mockRepo.On("GetCouponByCode", "SAVE20").Return(testCoupon, nil)
	mockRepo.On("ValidateCoupon", testCoupon, uint(1)).Return(nil)

	// Execute
	result, err := enrollmentService.ValidateCoupon(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Valid)
	assert.Equal(t, "SAVE20", result.CouponCode)
	assert.Equal(t, 20.0, result.DiscountPercentage)
	assert.NotNil(t, result.PriceCalculation)

	mockRepo.AssertExpectations(t)
}

// Test ValidateCoupon - Invalid coupon
func TestEnrollmentService_ValidateCoupon_InvalidCoupon(t *testing.T) {
	// Setup
	mockRepo := &MockEnrollmentRepository{}
	realMapper := mapper.NewEnrollmentMapper()
	enrollmentService := service.NewEnrollmentService(mockRepo, realMapper, &gorm.DB{})

	req := dto.CouponValidationRequest{
		CouponCode: "INVALID",
		PackageID:  1,
	}

	testPackage := createEnrollmentTestPackage()

	// Setup expectations
	mockRepo.On("GetPackageByID", uint(1)).Return(testPackage, nil)
	mockRepo.On("GetCouponByCode", "INVALID").Return(nil, repository.ErrCouponNotFound)

	// Execute
	result, err := enrollmentService.ValidateCoupon(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Valid)
	assert.Equal(t, "INVALID", result.CouponCode)
	assert.Contains(t, result.Message, "Coupon not found")

	mockRepo.AssertExpectations(t)
}

// Test ValidateCoupon - Coupon validation fails
func TestEnrollmentService_ValidateCoupon_CouponExpired(t *testing.T) {
	// Setup
	mockRepo := &MockEnrollmentRepository{}
	realMapper := mapper.NewEnrollmentMapper()
	enrollmentService := service.NewEnrollmentService(mockRepo, realMapper, &gorm.DB{})

	req := dto.CouponValidationRequest{
		CouponCode: "EXPIRED",
		PackageID:  1,
	}

	testPackage := createEnrollmentTestPackage()
	testCoupon := createEnrollmentTestCoupon()

	// Setup expectations
	mockRepo.On("GetPackageByID", uint(1)).Return(testPackage, nil)
	mockRepo.On("GetCouponByCode", "EXPIRED").Return(testCoupon, nil)
	mockRepo.On("ValidateCoupon", testCoupon, uint(1)).Return(repository.ErrCouponExpired)

	// Execute
	result, err := enrollmentService.ValidateCoupon(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Valid)
	assert.Contains(t, result.Message, "Coupon has expired")

	mockRepo.AssertExpectations(t)
}

// Test ValidateCoupon - Package not found
func TestEnrollmentService_ValidateCoupon_PackageNotFound(t *testing.T) {
	// Setup
	mockRepo := &MockEnrollmentRepository{}
	realMapper := mapper.NewEnrollmentMapper()
	enrollmentService := service.NewEnrollmentService(mockRepo, realMapper, &gorm.DB{})

	req := dto.CouponValidationRequest{
		CouponCode: "SAVE20",
		PackageID:  999,
	}

	// Setup expectations
	mockRepo.On("GetPackageByID", uint(999)).Return(nil, repository.ErrPackageNotFound)

	// Execute
	result, err := enrollmentService.ValidateCoupon(context.Background(), req)

	// Assert
	assert.NoError(t, err) // The service doesn't return an error, it returns a response with Valid: false
	assert.NotNil(t, result)
	assert.False(t, result.Valid)
	assert.Contains(t, result.Message, "Invalid package")

	mockRepo.AssertExpectations(t)
}

// Test CalculatePrice method
func TestEnrollmentService_CalculatePrice_WithoutCoupon(t *testing.T) {
	// Setup
	mockRepo := &MockEnrollmentRepository{}
	realMapper := mapper.NewEnrollmentMapper()
	enrollmentService := service.NewEnrollmentService(mockRepo, realMapper, &gorm.DB{})

	packagePrice := 99.99

	// Execute
	result := enrollmentService.CalculatePrice(packagePrice, nil)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, 99.99, result.OriginalPrice)
	assert.Equal(t, 0.0, result.DiscountPercentage)
	assert.Equal(t, 99.99, result.FinalPrice)
	assert.Nil(t, result.CouponCode)
}

// Test CalculatePrice with coupon
func TestEnrollmentService_CalculatePrice_WithCoupon(t *testing.T) {
	// Setup
	mockRepo := &MockEnrollmentRepository{}
	realMapper := mapper.NewEnrollmentMapper()
	enrollmentService := service.NewEnrollmentService(mockRepo, realMapper, &gorm.DB{})

	packagePrice := 99.99
	testCoupon := createEnrollmentTestCoupon()

	// Execute
	result := enrollmentService.CalculatePrice(packagePrice, testCoupon)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, 99.99, result.OriginalPrice)
	assert.Equal(t, 20.0, result.DiscountPercentage)
	assert.Equal(t, "SAVE20", *result.CouponCode)
	assert.InDelta(t, 79.99, result.FinalPrice, 0.01) // Allow for floating point precision
}

// Test interface compliance
func TestEnrollmentService_ImplementsInterface(t *testing.T) {
	// Setup
	mockRepo := &MockEnrollmentRepository{}
	realMapper := mapper.NewEnrollmentMapper()
	enrollmentService := service.NewEnrollmentService(mockRepo, realMapper, &gorm.DB{})

	// Verify that our service implements the expected interface
	var _ service.EnrollmentService = enrollmentService
}

// Test error scenarios for edge cases
func TestEnrollmentService_ValidateCoupon_CouponInvalid(t *testing.T) {
	// Setup
	mockRepo := &MockEnrollmentRepository{}
	realMapper := mapper.NewEnrollmentMapper()
	enrollmentService := service.NewEnrollmentService(mockRepo, realMapper, &gorm.DB{})

	req := dto.CouponValidationRequest{
		CouponCode: "INVALID_STATUS",
		PackageID:  1,
	}

	testPackage := createEnrollmentTestPackage()
	testCoupon := createEnrollmentTestCoupon()

	// Setup expectations
	mockRepo.On("GetPackageByID", uint(1)).Return(testPackage, nil)
	mockRepo.On("GetCouponByCode", "INVALID_STATUS").Return(testCoupon, nil)
	mockRepo.On("ValidateCoupon", testCoupon, uint(1)).Return(repository.ErrCouponInvalid)

	// Execute
	result, err := enrollmentService.ValidateCoupon(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Valid)
	assert.Contains(t, result.Message, "Coupon is not valid")

	mockRepo.AssertExpectations(t)
}

// Test CalculatePrice edge cases
func TestEnrollmentService_CalculatePrice_ZeroPrice(t *testing.T) {
	// Setup
	mockRepo := &MockEnrollmentRepository{}
	realMapper := mapper.NewEnrollmentMapper()
	enrollmentService := service.NewEnrollmentService(mockRepo, realMapper, &gorm.DB{})

	packagePrice := 0.0

	// Execute
	result := enrollmentService.CalculatePrice(packagePrice, nil)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, 0.0, result.OriginalPrice)
	assert.Equal(t, 0.0, result.DiscountPercentage)
	assert.Equal(t, 0.0, result.FinalPrice)
	assert.Nil(t, result.CouponCode)
}

// Test business logic validation
func TestEnrollmentService_CalculatePrice_HighDiscountCoupon(t *testing.T) {
	// Setup
	mockRepo := &MockEnrollmentRepository{}
	realMapper := mapper.NewEnrollmentMapper()
	enrollmentService := service.NewEnrollmentService(mockRepo, realMapper, &gorm.DB{})

	packagePrice := 50.0
	testCoupon := createEnrollmentTestCoupon()
	testCoupon.DiscountPercentage = 100.0 // 100% discount

	// Execute
	result := enrollmentService.CalculatePrice(packagePrice, testCoupon)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, 50.0, result.OriginalPrice)
	assert.Equal(t, 100.0, result.DiscountPercentage)
	assert.Equal(t, 0.0, result.FinalPrice) // Should be free
}

/*
NOTE: The EnrollInPackage method uses direct GORM database transactions and would require
integration testing with a real database connection. The transaction-based logic cannot be
easily unit tested without significant refactoring of the service architecture.

For comprehensive testing of EnrollInPackage, consider:
1. Creating integration tests that use a test database
2. Refactoring the service to use a transaction interface that can be mocked
3. Testing the business logic components separately

Example test scenarios that would require integration testing:
- TestEnrollmentService_EnrollInPackage_Success_WithoutCoupon
- TestEnrollmentService_EnrollInPackage_Success_WithCoupon
- TestEnrollmentService_EnrollInPackage_PackageNotFound
- TestEnrollmentService_EnrollInPackage_ActiveEnrollmentExists
- TestEnrollmentService_EnrollInPackage_InvalidCoupon
- TestEnrollmentService_EnrollInPackage_TransactionFails
*/
