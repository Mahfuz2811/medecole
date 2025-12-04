package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"github.com/Mahfuz2811/medecole/backend/internal/config"
	"github.com/Mahfuz2811/medecole/backend/internal/database"
	"github.com/Mahfuz2811/medecole/backend/internal/handlers"
	"github.com/Mahfuz2811/medecole/backend/internal/middleware"
	"github.com/Mahfuz2811/medecole/backend/internal/models"
	"github.com/Mahfuz2811/medecole/backend/internal/service"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TestApp holds the test application setup
type TestApp struct {
	Router      *gin.Engine
	DB          *gorm.DB
	AuthService *service.AuthService
	Config      *config.Config
}

// SetupTestApp creates a test application instance
func SetupTestApp(t *testing.T) *TestApp {
	// Set test environment
	os.Setenv("GIN_MODE", "test")
	gin.SetMode(gin.TestMode)

	// Create test configuration
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
		CORS: config.CORSConfig{
			FrontendURL: "http://localhost:3000",
		},
	}

	// Initialize test database
	db, err := database.New(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clean and migrate database
	CleanupTestDB(db.DB)
	if err := db.AutoMigrate(); err != nil {
		t.Fatalf("Failed to run test migrations: %v", err)
	}

	// Initialize services
	authService := service.NewAuthService(db.DB, cfg.JWT.Secret)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize router
	r := gin.New()

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.CORS.FrontendURL}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	// Setup routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Quizora Backend API is running",
		})
	})

	api := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
		}

		// Protected routes
		protected := api.Group("/auth")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret, authService))
		{
			protected.GET("/profile", authHandler.Profile)
		}
	}

	return &TestApp{
		Router:      r,
		DB:          db.DB,
		AuthService: authService,
		Config:      cfg,
	}
}

// CleanupTestDB cleans up test database
func CleanupTestDB(db *gorm.DB) {
	// Disable foreign key checks to allow dropping tables
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// Drop all tables in any order (foreign keys disabled)
	db.Exec("DROP TABLE IF EXISTS user_question_answers")
	db.Exec("DROP TABLE IF EXISTS user_exam_attempts")
	db.Exec("DROP TABLE IF EXISTS user_package_enrollments")
	db.Exec("DROP TABLE IF EXISTS package_exams")
	db.Exec("DROP TABLE IF EXISTS exams")
	db.Exec("DROP TABLE IF EXISTS packages")
	db.Exec("DROP TABLE IF EXISTS users")

	// Re-enable foreign key checks
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}

// TeardownTestApp cleans up after tests
func TeardownTestApp(app *TestApp) {
	CleanupTestDB(app.DB)
	sqlDB, _ := app.DB.DB()
	sqlDB.Close()
}

// CreateTestUser creates a test user in the database
func CreateTestUser(db *gorm.DB, name, msisdn, password string) (*models.User, error) {
	user := &models.User{
		Name:     name,
		MSISDN:   msisdn,
		Password: password, // This should be hashed in real implementation
		IsActive: true,
	}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// MakeRequest makes an HTTP request to the test server
func MakeRequest(router *gin.Engine, method, url string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer

	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, url, reqBody)
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// ParseJSONResponse parses JSON response into a struct
func ParseJSONResponse(t *testing.T, response *httptest.ResponseRecorder, v interface{}) {
	if err := json.Unmarshal(response.Body.Bytes(), v); err != nil {
		t.Fatalf("Failed to parse JSON response: %v\nResponse body: %s", err, response.Body.String())
	}
}

// AssertJSONResponse asserts the response matches expected JSON
func AssertJSONResponse(t *testing.T, response *httptest.ResponseRecorder, expectedStatus int, expectedBody interface{}) {
	if response.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d. Response: %s", expectedStatus, response.Code, response.Body.String())
	}

	if expectedBody != nil {
		expectedJSON, _ := json.Marshal(expectedBody)
		var actualJSON, expectedJSONNormalized interface{}

		json.Unmarshal(response.Body.Bytes(), &actualJSON)
		json.Unmarshal(expectedJSON, &expectedJSONNormalized)

		actualStr, _ := json.Marshal(actualJSON)
		expectedStr, _ := json.Marshal(expectedJSONNormalized)

		if string(actualStr) != string(expectedStr) {
			t.Errorf("Response body mismatch.\nExpected: %s\nActual: %s", string(expectedStr), string(actualStr))
		}
	}
}

// ExtractJWTToken extracts JWT token from auth response
func ExtractJWTToken(t *testing.T, response *httptest.ResponseRecorder) string {
	var authResponse models.AuthResponse
	ParseJSONResponse(t, response, &authResponse)

	if authResponse.Token == "" {
		t.Fatal("No token found in auth response")
	}

	return authResponse.Token
}

// GetAuthHeader creates authorization header with Bearer token
func GetAuthHeader(token string) map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
}
