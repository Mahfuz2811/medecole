package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"quizora-backend/internal/cache"
	"quizora-backend/internal/config"
	"quizora-backend/internal/models"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

const (
	stateKeyPrefix = "oauth_state:"
	stateTTL       = 5 * time.Minute
)

// OAuthService handles OAuth authentication logic
type OAuthService struct {
	googleConfig   *oauth2.Config
	facebookConfig *oauth2.Config
	cache          cache.CacheInterface
	frontendURL    string
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(cfg *config.Config, cacheInstance cache.CacheInterface) *OAuthService {
	googleConfig := &oauth2.Config{
		ClientID:     cfg.OAuth.Google.ClientID,
		ClientSecret: cfg.OAuth.Google.ClientSecret,
		RedirectURL:  cfg.OAuth.Google.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	facebookConfig := &oauth2.Config{
		ClientID:     cfg.OAuth.Facebook.AppID,
		ClientSecret: cfg.OAuth.Facebook.AppSecret,
		RedirectURL:  cfg.OAuth.Facebook.RedirectURL,
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     facebook.Endpoint,
	}

	return &OAuthService{
		googleConfig:   googleConfig,
		facebookConfig: facebookConfig,
		cache:          cacheInstance,
		frontendURL:    cfg.CORS.FrontendURL,
	}
}

// GenerateStateToken generates a secure state token for OAuth
func (s *OAuthService) GenerateStateToken(provider string) (string, error) {
	state := uuid.New().String()
	stateData := models.OAuthState{
		State:     state,
		Provider:  provider,
		CreatedAt: time.Now().Unix(),
	}

	// Store state in cache with TTL
	key := stateKeyPrefix + state
	if err := s.cache.Set(key, stateData, stateTTL); err != nil {
		return "", fmt.Errorf("failed to store state: %w", err)
	}

	return state, nil
}

// ValidateStateToken validates the OAuth state token
func (s *OAuthService) ValidateStateToken(state, expectedProvider string) error {
	key := stateKeyPrefix + state
	var stateData models.OAuthState

	if err := s.cache.Get(key, &stateData); err != nil {
		return errors.New("invalid or expired state token")
	}

	if stateData.Provider != expectedProvider {
		return errors.New("provider mismatch")
	}

	// Delete state after validation (one-time use)
	s.cache.Delete(key)

	return nil
}

// GetGoogleAuthURL generates Google OAuth authorization URL
func (s *OAuthService) GetGoogleAuthURL(state string) string {
	return s.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// HandleGoogleCallback handles Google OAuth callback
func (s *OAuthService) HandleGoogleCallback(code string) (*models.SocialUserInfo, error) {
	// Exchange code for token
	token, err := s.googleConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Fetch user info from Google
	client := s.googleConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		VerifiedEmail bool   `json:"verified_email"`
	}

	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &models.SocialUserInfo{
		ProviderUserID: googleUser.ID,
		Email:          googleUser.Email,
		Name:           googleUser.Name,
		ProfilePicture: googleUser.Picture,
		EmailVerified:  googleUser.VerifiedEmail,
	}, nil
}

// GetFacebookAuthURL generates Facebook OAuth authorization URL
func (s *OAuthService) GetFacebookAuthURL(state string) string {
	return s.facebookConfig.AuthCodeURL(state)
}

// HandleFacebookCallback handles Facebook OAuth callback
func (s *OAuthService) HandleFacebookCallback(code string) (*models.SocialUserInfo, error) {
	// Exchange code for token
	token, err := s.facebookConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Fetch user info from Facebook
	client := s.facebookConfig.Client(context.Background(), token)
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email,picture.type(large)")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var facebookUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}

	if err := json.Unmarshal(body, &facebookUser); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// Facebook doesn't provide email verification status via Graph API
	// We'll assume it's verified if email is present
	emailVerified := facebookUser.Email != ""

	return &models.SocialUserInfo{
		ProviderUserID: facebookUser.ID,
		Email:          facebookUser.Email,
		Name:           facebookUser.Name,
		ProfilePicture: facebookUser.Picture.Data.URL,
		EmailVerified:  emailVerified,
	}, nil
}

// VerifyGoogleIDToken verifies Google ID token (for client-side OAuth)
func (s *OAuthService) VerifyGoogleIDToken(idToken string) (*models.SocialUserInfo, error) {
	// Make request to Google's tokeninfo endpoint
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid token")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tokenInfo struct {
		Sub           string `json:"sub"`            // User ID
		Email         string `json:"email"`          // Email
		Name          string `json:"name"`           // Full name
		Picture       string `json:"picture"`        // Profile picture
		EmailVerified string `json:"email_verified"` // Email verified (string: "true" or "false")
		Aud           string `json:"aud"`            // Audience (should match client ID)
	}

	if err := json.Unmarshal(body, &tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to parse token info: %w", err)
	}

	// Verify audience matches our client ID
	if tokenInfo.Aud != s.googleConfig.ClientID {
		return nil, errors.New("invalid audience")
	}

	return &models.SocialUserInfo{
		ProviderUserID: tokenInfo.Sub,
		Email:          tokenInfo.Email,
		Name:           tokenInfo.Name,
		ProfilePicture: tokenInfo.Picture,
		EmailVerified:  tokenInfo.EmailVerified == "true",
	}, nil
}

// VerifyFacebookAccessToken verifies Facebook access token
func (s *OAuthService) VerifyFacebookAccessToken(accessToken string) (*models.SocialUserInfo, error) {
	// Verify token with Facebook
	verifyURL := fmt.Sprintf("https://graph.facebook.com/debug_token?input_token=%s&access_token=%s|%s",
		accessToken, s.facebookConfig.ClientID, s.facebookConfig.ClientSecret)

	resp, err := http.Get(verifyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid token")
	}

	// Fetch user info
	userResp, err := http.Get("https://graph.facebook.com/me?fields=id,name,email,picture.type(large)&access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer userResp.Body.Close()

	if userResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info, status: %d", userResp.StatusCode)
	}

	body, err := io.ReadAll(userResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var facebookUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}

	if err := json.Unmarshal(body, &facebookUser); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	emailVerified := facebookUser.Email != ""

	return &models.SocialUserInfo{
		ProviderUserID: facebookUser.ID,
		Email:          facebookUser.Email,
		Name:           facebookUser.Name,
		ProfilePicture: facebookUser.Picture.Data.URL,
		EmailVerified:  emailVerified,
	}, nil
}
