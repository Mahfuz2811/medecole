package errors

import "fmt"

// EnrollmentError represents base enrollment error type
type EnrollmentError struct {
	Code    string
	Message string
	Cause   error
}

func (e *EnrollmentError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *EnrollmentError) Unwrap() error {
	return e.Cause
}

// Specific enrollment error types
type PackageNotActiveError struct {
	*EnrollmentError
	PackageID uint
}

func NewPackageNotActiveError(packageID uint) *PackageNotActiveError {
	return &PackageNotActiveError{
		EnrollmentError: &EnrollmentError{
			Code:    "PACKAGE_NOT_ACTIVE",
			Message: "Package is not available for enrollment",
		},
		PackageID: packageID,
	}
}

type ActiveEnrollmentExistsError struct {
	*EnrollmentError
	UserID    uint
	PackageID uint
}

func NewActiveEnrollmentExistsError(userID, packageID uint) *ActiveEnrollmentExistsError {
	return &ActiveEnrollmentExistsError{
		EnrollmentError: &EnrollmentError{
			Code:    "ACTIVE_ENROLLMENT_EXISTS",
			Message: "You are already enrolled in this package",
		},
		UserID:    userID,
		PackageID: packageID,
	}
}

type CouponValidationError struct {
	*EnrollmentError
	CouponCode string
	Reason     string
}

func NewCouponValidationError(couponCode, reason string, cause error) *CouponValidationError {
	return &CouponValidationError{
		EnrollmentError: &EnrollmentError{
			Code:    "COUPON_VALIDATION_FAILED",
			Message: "Invalid or expired coupon",
			Cause:   cause,
		},
		CouponCode: couponCode,
		Reason:     reason,
	}
}

type PackageNotFoundError struct {
	*EnrollmentError
	PackageID uint
}

func NewPackageNotFoundError(packageID uint, cause error) *PackageNotFoundError {
	return &PackageNotFoundError{
		EnrollmentError: &EnrollmentError{
			Code:    "PACKAGE_NOT_FOUND",
			Message: "Package not found",
			Cause:   cause,
		},
		PackageID: packageID,
	}
}

type EnrollmentCreationError struct {
	*EnrollmentError
	UserID    uint
	PackageID uint
}

func NewEnrollmentCreationError(userID, packageID uint, cause error) *EnrollmentCreationError {
	return &EnrollmentCreationError{
		EnrollmentError: &EnrollmentError{
			Code:    "ENROLLMENT_CREATION_FAILED",
			Message: "Failed to create enrollment",
			Cause:   cause,
		},
		UserID:    userID,
		PackageID: packageID,
	}
}

type CouponProcessingError struct {
	*EnrollmentError
	CouponCode string
	Operation  string
}

func NewCouponProcessingError(couponCode, operation string, cause error) *CouponProcessingError {
	return &CouponProcessingError{
		EnrollmentError: &EnrollmentError{
			Code:    "COUPON_PROCESSING_FAILED",
			Message: "Failed to process coupon",
			Cause:   cause,
		},
		CouponCode: couponCode,
		Operation:  operation,
	}
}

type PackageStatsUpdateError struct {
	*EnrollmentError
	PackageID uint
}

func NewPackageStatsUpdateError(packageID uint, cause error) *PackageStatsUpdateError {
	return &PackageStatsUpdateError{
		EnrollmentError: &EnrollmentError{
			Code:    "PACKAGE_STATS_UPDATE_FAILED",
			Message: "Failed to update package statistics",
			Cause:   cause,
		},
		PackageID: packageID,
	}
}

type TransactionCommitError struct {
	*EnrollmentError
}

func NewTransactionCommitError(cause error) *TransactionCommitError {
	return &TransactionCommitError{
		EnrollmentError: &EnrollmentError{
			Code:    "TRANSACTION_COMMIT_FAILED",
			Message: "Failed to commit enrollment transaction",
			Cause:   cause,
		},
	}
}

type EnrollmentFetchError struct {
	*EnrollmentError
	EnrollmentID uint
}

func NewEnrollmentFetchError(enrollmentID uint, cause error) *EnrollmentFetchError {
	return &EnrollmentFetchError{
		EnrollmentError: &EnrollmentError{
			Code:    "ENROLLMENT_FETCH_FAILED",
			Message: "Failed to fetch enrollment details",
			Cause:   cause,
		},
		EnrollmentID: enrollmentID,
	}
}

// Helper functions to check error types
func IsPackageNotActiveError(err error) bool {
	_, ok := err.(*PackageNotActiveError)
	return ok
}

func IsActiveEnrollmentExistsError(err error) bool {
	_, ok := err.(*ActiveEnrollmentExistsError)
	return ok
}

func IsCouponValidationError(err error) bool {
	_, ok := err.(*CouponValidationError)
	return ok
}

func IsPackageNotFoundError(err error) bool {
	_, ok := err.(*PackageNotFoundError)
	return ok
}

func IsEnrollmentCreationError(err error) bool {
	_, ok := err.(*EnrollmentCreationError)
	return ok
}

func IsCouponProcessingError(err error) bool {
	_, ok := err.(*CouponProcessingError)
	return ok
}

func IsPackageStatsUpdateError(err error) bool {
	_, ok := err.(*PackageStatsUpdateError)
	return ok
}

func IsTransactionCommitError(err error) bool {
	_, ok := err.(*TransactionCommitError)
	return ok
}

func IsEnrollmentFetchError(err error) bool {
	_, ok := err.(*EnrollmentFetchError)
	return ok
}
