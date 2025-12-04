package unit

import (
	"errors"
	"quizora-backend/internal/models"
	"quizora-backend/internal/service"
	"quizora-backend/internal/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test GetDashboardSummary - Success
func TestDashboardService_GetDashboardSummary_Success(t *testing.T) {
	// Setup
	mockEnrollmentRepo := &MockDashboardEnrollmentRepository{}
	mockAttemptRepo := &MockUserExamAttemptRepository{}
	dashboardService := service.NewDashboardService(nil, mockEnrollmentRepo, mockAttemptRepo)

	userID := uint(1)
	testStats := createDashboardTestUserStats()
	testAttempts := createDashboardTestExamAttempts()

	// Setup expectations
	mockAttemptRepo.On("GetUserStats", userID).Return(testStats, nil)
	mockAttemptRepo.On("GetRecentActivity", userID, 10).Return(testAttempts, nil)

	// Execute
	result, err := dashboardService.GetDashboardSummary(userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Check user stats
	assert.Equal(t, 15, result.UserStats.TotalAttempts)
	assert.Equal(t, 120, result.UserStats.CorrectAnswers)
	assert.Equal(t, 80.0, result.UserStats.AccuracyRate)

	// Check recent activity
	assert.Len(t, result.RecentActivity, 2)
	assert.Equal(t, "JavaScript Basics Quiz", result.RecentActivity[0].ExamTitle)
	assert.Equal(t, 85.0, result.RecentActivity[0].Score)
	assert.Equal(t, 20, result.RecentActivity[0].TotalQuestions)
	assert.Equal(t, 17, result.RecentActivity[0].CorrectAnswers)
	assert.Equal(t, "30 min", result.RecentActivity[0].TimeTaken)

	mockAttemptRepo.AssertExpectations(t)
}

// Test GetDashboardSummary - GetUserStats Error
func TestDashboardService_GetDashboardSummary_GetUserStatsError(t *testing.T) {
	// Setup
	mockEnrollmentRepo := &MockDashboardEnrollmentRepository{}
	mockAttemptRepo := &MockUserExamAttemptRepository{}
	dashboardService := service.NewDashboardService(nil, mockEnrollmentRepo, mockAttemptRepo)

	userID := uint(1)

	// Setup expectations
	mockAttemptRepo.On("GetUserStats", userID).Return(nil, errors.New("database error"))

	// Execute
	result, err := dashboardService.GetDashboardSummary(userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get user stats")

	mockAttemptRepo.AssertExpectations(t)
}

// Test GetDashboardSummary - GetRecentActivity Error
func TestDashboardService_GetDashboardSummary_GetRecentActivityError(t *testing.T) {
	// Setup
	mockEnrollmentRepo := &MockDashboardEnrollmentRepository{}
	mockAttemptRepo := &MockUserExamAttemptRepository{}
	dashboardService := service.NewDashboardService(nil, mockEnrollmentRepo, mockAttemptRepo)

	userID := uint(1)
	testStats := createDashboardTestUserStats()

	// Setup expectations
	mockAttemptRepo.On("GetUserStats", userID).Return(testStats, nil)
	mockAttemptRepo.On("GetRecentActivity", userID, 10).Return([]models.UserExamAttempt{}, errors.New("database error"))

	// Execute
	result, err := dashboardService.GetDashboardSummary(userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get recent activity")

	mockAttemptRepo.AssertExpectations(t)
}

// Test accuracy rate calculation edge cases
func TestDashboardService_AccuracyRateCalculation_EdgeCases(t *testing.T) {
	// Setup
	mockEnrollmentRepo := &MockDashboardEnrollmentRepository{}
	mockAttemptRepo := &MockUserExamAttemptRepository{}
	dashboardService := service.NewDashboardService(nil, mockEnrollmentRepo, mockAttemptRepo)

	userID := uint(1)

	// Test case 1: Zero total questions
	testStats1 := &types.UserStatsData{
		TotalAttempts:  5,
		CorrectAnswers: 10,
		TotalQuestions: 0, // Zero total questions
		TotalTimeSpent: 1800,
		AverageScore:   0.0,
	}

	mockAttemptRepo.On("GetUserStats", userID).Return(testStats1, nil).Once()
	mockAttemptRepo.On("GetRecentActivity", userID, 10).Return([]models.UserExamAttempt{}, nil).Once()

	result1, err1 := dashboardService.GetDashboardSummary(userID)

	// Assert
	assert.NoError(t, err1)
	assert.Equal(t, 0.0, result1.UserStats.AccuracyRate) // Should be 0% when no questions

	// Test case 2: Perfect score
	testStats2 := &types.UserStatsData{
		TotalAttempts:  3,
		CorrectAnswers: 50,
		TotalQuestions: 50, // All correct
		TotalTimeSpent: 3600,
		AverageScore:   100.0,
	}

	mockAttemptRepo.On("GetUserStats", userID).Return(testStats2, nil).Once()
	mockAttemptRepo.On("GetRecentActivity", userID, 10).Return([]models.UserExamAttempt{}, nil).Once()

	result2, err2 := dashboardService.GetDashboardSummary(userID)

	// Assert
	assert.NoError(t, err2)
	assert.Equal(t, 100.0, result2.UserStats.AccuracyRate) // Should be 100%

	mockAttemptRepo.AssertExpectations(t)
}

// Test relative time calculation
func TestDashboardService_RelativeTimeCalculation(t *testing.T) {
	// Setup
	mockEnrollmentRepo := &MockDashboardEnrollmentRepository{}
	mockAttemptRepo := &MockUserExamAttemptRepository{}
	dashboardService := service.NewDashboardService(nil, mockEnrollmentRepo, mockAttemptRepo)

	userID := uint(1)
	now := time.Now()

	// Create test attempts with different time stamps
	testAttempts := []models.UserExamAttempt{
		{
			ID:              1,
			UserID:          1,
			ExamID:          1,
			Status:          models.AttemptStatusCompleted,
			StartedAt:       now.Add(-30 * time.Minute), // 30 minutes ago
			ActualTimeSpent: 1800,
			TotalQuestions:  20,
			Score:           ptrDashboardFloat64(85.0),
			CorrectAnswers:  &[]int{17}[0],
			Exam: models.Exam{
				ID:    1,
				Title: "Recent Quiz",
			},
		},
		{
			ID:              2,
			UserID:          1,
			ExamID:          2,
			Status:          models.AttemptStatusCompleted,
			StartedAt:       now.Add(-2 * time.Hour), // 2 hours ago
			ActualTimeSpent: 2400,
			TotalQuestions:  25,
			Score:           ptrDashboardFloat64(92.0),
			CorrectAnswers:  &[]int{23}[0],
			Exam: models.Exam{
				ID:    2,
				Title: "Hours Ago Quiz",
			},
		},
		{
			ID:              3,
			UserID:          1,
			ExamID:          3,
			Status:          models.AttemptStatusCompleted,
			StartedAt:       now.Add(-25 * time.Hour), // Yesterday
			ActualTimeSpent: 3000,
			TotalQuestions:  30,
			Score:           ptrDashboardFloat64(78.0),
			CorrectAnswers:  &[]int{23}[0],
			Exam: models.Exam{
				ID:    3,
				Title: "Yesterday Quiz",
			},
		},
	}

	testStats := createDashboardTestUserStats()

	// Setup expectations
	mockAttemptRepo.On("GetUserStats", userID).Return(testStats, nil)
	mockAttemptRepo.On("GetRecentActivity", userID, 10).Return(testAttempts, nil)

	// Execute
	result, err := dashboardService.GetDashboardSummary(userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.RecentActivity, 3)

	// Check relative time formatting (these are approximations)
	activity1 := result.RecentActivity[0]
	assert.Contains(t, activity1.Date, "30 minutes ago")

	activity2 := result.RecentActivity[1]
	assert.Contains(t, activity2.Date, "2 hours ago")

	activity3 := result.RecentActivity[2]
	assert.Contains(t, activity3.Date, "Yesterday")

	mockAttemptRepo.AssertExpectations(t)
}

// Test interface compliance
func TestDashboardService_ImplementsInterface(t *testing.T) {
	// Setup
	mockEnrollmentRepo := &MockDashboardEnrollmentRepository{}
	mockAttemptRepo := &MockUserExamAttemptRepository{}
	dashboardService := service.NewDashboardService(nil, mockEnrollmentRepo, mockAttemptRepo)

	// Verify that our service implements the expected interface
	var _ service.DashboardService = dashboardService
}

// Test time taken formatting edge cases
func TestDashboardService_TimeTakenFormatting_EdgeCases(t *testing.T) {
	// Setup
	mockEnrollmentRepo := &MockDashboardEnrollmentRepository{}
	mockAttemptRepo := &MockUserExamAttemptRepository{}
	dashboardService := service.NewDashboardService(nil, mockEnrollmentRepo, mockAttemptRepo)

	userID := uint(1)
	now := time.Now()

	// Create test attempts with different time durations
	testAttempts := []models.UserExamAttempt{
		{
			ID:              1,
			UserID:          1,
			ExamID:          1,
			Status:          models.AttemptStatusCompleted,
			StartedAt:       now.Add(-1 * time.Hour),
			ActualTimeSpent: 59, // 59 seconds (less than 1 minute)
			TotalQuestions:  10,
			Score:           ptrDashboardFloat64(80.0),
			CorrectAnswers:  &[]int{8}[0],
			Exam: models.Exam{
				ID:    1,
				Title: "Quick Quiz",
			},
		},
		{
			ID:              2,
			UserID:          1,
			ExamID:          2,
			Status:          models.AttemptStatusCompleted,
			StartedAt:       now.Add(-2 * time.Hour),
			ActualTimeSpent: 3661, // 1 hour and 1 minute
			TotalQuestions:  50,
			Score:           ptrDashboardFloat64(90.0),
			CorrectAnswers:  &[]int{45}[0],
			Exam: models.Exam{
				ID:    2,
				Title: "Long Quiz",
			},
		},
	}

	testStats := createDashboardTestUserStats()

	// Setup expectations
	mockAttemptRepo.On("GetUserStats", userID).Return(testStats, nil)
	mockAttemptRepo.On("GetRecentActivity", userID, 10).Return(testAttempts, nil)

	// Execute
	result, err := dashboardService.GetDashboardSummary(userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.RecentActivity, 2)

	// Check time formatting (59 seconds should be 0 min, 3661 seconds should be 61 min)
	assert.Equal(t, "0 min", result.RecentActivity[0].TimeTaken)
	assert.Equal(t, "61 min", result.RecentActivity[1].TimeTaken)

	mockAttemptRepo.AssertExpectations(t)
}