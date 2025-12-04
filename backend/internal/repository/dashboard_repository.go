package repository

import (
	"quizora-backend/internal/models"
	"quizora-backend/internal/types"

	"gorm.io/gorm"
)

// UserExamAttemptRepository interface for exam attempt operations
type UserExamAttemptRepository interface {
	GetUserStats(userID uint) (*types.UserStatsData, error)
	GetRecentActivity(userID uint, limit int) ([]models.UserExamAttempt, error)
	GetUserAttemptsWithFilters(userID uint, filters types.AttemptFilters) ([]models.UserExamAttempt, error)
}

// userExamAttemptRepository implements UserExamAttemptRepository
type userExamAttemptRepository struct {
	db *gorm.DB
}

// NewUserExamAttemptRepository creates a new user exam attempt repository
func NewUserExamAttemptRepository(db *gorm.DB) UserExamAttemptRepository {
	return &userExamAttemptRepository{
		db: db,
	}
}

// GetUserStats retrieves aggregated statistics for a user
func (r *userExamAttemptRepository) GetUserStats(userID uint) (*types.UserStatsData, error) {
	var stats types.UserStatsData

	err := r.db.Model(&models.UserExamAttempt{}).
		Select(`
			COUNT(*) as total_attempts,
			COALESCE(SUM(correct_answers), 0) as correct_answers,
			COALESCE(SUM(total_questions), 0) as total_questions
		`).
		Where("user_id = ?", userID).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetRecentActivity retrieves recent exam activities
func (r *userExamAttemptRepository) GetRecentActivity(userID uint, limit int) ([]models.UserExamAttempt, error) {
	var attempts []models.UserExamAttempt
	err := r.db.Where("user_id = ?", userID).
		Preload("Exam").
		Order("started_at DESC").
		Limit(limit).
		Find(&attempts).Error
	return attempts, err
}

// GetUserAttemptsWithFilters retrieves filtered attempts for a user
func (r *userExamAttemptRepository) GetUserAttemptsWithFilters(userID uint, filters types.AttemptFilters) ([]models.UserExamAttempt, error) {
	query := r.db.Where("user_id = ?", userID)

	// Apply filters based on the AttemptFilters struct
	if filters.PackageID != nil {
		query = query.Joins("JOIN exams ON user_exam_attempts.exam_id = exams.id").
			Joins("JOIN package_exams ON exams.id = package_exams.exam_id").
			Where("package_exams.package_id = ?", *filters.PackageID)
	}

	if filters.StartDate != nil {
		query = query.Where("started_at >= ?", *filters.StartDate)
	}

	if filters.EndDate != nil {
		query = query.Where("started_at <= ?", *filters.EndDate)
	}

	var attempts []models.UserExamAttempt
	err := query.Preload("Exam").Order("started_at DESC").Find(&attempts).Error
	return attempts, err
}
