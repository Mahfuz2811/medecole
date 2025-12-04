package integration

import (
	"fmt"
	"net/http"
	"os"
	"quizora-backend/internal/models"
	"quizora-backend/tests/helpers"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Setup test environment
	os.Setenv("GIN_MODE", "test")
	code := m.Run()
	os.Exit(code)
}

func TestAuthFlow_Complete(t *testing.T) {
	app := helpers.SetupTestApp(t)
	defer helpers.TeardownTestApp(app)

	// Test data
	registerReq := models.RegisterRequest{
		Name:     "Integration Test User",
		MSISDN:   "01712345678",
		Password: "password123",
	}

	loginReq := models.LoginRequest{
		MSISDN:   registerReq.MSISDN,
		Password: registerReq.Password,
	}

	// Step 1: Register new user
	t.Run("Register new user", func(t *testing.T) {
		response := helpers.MakeRequest(app.Router, "POST", "/api/v1/auth/register", registerReq, nil)

		assert.Equal(t, http.StatusCreated, response.Code)

		var authResp models.AuthResponse
		helpers.ParseJSONResponse(t, response, &authResp)

		assert.NotEmpty(t, authResp.Token)
		assert.Equal(t, registerReq.Name, authResp.User.Name)
		// MSISDN gets normalized to +88 format
		assert.Contains(t, authResp.User.MSISDN, "01712345678")
		assert.True(t, authResp.User.IsActive)
		assert.NotZero(t, authResp.User.ID)
		assert.NotZero(t, authResp.User.CreatedAt)
	})

	// Step 2: Try to register duplicate user
	t.Run("Register duplicate user should fail", func(t *testing.T) {
		response := helpers.MakeRequest(app.Router, "POST", "/api/v1/auth/register", registerReq, nil)

		assert.Equal(t, http.StatusConflict, response.Code)

		var errorResp models.ErrorResponse
		helpers.ParseJSONResponse(t, response, &errorResp)

		assert.Equal(t, "Registration Failed", errorResp.Error)
		assert.Contains(t, errorResp.Message, "already exists")
	})

	// Step 3: Login with correct credentials
	t.Run("Login with correct credentials", func(t *testing.T) {
		response := helpers.MakeRequest(app.Router, "POST", "/api/v1/auth/login", loginReq, nil)

		assert.Equal(t, http.StatusOK, response.Code)

		var authResp models.AuthResponse
		helpers.ParseJSONResponse(t, response, &authResp)

		assert.NotEmpty(t, authResp.Token)
		assert.Equal(t, registerReq.Name, authResp.User.Name)
		assert.Contains(t, authResp.User.MSISDN, "01712345678")
	})

	// Step 4: Login with wrong password
	t.Run("Login with wrong password should fail", func(t *testing.T) {
		wrongLoginReq := models.LoginRequest{
			MSISDN:   registerReq.MSISDN,
			Password: "wrongpassword",
		}

		response := helpers.MakeRequest(app.Router, "POST", "/api/v1/auth/login", wrongLoginReq, nil)

		assert.Equal(t, http.StatusUnauthorized, response.Code)

		var errorResp models.ErrorResponse
		helpers.ParseJSONResponse(t, response, &errorResp)

		assert.Equal(t, "Login Failed", errorResp.Error)
		assert.Contains(t, errorResp.Message, "invalid credentials")
	})

	// Step 5: Login with non-existent user
	t.Run("Login with non-existent user should fail", func(t *testing.T) {
		nonExistentLoginReq := models.LoginRequest{
			MSISDN:   "01799999999",
			Password: "password123",
		}

		response := helpers.MakeRequest(app.Router, "POST", "/api/v1/auth/login", nonExistentLoginReq, nil)

		assert.Equal(t, http.StatusUnauthorized, response.Code)

		var errorResp models.ErrorResponse
		helpers.ParseJSONResponse(t, response, &errorResp)

		assert.Equal(t, "Login Failed", errorResp.Error)
		assert.Contains(t, errorResp.Message, "invalid credentials")
	})

	// Step 6: Access profile with valid token
	t.Run("Access profile with valid token", func(t *testing.T) {
		// First login to get token
		loginResponse := helpers.MakeRequest(app.Router, "POST", "/api/v1/auth/login", loginReq, nil)
		require.Equal(t, http.StatusOK, loginResponse.Code)

		token := helpers.ExtractJWTToken(t, loginResponse)
		headers := helpers.GetAuthHeader(token)

		// Access profile
		response := helpers.MakeRequest(app.Router, "GET", "/api/v1/auth/profile", nil, headers)

		assert.Equal(t, http.StatusOK, response.Code)

		var user models.UserResponse
		helpers.ParseJSONResponse(t, response, &user)

		assert.Equal(t, registerReq.Name, user.Name)
		assert.Contains(t, user.MSISDN, "01712345678")
		assert.True(t, user.IsActive)
		assert.NotZero(t, user.ID)
	})

	// Step 7: Access profile without token
	t.Run("Access profile without token should fail", func(t *testing.T) {
		response := helpers.MakeRequest(app.Router, "GET", "/api/v1/auth/profile", nil, nil)

		assert.Equal(t, http.StatusUnauthorized, response.Code)

		var errResp models.ErrorResponse
		helpers.ParseJSONResponse(t, response, &errResp)
		assert.Contains(t, errResp.Message, "Authorization header is required")
	})

	// Step 8: Access profile with invalid token
	t.Run("Access profile with invalid token should fail", func(t *testing.T) {
		headers := map[string]string{
			"Authorization": "Bearer invalid-token",
		}
		response := helpers.MakeRequest(app.Router, "GET", "/api/v1/auth/profile", nil, headers)

		assert.Equal(t, http.StatusUnauthorized, response.Code)

		var errResp models.ErrorResponse
		helpers.ParseJSONResponse(t, response, &errResp)
		assert.Contains(t, errResp.Message, "Invalid or expired token")
	})

	// Step 9: Logout with valid token
	t.Run("Logout with valid token", func(t *testing.T) {
		// First login to get token
		loginResponse := helpers.MakeRequest(app.Router, "POST", "/api/v1/auth/login", loginReq, nil)
		require.Equal(t, http.StatusOK, loginResponse.Code)

		token := helpers.ExtractJWTToken(t, loginResponse)
		headers := helpers.GetAuthHeader(token)

		// Logout
		response := helpers.MakeRequest(app.Router, "POST", "/api/v1/auth/logout", nil, headers)

		assert.Equal(t, http.StatusOK, response.Code)

		var successResp models.SuccessResponse
		helpers.ParseJSONResponse(t, response, &successResp)

		assert.Contains(t, successResp.Message, "Logged out successfully")
	})

	// Step 10: Logout without token
	t.Run("Logout without token should succeed", func(t *testing.T) {
		response := helpers.MakeRequest(app.Router, "POST", "/api/v1/auth/logout", nil, nil)

		assert.Equal(t, http.StatusOK, response.Code)
	})
}

