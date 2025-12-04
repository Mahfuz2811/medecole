package models

import (
	"time"

	"gorm.io/gorm"
)

// PaymentStatus enum for different payment states
type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "PENDING"
	PaymentStatusPaid     PaymentStatus = "PAID"
	PaymentStatusFailed   PaymentStatus = "FAILED"
	PaymentStatusRefunded PaymentStatus = "REFUNDED"
	PaymentStatusFree     PaymentStatus = "FREE"     // For free enrollments
	PaymentStatusExpired  PaymentStatus = "EXPIRED"  // Free trial expired, payment required
	PaymentStatusUpgraded PaymentStatus = "UPGRADED" // Upgraded from free to paid
)

// EnrollmentType enum for different enrollment types
type EnrollmentType string

const (
	EnrollmentTypeTrial   EnrollmentType = "TRIAL"   // Free trial enrollment
	EnrollmentTypeFull    EnrollmentType = "FULL"    // Full paid enrollment
	EnrollmentTypeUpgrade EnrollmentType = "UPGRADE" // Upgrade from trial to paid
)

// UserPackageEnrollment represents user enrollment in packages
type UserPackageEnrollment struct {
	ID        uint `json:"id" gorm:"primarykey"`
	UserID    uint `json:"user_id" gorm:"not null;index:idx_user_id"`
	PackageID uint `json:"package_id" gorm:"not null;index:idx_package_id"`

	// Enrollment Details
	EnrollmentType EnrollmentType `json:"enrollment_type" gorm:"type:enum('TRIAL','FULL','UPGRADE');default:'FULL';index:idx_enrollment_type"`
	EnrolledAt     time.Time      `json:"enrolled_at"`
	ExpiresAt      *time.Time     `json:"expires_at" gorm:"index:idx_expires_at"`

	// Trial Management
	IsTrialUsed     bool       `json:"is_trial_used" gorm:"default:false;comment:'Whether user has used trial for this package'"`
	TrialExpiresAt  *time.Time `json:"trial_expires_at" gorm:"comment:'When trial access expires'"`
	TrialExtendedAt *time.Time `json:"trial_extended_at" gorm:"comment:'If trial was extended'"`

	// Package State at Enrollment (snapshot for pricing transitions)
	EnrolledPackageType PackageType `json:"enrolled_package_type" gorm:"type:enum('FREE','PREMIUM');not null;comment:'Package type when user enrolled'"`
	EnrolledPrice       float64     `json:"enrolled_price" gorm:"type:decimal(10,2);default:0.00;comment:'Price when user enrolled'"`

	// Payment Details
	PaymentStatus    PaymentStatus `json:"payment_status" gorm:"type:enum('PENDING','PAID','FAILED','REFUNDED','FREE','EXPIRED','UPGRADED');default:'FREE';index:idx_payment_status"`
	PaymentAmount    *float64      `json:"payment_amount" gorm:"type:decimal(10,2);comment:'Actual amount paid'"`
	PaymentReference *string       `json:"payment_reference" gorm:"size:100;comment:'Payment gateway reference'"`
	PaymentDate      *time.Time    `json:"payment_date" gorm:"comment:'When payment was completed'"`

	// Coupon Details (snapshot at enrollment)
	CouponID           *uint    `json:"coupon_id" gorm:"index:idx_coupon_id;comment:'Coupon used for this enrollment'"`
	CouponCode         *string  `json:"coupon_code" gorm:"size:50;comment:'Coupon code snapshot'"`
	OriginalPrice      *float64 `json:"original_price" gorm:"type:decimal(10,2);comment:'Price before coupon discount'"`
	DiscountPercentage *float64 `json:"discount_percentage" gorm:"type:decimal(5,2);comment:'Discount percentage applied'"`
	DiscountAmount     *float64 `json:"discount_amount" gorm:"type:decimal(10,2);comment:'Discount amount applied'"`
	FinalPrice         *float64 `json:"final_price" gorm:"type:decimal(10,2);comment:'Final price after discount'"`

	// Status
	IsActive bool `json:"is_active" gorm:"default:true;index:idx_active"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Package Package `json:"package,omitempty" gorm:"foreignKey:PackageID"`
	Coupon  *Coupon `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`
}

// TableName specifies the table name for UserPackageEnrollment
func (UserPackageEnrollment) TableName() string {
	return "user_package_enrollments"
}

// IsFreePurchase checks if this enrollment was free
func (u *UserPackageEnrollment) IsFreePurchase() bool {
	return u.PaymentStatus == PaymentStatusFree || u.EnrolledPrice == 0.00
}

// IsPaidPurchase checks if this enrollment required payment
func (u *UserPackageEnrollment) IsPaidPurchase() bool {
	return u.EnrolledPrice > 0.00 && u.PaymentStatus == PaymentStatusPaid
}

// IsExpired checks if the enrollment has expired
func (u *UserPackageEnrollment) IsExpired() bool {
	if u.ExpiresAt == nil {
		return false // No expiry = never expires
	}
	return time.Now().After(*u.ExpiresAt)
}

// IsTrialExpired checks if the trial period has expired
func (u *UserPackageEnrollment) IsTrialExpired() bool {
	if u.TrialExpiresAt == nil {
		return false // No trial expiry set
	}
	return time.Now().After(*u.TrialExpiresAt)
}

// CanAccessContent checks if user can currently access package content
func (u *UserPackageEnrollment) CanAccessContent() bool {
	if !u.IsActive {
		return false
	}

	// If it's a paid enrollment, check main expiry
	if u.IsPaidPurchase() {
		return !u.IsExpired()
	}

	// If it's a trial enrollment, check trial expiry
	if u.EnrollmentType == EnrollmentTypeTrial {
		return !u.IsTrialExpired()
	}

	// For other free enrollments, check main expiry
	return !u.IsExpired()
}

// NeedsPaymentToAccess checks if user needs to pay to access content
func (u *UserPackageEnrollment) NeedsPaymentToAccess() bool {
	return u.IsTrialExpired() && u.PaymentStatus != PaymentStatusPaid
}

// GetEffectiveStatus returns the current effective status considering expiry
func (u *UserPackageEnrollment) GetEffectiveStatus() string {
	if !u.IsActive {
		return "INACTIVE"
	}

	if u.EnrollmentType == EnrollmentTypeTrial {
		if u.IsTrialExpired() {
			return "TRIAL_EXPIRED_PAYMENT_REQUIRED"
		}
		return "TRIAL_ACTIVE"
	}

	if u.IsExpired() {
		return "EXPIRED"
	}

	if u.PaymentStatus == PaymentStatusPending {
		return "PENDING_PAYMENT"
	}

	if u.IsPaidPurchase() {
		return "PAID_ACTIVE"
	}

	return "ACTIVE"
}
