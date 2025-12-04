package models

import (
	"time"

	"gorm.io/gorm"
)

// PackageExam represents the many-to-many relationship between packages and exams
type PackageExam struct {
	ID        uint `json:"id" gorm:"primarykey"`
	PackageID uint `json:"package_id" gorm:"not null;index:idx_package_id"`
	ExamID    uint `json:"exam_id" gorm:"not null;index:idx_exam_id"`

	// Ordering within package
	SortOrder int `json:"sort_order" gorm:"default:0"`

	// Status
	IsActive bool `json:"is_active" gorm:"default:true"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Package Package `json:"package,omitempty" gorm:"foreignKey:PackageID"`
	Exam    Exam    `json:"exam,omitempty" gorm:"foreignKey:ExamID"`
}

// TableName specifies the table name for PackageExam
func (PackageExam) TableName() string {
	return "package_exams"
}
