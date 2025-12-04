package repository

import (
	"errors"
	"quizora-backend/internal/models"

	"gorm.io/gorm"
)

// PackageRepository handles database operations for packages
type PackageRepository interface {
	GetActivePackages() ([]models.Package, error)
	GetBySlugWithExams(slug string) (*models.Package, error)
}

// packageRepository implements PackageRepository
type packageRepository struct {
	db *gorm.DB
}

// NewPackageRepository creates a new package repository
func NewPackageRepository(db *gorm.DB) PackageRepository {
	return &packageRepository{db: db}
}

// GetActivePackages retrieves all active packages ordered by sort_order
func (r *packageRepository) GetActivePackages() ([]models.Package, error) {
	var packages []models.Package

	err := r.db.Where("is_active = ?", true).
		Order("sort_order ASC").
		Find(&packages).Error

	if err != nil {
		return nil, err
	}

	return packages, nil
}

// GetBySlugWithExams retrieves a package by slug with exam schedule data
func (r *packageRepository) GetBySlugWithExams(slug string) (*models.Package, error) {
	var pkg models.Package

	err := r.db.Preload("PackageExams", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = ?", true).Order("sort_order ASC")
	}).Preload("PackageExams.Exam", func(db *gorm.DB) *gorm.DB {
		// Only load basic exam fields, not questions_data
		return db.Select("id, title, slug, description, exam_type, total_questions, duration_minutes, passing_score, max_attempts, scheduled_start_date, scheduled_end_date, status, is_active, attempt_count, completed_attempt_count, average_score, pass_rate, last_attempt_at, created_at, updated_at").
			Where("is_active = ?", true)
	}).Where("slug = ? AND is_active = ?", slug, true).First(&pkg).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPackageNotFound
		}
		return nil, err
	}

	return &pkg, nil
}

// Custom errors
var (
	ErrPackageNotFound = errors.New("package not found")
)
