package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the user model in the database
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null;size:100"`
	MSISDN    string         `json:"msisdn" gorm:"size:20;uniqueIndex"` // Nullable for social auth users, unique constraint
	Password  string         `json:"-" gorm:"type:varchar(255)"`        // Nullable for social auth users, hidden from JSON
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete

	// Social Authentication Fields
	Email          string `json:"email" gorm:"size:255;index"`                           // Email for social auth users
	AuthProvider   string `json:"auth_provider" gorm:"size:20;default:'local';not null"` // "local", "google", "facebook"
	ProviderUserID string `json:"provider_user_id" gorm:"size:255;index"`                // Google/Facebook user ID
	ProfilePicture string `json:"profile_picture" gorm:"type:text"`                      // Avatar URL from social provider
	EmailVerified  bool   `json:"email_verified" gorm:"default:false"`                   // Email verification status
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	MSISDN         string    `json:"msisdn,omitempty"`          // Omit if empty for social users
	Email          string    `json:"email,omitempty"`           // Include for social users
	AuthProvider   string    `json:"auth_provider"`             // Show auth method
	ProfilePicture string    `json:"profile_picture,omitempty"` // Include if available
	EmailVerified  bool      `json:"email_verified"`            // Email verification status
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

// ToResponse converts User model to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:             u.ID,
		Name:           u.Name,
		MSISDN:         u.MSISDN,
		Email:          u.Email,
		AuthProvider:   u.AuthProvider,
		ProfilePicture: u.ProfilePicture,
		EmailVerified:  u.EmailVerified,
		IsActive:       u.IsActive,
		CreatedAt:      u.CreatedAt,
	}
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
