package unit

import (
	"github.com/Mahfuz2811/medecole/backend/internal/models"
	"github.com/Mahfuz2811/medecole/backend/internal/repository"
	"github.com/Mahfuz2811/medecole/backend/internal/types"
	"time"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDashboardEnrollmentRepository for dashboard service tests
type MockDashboardEnrollmentRepository struct {
	mock.Mock
}

func (m *MockDashboardEnrollmentRepository) CreateEnrollment(enrollment *models.UserPackageEnrollment) error {
	args := m.Called(enrollment)
	return args.Error(0)
}

func (m *MockDashboardEnrollmentRepository) GetEnrollmentByID(id uint) (*models.UserPackageEnrollment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserPackageEnrollment), args.Error(1)
}

func (m *MockDashboardEnrollmentRepository) GetUserEnrollments(userID uint) ([]models.UserPackageEnrollment, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.UserPackageEnrollment), args.Error(1)
}

func (m *MockDashboardEnrollmentRepository) GetActiveEnrollment(userID, packageID uint) (*models.UserPackageEnrollment, error) {
	args := m.Called(userID, packageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserPackageEnrollment), args.Error(1)
}

func (m *MockDashboardEnrollmentRepository) GetPackageByID(packageID uint) (*models.Package, error) {
	args := m.Called(packageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Package), args.Error(1)
}

func (m *MockDashboardEnrollmentRepository) GetCouponByCode(code string) (*models.Coupon, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Coupon), args.Error(1)
}

func (m *MockDashboardEnrollmentRepository) ValidateCoupon(coupon *models.Coupon, packageID uint) error {
	args := m.Called(coupon, packageID)
	return args.Error(0)
}

func (m *MockDashboardEnrollmentRepository) IncrementCouponUsage(couponID uint) error {
	args := m.Called(couponID)
	return args.Error(0)
}

func (m *MockDashboardEnrollmentRepository) CreateCouponUsage(usage *models.CouponUsage) error {
	args := m.Called(usage)
	return args.Error(0)
}

func (m *MockDashboardEnrollmentRepository) WithTransaction(tx *gorm.DB) repository.EnrollmentRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.EnrollmentRepository)
}

// IsUserEnrolledInPackage mock implementation for Phase 2 optimization
func (m *MockDashboardEnrollmentRepository) IsUserEnrolledInPackage(userID, packageID uint) (bool, error) {
	args := m.Called(userID, packageID)
	return args.Bool(0), args.Error(1)
}

// MockUserExamAttemptRepository mocks the UserExamAttemptRepository interface
type MockUserExamAttemptRepository struct {
	mock.Mock
}

func (m *MockUserExamAttemptRepository) GetUserStats(userID uint) (*types.UserStatsData, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.UserStatsData), args.Error(1)
}

func (m *MockUserExamAttemptRepository) GetRecentActivity(userID uint, limit int) ([]models.UserExamAttempt, error) {
	args := m.Called(userID, limit)
	return args.Get(0).([]models.UserExamAttempt), args.Error(1)
}

// MockDB mocks GORM DB for specific queries
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Model(value interface{}) *MockDB {
	m.Called(value)
	return m
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *MockDB {
	m.Called(query, args)
	return m
}

func (m *MockDB) Joins(query string, args ...interface{}) *MockDB {
	m.Called(query, args)
	return m
}

func (m *MockDB) Count(count *int64) *MockDB {
	args := m.Called(count)
	if args.Get(0) != nil {
		*count = int64(args.Int(0))
	}
	return m
}

func (m *MockDB) Distinct(args ...interface{}) *MockDB {
	m.Called(args)
	return m
}

// Test data generators for dashboard service tests
func createDashboardTestUserStats() *types.UserStatsData {
	return &types.UserStatsData{
		TotalAttempts:  15,
		CorrectAnswers: 120,
		TotalQuestions: 150,
		TotalTimeSpent: 3600, // 1 hour
		AverageScore:   80.0,
	}
}

func createDashboardTestExamAttempts() []models.UserExamAttempt {
	now := time.Now()
	correctAnswers1 := 17
	correctAnswers2 := 23

	return []models.UserExamAttempt{
		{
			ID:              1,
			UserID:          1,
			ExamID:          1,
			Status:          models.AttemptStatusCompleted,
			StartedAt:       now.Add(-2 * time.Hour),
			CompletedAt:     &now,
			ActualTimeSpent: 1800, // 30 minutes
			TotalQuestions:  20,
			Score:           ptrDashboardFloat64(85.0),
			CorrectAnswers:  &correctAnswers1,
			Exam: models.Exam{
				ID:    1,
				Title: "JavaScript Basics Quiz",
			},
		},
		{
			ID:              2,
			UserID:          1,
			ExamID:          2,
			Status:          models.AttemptStatusCompleted,
			StartedAt:       now.Add(-1 * time.Hour),
			CompletedAt:     &now,
			ActualTimeSpent: 2400, // 40 minutes
			TotalQuestions:  25,
			Score:           ptrDashboardFloat64(92.0),
			CorrectAnswers:  &correctAnswers2,
			Exam: models.Exam{
				ID:    2,
				Title: "React Fundamentals",
			},
		},
	}
}

// Helper functions for dashboard service tests
func ptrDashboardFloat64(f float64) *float64 {
	return &f
}
