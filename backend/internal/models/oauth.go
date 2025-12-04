package models

// SocialUserInfo represents user information retrieved from OAuth providers
type SocialUserInfo struct {
	ProviderUserID string `json:"provider_user_id"` // Unique ID from provider
	Email          string `json:"email"`            // User's email
	Name           string `json:"name"`             // User's full name
	ProfilePicture string `json:"profile_picture"`  // Avatar URL
	EmailVerified  bool   `json:"email_verified"`   // Email verification status
}

// OAuthCallbackRequest represents the OAuth callback request
type OAuthCallbackRequest struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state" binding:"required"`
}

// OAuthState represents the OAuth state stored in Redis
type OAuthState struct {
	State     string `json:"state"`
	Provider  string `json:"provider"`
	CreatedAt int64  `json:"created_at"`
}

// GoogleAuthRequest represents Google credential authentication request
type GoogleAuthRequest struct {
	Credential string `json:"credential" binding:"required"` // Google ID token
}

// FacebookAuthRequest represents Facebook authentication request
type FacebookAuthRequest struct {
	AccessToken string `json:"access_token" binding:"required"` // Facebook access token
}
