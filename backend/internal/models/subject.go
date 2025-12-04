package models

import (
	"time"

	"gorm.io/gorm"
)

// Subject represents the subject table
type Subject struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Slug        string         `json:"slug" gorm:"size:100;uniqueIndex;not null"`
	Description *string        `json:"description" gorm:"type:text"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Systems []System `json:"systems,omitempty" gorm:"foreignKey:SubjectID"`
}

// TableName specifies the table name for Subject
func (Subject) TableName() string {
	return "subjects"
}
