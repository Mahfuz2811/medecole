package service

import (
	"context"
	"fmt"
	"math"
	"quizora-backend/internal/dto"
	"quizora-backend/internal/errors"
	"quizora-backend/internal/logger"
	"quizora-backend/internal/mapper"
	"quizora-backend/internal/models"
	"quizora-backend/internal/repository"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// EnrollmentService handles enrollment business logic
type EnrollmentService interface {
	// Core enrollment operations
	EnrollInPackage(ctx context.Context, userID uint, req dto.EnrollmentRequest) (*dto.EnrollmentResponse, error)
	CheckEnrollmentStatus(ctx context.Context, userID, packageID uint) (*dto.EnrollmentStatusResponse, error)

	// Coupon operations
	ValidateCoupon(ctx context.Context, req dto.CouponValidationRequest) (*dto.CouponValidationResponse, error)
	CalculatePrice(packagePrice float64, coupon *models.Coupon) *dto.PriceCalculationResult
}

// enrollmentService implements EnrollmentService
type enrollmentService struct {
	repo   repository.EnrollmentRepository
	mapper *mapper.EnrollmentMapper
	db     *gorm.DB
}

// NewEnrollmentService creates a new enrollment service
func NewEnrollmentService(repo repository.EnrollmentRepository, mapper *mapper.EnrollmentMapper, db *gorm.DB) EnrollmentService {
	return &enrollmentService{
		repo:   repo,
		mapper: mapper,
		db:     db,
	}
}

// EnrollInPackage handles package enrollment with full business logic
func (s *enrollmentService) EnrollInPackage(ctx context.Context, userID uint, req dto.EnrollmentRequest) (*dto.EnrollmentResponse, error) {
	// Add operation context for logging
	ctx = logger.AddOperationToContext(ctx, "EnrollInPackage")
	ctx = logger.AddServiceToContext(ctx, "enrollment")

	log := logger.WithContext(ctx).WithFields(logrus.Fields{
		"user_id":    userID,
		"package_id": req.PackageID,
		"operation":  "EnrollInPackage",
		"service":    "enrollment",
	})

	log.Info("Starting package enrollment process")

	// Start transaction for data consistency
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.WithField("panic", r).Error("Panic occurred during enrollment, rolling back transaction")
			tx.Rollback()
		}
	}()

	repoTx := s.repo.WithTransaction(tx)

	// 1. Validate package exists and is active
	pkg, err := repoTx.GetPackageByID(req.PackageID)
	if err != nil {
		log.WithError(err).Error("Failed to fetch package from database")
		tx.Rollback()
		return nil, errors.NewPackageNotFoundError(req.PackageID, err)
	}

	log.WithFields(logrus.Fields{
		"package_name":   pkg.Name,
		"package_type":   pkg.PackageType,
		"package_price":  pkg.Price,
		"package_active": pkg.IsActive,
	}).Debug("Package details retrieved")

	if !pkg.IsActive {
		log.Warn("Attempted enrollment in inactive package")
		tx.Rollback()
		return nil, errors.NewPackageNotActiveError(req.PackageID)
	}

	// 2. Check for existing active enrollment (prevent duplicates)
	existingEnrollment, err := repoTx.GetActiveEnrollment(userID, req.PackageID)
	if err != nil {
		log.WithError(err).Error("Failed to check existing enrollment")
		tx.Rollback()
		return nil, fmt.Errorf("failed to check existing enrollment: %w", err)
	}

	if existingEnrollment != nil && existingEnrollment.CanAccessContent() {
		log.WithField("existing_enrollment_id", existingEnrollment.ID).Warn("User already has active enrollment for this package")
		tx.Rollback()
		return nil, errors.NewActiveEnrollmentExistsError(userID, req.PackageID)
	}

	// 3. Validate and process coupon if provided
	var coupon *models.Coupon
	if req.CouponCode != nil && *req.CouponCode != "" {
		log.WithField("coupon_code", *req.CouponCode).Info("Processing coupon validation")
		coupon, err = repoTx.GetCouponByCode(*req.CouponCode)
		if err != nil {
			log.WithError(err).WithField("coupon_code", *req.CouponCode).Error("Failed to retrieve coupon")
			tx.Rollback()
			return nil, errors.NewCouponValidationError(*req.CouponCode, "coupon not found", err)
		}

		log.WithFields(logrus.Fields{
			"coupon_id":               coupon.ID,
			"coupon_discount_percent": coupon.DiscountPercentage,
			"coupon_usage_limit":      coupon.UsageLimit,
		}).Debug("Coupon details retrieved")

		if err := repoTx.ValidateCoupon(coupon, req.PackageID); err != nil {
			log.WithError(err).WithField("coupon_code", *req.CouponCode).Error("Coupon validation failed")
			tx.Rollback()
			return nil, errors.NewCouponValidationError(*req.CouponCode, "coupon validation failed", err)
		}
		log.WithField("coupon_code", *req.CouponCode).Info("Coupon validated successfully")
	}

	// 4. Calculate pricing
	log.Debug("Calculating pricing with potential coupon discount")
	priceCalc := s.CalculatePrice(pkg.Price, coupon)

	log.WithFields(logrus.Fields{
		"original_price":      priceCalc.OriginalPrice,
		"discount_percentage": priceCalc.DiscountPercentage,
		"discount_amount":     priceCalc.DiscountAmount,
		"final_price":         priceCalc.FinalPrice,
	}).Info("Price calculation completed")

	// 5. Create enrollment record
	enrollment := s.buildEnrollment(userID, pkg, coupon, priceCalc)

	if err := repoTx.CreateEnrollment(enrollment); err != nil {
		log.WithError(err).Error("Failed to create enrollment record")
		tx.Rollback()
		return nil, errors.NewEnrollmentCreationError(userID, req.PackageID, err)
	}

	log.WithField("enrollment_id", enrollment.ID).Info("Enrollment record created successfully")

	// 6. Process coupon usage if applicable
	if coupon != nil {
		log.WithField("coupon_code", coupon.Code).Debug("Processing coupon usage tracking")
		if err := s.processCouponUsage(repoTx, coupon, enrollment, priceCalc); err != nil {
			log.WithError(err).Error("Failed to process coupon usage")
			tx.Rollback()
			return nil, errors.NewCouponProcessingError(coupon.Code, "usage tracking", err)
		}
		log.WithField("coupon_code", coupon.Code).Info("Coupon usage processed successfully")
	}

	// 7. Update package enrollment statistics
	if err := s.updatePackageStats(tx, pkg.ID); err != nil {
		log.WithError(err).Error("Failed to update package statistics")
		tx.Rollback()
		return nil, errors.NewPackageStatsUpdateError(pkg.ID, err)
	}

	// 8. Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.WithError(err).Error("Failed to commit transaction")
		return nil, errors.NewTransactionCommitError(err)
	}

	log.Info("Transaction committed successfully")

	// 9. Fetch and return complete enrollment data
	finalEnrollment, err := s.repo.GetEnrollmentByID(enrollment.ID)
	if err != nil {
		log.WithError(err).WithField("enrollment_id", enrollment.ID).Error("Failed to fetch enrollment details")
		return nil, errors.NewEnrollmentFetchError(enrollment.ID, err)
	}

	response := s.mapper.ToEnrollmentResponse(finalEnrollment)
	response.PriceCalculation = priceCalc

	log.WithFields(logrus.Fields{
		"enrollment_id":      finalEnrollment.ID,
		"final_price":        priceCalc.FinalPrice,
		"payment_status":     finalEnrollment.PaymentStatus,
		"can_access_content": finalEnrollment.CanAccessContent(),
	}).Info("Package enrollment completed successfully")

	return response, nil
}

