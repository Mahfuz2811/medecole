package unit

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"quizora-backend/internal/dto"
	"quizora-backend/internal/models"
	"quizora-backend/internal/repository"
	"quizora-backend/internal/service"
)

// MockExamRepository is a mock implementation of repository.ExamRepository
type MockExamRepository struct {
	mock.Mock
}

func (m *MockExamRepository) GetExamsByPackageSlug(packageSlug string, userID uint) ([]repository.ExamWithUserData, error) {
	args := m.Called(packageSlug, userID)
	return args.Get(0).([]repository.ExamWithUserData), args.Error(1)
}

func (m *MockExamRepository) GetPackageWithExamsBySlug(packageSlug string, userID uint) (*repository.PackageWithExamsData, error) {
	args := m.Called(packageSlug, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PackageWithExamsData), args.Error(1)
}

func (m *MockExamRepository) GetExamBySlug(examSlug string) (*models.Exam, error) {
	args := m.Called(examSlug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Exam), args.Error(1)
}

func (m *MockExamRepository) GetActiveAttemptByUserAndExam(userID uint, examID uint) (*models.UserExamAttempt, error) {
	args := m.Called(userID, examID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserExamAttempt), args.Error(1)
}

func (m *MockExamRepository) GetUserAttemptsByExam(userID uint, examID uint) ([]models.UserExamAttempt, error) {
	args := m.Called(userID, examID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserExamAttempt), args.Error(1)
}

func (m *MockExamRepository) CreateExamAttempt(userID uint, examID uint, packageID uint, deviceInfo map[string]string) (*models.UserExamAttempt, error) {
	args := m.Called(userID, examID, packageID, deviceInfo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserExamAttempt), args.Error(1)
}

func (m *MockExamRepository) GetActiveSessionByID(sessionID string) (*repository.SessionWithExamData, error) {
	args := m.Called(sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.SessionWithExamData), args.Error(1)
}

func (m *MockExamRepository) GetCompletedSessionByID(sessionID string) (*repository.SessionWithExamData, error) {
	args := m.Called(sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.SessionWithExamData), args.Error(1)
}

func (m *MockExamRepository) SyncSessionAnswers(sessionID string, answers map[uint]string) error {
	args := m.Called(sessionID, answers)
	return args.Error(0)
}

func (m *MockExamRepository) GetSessionAnswers(sessionID string) (map[uint]string, error) {
	args := m.Called(sessionID)
	if args.Get(0) == nil {
		return make(map[uint]string), args.Error(1)
	}
	return args.Get(0).(map[uint]string), args.Error(1)
}

func (m *MockExamRepository) CompleteExamAttempt(attemptID uint, score float64, passed bool) error {
	args := m.Called(attemptID, score, passed)
	return args.Error(0)
}

func (m *MockExamRepository) CompleteExamAttemptWithAnswers(attemptID uint, score float64, passed bool, answersData string, correctAnswers int) error {
	args := m.Called(attemptID, score, passed, answersData, correctAnswers)
	return args.Error(0)
}

func (m *MockExamRepository) GetAttemptBySessionAndUser(sessionID string, userID uint) (*models.UserExamAttempt, error) {
	args := m.Called(sessionID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserExamAttempt), args.Error(1)
}

// Phase 1 & 2 optimization methods
func (m *MockExamRepository) GetUserAttemptForExam(userID uint, examID uint) (*models.UserExamAttempt, error) {
	args := m.Called(userID, examID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserExamAttempt), args.Error(1)
}

func (m *MockExamRepository) GetUserAttemptForExamInPackage(userID uint, examID uint, packageID uint) (*models.UserExamAttempt, error) {
	args := m.Called(userID, examID, packageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserExamAttempt), args.Error(1)
}

func (m *MockExamRepository) CreateExamAttemptWithExam(userID uint, exam *models.Exam, packageID uint, deviceInfo map[string]string) (*models.UserExamAttempt, error) {
	args := m.Called(userID, exam, packageID, deviceInfo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserExamAttempt), args.Error(1)
}

func (m *MockExamRepository) GetPackageIDForExam(examID uint) (uint, error) {
	args := m.Called(examID)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockExamRepository) MarkExpiredSessionsAsAbandoned(currentTime time.Time, gracePeriodSeconds int) (int64, error) {
	args := m.Called(currentTime, gracePeriodSeconds)
	return args.Get(0).(int64), args.Error(1)
}

// MockExamMapper is a mock implementation of mapper.ExamMapper
type MockExamMapper struct {
	mock.Mock
}

func (m *MockExamMapper) ToExamResponse(exam repository.ExamWithUserData) dto.ExamResponse {
	args := m.Called(exam)
	return args.Get(0).(dto.ExamResponse)
}

func (m *MockExamMapper) ToExamListResponse(exams []repository.ExamWithUserData) dto.ExamListResponse {
	args := m.Called(exams)
	return args.Get(0).(dto.ExamListResponse)
}

func (m *MockExamMapper) ToExamListResponseWithPackage(pkg models.Package, exams []repository.ExamWithUserData) dto.ExamListResponse {
	args := m.Called(pkg, exams)
	return args.Get(0).(dto.ExamListResponse)
}

func (m *MockExamMapper) ToExamContentResponse(exam repository.ExamWithUserData) dto.ExamContentResponse {
	args := m.Called(exam)
	return args.Get(0).(dto.ExamContentResponse)
}

func (m *MockExamMapper) ToExamMetaResponse(exam models.Exam) dto.ExamMetaResponse {
	args := m.Called(exam)
	return args.Get(0).(dto.ExamMetaResponse)
}

func (m *MockExamMapper) ToExamSessionResponse(attempt models.UserExamAttempt, exam models.Exam) dto.ExamSessionResponse {
	args := m.Called(attempt, exam)
	return args.Get(0).(dto.ExamSessionResponse)
}

func (m *MockExamMapper) ToExamSessionResponseWithAnswers(attempt models.UserExamAttempt, exam models.Exam, savedAnswers []dto.UserAnswerResponse) dto.ExamSessionResponse {
	args := m.Called(attempt, exam, savedAnswers)
	return args.Get(0).(dto.ExamSessionResponse)
}

// Test data helpers for exam list service
func createTestPackageForExamList() models.Package {
	description := "Frontend Development Bootcamp Package"
	validityDays := 30

	return models.Package{
		ID:              1,
		Name:            "Frontend Bootcamp",
		Slug:            "frontend-bootcamp",
		Description:     &description,
		PackageType:     models.PackageTypeFree,
		Price:           0.0,
		ValidityType:    models.ValidityTypeRelative,
		ValidityDays:    &validityDays,
		TotalExams:      2,
		EnrollmentCount: 150,
		IsActive:        true,
		SortOrder:       1,
	}
}

func createTestPackageWithExamsData() *repository.PackageWithExamsData {
	pkg := createTestPackageForExamList()
	exams := createTestExamWithUserData()

	return &repository.PackageWithExamsData{
		Package: pkg,
		Exams:   exams,
	}
}

