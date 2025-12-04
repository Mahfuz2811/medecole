package models

import (
	"time"

	"gorm.io/gorm"
)

// CouponStatus enum for coupon states
type CouponStatus string

const (
	CouponStatusActive    CouponStatus = "ACTIVE"
	CouponStatusInactive  CouponStatus = "INACTIVE"
	CouponStatusExpired   CouponStatus = "EXPIRED"
	CouponStatusExhausted CouponStatus = "EXHAUSTED" // Usage limit reached
)

// Coupon represents the coupons table (simplified)
type Coupon struct {
	ID          uint    `json:"id" gorm:"primarykey"`
	Code        string  `json:"code" gorm:"size:50;uniqueIndex;not null;comment:'Unique coupon code'"`
	Name        string  `json:"name" gorm:"size:200;not null;comment:'Display name for admin'"`
	Description *string `json:"description" gorm:"type:text;comment:'Coupon description'"`

	// Percentage Discount Only
	DiscountPercentage float64 `json:"discount_percentage" gorm:"type:decimal(5,2);not null;comment:'Discount percentage (0-100)'"`

	// Usage Limits
	UsageLimit *int `json:"usage_limit" gorm:"comment:'Total usage limit (null = unlimited)'"`
	UsageCount int  `json:"usage_count" gorm:"default:0;index:idx_usage_count;comment:'How many times used'"`

	// Validity Period
	ValidFrom  time.Time  `json:"valid_from" gorm:"index:idx_validity"`
	ValidUntil *time.Time `json:"valid_until" gorm:"index:idx_validity"`

	// Status
	Status   CouponStatus `json:"status" gorm:"type:enum('ACTIVE','INACTIVE','EXPIRED','EXHAUSTED');default:'ACTIVE';index:idx_status"`
	IsActive bool         `json:"is_active" gorm:"default:true;index:idx_active"`

	CreatedBy uint           `json:"created_by" gorm:"not null;index:idx_created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName specifies the table name for Coupon
func (Coupon) TableName() string {
	return "coupons"
}

// IsValid checks if coupon is currently valid and usable
func (c *Coupon) IsValid() bool {
	now := time.Now()

	// Check basic status
	if !c.IsActive || c.Status != CouponStatusActive {
		return false
	}

	// Check validity period
	if now.Before(c.ValidFrom) {
		return false
	}

	if c.ValidUntil != nil && now.After(*c.ValidUntil) {
		return false
	}

	// Check usage limits
	if c.UsageLimit != nil && c.UsageCount >= *c.UsageLimit {
		return false
	}

	return true
}

// CalculateDiscount calculates the discount amount for a given package price
func (c *Coupon) CalculateDiscount(packagePrice float64) float64 {
	if !c.IsValid() {
		return 0.0
	}

	return (c.DiscountPercentage / 100.0) * packagePrice
}

// IncrementUsage increments the usage count
func (c *Coupon) IncrementUsage(db *gorm.DB) error {
	return db.Model(c).UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1)).Error
}