func TestRegistrationValidation(t *testing.T) {
	app := helpers.SetupTestApp(t)
	defer helpers.TeardownTestApp(app)

	tests := []struct {
		name           string
		request        interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Missing name",
			request: map[string]interface{}{
				"msisdn":   "01712345678",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Name",
		},
		{
			name: "Missing MSISDN",
			request: map[string]interface{}{
				"name":     "Test User",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "MSISDN",
		},
		{
			name: "Missing password",
			request: map[string]interface{}{
				"name":   "Test User",
				"msisdn": "01712345678",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Password",
		},
		{
			name: "Empty name",
			request: models.RegisterRequest{
				Name:     "",
				MSISDN:   "01712345678",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Name",
		},
		{
			name: "Empty MSISDN",
			request: models.RegisterRequest{
				Name:     "Test User",
				MSISDN:   "",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "MSISDN",
		},
		{
			name: "Empty password",
			request: models.RegisterRequest{
				Name:     "Test User",
				MSISDN:   "01712345678",
				Password: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Password",
		},
		{
			name: "Invalid MSISDN format",
			request: models.RegisterRequest{
				Name:     "Test User",
				MSISDN:   "123456789",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid MSISDN format",
		},
		{
			name: "Short password",
			request: models.RegisterRequest{
				Name:     "Test User",
				MSISDN:   "01712345678",
				Password: "123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "min",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := helpers.MakeRequest(app.Router, "POST", "/api/v1/auth/register", tt.request, nil)

			assert.Equal(t, tt.expectedStatus, response.Code)

			var errorResp models.ErrorResponse
			helpers.ParseJSONResponse(t, response, &errorResp)

			if tt.name == "Invalid MSISDN format" {
				// Service layer validation returns "Registration Failed"
				assert.Equal(t, "Registration Failed", errorResp.Error)
				assert.Contains(t, errorResp.Message, tt.expectedError)
			} else {
				// Binding layer validation returns "Bad Request"
				assert.Equal(t, "Bad Request", errorResp.Error)
				assert.Contains(t, errorResp.Message, tt.expectedError)
			}
		})
	}
}

func TestLoginValidation(t *testing.T) {
	app := helpers.SetupTestApp(t)
	defer helpers.TeardownTestApp(app)

	tests := []struct {
		name           string
		request        interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Missing MSISDN",
			request: map[string]interface{}{
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "MSISDN",
		},
		{
			name: "Missing password",
			request: map[string]interface{}{
				"msisdn": "01712345678",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Password",
		},
		{
			name: "Empty MSISDN",
			request: models.LoginRequest{
				MSISDN:   "",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "MSISDN",
		},
		{
			name: "Empty password",
			request: models.LoginRequest{
				MSISDN:   "01712345678",
				Password: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := helpers.MakeRequest(app.Router, "POST", "/api/v1/auth/login", tt.request, nil)

			assert.Equal(t, tt.expectedStatus, response.Code)

			var errorResp models.ErrorResponse
			helpers.ParseJSONResponse(t, response, &errorResp)

			assert.Equal(t, "Bad Request", errorResp.Error)
			assert.Contains(t, errorResp.Message, tt.expectedError)
		})
	}
}

func TestConcurrentRegistrations(t *testing.T) {
	app := helpers.SetupTestApp(t)
	defer helpers.TeardownTestApp(app)

	msisdn := "01712345678"
	numGoroutines := 5

	results := make(chan int, numGoroutines)

	// Start multiple goroutines trying to register the same MSISDN
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			req := models.RegisterRequest{
				Name:     fmt.Sprintf("User %c", 'A'+id),
				MSISDN:   msisdn,
				Password: "password123",
			}

			response := helpers.MakeRequest(app.Router, "POST", "/api/v1/auth/register", req, nil)
			results <- response.Code
		}(i)
	}

	// Collect results
	var successCount, errorCount int
	for i := 0; i < numGoroutines; i++ {
		statusCode := <-results
		switch statusCode {
		case http.StatusCreated:
			successCount++
		case http.StatusConflict, http.StatusInternalServerError, http.StatusBadRequest:
			errorCount++
		default:
			t.Logf("Unexpected status code: %d", statusCode)
			errorCount++
		}
	}

	// Exactly one should succeed, others should fail
	assert.Equal(t, 1, successCount, "Exactly one registration should succeed")
	assert.Equal(t, numGoroutines-1, errorCount, "Other registrations should fail")
}

func TestHealthEndpoint(t *testing.T) {
	app := helpers.SetupTestApp(t)
	defer helpers.TeardownTestApp(app)

	response := helpers.MakeRequest(app.Router, "GET", "/health", nil, nil)

	assert.Equal(t, http.StatusOK, response.Code)

	var healthResp map[string]interface{}
	helpers.ParseJSONResponse(t, response, &healthResp)

	assert.Equal(t, "ok", healthResp["status"])
	assert.Contains(t, healthResp["message"], "Quizora Backend API is running")
}

func TestCORSHeaders(t *testing.T) {
	app := helpers.SetupTestApp(t)
	defer helpers.TeardownTestApp(app)

	// Test OPTIONS request
	response := helpers.MakeRequest(app.Router, "OPTIONS", "/api/v1/auth/register", nil, map[string]string{
		"Origin":                         "http://localhost:3000",
		"Access-Control-Request-Method":  "POST",
		"Access-Control-Request-Headers": "Content-Type,Authorization",
	})

	assert.Equal(t, http.StatusNoContent, response.Code)
	assert.Equal(t, "http://localhost:3000", response.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, response.Header().Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, response.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
	assert.Contains(t, response.Header().Get("Access-Control-Allow-Headers"), "Authorization")
	assert.Equal(t, "true", response.Header().Get("Access-Control-Allow-Credentials"))
}

func TestTokenPersistence(t *testing.T) {
	app := helpers.SetupTestApp(t)
	defer helpers.TeardownTestApp(app)

	// Register and login
	registerReq := models.RegisterRequest{
		Name:     "Token Test User",
		MSISDN:   "01712345678",
		Password: "password123",
	}

	regResponse := helpers.MakeRequest(app.Router, "POST", "/api/v1/auth/register", registerReq, nil)
	require.Equal(t, http.StatusCreated, regResponse.Code)

	loginReq := models.LoginRequest{
		MSISDN:   registerReq.MSISDN,
		Password: registerReq.Password,
	}

	loginResponse := helpers.MakeRequest(app.Router, "POST", "/api/v1/auth/login", loginReq, nil)
	require.Equal(t, http.StatusOK, loginResponse.Code)

	token := helpers.ExtractJWTToken(t, loginResponse)
	headers := helpers.GetAuthHeader(token)

	// Use token multiple times
	for i := 0; i < 3; i++ {
		response := helpers.MakeRequest(app.Router, "GET", "/api/v1/auth/profile", nil, headers)
		assert.Equal(t, http.StatusOK, response.Code, "Token should be valid for multiple requests")

		// Add small delay to test token persistence over time
		time.Sleep(100 * time.Millisecond)
	}
}

func TestDifferentMSISDNFormats(t *testing.T) {
	app := helpers.SetupTestApp(t)
	defer helpers.TeardownTestApp(app)

	validMSISDNs := []string{
		"01712345678", // Grameenphone (013,014,015,016,017,018,019)
		"01812345679", // Robi (018)
		"01912345680", // Banglalink (019,014)
		"01612345681", // Airtel (016)
		"01512345682", // Teletalk (015)
	}

	for i, msisdn := range validMSISDNs {
		t.Run(fmt.Sprintf("MSISDN_%s", msisdn), func(t *testing.T) {
			names := []string{"Alice", "Bob", "Charlie", "Diana", "Eve"}
			req := models.RegisterRequest{
				Name:     names[i],
				MSISDN:   msisdn,
				Password: "password123",
			}

			response := helpers.MakeRequest(app.Router, "POST", "/api/v1/auth/register", req, nil)
			assert.Equal(t, http.StatusCreated, response.Code)
		})
	}
}