func createExpectedEnhancedExamListResponse() dto.ExamListResponse {
	scheduledStartDate := "2025-08-08T00:00:00Z"

	return dto.ExamListResponse{
		Package: dto.PackageInfoResponse{
			ID:                    1,
			Name:                  "Frontend Bootcamp",
			Slug:                  "frontend-bootcamp",
			Description:           stringPtr("Frontend Development Bootcamp Package"),
			PackageType:           models.PackageTypeFree,
			Price:                 0.0,
			ValidityType:          models.ValidityTypeRelative,
			ValidityDays:          intPtr(30),
			TotalExams:            2,
			EnrollmentCount:       150,
			ActiveEnrollmentCount: 0,
		},
		Exams: []dto.ExamResponse{
			{
				ID:                 1,
				Title:              "JavaScript Fundamentals",
				Slug:               "javascript-fundamentals",
				Description:        "Test exam description",
				ExamType:           "DAILY",
				TotalQuestions:     20,
				DurationMinutes:    60,
				PassingScore:       70.0,
				ScheduledStartDate: &scheduledStartDate,
				ScheduledEndDate:   nil,
				AttemptCount:       150,
				AverageScore:       85.5,
				PassRate:           75.0,
				ComputedStatus:     "AVAILABLE",
				SortOrder:          1,
				HasAttempted:       true,
			},
			{
				ID:                 2,
				Title:              "React Components",
				Slug:               "react-components",
				Description:        "",
				ExamType:           "MOCK",
				TotalQuestions:     25,
				DurationMinutes:    90,
				PassingScore:       80.0,
				ScheduledStartDate: nil,
				ScheduledEndDate:   nil,
				AttemptCount:       85,
				AverageScore:       0.0,
				PassRate:           0.0,
				ComputedStatus:     "UPCOMING",
				SortOrder:          2,
				HasAttempted:       false,
			},
		},
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func createTestExamWithUserData() []repository.ExamWithUserData {
	now := time.Now()
	description := "Test exam description"
	averageScore := 85.5
	passRate := 75.0

	return []repository.ExamWithUserData{
		{
			Exam: models.Exam{
				ID:                 1,
				Title:              "JavaScript Fundamentals",
				Slug:               "javascript-fundamentals",
				Description:        &description,
				ExamType:           models.ExamTypeDaily,
				TotalQuestions:     20,
				DurationMinutes:    60,
				PassingScore:       70.0,
				ScheduledStartDate: &now,
				ScheduledEndDate:   nil,
				AttemptCount:       150,
				AverageScore:       &averageScore,
				PassRate:           &passRate,
			},
			HasAttempted:   true,
			SortOrder:      1,
			ComputedStatus: "AVAILABLE",
		},
		{
			Exam: models.Exam{
				ID:                 2,
				Title:              "React Components",
				Slug:               "react-components",
				Description:        nil,
				ExamType:           models.ExamTypeMock,
				TotalQuestions:     25,
				DurationMinutes:    90,
				PassingScore:       80.0,
				ScheduledStartDate: nil,
				ScheduledEndDate:   nil,
				AttemptCount:       85,
				AverageScore:       nil,
				PassRate:           nil,
			},
			HasAttempted:   false,
			SortOrder:      2,
			ComputedStatus: "UPCOMING",
		},
	}
}

// Test GetPackageExamsBySlug - Success (Updated for Enhanced API)
func TestExamService_GetPackageExamsBySlug_Success(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Test data
	packageSlug := "frontend-bootcamp"
	userID := uint(1)
	testPackageWithExams := createTestPackageWithExamsData()
	expectedResponse := createExpectedEnhancedExamListResponse()

	// Setup expectations - now using the enhanced API method
	mockExamRepo.On("GetPackageWithExamsBySlug", packageSlug, userID).Return(testPackageWithExams, nil)
	mockExamMapper.On("ToExamListResponseWithPackage", testPackageWithExams.Package, testPackageWithExams.Exams).Return(expectedResponse)

	// Execute
	result, err := examService.GetPackageExamsBySlug(packageSlug, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)

	// Verify package data is included
	assert.Equal(t, uint(1), result.Package.ID)
	assert.Equal(t, "Frontend Bootcamp", result.Package.Name)
	assert.Equal(t, "frontend-bootcamp", result.Package.Slug)
	assert.Equal(t, models.PackageTypeFree, result.Package.PackageType)

	// Verify exam data
	assert.Len(t, result.Exams, 2)

	// Verify first exam data
	firstExam := result.Exams[0]
	assert.Equal(t, 1, firstExam.ID)
	assert.Equal(t, "JavaScript Fundamentals", firstExam.Title)
	assert.Equal(t, "javascript-fundamentals", firstExam.Slug)
	assert.Equal(t, "Test exam description", firstExam.Description)
	assert.Equal(t, "DAILY", firstExam.ExamType)
	assert.Equal(t, 20, firstExam.TotalQuestions)
	assert.Equal(t, 60, firstExam.DurationMinutes)
	assert.Equal(t, 70.0, firstExam.PassingScore)
	assert.Equal(t, 150, firstExam.AttemptCount)
	assert.Equal(t, 85.5, firstExam.AverageScore)
	assert.Equal(t, 75.0, firstExam.PassRate)
	assert.Equal(t, "AVAILABLE", firstExam.ComputedStatus)
	assert.Equal(t, 1, firstExam.SortOrder)
	assert.True(t, firstExam.HasAttempted)
	assert.NotNil(t, firstExam.ScheduledStartDate)

	// Verify second exam
	secondExam := result.Exams[1]
	assert.Equal(t, 2, secondExam.ID)
	assert.Equal(t, "React Components", secondExam.Title)
	assert.Equal(t, "react-components", secondExam.Slug)
	assert.Equal(t, "", secondExam.Description) // Mapper should handle nil description
	assert.Equal(t, "MOCK", secondExam.ExamType)
	assert.Equal(t, 25, secondExam.TotalQuestions)
	assert.Equal(t, 90, secondExam.DurationMinutes)
	assert.Equal(t, 80.0, secondExam.PassingScore)
	assert.Equal(t, 85, secondExam.AttemptCount)
	assert.Equal(t, 0.0, secondExam.AverageScore) // Mapper should handle nil values
	assert.Equal(t, 0.0, secondExam.PassRate)     // Mapper should handle nil values
	assert.Equal(t, "UPCOMING", secondExam.ComputedStatus)
	assert.Equal(t, 2, secondExam.SortOrder)
	assert.False(t, secondExam.HasAttempted)
	assert.Nil(t, secondExam.ScheduledStartDate)
	assert.Nil(t, secondExam.ScheduledEndDate)

	// Verify mock expectations
	mockExamRepo.AssertExpectations(t)
	mockExamMapper.AssertExpectations(t)
}

// Test GetPackageExamsBySlug - Repository Error (Updated for Enhanced API)
func TestExamService_GetPackageExamsBySlug_RepositoryError(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Test data
	packageSlug := "invalid-package"
	userID := uint(1)
	expectedError := errors.New("failed to fetch package: package not found")

	// Setup expectations
	mockExamRepo.On("GetPackageWithExamsBySlug", packageSlug, userID).Return(nil, expectedError)

	// Execute
	result, err := examService.GetPackageExamsBySlug(packageSlug, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, dto.ExamListResponse{}, result)
	assert.Contains(t, err.Error(), "failed to fetch package")

	// Verify mock expectations
	mockExamRepo.AssertExpectations(t)
	// Mapper should not be called when repository fails
	mockExamMapper.AssertNotCalled(t, "ToExamListResponseWithPackage")
}

// Test GetPackageExamsBySlug - Empty Results (Updated for Enhanced API)
func TestExamService_GetPackageExamsBySlug_EmptyResults(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Test data
	packageSlug := "empty-package"
	userID := uint(1)

	emptyPackageData := &repository.PackageWithExamsData{
		Package: models.Package{
			ID:              1,
			Name:            "Empty Package",
			Slug:            "empty-package",
			PackageType:     models.PackageTypeFree,
			Price:           0.0,
			ValidityType:    models.ValidityTypeRelative,
			TotalExams:      0,
			EnrollmentCount: 0,
			IsActive:        true,
		},
		Exams: []repository.ExamWithUserData{},
	}

	emptyResponse := dto.ExamListResponse{
		Package: dto.PackageInfoResponse{
			ID:                    1,
			Name:                  "Empty Package",
			Slug:                  "empty-package",
			PackageType:           models.PackageTypeFree,
			Price:                 0.0,
			ValidityType:          models.ValidityTypeRelative,
			TotalExams:            0,
			EnrollmentCount:       0,
			ActiveEnrollmentCount: 0,
		},
		Exams: []dto.ExamResponse{},
	}

	// Setup expectations
	mockExamRepo.On("GetPackageWithExamsBySlug", packageSlug, userID).Return(emptyPackageData, nil)
	mockExamMapper.On("ToExamListResponseWithPackage", emptyPackageData.Package, emptyPackageData.Exams).Return(emptyResponse)

	// Execute
	result, err := examService.GetPackageExamsBySlug(packageSlug, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, result)
	assert.Equal(t, "Empty Package", result.Package.Name)
	assert.Len(t, result.Exams, 0)

	// Verify mock expectations
	mockExamRepo.AssertExpectations(t)
	mockExamMapper.AssertExpectations(t)
}

// Test GetPackageExamsBySlug - Different User IDs (Updated for Enhanced API)
func TestExamService_GetPackageExamsBySlug_DifferentUserIDs(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Test data
	packageSlug := "frontend-bootcamp"
	user1ID := uint(1)
	user2ID := uint(2)

	// Package data (same for both users)
	pkg := models.Package{
		ID:              1,
		Name:            "Frontend Bootcamp",
		Slug:            packageSlug,
		PackageType:     models.PackageTypeFree,
		Price:           0.0,
		ValidityType:    models.ValidityTypeRelative,
		TotalExams:      1,
		EnrollmentCount: 100,
		IsActive:        true,
	}

	// User 1 has attempted exams
	user1Exams := []repository.ExamWithUserData{
		{
			Exam: models.Exam{
				ID:              1,
				Title:           "JavaScript Fundamentals",
				Slug:            "javascript-fundamentals",
				ExamType:        models.ExamTypeDaily,
				TotalQuestions:  20,
				DurationMinutes: 60,
				PassingScore:    70.0,
			},
			HasAttempted:   true,
			SortOrder:      1,
			ComputedStatus: "AVAILABLE",
		},
	}

	// User 2 has not attempted any exams
	user2Exams := []repository.ExamWithUserData{
		{
			Exam: models.Exam{
				ID:              1,
				Title:           "JavaScript Fundamentals",
				Slug:            "javascript-fundamentals",
				ExamType:        models.ExamTypeDaily,
				TotalQuestions:  20,
				DurationMinutes: 60,
				PassingScore:    70.0,
			},
			HasAttempted:   false,
			SortOrder:      1,
			ComputedStatus: "AVAILABLE",
		},
	}

	// Package with exams data for user 1
	user1PackageData := &repository.PackageWithExamsData{
		Package: pkg,
		Exams:   user1Exams,
	}

	// Package with exams data for user 2
	user2PackageData := &repository.PackageWithExamsData{
		Package: pkg,
		Exams:   user2Exams,
	}

	user1Response := dto.ExamListResponse{
		Package: dto.PackageInfoResponse{
			ID:                    1,
			Name:                  "Frontend Bootcamp",
			Slug:                  packageSlug,
			PackageType:           models.PackageTypeFree,
			Price:                 0.0,
			ValidityType:          models.ValidityTypeRelative,
			TotalExams:            1,
			EnrollmentCount:       100,
			ActiveEnrollmentCount: 0,
		},
		Exams: []dto.ExamResponse{
			{
				ID:              1,
				Title:           "JavaScript Fundamentals",
				Slug:            "javascript-fundamentals",
				ExamType:        "DAILY",
				TotalQuestions:  20,
				DurationMinutes: 60,
				PassingScore:    70.0,
				ComputedStatus:  "AVAILABLE",
				SortOrder:       1,
				HasAttempted:    true,
			},
		},
	}

	user2Response := dto.ExamListResponse{
		Package: dto.PackageInfoResponse{
			ID:                    1,
			Name:                  "Frontend Bootcamp",
			Slug:                  packageSlug,
			PackageType:           models.PackageTypeFree,
			Price:                 0.0,
			ValidityType:          models.ValidityTypeRelative,
			TotalExams:            1,
			EnrollmentCount:       100,
			ActiveEnrollmentCount: 0,
		},
		Exams: []dto.ExamResponse{
			{
				ID:              1,
				Title:           "JavaScript Fundamentals",
				Slug:            "javascript-fundamentals",
				ExamType:        "DAILY",
				TotalQuestions:  20,
				DurationMinutes: 60,
				PassingScore:    70.0,
				ComputedStatus:  "AVAILABLE",
				SortOrder:       1,
				HasAttempted:    false,
			},
		},
	}

	// Setup expectations for user 1
	mockExamRepo.On("GetPackageWithExamsBySlug", packageSlug, user1ID).Return(user1PackageData, nil)
	mockExamMapper.On("ToExamListResponseWithPackage", pkg, user1Exams).Return(user1Response)

	// Setup expectations for user 2
	mockExamRepo.On("GetPackageWithExamsBySlug", packageSlug, user2ID).Return(user2PackageData, nil)
	mockExamMapper.On("ToExamListResponseWithPackage", pkg, user2Exams).Return(user2Response)

	// Execute for user 1
	result1, err1 := examService.GetPackageExamsBySlug(packageSlug, user1ID)

	// Execute for user 2
	result2, err2 := examService.GetPackageExamsBySlug(packageSlug, user2ID)

	// Assert user 1 results
	assert.NoError(t, err1)
	assert.True(t, result1.Exams[0].HasAttempted)
	assert.Equal(t, "Frontend Bootcamp", result1.Package.Name)

	// Assert user 2 results
	assert.NoError(t, err2)
	assert.False(t, result2.Exams[0].HasAttempted)
	assert.Equal(t, "Frontend Bootcamp", result2.Package.Name)

	// Verify mock expectations
	mockExamRepo.AssertExpectations(t)
	mockExamMapper.AssertExpectations(t)
}

// Test GetPackageExamsBySlug - Various Exam Types and Statuses (Updated for Enhanced API)
func TestExamService_GetPackageExamsBySlug_VariousExamTypesAndStatuses(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Test data with different exam types and statuses
	packageSlug := "comprehensive-package"
	userID := uint(1)

	// Package data
	pkg := models.Package{
		ID:              1,
		Name:            "Comprehensive Package",
		Slug:            packageSlug,
		PackageType:     models.PackageTypePremium,
		Price:           99.99,
		ValidityType:    models.ValidityTypeRelative,
		TotalExams:      4,
		EnrollmentCount: 250,
		IsActive:        true,
	}

	variedExams := []repository.ExamWithUserData{
		{
			Exam: models.Exam{
				ID:              1,
				Title:           "Daily Practice",
				ExamType:        models.ExamTypeDaily,
				TotalQuestions:  10,
				DurationMinutes: 30,
			},
			ComputedStatus: "AVAILABLE",
			SortOrder:      1,
			HasAttempted:   false,
		},
		{
			Exam: models.Exam{
				ID:              2,
				Title:           "Mock Test",
				ExamType:        models.ExamTypeMock,
				TotalQuestions:  50,
				DurationMinutes: 120,
			},
			ComputedStatus: "UPCOMING",
			SortOrder:      2,
			HasAttempted:   false,
		},
		{
			Exam: models.Exam{
				ID:              3,
				Title:           "Review Session",
				ExamType:        models.ExamTypeReview,
				TotalQuestions:  15,
				DurationMinutes: 45,
			},
			ComputedStatus: "LIVE",
			SortOrder:      3,
			HasAttempted:   true,
		},
		{
			Exam: models.Exam{
				ID:              4,
				Title:           "Final Exam",
				ExamType:        models.ExamTypeFinal,
				TotalQuestions:  100,
				DurationMinutes: 180,
			},
			ComputedStatus: "COMPLETED",
			SortOrder:      4,
			HasAttempted:   true,
		},
	}

	// Package with exams data
	packageWithExams := &repository.PackageWithExamsData{
		Package: pkg,
		Exams:   variedExams,
	}

	expectedResponse := dto.ExamListResponse{
		Package: dto.PackageInfoResponse{
			ID:                    1,
			Name:                  "Comprehensive Package",
			Slug:                  packageSlug,
			PackageType:           models.PackageTypePremium,
			Price:                 99.99,
			ValidityType:          models.ValidityTypeRelative,
			TotalExams:            4,
			EnrollmentCount:       250,
			ActiveEnrollmentCount: 0,
		},
		Exams: []dto.ExamResponse{
			{
				ID:              1,
				Title:           "Daily Practice",
				ExamType:        "DAILY",
				TotalQuestions:  10,
				DurationMinutes: 30,
				ComputedStatus:  "AVAILABLE",
				SortOrder:       1,
				HasAttempted:    false,
			},
			{
				ID:              2,
				Title:           "Mock Test",
				ExamType:        "MOCK",
				TotalQuestions:  50,
				DurationMinutes: 120,
				ComputedStatus:  "UPCOMING",
				SortOrder:       2,
				HasAttempted:    false,
			},
			{
				ID:              3,
				Title:           "Review Session",
				ExamType:        "REVIEW",
				TotalQuestions:  15,
				DurationMinutes: 45,
				ComputedStatus:  "LIVE",
				SortOrder:       3,
				HasAttempted:    true,
			},
			{
				ID:              4,
				Title:           "Final Exam",
				ExamType:        "FINAL",
				TotalQuestions:  100,
				DurationMinutes: 180,
				ComputedStatus:  "COMPLETED",
				SortOrder:       4,
				HasAttempted:    true,
			},
		},
	}

	// Setup expectations
	mockExamRepo.On("GetPackageWithExamsBySlug", packageSlug, userID).Return(packageWithExams, nil)
	mockExamMapper.On("ToExamListResponseWithPackage", pkg, variedExams).Return(expectedResponse)

	// Execute
	result, err := examService.GetPackageExamsBySlug(packageSlug, userID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result.Exams, 4)
	assert.Equal(t, "Comprehensive Package", result.Package.Name)
	assert.Equal(t, models.PackageTypePremium, result.Package.PackageType)

	// Verify different exam types are handled correctly
	examTypes := []string{"DAILY", "MOCK", "REVIEW", "FINAL"}
	statuses := []string{"AVAILABLE", "UPCOMING", "LIVE", "COMPLETED"}
	attemptedStates := []bool{false, false, true, true}

	for i, exam := range result.Exams {
		assert.Equal(t, examTypes[i], exam.ExamType)
		assert.Equal(t, statuses[i], exam.ComputedStatus)
		assert.Equal(t, attemptedStates[i], exam.HasAttempted)
		assert.Equal(t, i+1, exam.SortOrder) // Verify sorting order
	}

	// Verify mock expectations
	mockExamRepo.AssertExpectations(t)
	mockExamMapper.AssertExpectations(t)
}

// Test GetPackageExamsBySlug - Edge Cases with Package Slug (Updated for Enhanced API)
func TestExamService_GetPackageExamsBySlug_EdgeCasesPackageSlug(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	testCases := []struct {
		name        string
		packageSlug string
		userID      uint
		expectError bool
	}{
		{
			name:        "Empty package slug",
			packageSlug: "",
			userID:      1,
			expectError: false, // Service doesn't validate, repository handles it
		},
		{
			name:        "Package slug with special characters",
			packageSlug: "package-with-special_chars.123",
			userID:      1,
			expectError: false,
		},
		{
			name:        "Very long package slug",
			packageSlug: "very-long-package-slug-that-might-cause-issues-in-some-systems-but-should-be-handled-gracefully",
			userID:      1,
			expectError: false,
		},
		{
			name:        "Zero user ID",
			packageSlug: "valid-package",
			userID:      0,
			expectError: false, // Service doesn't validate, repository handles it
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create empty package data for edge cases
			emptyPkg := models.Package{
				ID:              1,
				Name:            "Edge Case Package",
				Slug:            tc.packageSlug,
				PackageType:     models.PackageTypeFree,
				Price:           0.0,
				ValidityType:    models.ValidityTypeRelative,
				TotalExams:      0,
				EnrollmentCount: 0,
				IsActive:        true,
			}

			emptyPackageData := &repository.PackageWithExamsData{
				Package: emptyPkg,
				Exams:   []repository.ExamWithUserData{},
			}

			emptyResponse := dto.ExamListResponse{
				Package: dto.PackageInfoResponse{
					ID:                    1,
					Name:                  "Edge Case Package",
					Slug:                  tc.packageSlug,
					PackageType:           models.PackageTypeFree,
					Price:                 0.0,
					ValidityType:          models.ValidityTypeRelative,
					TotalExams:            0,
					EnrollmentCount:       0,
					ActiveEnrollmentCount: 0,
				},
				Exams: []dto.ExamResponse{},
			}

			// Setup expectations for each test case
			mockExamRepo.On("GetPackageWithExamsBySlug", tc.packageSlug, tc.userID).Return(emptyPackageData, nil).Once()
			mockExamMapper.On("ToExamListResponseWithPackage", emptyPkg, []repository.ExamWithUserData{}).Return(emptyResponse).Once()

			// Execute
			result, err := examService.GetPackageExamsBySlug(tc.packageSlug, tc.userID)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, emptyResponse, result)
				assert.Equal(t, tc.packageSlug, result.Package.Slug)
			}
		})
	}

	// Verify all mock expectations
	mockExamRepo.AssertExpectations(t)
	mockExamMapper.AssertExpectations(t)
}

