package handlers

import (
	"net/http"
	"quizora-backend/internal/models"
	"quizora-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// OAuthHandler handles OAuth authentication endpoints
type OAuthHandler struct {
	oauthService *service.OAuthService
	authService  *service.AuthService
	frontendURL  string
}

// NewOAuthHandler creates a new OAuth handler
func NewOAuthHandler(oauthService *service.OAuthService, authService *service.AuthService, frontendURL string) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
		authService:  authService,
		frontendURL:  frontendURL,
	}
}

// GoogleLogin initiates Google OAuth flow
// @Summary Initiate Google OAuth
// @Description Redirects user to Google login page
// @Tags oauth
// @Accept json
// @Produce json
// @Success 302 {string} string "Redirect to Google"
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/google [get]
func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
	// Generate state token
	state, err := h.oauthService.GenerateStateToken("google")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to generate state token",
		})
		return
	}

	// Get Google auth URL
	authURL := h.oauthService.GetGoogleAuthURL(state)

	// Redirect to Google
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GoogleCallback handles Google OAuth callback
// @Summary Google OAuth callback
// @Description Handles callback from Google OAuth
// @Tags oauth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code"
// @Param state query string true "State token"
// @Success 302 {string} string "Redirect to frontend"
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/google/callback [get]
func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth?error=missing_parameters")
		return
	}

	// Validate state token
	if err := h.oauthService.ValidateStateToken(state, "google"); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth?error=invalid_state")
		return
	}

	// Handle Google callback
	userInfo, err := h.oauthService.HandleGoogleCallback(code)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth?error=google_auth_failed")
		return
	}

	// Authenticate or create user
	authResponse, err := h.authService.SocialAuth("google", userInfo)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth?error="+err.Error())
		return
	}

	// Redirect to frontend with token
	c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth/callback?token="+authResponse.Token)
}

// FacebookLogin initiates Facebook OAuth flow
// @Summary Initiate Facebook OAuth
// @Description Redirects user to Facebook login page
// @Tags oauth
// @Accept json
// @Produce json
// @Success 302 {string} string "Redirect to Facebook"
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/facebook [get]
func (h *OAuthHandler) FacebookLogin(c *gin.Context) {
	// Generate state token
	state, err := h.oauthService.GenerateStateToken("facebook")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to generate state token",
		})
		return
	}

	// Get Facebook auth URL
	authURL := h.oauthService.GetFacebookAuthURL(state)

	// Redirect to Facebook
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// FacebookCallback handles Facebook OAuth callback
// @Summary Facebook OAuth callback
// @Description Handles callback from Facebook OAuth
// @Tags oauth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code"
// @Param state query string true "State token"
// @Success 302 {string} string "Redirect to frontend"
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/facebook/callback [get]
func (h *OAuthHandler) FacebookCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth?error=missing_parameters")
		return
	}

	// Validate state token
	if err := h.oauthService.ValidateStateToken(state, "facebook"); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth?error=invalid_state")
		return
	}

	// Handle Facebook callback
	userInfo, err := h.oauthService.HandleFacebookCallback(code)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth?error=facebook_auth_failed")
		return
	}

	// Authenticate or create user
	authResponse, err := h.authService.SocialAuth("facebook", userInfo)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth?error="+err.Error())
		return
	}

	// Redirect to frontend with token
	c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth/callback?token="+authResponse.Token)
}

// GoogleAuthWithCredential handles Google authentication with ID token (client-side flow)
// @Summary Google authentication with credential
// @Description Authenticate with Google ID token from client-side
// @Tags oauth
// @Accept json
// @Produce json
// @Param request body models.GoogleAuthRequest true "Google credential"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/google/credential [post]
func (h *OAuthHandler) GoogleAuthWithCredential(c *gin.Context) {
	var req models.GoogleAuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	// Verify Google ID token
	userInfo, err := h.oauthService.VerifyGoogleIDToken(req.Credential)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid Google credential",
		})
		return
	}

	// Authenticate or create user
	authResponse, err := h.authService.SocialAuth("google", userInfo)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "email already registered with different provider" {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, models.ErrorResponse{
			Error:   "Authentication Failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// FacebookAuthWithToken handles Facebook authentication with access token (client-side flow)
// @Summary Facebook authentication with access token
// @Description Authenticate with Facebook access token from client-side
// @Tags oauth
// @Accept json
// @Produce json
// @Param request body models.FacebookAuthRequest true "Facebook access token"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/facebook/token [post]
func (h *OAuthHandler) FacebookAuthWithToken(c *gin.Context) {
	var req models.FacebookAuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	// Verify Facebook access token
	userInfo, err := h.oauthService.VerifyFacebookAccessToken(req.AccessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid Facebook access token",
		})
		return
	}

	// Authenticate or create user
	authResponse, err := h.authService.SocialAuth("facebook", userInfo)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "email already registered with different provider" {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, models.ErrorResponse{
			Error:   "Authentication Failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}
