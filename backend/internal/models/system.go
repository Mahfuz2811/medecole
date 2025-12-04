package models

import (
	"time"

	"gorm.io/gorm"
)

// System represents the system table (body systems like respiratory, cardiovascular, etc.)
type System struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	SubjectID   uint           `json:"subject_id" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Slug        string         `json:"slug" gorm:"size:100;not null"`
	Description *string        `json:"description" gorm:"type:text"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Subject   Subject    `json:"subject,omitempty" gorm:"foreignKey:SubjectID"`
	Questions []Question `json:"questions,omitempty" gorm:"foreignKey:SystemID"`
}

// TableName specifies the table name for System
func (System) TableName() string {
	return "systems"
}