// buildEnrollment creates an enrollment model with calculated values
func (s *enrollmentService) buildEnrollment(userID uint, pkg *models.Package, coupon *models.Coupon, priceCalc *dto.PriceCalculationResult) *models.UserPackageEnrollment {
	now := time.Now()

	enrollment := &models.UserPackageEnrollment{
		UserID:              userID,
		PackageID:           pkg.ID,
		EnrollmentType:      models.EnrollmentTypeFull, // Always FULL for now (TRIAL will be added later)
		EnrolledAt:          now,
		EnrolledPackageType: pkg.PackageType,
		EnrolledPrice:       priceCalc.FinalPrice,
		IsActive:            true,
	}

	// Set payment status based on price
	if priceCalc.FinalPrice == 0 {
		enrollment.PaymentStatus = models.PaymentStatusFree
	} else {
		enrollment.PaymentStatus = models.PaymentStatusPending // Will be updated when payment is processed
	}

	// Calculate expiration date based on validity type
	enrollment.ExpiresAt = s.calculateExpirationDate(pkg)

	// Set coupon details if used
	if coupon != nil {
		enrollment.CouponID = &coupon.ID
		enrollment.CouponCode = &coupon.Code
		enrollment.OriginalPrice = &priceCalc.OriginalPrice
		enrollment.DiscountPercentage = &priceCalc.DiscountPercentage
		enrollment.DiscountAmount = &priceCalc.DiscountAmount
		enrollment.FinalPrice = &priceCalc.FinalPrice
	}

	return enrollment
}