// Test Service Interface Compliance
func TestExamService_ImplementsInterface(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}

	// Verify that the service implements the interface
	var _ service.ExamService = service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Test passes if compilation succeeds
	assert.True(t, true, "Service implements ExamService interface")
}

// Test Constructor
func TestNewExamService(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}

	// Execute
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Assert
	assert.NotNil(t, examService)
	assert.Implements(t, (*service.ExamService)(nil), examService)
}

// =============================================================================
// NEW TESTS FOR ENHANCED API WITH PACKAGE DATA
// =============================================================================

// Test GetPackageExamsBySlug - Enhanced API with Package Data Success
func TestExamService_GetPackageExamsBySlug_WithPackageData_Success(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Test data
	packageSlug := "frontend-bootcamp"
	userID := uint(1)
	testPackageWithExams := createTestPackageWithExamsData()
	expectedResponse := createExpectedEnhancedExamListResponse()

	// Setup expectations - the service now uses GetPackageWithExamsBySlug
	mockExamRepo.On("GetPackageWithExamsBySlug", packageSlug, userID).Return(testPackageWithExams, nil)
	mockExamMapper.On("ToExamListResponseWithPackage", testPackageWithExams.Package, testPackageWithExams.Exams).Return(expectedResponse)

	// Execute
	result, err := examService.GetPackageExamsBySlug(packageSlug, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)

	// Verify package data is included
	assert.Equal(t, uint(1), result.Package.ID)
	assert.Equal(t, "Frontend Bootcamp", result.Package.Name)
	assert.Equal(t, "frontend-bootcamp", result.Package.Slug)
	assert.Equal(t, models.PackageTypeFree, result.Package.PackageType)
	assert.Equal(t, 0.0, result.Package.Price)
	assert.Equal(t, models.ValidityTypeRelative, result.Package.ValidityType)
	assert.Equal(t, 30, *result.Package.ValidityDays)
	assert.Equal(t, 2, result.Package.TotalExams)
	assert.Equal(t, 150, result.Package.EnrollmentCount)
	assert.Equal(t, "Frontend Development Bootcamp Package", *result.Package.Description)

	// Verify exam data is included
	assert.Len(t, result.Exams, 2)

	// Verify first exam
	firstExam := result.Exams[0]
	assert.Equal(t, 1, firstExam.ID)
	assert.Equal(t, "JavaScript Fundamentals", firstExam.Title)
	assert.Equal(t, "javascript-fundamentals", firstExam.Slug)
	assert.Equal(t, "Test exam description", firstExam.Description)
	assert.Equal(t, "DAILY", firstExam.ExamType)
	assert.Equal(t, 20, firstExam.TotalQuestions)
	assert.Equal(t, 60, firstExam.DurationMinutes)
	assert.Equal(t, 70.0, firstExam.PassingScore)
	assert.Equal(t, 150, firstExam.AttemptCount)
	assert.Equal(t, 85.5, firstExam.AverageScore)
	assert.Equal(t, 75.0, firstExam.PassRate)
	assert.Equal(t, "AVAILABLE", firstExam.ComputedStatus)
	assert.Equal(t, 1, firstExam.SortOrder)
	assert.True(t, firstExam.HasAttempted)

	// Verify mock expectations
	mockExamRepo.AssertExpectations(t)
	mockExamMapper.AssertExpectations(t)
}

