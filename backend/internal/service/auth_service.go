package service

import (
	"errors"
	"quizora-backend/internal/models"
	"quizora-backend/internal/utils"

	"gorm.io/gorm"
)

// AuthService handles authentication business logic
type AuthService struct {
	db        *gorm.DB
	jwtSecret string
}

// NewAuthService creates a new auth service
func NewAuthService(db *gorm.DB, jwtSecret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

// Register registers a new user
func (s *AuthService) Register(req models.RegisterRequest) (*models.AuthResponse, error) {
	// Validate input
	if !utils.ValidateName(req.Name) {
		return nil, errors.New("invalid name format")
	}

	if !utils.ValidateMSISDN(req.MSISDN) {
		return nil, errors.New("invalid MSISDN format")
	}

	if !utils.ValidatePassword(req.Password) {
		return nil, errors.New("password must be at least 6 characters long")
	}

	// Normalize MSISDN
	normalizedMSISDN := utils.NormalizeMSISDN(req.MSISDN)

	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("msisdn = ?", normalizedMSISDN).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists with this MSISDN")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := models.User{
		Name:     req.Name,
		MSISDN:   normalizedMSISDN,
		Password: hashedPassword,
		IsActive: true,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.MSISDN, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &models.AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	// Validate input
	if !utils.ValidateMSISDN(req.MSISDN) {
		return nil, errors.New("invalid MSISDN format")
	}

	// Normalize MSISDN
	normalizedMSISDN := utils.NormalizeMSISDN(req.MSISDN)

	// Find user
	var user models.User
	if err := s.db.Where("msisdn = ? AND is_active = ?", normalizedMSISDN, true).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, errors.New("database error")
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.MSISDN, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &models.AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ? AND is_active = ?", userID, true).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("database error")
	}

	return &user, nil
}

// SocialAuth handles social authentication (login or registration)
func (s *AuthService) SocialAuth(provider string, userInfo *models.SocialUserInfo) (*models.AuthResponse, error) {
	// Validate provider
	if provider != "google" && provider != "facebook" {
		return nil, errors.New("invalid auth provider")
	}

	// Check if user exists with this provider and provider user ID
	var user models.User
	err := s.db.Where("auth_provider = ? AND provider_user_id = ?", provider, userInfo.ProviderUserID).
		First(&user).Error

	if err == nil {
		// User exists - login
		if !user.IsActive {
			return nil, errors.New("user account is inactive")
		}

		// Update user info in case it changed
		user.Name = userInfo.Name
		user.Email = userInfo.Email
		user.ProfilePicture = userInfo.ProfilePicture
		user.EmailVerified = userInfo.EmailVerified

		if err := s.db.Save(&user).Error; err != nil {
			return nil, errors.New("failed to update user info")
		}

		// Generate JWT token
		token, err := utils.GenerateJWT(user.ID, user.MSISDN, s.jwtSecret)
		if err != nil {
			return nil, errors.New("failed to generate token")
		}

		return &models.AuthResponse{
			User:  user.ToResponse(),
			Token: token,
		}, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("database error")
	}

	// User doesn't exist - check if email is already used by another account
	if userInfo.Email != "" {
		var existingUser models.User
		if err := s.db.Where("email = ?", userInfo.Email).First(&existingUser).Error; err == nil {
			// Email exists with different provider
			return nil, errors.New("email already registered with different provider")
		}
	}

	// Create new user
	newUser := models.User{
		Name:           userInfo.Name,
		Email:          userInfo.Email,
		AuthProvider:   provider,
		ProviderUserID: userInfo.ProviderUserID,
		ProfilePicture: userInfo.ProfilePicture,
		EmailVerified:  userInfo.EmailVerified,
		IsActive:       true,
		MSISDN:         "", // No phone number for social auth users
		Password:       "", // No password for social auth users
	}

	if err := s.db.Create(&newUser).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(newUser.ID, newUser.MSISDN, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &models.AuthResponse{
		User:  newUser.ToResponse(),
		Token: token,
	}, nil
}

// LinkSocialAccount links a social account to an existing user
func (s *AuthService) LinkSocialAccount(userID uint, provider string, userInfo *models.SocialUserInfo) error {
	// Validate provider
	if provider != "google" && provider != "facebook" {
		return errors.New("invalid auth provider")
	}

	// Get existing user
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("database error")
	}

	// Check if this provider account is already linked to another user
	var existingUser models.User
	err := s.db.Where("auth_provider = ? AND provider_user_id = ? AND id != ?",
		provider, userInfo.ProviderUserID, userID).First(&existingUser).Error

	if err == nil {
		return errors.New("this social account is already linked to another user")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("database error")
	}

	// Update user with social account info
	user.Email = userInfo.Email
	user.ProfilePicture = userInfo.ProfilePicture
	user.EmailVerified = userInfo.EmailVerified

	if err := s.db.Save(&user).Error; err != nil {
		return errors.New("failed to link social account")
	}

	return nil
}

// GetUserByEmail retrieves a user by email
func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ? AND is_active = ?", email, true).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("database error")
	}

	return &user, nil
}