// calculateExpirationDate calculates when the enrollment expires
func (s *enrollmentService) calculateExpirationDate(pkg *models.Package) *time.Time {
	now := time.Now()
	var expiresAt time.Time

	switch pkg.ValidityType {
	case models.ValidityTypeFixed:
		if pkg.ValidityDate != nil {
			expiresAt = *pkg.ValidityDate
		} else {
			// Fallback: 1 year from now if no fixed date set
			expiresAt = now.AddDate(1, 0, 0)
		}
	case models.ValidityTypeRelative:
		if pkg.ValidityDays != nil {
			expiresAt = now.AddDate(0, 0, *pkg.ValidityDays)
		} else {
			// Fallback: 1 year from now if no days set
			expiresAt = now.AddDate(1, 0, 0)
		}
	default:
		// Default: 1 year from now
		expiresAt = now.AddDate(1, 0, 0)
	}

	return &expiresAt
}

// processCouponUsage handles coupon usage tracking
func (s *enrollmentService) processCouponUsage(repo repository.EnrollmentRepository, coupon *models.Coupon, enrollment *models.UserPackageEnrollment, priceCalc *dto.PriceCalculationResult) error {
	// Increment coupon usage count
	if err := repo.IncrementCouponUsage(coupon.ID); err != nil {
		return fmt.Errorf("failed to increment coupon usage: %w", err)
	}

	// Create coupon usage record for tracking
	usage := &models.CouponUsage{
		CouponID:           coupon.ID,
		UserID:             enrollment.UserID,
		EnrollmentID:       enrollment.ID,
		PackageID:          enrollment.PackageID,
		OriginalPrice:      priceCalc.OriginalPrice,
		DiscountPercentage: priceCalc.DiscountPercentage,
		DiscountAmount:     priceCalc.DiscountAmount,
		FinalPrice:         priceCalc.FinalPrice,
		CouponCode:         coupon.Code,
		UsedAt:             time.Now(),
	}

	return repo.CreateCouponUsage(usage)
}

// updatePackageStats updates package enrollment statistics
func (s *enrollmentService) updatePackageStats(tx *gorm.DB, packageID uint) error {
	return tx.Model(&models.Package{}).Where("id = ?", packageID).Updates(map[string]interface{}{
		"enrollment_count":   gorm.Expr("enrollment_count + 1"),
		"last_enrollment_at": time.Now(),
	}).Error
}