// Test GetPackageExamsBySlug - Enhanced API Repository Error
func TestExamService_GetPackageExamsBySlug_WithPackageData_RepositoryError(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Test data
	packageSlug := "invalid-package"
	userID := uint(1)
	expectedError := errors.New("failed to fetch package: package not found")

	// Setup expectations
	mockExamRepo.On("GetPackageWithExamsBySlug", packageSlug, userID).Return(nil, expectedError)

	// Execute
	result, err := examService.GetPackageExamsBySlug(packageSlug, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, dto.ExamListResponse{}, result)
	assert.Contains(t, err.Error(), "failed to fetch package")

	// Verify mock expectations
	mockExamRepo.AssertExpectations(t)
	// Mapper should not be called when repository fails
	mockExamMapper.AssertNotCalled(t, "ToExamListResponseWithPackage")
}

// Test GetPackageExamsBySlug - Enhanced API Empty Package
func TestExamService_GetPackageExamsBySlug_WithPackageData_EmptyPackage(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Test data - package with no exams
	packageSlug := "empty-package"
	userID := uint(1)

	emptyPackageData := &repository.PackageWithExamsData{
		Package: models.Package{
			ID:              1,
			Name:            "Empty Package",
			Slug:            "empty-package",
			PackageType:     models.PackageTypeFree,
			Price:           0.0,
			ValidityType:    models.ValidityTypeRelative,
			TotalExams:      0,
			EnrollmentCount: 0,
			IsActive:        true,
		},
		Exams: []repository.ExamWithUserData{},
	}

	emptyResponse := dto.ExamListResponse{
		Package: dto.PackageInfoResponse{
			ID:                    1,
			Name:                  "Empty Package",
			Slug:                  "empty-package",
			PackageType:           models.PackageTypeFree,
			Price:                 0.0,
			ValidityType:          models.ValidityTypeRelative,
			TotalExams:            0,
			EnrollmentCount:       0,
			ActiveEnrollmentCount: 0,
		},
		Exams: []dto.ExamResponse{},
	}

	// Setup expectations
	mockExamRepo.On("GetPackageWithExamsBySlug", packageSlug, userID).Return(emptyPackageData, nil)
	mockExamMapper.On("ToExamListResponseWithPackage", emptyPackageData.Package, emptyPackageData.Exams).Return(emptyResponse)

	// Execute
	result, err := examService.GetPackageExamsBySlug(packageSlug, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, emptyResponse, result)
	assert.Equal(t, "Empty Package", result.Package.Name)
	assert.Len(t, result.Exams, 0)

	// Verify mock expectations
	mockExamRepo.AssertExpectations(t)
	mockExamMapper.AssertExpectations(t)
}

