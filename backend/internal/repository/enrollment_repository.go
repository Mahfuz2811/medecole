package repository

import (
	"errors"
	"fmt"
	"github.com/Mahfuz2811/medecole/backend/internal/database"
	"github.com/Mahfuz2811/medecole/backend/internal/models"
	"time"

	"gorm.io/gorm"
)

var (
	ErrEnrollmentNotFound     = errors.New("enrollment not found")
	ErrActiveEnrollmentExists = errors.New("active enrollment already exists")
	ErrCouponNotFound         = errors.New("coupon not found")
	ErrCouponInvalid          = errors.New("coupon is invalid")
	ErrCouponExpired          = errors.New("coupon has expired")
	ErrCouponExhausted        = errors.New("coupon usage limit exceeded")
)

// EnrollmentRepository handles enrollment data operations
type EnrollmentRepository interface {
	// Enrollment CRUD
	CreateEnrollment(enrollment *models.UserPackageEnrollment) error
	GetEnrollmentByID(id uint) (*models.UserPackageEnrollment, error)
	GetUserEnrollments(userID uint) ([]models.UserPackageEnrollment, error)
	GetActiveEnrollment(userID, packageID uint) (*models.UserPackageEnrollment, error)

	// Package operations
	GetPackageByID(packageID uint) (*models.Package, error)

	// Coupon operations
	GetCouponByCode(code string) (*models.Coupon, error)
	ValidateCoupon(coupon *models.Coupon, packageID uint) error
	IncrementCouponUsage(couponID uint) error
	CreateCouponUsage(usage *models.CouponUsage) error

	// Transaction support
	WithTransaction(tx *gorm.DB) EnrollmentRepository

	// Phase 2: Enrollment validation
	IsUserEnrolledInPackage(userID, packageID uint) (bool, error)
}

// enrollmentRepository implements EnrollmentRepository
type enrollmentRepository struct {
	db *database.Database
	tx *gorm.DB // For transaction support
}

// NewEnrollmentRepository creates a new enrollment repository
func NewEnrollmentRepository(db *database.Database) EnrollmentRepository {
	return &enrollmentRepository{
		db: db,
	}
}

// WithTransaction creates a repository instance with transaction
func (r *enrollmentRepository) WithTransaction(tx *gorm.DB) EnrollmentRepository {
	return &enrollmentRepository{
		db: r.db,
		tx: tx,
	}
}

// getDB returns the appropriate database connection
func (r *enrollmentRepository) getDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db.DB
}

// CreateEnrollment creates a new enrollment
func (r *enrollmentRepository) CreateEnrollment(enrollment *models.UserPackageEnrollment) error {
	return r.getDB().Create(enrollment).Error
}

// GetEnrollmentByID retrieves enrollment by ID
func (r *enrollmentRepository) GetEnrollmentByID(id uint) (*models.UserPackageEnrollment, error) {
	var enrollment models.UserPackageEnrollment
	err := r.getDB().Preload("Package").Preload("Coupon").First(&enrollment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEnrollmentNotFound
		}
		return nil, err
	}
	return &enrollment, nil
}

// GetUserEnrollments retrieves all enrollments for a user
func (r *enrollmentRepository) GetUserEnrollments(userID uint) ([]models.UserPackageEnrollment, error) {
	var enrollments []models.UserPackageEnrollment
	err := r.getDB().Where("user_id = ?", userID).
		Preload("Package").
		Preload("Coupon").
		Order("created_at DESC").
		Find(&enrollments).Error
	return enrollments, err
}

// GetActiveEnrollment checks if user has active enrollment for package
func (r *enrollmentRepository) GetActiveEnrollment(userID, packageID uint) (*models.UserPackageEnrollment, error) {
	var enrollment models.UserPackageEnrollment
	err := r.getDB().Where("user_id = ? AND package_id = ? AND is_active = ?", userID, packageID, true).
		First(&enrollment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No active enrollment found (not an error)
		}
		return nil, err
	}
	return &enrollment, nil
}

// GetPackageByID retrieves package by ID
func (r *enrollmentRepository) GetPackageByID(packageID uint) (*models.Package, error) {
	var pkg models.Package
	err := r.getDB().First(&pkg, packageID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPackageNotFound
		}
		return nil, err
	}
	return &pkg, nil
}

// GetCouponByCode retrieves coupon by code
func (r *enrollmentRepository) GetCouponByCode(code string) (*models.Coupon, error) {
	var coupon models.Coupon
	err := r.getDB().Where("code = ?", code).First(&coupon).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCouponNotFound
		}
		return nil, err
	}
	return &coupon, nil
}

// ValidateCoupon validates if coupon can be used
func (r *enrollmentRepository) ValidateCoupon(coupon *models.Coupon, packageID uint) error {
	now := time.Now()

	// Check if coupon is active
	if !coupon.IsActive || coupon.Status != models.CouponStatusActive {
		return ErrCouponInvalid
	}

	// Check validity period
	if coupon.ValidFrom.After(now) {
		return ErrCouponInvalid
	}

	if coupon.ValidUntil != nil && coupon.ValidUntil.Before(now) {
		return ErrCouponExpired
	}

	// Check usage limit
	if coupon.UsageLimit != nil && coupon.UsageCount >= *coupon.UsageLimit {
		return ErrCouponExhausted
	}

	return nil
}

// IncrementCouponUsage increments the coupon usage count
func (r *enrollmentRepository) IncrementCouponUsage(couponID uint) error {
	return r.getDB().Model(&models.Coupon{}).Where("id = ?", couponID).
		Update("usage_count", gorm.Expr("usage_count + 1")).Error
}

// CreateCouponUsage creates a coupon usage record
func (r *enrollmentRepository) CreateCouponUsage(usage *models.CouponUsage) error {
	return r.getDB().Create(usage).Error
}

// IsUserEnrolledInPackage checks if user has active enrollment in package
// Phase 2: Optimized enrollment validation for exam access
func (r *enrollmentRepository) IsUserEnrolledInPackage(userID, packageID uint) (bool, error) {
	var enrollment models.UserPackageEnrollment
	err := r.getDB().Where("user_id = ? AND package_id = ? AND is_active = true", userID, packageID).
		First(&enrollment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // No enrollment found
		}
		return false, fmt.Errorf("failed to check enrollment status: %w", err)
	}

	// Use the model's business logic to check if user can access content
	return enrollment.CanAccessContent(), nil
}
