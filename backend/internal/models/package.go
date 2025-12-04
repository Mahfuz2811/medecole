package models

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// PackageType enum for different package types
type PackageType string

const (
	PackageTypeFree    PackageType = "FREE"
	PackageTypePremium PackageType = "PREMIUM"
)

// ValidityType enum for different validity approaches
type ValidityType string

const (
	ValidityTypeFixed    ValidityType = "FIXED"    // Expires on a specific date
	ValidityTypeRelative ValidityType = "RELATIVE" // Days from enrollment date
)

// Package represents the packages table
type Package struct {
	ID          uint        `json:"id" gorm:"primarykey"`
	Name        string      `json:"name" gorm:"size:200;not null"`
	Slug        string      `json:"slug" gorm:"size:200;uniqueIndex;not null"`
	Description *string     `json:"description" gorm:"type:text"`
	PackageType PackageType `json:"package_type" gorm:"type:enum('FREE','PREMIUM');default:'FREE';index:idx_package_type"`

	// Pricing
	Price float64 `json:"price" gorm:"type:decimal(10,2);default:0.00"`

	// Image Support (responsive design ready)
	ImageURL      *string `json:"image_url" gorm:"size:500;comment:'Primary package image URL'"`
	ImageAlt      *string `json:"image_alt" gorm:"size:200;comment:'Alt text for accessibility'"`
	ThumbnailURL  *string `json:"thumbnail_url" gorm:"size:500;comment:'Small thumbnail (optional, can be generated from ImageURL)'"`
	ImageMetadata *string `json:"image_metadata" gorm:"type:text;comment:'JSON metadata: dimensions, file size, format, etc.'"`

	// Coupon Support
	CouponCode *string `json:"coupon_code" gorm:"size:50;index:idx_coupon_code;comment:'Optional default coupon for this package'"`

	// Validity Configuration
	ValidityType ValidityType `json:"validity_type" gorm:"type:enum('FIXED','RELATIVE');default:'RELATIVE';index:idx_validity_type"`
	ValidityDays *int         `json:"validity_days" gorm:"comment:'Days from enrollment (for RELATIVE type)'"`
	ValidityDate *time.Time   `json:"validity_date" gorm:"comment:'Fixed expiry date (for FIXED type)'"`

	// Metadata
	TotalExams int `json:"total_exams" gorm:"default:0"`

	// Analytics & Statistics (denormalized for performance)
	EnrollmentCount  int        `json:"enrollment_count" gorm:"default:0;index:idx_enrollment_count;comment:'Total number of users enrolled (includes free and paid)'"`
	LastEnrollmentAt *time.Time `json:"last_enrollment_at" gorm:"comment:'When the last user enrolled'"`

	// Status
	IsActive  bool `json:"is_active" gorm:"default:true;index:idx_active"`
	SortOrder int  `json:"sort_order" gorm:"default:0;index:idx_sort_order"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	PackageExams []PackageExam `json:"package_exams,omitempty" gorm:"foreignKey:PackageID"`
}

// TableName specifies the table name for Package
func (Package) TableName() string {
	return "packages"
}

// CalculateExpiryDate calculates the expiry date based on validity type
func (p *Package) CalculateExpiryDate(enrollmentDate time.Time) *time.Time {
	switch p.ValidityType {
	case ValidityTypeFixed:
		// Return the fixed validity date
		return p.ValidityDate
	case ValidityTypeRelative:
		// Calculate based on enrollment date + validity days
		if p.ValidityDays != nil {
			expiryDate := enrollmentDate.AddDate(0, 0, *p.ValidityDays)
			return &expiryDate
		}
		return nil
	default:
		return nil
	}
}

// IsValidConfiguration checks if the package validity configuration is valid
func (p *Package) IsValidConfiguration() bool {
	switch p.ValidityType {
	case ValidityTypeFixed:
		return p.ValidityDate != nil
	case ValidityTypeRelative:
		return p.ValidityDays != nil && *p.ValidityDays > 0
	default:
		return false
	}
}

// ImageMetadata represents the structure for image metadata JSON
type ImageMetadata struct {
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	FileSize    int64  `json:"file_size"`
	Format      string `json:"format"`
	OriginalURL string `json:"original_url,omitempty"`
}

// GetImageMetadata parses and returns the image metadata
func (p *Package) GetImageMetadata() (*ImageMetadata, error) {
	if p.ImageMetadata == nil || *p.ImageMetadata == "" {
		return nil, nil
	}

	var metadata ImageMetadata
	if err := json.Unmarshal([]byte(*p.ImageMetadata), &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// SetImageMetadata sets the image metadata as JSON string
func (p *Package) SetImageMetadata(metadata ImageMetadata) error {
	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	metadataStr := string(data)
	p.ImageMetadata = &metadataStr
	return nil
}

// GetOriginalImageURL returns the original image URL for Next.js Image component
// Next.js will handle all optimization automatically (WebP conversion, responsive sizes, etc.)
func (p *Package) GetOriginalImageURL() string {
	if p.ImageURL == nil || *p.ImageURL == "" {
		return ""
	}
	return *p.ImageURL
}

// GetNextJSImageSrc returns the image source optimized for Next.js Image component
// This works with local storage, S3, or any static file hosting
func (p *Package) GetNextJSImageSrc() string {
	if p.ImageURL == nil || *p.ImageURL == "" {
		// Return a placeholder if no image is available
		return "/images/package-placeholder.jpg"
	}

	// Return the URL as-is - Next.js Image component will handle optimization
	return *p.ImageURL
}

// GetImageURLForDevice returns the original image URL since Next.js handles optimization
// Next.js Image component automatically generates responsive sizes based on device
func (p *Package) GetImageURLForDevice(deviceType string) string {
	if p.ImageURL == nil {
		return "/images/package-placeholder.jpg"
	}

	// With Next.js Image component, we just return the original URL
	// Next.js automatically optimizes for device type, screen density, and format
	return *p.ImageURL
}

// HasImage checks if the package has an image configured
func (p *Package) HasImage() bool {
	return p.ImageURL != nil && *p.ImageURL != ""
}

// GetDisplayImageURL returns the best image URL for display
// Falls back to thumbnail if main image is not available
func (p *Package) GetDisplayImageURL() string {
	if p.HasImage() {
		return *p.ImageURL
	}
	if p.ThumbnailURL != nil && *p.ThumbnailURL != "" {
		return *p.ThumbnailURL
	}
	return "" // No image available
}

// GetImageAltText returns alt text for accessibility
func (p *Package) GetImageAltText() string {
	if p.ImageAlt != nil && *p.ImageAlt != "" {
		return *p.ImageAlt
	}
	// Fallback to package name
	return fmt.Sprintf("Image for %s package", p.Name)
}

// UpdateEnrollmentCount increments the enrollment counter
func (p *Package) UpdateEnrollmentCount(db *gorm.DB, increment int) error {
	now := time.Now()
	return db.Model(p).Updates(map[string]interface{}{
		"enrollment_count":   gorm.Expr("enrollment_count + ?", increment),
		"last_enrollment_at": now,
	}).Error
}

// UpdateActiveEnrollmentCount updates the active enrollment counter
func (p *Package) UpdateActiveEnrollmentCount(db *gorm.DB, increment int) error {
	return db.Model(p).Update("active_enrollment_count", gorm.Expr("active_enrollment_count + ?", increment)).Error
}

// RecalculateEnrollmentStats recalculates enrollment statistics from actual data
func (p *Package) RecalculateEnrollmentStats(db *gorm.DB) error {
	var totalCount int64
	var activeCount int64
	var lastEnrollment time.Time

	// Count total enrollments
	if err := db.Model(&UserPackageEnrollment{}).Where("package_id = ?", p.ID).Count(&totalCount).Error; err != nil {
		return err
	}

	// Count active enrollments (not expired)
	now := time.Now()
	if err := db.Model(&UserPackageEnrollment{}).
		Where("package_id = ? AND (expires_at IS NULL OR expires_at > ?)", p.ID, now).
		Count(&activeCount).Error; err != nil {
		return err
	}

	// Get last enrollment date
	var enrollment UserPackageEnrollment
	if err := db.Where("package_id = ?", p.ID).
		Order("enrolled_at DESC").
		First(&enrollment).Error; err == nil {
		lastEnrollment = enrollment.EnrolledAt
	}

	// Update the package
	updates := map[string]interface{}{
		"enrollment_count":        totalCount,
		"active_enrollment_count": activeCount,
	}

	if !lastEnrollment.IsZero() {
		updates["last_enrollment_at"] = lastEnrollment
	}

	return db.Model(p).Updates(updates).Error
}