// Test GetPackageExamsBySlug - Enhanced API Different Package Types
func TestExamService_GetPackageExamsBySlug_WithPackageData_DifferentPackageTypes(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	testCases := []struct {
		name        string
		packageType models.PackageType
		price       float64
	}{
		{
			name:        "Free Package",
			packageType: models.PackageTypeFree,
			price:       0.0,
		},
		{
			name:        "Premium Package",
			packageType: models.PackageTypePremium,
			price:       99.99,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			packageSlug := "test-package"
			userID := uint(1)

			packageData := &repository.PackageWithExamsData{
				Package: models.Package{
					ID:              1,
					Name:            tc.name,
					Slug:            packageSlug,
					PackageType:     tc.packageType,
					Price:           tc.price,
					ValidityType:    models.ValidityTypeRelative,
					TotalExams:      1,
					EnrollmentCount: 50,
					IsActive:        true,
				},
				Exams: []repository.ExamWithUserData{},
			}

			expectedResponse := dto.ExamListResponse{
				Package: dto.PackageInfoResponse{
					ID:                    1,
					Name:                  tc.name,
					Slug:                  packageSlug,
					PackageType:           tc.packageType,
					Price:                 tc.price,
					ValidityType:          models.ValidityTypeRelative,
					TotalExams:            1,
					EnrollmentCount:       50,
					ActiveEnrollmentCount: 0,
				},
				Exams: []dto.ExamResponse{},
			}

			// Setup expectations
			mockExamRepo.On("GetPackageWithExamsBySlug", packageSlug, userID).Return(packageData, nil).Once()
			mockExamMapper.On("ToExamListResponseWithPackage", packageData.Package, packageData.Exams).Return(expectedResponse).Once()

			// Execute
			result, err := examService.GetPackageExamsBySlug(packageSlug, userID)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.packageType, result.Package.PackageType)
			assert.Equal(t, tc.price, result.Package.Price)
		})
	}

	// Verify mock expectations
	mockExamRepo.AssertExpectations(t)
	mockExamMapper.AssertExpectations(t)
}

