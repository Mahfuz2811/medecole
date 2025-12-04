package unit

import (
	"os"
	"quizora-backend/internal/config"
	"quizora-backend/internal/database"
	"quizora-backend/internal/models"
	"quizora-backend/internal/service"
	"quizora-backend/tests/helpers"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	// Setup test environment
	os.Setenv("GIN_MODE", "test")
	code := m.Run()
	os.Exit(code)
}

func setupTestAuthService(t *testing.T) (*service.AuthService, *gorm.DB, func()) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     "3306",
			User:     "root",
			Password: "secret",
			Name:     "quizora_test",
		},
		JWT: config.JWTConfig{
			Secret: "test-secret-key",
		},
	}

	db, err := database.New(cfg)
	require.NoError(t, err)

	// Clean and migrate database
	helpers.CleanupTestDB(db.DB)
	err = db.AutoMigrate()
	require.NoError(t, err)

	authService := service.NewAuthService(db.DB, cfg.JWT.Secret)

	cleanup := func() {
		helpers.CleanupTestDB(db.DB)
		sqlDB, _ := db.DB.DB()
		sqlDB.Close()
	}

	return authService, db.DB, cleanup
}

func TestAuthService_Register(t *testing.T) {
	authService, db, cleanup := setupTestAuthService(t)
	defer cleanup()

	tests := []struct {
		name        string
		input       models.RegisterRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid registration",
			input: models.RegisterRequest{
				Name:     "John Doe",
				MSISDN:   "01712345678",
				Password: "password123",
			},
			expectError: false,
		},
		{
			name: "Duplicate MSISDN",
			input: models.RegisterRequest{
				Name:     "Jane Doe",
				MSISDN:   "01712345678", // Same as above
				Password: "password456",
			},
			expectError: true,
			errorMsg:    "already exists",
		},
		{
			name: "Empty name",
			input: models.RegisterRequest{
				Name:     "",
				MSISDN:   "01787654321",
				Password: "password123",
			},
			expectError: true,
		},
		{
			name: "Invalid MSISDN format",
			input: models.RegisterRequest{
				Name:     "Test User",
				MSISDN:   "123456789", // Invalid format
				Password: "password123",
			},
			expectError: true,
		},
		{
			name: "Short password",
			input: models.RegisterRequest{
				Name:     "Test User",
				MSISDN:   "01798765432",
				Password: "123", // Too short
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authResp, err := authService.Register(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, authResp)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, authResp)
				assert.Equal(t, tt.input.Name, authResp.User.Name)
				// MSISDN gets normalized to +88 format
				assert.True(t, strings.Contains(authResp.User.MSISDN, tt.input.MSISDN) ||
					strings.HasSuffix(authResp.User.MSISDN, tt.input.MSISDN))
				assert.NotEmpty(t, authResp.User.ID)
				assert.True(t, authResp.User.IsActive)
				assert.NotEmpty(t, authResp.User.CreatedAt)
				assert.NotEmpty(t, authResp.Token)

				// Verify user exists in database
				var dbUser models.User
				err := db.First(&dbUser, authResp.User.ID).Error
				assert.NoError(t, err)
				assert.Equal(t, authResp.User.ID, dbUser.ID)

				// Verify password is hashed
				assert.NotEqual(t, tt.input.Password, dbUser.Password)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	authService, _, cleanup := setupTestAuthService(t)
	defer cleanup()

	// Create a test user first
	registerReq := models.RegisterRequest{
		Name:     "Test User",
		MSISDN:   "01712345678",
		Password: "password123",
	}
	registerResp, err := authService.Register(registerReq)
	require.NoError(t, err)
	require.NotNil(t, registerResp)

	tests := []struct {
		name        string
		input       models.LoginRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid login",
			input: models.LoginRequest{
				MSISDN:   "01712345678",
				Password: "password123",
			},
			expectError: false,
		},
		{
			name: "Invalid MSISDN",
			input: models.LoginRequest{
				MSISDN:   "01799999999",
				Password: "password123",
			},
			expectError: true,
			errorMsg:    "invalid credentials",
		},
		{
			name: "Invalid password",
			input: models.LoginRequest{
				MSISDN:   "01712345678",
				Password: "wrongpassword",
			},
			expectError: true,
			errorMsg:    "invalid credentials",
		},
		{
			name: "Empty MSISDN",
			input: models.LoginRequest{
				MSISDN:   "",
				Password: "password123",
			},
			expectError: true,
		},
		{
			name: "Empty password",
			input: models.LoginRequest{
				MSISDN:   "01712345678",
				Password: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authResp, err := authService.Login(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, authResp)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, authResp)
				assert.Equal(t, registerResp.User.ID, authResp.User.ID)
				assert.Equal(t, registerResp.User.Name, authResp.User.Name)
				assert.Equal(t, registerResp.User.MSISDN, authResp.User.MSISDN)
				assert.NotEmpty(t, authResp.Token)
			}
		})
	}
}

func TestAuthService_GetUserByID(t *testing.T) {
	authService, _, cleanup := setupTestAuthService(t)
	defer cleanup()

	// Create a test user
	registerReq := models.RegisterRequest{
		Name:     "Test User",
		MSISDN:   "01712345678",
		Password: "password123",
	}
	registerResp, err := authService.Register(registerReq)
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      uint
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid user ID",
			userID:      registerResp.User.ID,
			expectError: false,
		},
		{
			name:        "Invalid user ID",
			userID:      99999,
			expectError: true,
			errorMsg:    "not found",
		},
		{
			name:        "Zero user ID",
			userID:      0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundUser, err := authService.GetUserByID(tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, foundUser)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, foundUser)
				assert.Equal(t, registerResp.User.ID, foundUser.ID)
				assert.Equal(t, registerResp.User.Name, foundUser.Name)
				assert.Equal(t, registerResp.User.MSISDN, foundUser.MSISDN)
				assert.True(t, foundUser.IsActive)
			}
		})
	}
}

func TestAuthService_ConcurrentRegistrations(t *testing.T) {
	authService, _, cleanup := setupTestAuthService(t)
	defer cleanup()

	// Test concurrent registrations with the same MSISDN
	// Only one should succeed
	msisdn := "01712345678"

	type result struct {
		authResp *models.AuthResponse
		err      error
	}

	results := make(chan result, 2)

	// Start two goroutines trying to register the same MSISDN
	for i := 0; i < 2; i++ {
		go func(id int) {
			req := models.RegisterRequest{
				Name:     "User " + string(rune(id+'A')), // Use letter instead of number
				MSISDN:   msisdn,
				Password: "password123",
			}
			authResp, err := authService.Register(req)
			results <- result{authResp: authResp, err: err}
		}(i)
	}

	// Collect results
	var successCount, errorCount int
	for i := 0; i < 2; i++ {
		res := <-results
		if res.err != nil {
			errorCount++
			// Could be either "already exists" or "failed to create user" due to race condition
			assert.True(t, strings.Contains(res.err.Error(), "already exists") ||
				strings.Contains(res.err.Error(), "failed to create user"))
		} else {
			successCount++
			assert.NotNil(t, res.authResp)
		}
	}

	// Exactly one should succeed, one should fail
	assert.Equal(t, 1, successCount, "Exactly one registration should succeed")
	assert.Equal(t, 1, errorCount, "Exactly one registration should fail")
}