// CalculatePrice calculates final price with coupon discount
func (s *enrollmentService) CalculatePrice(packagePrice float64, coupon *models.Coupon) *dto.PriceCalculationResult {
	result := &dto.PriceCalculationResult{
		OriginalPrice:      packagePrice,
		DiscountPercentage: 0,
		DiscountAmount:     0,
		FinalPrice:         packagePrice,
	}

	if coupon != nil {
		result.CouponCode = &coupon.Code
		result.DiscountPercentage = coupon.DiscountPercentage
		result.DiscountAmount = math.Round((packagePrice*coupon.DiscountPercentage/100)*100) / 100
		result.FinalPrice = math.Max(0, packagePrice-result.DiscountAmount)
	}

	return result
}

// ValidateCoupon validates a coupon for a specific package
func (s *enrollmentService) ValidateCoupon(ctx context.Context, req dto.CouponValidationRequest) (*dto.CouponValidationResponse, error) {
	ctx = logger.AddOperationToContext(ctx, "ValidateCoupon")
	log := logger.WithContext(ctx).WithFields(logrus.Fields{
		"coupon_code": req.CouponCode,
		"package_id":  req.PackageID,
		"operation":   "ValidateCoupon",
	})

	log.Debug("Starting coupon validation")

	// Get package to validate pricing
	pkg, err := s.repo.GetPackageByID(req.PackageID)
	if err != nil {
		log.WithError(err).Error("Failed to retrieve package for coupon validation")
		return &dto.CouponValidationResponse{
			Valid:   false,
			Message: "Invalid package",
		}, nil
	}

	// Get and validate coupon
	coupon, err := s.repo.GetCouponByCode(req.CouponCode)
	if err != nil {
		log.WithError(err).Warn("Coupon not found during validation")
		return &dto.CouponValidationResponse{
			Valid:      false,
			CouponCode: req.CouponCode,
			Message:    "Coupon not found",
		}, nil
	}

	// Validate coupon
	if err := s.repo.ValidateCoupon(coupon, req.PackageID); err != nil {
		var message string
		switch err {
		case repository.ErrCouponExpired:
			message = "Coupon has expired"
		case repository.ErrCouponExhausted:
			message = "Coupon usage limit exceeded"
		case repository.ErrCouponInvalid:
			message = "Coupon is not valid"
		default:
			message = "Coupon validation failed"
		}

		return &dto.CouponValidationResponse{
			Valid:      false,
			CouponCode: req.CouponCode,
			Message:    message,
		}, nil
	}

	// Calculate price with coupon
	priceCalc := s.CalculatePrice(pkg.Price, coupon)

	return &dto.CouponValidationResponse{
		Valid:              true,
		CouponCode:         req.CouponCode,
		DiscountPercentage: coupon.DiscountPercentage,
		Message:            fmt.Sprintf("Coupon applied! Save ৳%.2f", priceCalc.DiscountAmount),
		PriceCalculation:   priceCalc,
	}, nil
}

// CheckEnrollmentStatus checks if user has active enrollment for package
func (s *enrollmentService) CheckEnrollmentStatus(ctx context.Context, userID, packageID uint) (*dto.EnrollmentStatusResponse, error) {
	ctx = logger.AddOperationToContext(ctx, "CheckEnrollmentStatus")
	log := logger.WithContext(ctx).WithFields(logrus.Fields{
		"user_id":    userID,
		"package_id": packageID,
		"operation":  "CheckEnrollmentStatus",
	})

	enrollment, err := s.repo.GetActiveEnrollment(userID, packageID)
	if err != nil {
		log.WithError(err).Error("Failed to check enrollment status")
		return nil, fmt.Errorf("failed to check enrollment status: %w", err)
	}

	hasActiveEnrollment := enrollment != nil && enrollment.CanAccessContent()

	response := &dto.EnrollmentStatusResponse{
		HasActiveEnrollment: hasActiveEnrollment,
	}

	if response.HasActiveEnrollment {
		response.Enrollment = s.mapper.ToEnrollmentResponse(enrollment)
	}

	return response, nil
}