// Test GetPackageExamsBySlug - Enhanced API Validity Types
func TestExamService_GetPackageExamsBySlug_WithPackageData_ValidityTypes(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	validityDate := time.Now().AddDate(0, 1, 0) // 1 month from now
	validityDays := 30

	testCases := []struct {
		name         string
		validityType models.ValidityType
		validityDays *int
		validityDate *time.Time
	}{
		{
			name:         "Relative Validity",
			validityType: models.ValidityTypeRelative,
			validityDays: &validityDays,
			validityDate: nil,
		},
		{
			name:         "Fixed Validity",
			validityType: models.ValidityTypeFixed,
			validityDays: nil,
			validityDate: &validityDate,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			packageSlug := "validity-test-package"
			userID := uint(1)

			packageData := &repository.PackageWithExamsData{
				Package: models.Package{
					ID:           1,
					Name:         tc.name,
					Slug:         packageSlug,
					PackageType:  models.PackageTypeFree,
					Price:        0.0,
					ValidityType: tc.validityType,
					ValidityDays: tc.validityDays,
					ValidityDate: tc.validityDate,
					TotalExams:   1,
					IsActive:     true,
				},
				Exams: []repository.ExamWithUserData{},
			}

			var validityDateStr *string
			if tc.validityDate != nil {
				str := tc.validityDate.Format("2006-01-02")
				validityDateStr = &str
			}

			expectedResponse := dto.ExamListResponse{
				Package: dto.PackageInfoResponse{
					ID:                    1,
					Name:                  tc.name,
					Slug:                  packageSlug,
					PackageType:           models.PackageTypeFree,
					Price:                 0.0,
					ValidityType:          tc.validityType,
					ValidityDays:          tc.validityDays,
					ValidityDate:          validityDateStr,
					TotalExams:            1,
					EnrollmentCount:       0,
					ActiveEnrollmentCount: 0,
				},
				Exams: []dto.ExamResponse{},
			}

			// Setup expectations
			mockExamRepo.On("GetPackageWithExamsBySlug", packageSlug, userID).Return(packageData, nil).Once()
			mockExamMapper.On("ToExamListResponseWithPackage", packageData.Package, packageData.Exams).Return(expectedResponse).Once()

			// Execute
			result, err := examService.GetPackageExamsBySlug(packageSlug, userID)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.validityType, result.Package.ValidityType)
			assert.Equal(t, tc.validityDays, result.Package.ValidityDays)
			assert.Equal(t, validityDateStr, result.Package.ValidityDate)
		})
	}

	// Verify mock expectations
	mockExamRepo.AssertExpectations(t)
	mockExamMapper.AssertExpectations(t)
}

// =============================================================================
// TESTS FOR SUBMIT EXAM FUNCTIONALITY
// =============================================================================

// Test SubmitExam - Success
func TestExamService_SubmitExam_Success(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Test data
	sessionID := "test_session_123"
	userID := uint(1)

	// Create test exam with questions data
	questionsJSON := `[{
		"id": 1,
		"question_text": "What is JavaScript?",
		"question_type": "SBA",
		"options": {
			"a": {"text": "A programming language", "is_correct": true},
			"b": {"text": "A markup language", "is_correct": false},
			"c": {"text": "A database", "is_correct": false}
		},
		"points": 1
	}]`

	// Create session data
	sessionData := &repository.SessionWithExamData{
		Attempt: models.UserExamAttempt{
			ID:        1,
			UserID:    userID,
			ExamID:    1,
			PackageID: 1,
			Status:    models.AttemptStatusStarted,
		},
		Exam: models.Exam{
			ID:            1,
			Title:         "Test Exam",
			PassingScore:  70.0,
			QuestionsData: questionsJSON,
		},
	}

	// User answers - correct answer
	userAnswers := map[uint]string{
		1: "a", // Correct answer
	}

	// Setup expectations
	mockExamRepo.On("GetActiveSessionByID", sessionID).Return(sessionData, nil)
	mockExamRepo.On("GetSessionAnswers", sessionID).Return(userAnswers, nil)
	mockExamRepo.On("CompleteExamAttemptWithAnswers", uint(1), 1.0, false, mock.AnythingOfType("string"), 1).Return(nil)

	// Execute
	result, err := examService.SubmitExam(sessionID, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, sessionID, result.SessionID)
	assert.Equal(t, 1.0, result.Score)
	assert.False(t, result.Passed) // 1 point is less than 70 passing score
	assert.Equal(t, 1, result.TotalQuestions)
	assert.Equal(t, 1, result.CorrectAnswers)

	// Verify mock expectations
	mockExamRepo.AssertExpectations(t)
}

// Test SubmitExam - Session Not Found
func TestExamService_SubmitExam_SessionNotFound(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Test data
	sessionID := "invalid_session"
	userID := uint(1)

	// Setup expectations
	mockExamRepo.On("GetActiveSessionByID", sessionID).Return(nil, repository.ErrAttemptNotFound)

	// Execute
	result, err := examService.SubmitExam(sessionID, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, dto.SubmitExamResponse{}, result)
	assert.Equal(t, repository.ErrAttemptNotFound, err)

	// Verify mock expectations
	mockExamRepo.AssertExpectations(t)
}

// =============================================================================
// COMPREHENSIVE SCORING TESTS FOR SUBMIT EXAM FUNCTIONALITY
// =============================================================================

// Test SBA Question Scoring - Correct Answer with Real Structure
func TestExamService_SubmitExam_SBA_CorrectAnswer_RealStructure(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Test data - SBA question with correct structure (matching your sample data)
	questionsJSON := `[{
		"id": 2,
		"question_text": "Which liver enzyme is most specific for hepatocellular damage?",
		"question_type": "SBA",
		"options": {
			"a": {"text": "ALP (Alkaline Phosphatase)", "is_correct": false},
			"b": {"text": "AST (Aspartate Aminotransferase)", "is_correct": false},
			"c": {"text": "ALT (Alanine Aminotransferase)", "is_correct": true},
			"d": {"text": "PT (Prothrombin Time)", "is_correct": false},
			"e": {"text": "Bilirubin", "is_correct": false}
		},
		"points": 1
	}]`

	sessionData := &repository.SessionWithExamData{
		Attempt: models.UserExamAttempt{
			ID:     1,
			UserID: 1,
			ExamID: 1,
			Status: models.AttemptStatusStarted,
		},
		Exam: models.Exam{
			ID:            1,
			Title:         "Test Exam",
			PassingScore:  70.0,
			QuestionsData: questionsJSON,
		},
	}

	// User answers - correct answer
	userAnswers := map[uint]string{
		2: "c", // Correct answer (ALT)
	}

	// Setup expectations
	mockExamRepo.On("GetActiveSessionByID", "test_session").Return(sessionData, nil)
	mockExamRepo.On("GetSessionAnswers", "test_session").Return(userAnswers, nil)
	mockExamRepo.On("CompleteExamAttemptWithAnswers", uint(1), 1.0, false, mock.AnythingOfType("string"), 1).Return(nil)

	// Execute
	result, err := examService.SubmitExam("test_session", uint(1))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1.0, result.Score)
	assert.False(t, result.Passed) // 1 point < 70 passing score
	assert.Equal(t, 1, result.CorrectAnswers)

	// Verify the answers data contains correct information
	call := mockExamRepo.Calls[2] // CompleteExamAttemptWithAnswers call
	answersDataJSON := call.Arguments[3].(string)

	var answersData map[string]interface{}
	err = json.Unmarshal([]byte(answersDataJSON), &answersData)
	assert.NoError(t, err)

	answers := answersData["answers"].([]interface{})
	firstAnswer := answers[0].(map[string]interface{})

	assert.Equal(t, "c", firstAnswer["correct_answer"].(string))
	assert.True(t, firstAnswer["is_correct"].(bool))
	assert.Equal(t, float64(1), firstAnswer["points_earned"].(float64))

	mockExamRepo.AssertExpectations(t)
}

