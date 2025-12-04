package models

import (
	"time"

	"gorm.io/gorm"
)

// CouponUsage represents the coupon_usages table (simplified)
type CouponUsage struct {
	ID       uint `json:"id" gorm:"primarykey"`
	CouponID uint `json:"coupon_id" gorm:"not null;index:idx_coupon_id"`
	UserID   uint `json:"user_id" gorm:"not null;index:idx_user_id"`

	// Enrollment Reference
	EnrollmentID uint `json:"enrollment_id" gorm:"not null;index:idx_enrollment_id"`
	PackageID    uint `json:"package_id" gorm:"not null;index:idx_package_id"`

	// Discount Applied
	OriginalPrice      float64 `json:"original_price" gorm:"type:decimal(10,2);not null"`
	DiscountPercentage float64 `json:"discount_percentage" gorm:"type:decimal(5,2);not null"`
	DiscountAmount     float64 `json:"discount_amount" gorm:"type:decimal(10,2);not null"`
	FinalPrice         float64 `json:"final_price" gorm:"type:decimal(10,2);not null"`

	// Usage Details
	CouponCode string    `json:"coupon_code" gorm:"size:50;not null;comment:'Snapshot of coupon code'"`
	UsedAt     time.Time `json:"used_at" gorm:"index:idx_used_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Coupon     Coupon                `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`
	User       User                  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Package    Package               `json:"package,omitempty" gorm:"foreignKey:PackageID"`
	Enrollment UserPackageEnrollment `json:"enrollment,omitempty" gorm:"foreignKey:EnrollmentID"`
}

// TableName specifies the table name for CouponUsage
func (CouponUsage) TableName() string {
	return "coupon_usages"
}