// Test TRUE_FALSE Question Scoring - Perfect Answer (matching your sample)
func TestExamService_SubmitExam_TrueFalse_PerfectAnswer_RealStructure(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Test data - TRUE_FALSE question matching your sample
	questionsJSON := `[{
		"id": 1,
		"question_text": "Which of the following drugs are known teratogens that can cause birth defects?",
		"question_type": "TRUE_FALSE",
		"options": {
			"a": {"text": "Isotretinoin", "is_correct": true},
			"b": {"text": "OCP (Oral Contraceptive Pills)", "is_correct": true},
			"c": {"text": "Valproate", "is_correct": true},
			"d": {"text": "Sulfasalazine", "is_correct": false},
			"e": {"text": "Amiodarone", "is_correct": false}
		},
		"points": 2
	}]`

	sessionData := &repository.SessionWithExamData{
		Attempt: models.UserExamAttempt{
			ID:     1,
			UserID: 1,
			ExamID: 1,
			Status: models.AttemptStatusStarted,
		},
		Exam: models.Exam{
			ID:            1,
			Title:         "Test Exam",
			PassingScore:  70.0,
			QuestionsData: questionsJSON,
		},
	}

	// User answers - perfect answer (all options marked correctly)
	userAnswers := map[uint]string{
		1: `["a:true","b:true","c:true","d:false","e:false"]`, // Perfect answer - all options marked correctly
	}

	// Setup expectations
	mockExamRepo.On("GetActiveSessionByID", "test_session").Return(sessionData, nil)
	mockExamRepo.On("GetSessionAnswers", "test_session").Return(userAnswers, nil)
	mockExamRepo.On("CompleteExamAttemptWithAnswers", uint(1), 2.0, false, mock.AnythingOfType("string"), 1).Return(nil)

	// Execute
	result, err := examService.SubmitExam("test_session", uint(1))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2.0, result.Score)
	assert.False(t, result.Passed) // 2 points < 70 passing score
	assert.Equal(t, 1, result.CorrectAnswers)

	// Verify the answers data
	call := mockExamRepo.Calls[2]
	answersDataJSON := call.Arguments[3].(string)

	var answersData map[string]interface{}
	err = json.Unmarshal([]byte(answersDataJSON), &answersData)
	assert.NoError(t, err)

	answers := answersData["answers"].([]interface{})
	firstAnswer := answers[0].(map[string]interface{})

	// For TRUE_FALSE questions, correct_answer is now a map[string]bool
	correctAnswer := firstAnswer["correct_answer"].(map[string]interface{})
	assert.Equal(t, true, correctAnswer["a"])
	assert.Equal(t, true, correctAnswer["b"])
	assert.Equal(t, true, correctAnswer["c"])
	assert.Equal(t, false, correctAnswer["d"])
	assert.Equal(t, false, correctAnswer["e"])
	assert.True(t, firstAnswer["is_correct"].(bool))
	assert.Equal(t, float64(2), firstAnswer["points_earned"].(float64))

	mockExamRepo.AssertExpectations(t)
}

// Test TRUE_FALSE Question Scoring - Partial Answer like your ACE inhibitors question
func TestExamService_SubmitExam_TrueFalse_PartialAnswer_YourSample(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Test data - TRUE_FALSE question like your ACE inhibitors
	questionsJSON := `[{
		"id": 3,
		"question_text": "Which of the following are common side effects of ACE inhibitors?",
		"question_type": "TRUE_FALSE",
		"options": {
			"a": {"text": "Dry cough", "is_correct": true},
			"b": {"text": "Hyperkalemia", "is_correct": true},
			"c": {"text": "Angioedema", "is_correct": true},
			"d": {"text": "Hyponatremia", "is_correct": false},
			"e": {"text": "Weight gain", "is_correct": false}
		},
		"points": 2
	}]`

	sessionData := &repository.SessionWithExamData{
		Attempt: models.UserExamAttempt{
			ID:     1,
			UserID: 1,
			ExamID: 1,
			Status: models.AttemptStatusStarted,
		},
		Exam: models.Exam{
			ID:            1,
			Title:         "Test Exam",
			PassingScore:  70.0,
			QuestionsData: questionsJSON,
		},
	}

	// User answers - partial answer (some options marked correctly, some missing/wrong)
	userAnswers := map[uint]string{
		3: `["a:true","b:true","c:false","d:false","e:false"]`, // 2 correct, 1 wrong (c should be true)
	}

	// Calculate expected score: (2/3) * 2 = 1.33 points
	expectedPartialPoints := (4.0 / 5.0) * 2.0 // 1.6 points (4 out of 5 correct)

	// Setup expectations
	mockExamRepo.On("GetActiveSessionByID", "test_session").Return(sessionData, nil)
	mockExamRepo.On("GetSessionAnswers", "test_session").Return(userAnswers, nil)
	mockExamRepo.On("CompleteExamAttemptWithAnswers", uint(1), expectedPartialPoints, false, mock.AnythingOfType("string"), 0).Return(nil)

	// Execute
	result, err := examService.SubmitExam("test_session", uint(1))

	// Assert
	assert.NoError(t, err)
	assert.InDelta(t, expectedPartialPoints, result.Score, 0.01)
	assert.Equal(t, 0, result.CorrectAnswers) // Not counted as fully correct

	// Verify the answers data
	call := mockExamRepo.Calls[2]
	answersDataJSON := call.Arguments[3].(string)

	var answersData map[string]interface{}
	err = json.Unmarshal([]byte(answersDataJSON), &answersData)
	assert.NoError(t, err)

	answers := answersData["answers"].([]interface{})
	firstAnswer := answers[0].(map[string]interface{})

	assert.False(t, firstAnswer["is_correct"].(bool))
	assert.InDelta(t, expectedPartialPoints, firstAnswer["points_earned"].(float64), 0.01)

	mockExamRepo.AssertExpectations(t)
}

// Test TRUE_FALSE Question Scoring - Wrong Answer with Penalty (like your DVT question)
func TestExamService_SubmitExam_TrueFalse_WrongAnswerWithPenalty_YourSample(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Test data - TRUE_FALSE question like your DVT question
	questionsJSON := `[{
		"id": 7,
		"question_text": "Which of the following are risk factors for deep vein thrombosis (DVT)?",
		"question_type": "TRUE_FALSE",
		"options": {
			"a": {"text": "Prolonged immobilization", "is_correct": true},
			"b": {"text": "Oral contraceptive use", "is_correct": true},
			"c": {"text": "Recent surgery", "is_correct": true},
			"d": {"text": "Young age (under 20)", "is_correct": false},
			"e": {"text": "Low body weight", "is_correct": false}
		},
		"points": 2
	}]`

	sessionData := &repository.SessionWithExamData{
		Attempt: models.UserExamAttempt{
			ID:     1,
			UserID: 1,
			ExamID: 1,
			Status: models.AttemptStatusStarted,
		},
		Exam: models.Exam{
			ID:            1,
			Title:         "Test Exam",
			PassingScore:  70.0,
			QuestionsData: questionsJSON,
		},
	}

	// User answers - some correct, some wrong
	userAnswers := map[uint]string{
		7: `["a:true","b:false","c:true","d:false","e:true"]`, // 2 correct (a,c), 2 wrong (b should be true, e should be false)
	}

	// Calculate expected score: (2/3) * 2 - 0.1 * 2 = 1.33 - 0.2 = 1.13 points
	expectedPartialPoints := (3.0 / 5.0) * 2.0 // 1.2 points (3 out of 5 correct)

	// Setup expectations
	mockExamRepo.On("GetActiveSessionByID", "test_session").Return(sessionData, nil)
	mockExamRepo.On("GetSessionAnswers", "test_session").Return(userAnswers, nil)
	mockExamRepo.On("CompleteExamAttemptWithAnswers", uint(1), expectedPartialPoints, false, mock.AnythingOfType("string"), 0).Return(nil)

	// Execute
	result, err := examService.SubmitExam("test_session", uint(1))

	// Assert
	assert.NoError(t, err)
	assert.InDelta(t, expectedPartialPoints, result.Score, 0.01)
	assert.Equal(t, 0, result.CorrectAnswers) // Not counted as fully correct due to wrong selection

	mockExamRepo.AssertExpectations(t)
}

// Test Your Exact Exam Scenario - Comprehensive Real-World Test
func TestExamService_SubmitExam_YourExactExamScenario(t *testing.T) {
	// Setup
	mockExamRepo := &MockExamRepository{}
	mockExamMapper := &MockExamMapper{}
	examService := service.NewExamService(mockExamRepo, &MockEnrollmentRepository{}, mockExamMapper)

	// Your exact exam data (first 5 questions for testing)
	questionsJSON := `[
		{
			"id": 1,
			"question_text": "Which of the following drugs are known teratogens that can cause birth defects?",
			"question_type": "TRUE_FALSE",
			"options": {
				"a": {"text": "Isotretinoin", "is_correct": true},
				"b": {"text": "OCP (Oral Contraceptive Pills)", "is_correct": true},
				"c": {"text": "Valproate", "is_correct": true},
				"d": {"text": "Sulfasalazine", "is_correct": false},
				"e": {"text": "Amiodarone", "is_correct": false}
			},
			"points": 2
		},
		{
			"id": 2,
			"question_text": "Which liver enzyme is most specific for hepatocellular damage?",
			"question_type": "SBA",
			"options": {
				"a": {"text": "ALP (Alkaline Phosphatase)", "is_correct": false},
				"b": {"text": "AST (Aspartate Aminotransferase)", "is_correct": false},
				"c": {"text": "ALT (Alanine Aminotransferase)", "is_correct": true},
				"d": {"text": "PT (Prothrombin Time)", "is_correct": false},
				"e": {"text": "Bilirubin", "is_correct": false}
			},
			"points": 1
		},
		{
			"id": 4,
			"question_text": "What is the most appropriate first-line treatment for Type 2 diabetes mellitus?",
			"question_type": "SBA",
			"options": {
				"a": {"text": "Insulin", "is_correct": false},
				"b": {"text": "Metformin", "is_correct": true},
				"c": {"text": "Sulfonylureas", "is_correct": false},
				"d": {"text": "DPP-4 inhibitors", "is_correct": false},
				"e": {"text": "GLP-1 agonists", "is_correct": false}
			},
			"points": 1
		},
		{
			"id": 5,
			"question_text": "Which of the following are contraindications for MRI scanning?",
			"question_type": "TRUE_FALSE",
			"options": {
				"a": {"text": "Pacemaker (non-MRI compatible)", "is_correct": true},
				"b": {"text": "Cochlear implants (older models)", "is_correct": true},
				"c": {"text": "Metallic foreign body in the eye", "is_correct": true},
				"d": {"text": "Dental fillings", "is_correct": false},
				"e": {"text": "Surgical clips (titanium)", "is_correct": false}
			},
			"points": 2
		},
		{
			"id": 6,
			"question_text": "Which antibiotic is the drug of choice for treating MRSA (Methicillin-Resistant Staphylococcus Aureus)?",
			"question_type": "SBA",
			"options": {
				"a": {"text": "Penicillin", "is_correct": false},
				"b": {"text": "Cephalexin", "is_correct": false},
				"c": {"text": "Vancomycin", "is_correct": true},
				"d": {"text": "Amoxicillin", "is_correct": false},
				"e": {"text": "Ciprofloxacin", "is_correct": false}
			},
			"points": 1
		}
	]`

	sessionData := &repository.SessionWithExamData{
		Attempt: models.UserExamAttempt{
			ID:     1,
			UserID: 1,
			ExamID: 1,
			Status: models.AttemptStatusStarted,
		},
		Exam: models.Exam{
			ID:            1,
			Title:         "Medical Knowledge Exam",
			PassingScore:  70.0,
			QuestionsData: questionsJSON,
		},
	}

	// Your exact answers but updated to new TRUE_FALSE format
	userAnswers := map[uint]string{
		1: `["a:true","b:true","c:true","d:false","e:false"]`, // Perfect TRUE_FALSE answer (2 points)
		2: "c",                                                // Correct SBA answer (1 point)
		4: "b",                                                // Correct SBA answer (1 point)
		5: `["a:true","b:true","c:true","d:false","e:false"]`, // Perfect TRUE_FALSE answer (2 points)
		6: "c",                                                // Correct SBA answer (1 point)
	}

	// Expected calculation:
	// Q1: 2/2 points (perfect TRUE_FALSE)
	// Q2: 1/1 points (correct SBA)
	// Q4: 1/1 points (correct SBA)
	// Q5: 2/2 points (perfect TRUE_FALSE)
	// Q6: 1/1 points (correct SBA)
	// Total: 7/7 points (not percentage)
	expectedScore := 7.0
	expectedCorrectAnswers := 5

	// Setup expectations
	mockExamRepo.On("GetActiveSessionByID", "test_session").Return(sessionData, nil)
	mockExamRepo.On("GetSessionAnswers", "test_session").Return(userAnswers, nil)
	mockExamRepo.On("CompleteExamAttemptWithAnswers", uint(1), expectedScore, false, mock.AnythingOfType("string"), expectedCorrectAnswers).Return(nil)

	// Execute
	result, err := examService.SubmitExam("test_session", uint(1))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedScore, result.Score)
	assert.False(t, result.Passed) // 7 points < 70 passing score
	assert.Equal(t, expectedCorrectAnswers, result.CorrectAnswers)
	assert.Equal(t, 5, result.TotalQuestions)

	// Verify the answers data contains proper structure
	call := mockExamRepo.Calls[2]
	answersDataJSON := call.Arguments[3].(string)

	var answersData map[string]interface{}
	err = json.Unmarshal([]byte(answersDataJSON), &answersData)
	assert.NoError(t, err)

	answers := answersData["answers"].([]interface{})
	assert.Len(t, answers, 5)

	// Check first answer (TRUE_FALSE) - correct_answer is now a map
	firstAnswer := answers[0].(map[string]interface{})
	assert.True(t, firstAnswer["is_correct"].(bool))
	assert.Equal(t, float64(2), firstAnswer["points_earned"].(float64))
	correctAnswers := firstAnswer["correct_answer"].(map[string]interface{})
	assert.Equal(t, true, correctAnswers["a"])
	assert.Equal(t, true, correctAnswers["b"])
	assert.Equal(t, true, correctAnswers["c"])
	assert.Equal(t, false, correctAnswers["d"])
	assert.Equal(t, false, correctAnswers["e"])

	// Check second answer (SBA)
	secondAnswer := answers[1].(map[string]interface{})
	assert.True(t, secondAnswer["is_correct"].(bool))
	assert.Equal(t, float64(1), secondAnswer["points_earned"].(float64))
	assert.Equal(t, "c", secondAnswer["correct_answer"].(string))

	mockExamRepo.AssertExpectations(t)
}
