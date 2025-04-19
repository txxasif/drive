package service

import (
	"context"
	"drive/internal/model"
	"drive/internal/repository"
	"drive/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	// ErrInvalidOAuthToken indicates that the token provided for OAuth is invalid
	ErrInvalidOAuthToken = errors.New("invalid oauth token")
	// ErrUnsupportedProvider indicates an unsupported OAuth provider
	ErrUnsupportedProvider = errors.New("unsupported oauth provider")
)

// OAuthProvider defines the interface for OAuth providers
type OAuthProvider interface {
	// GetUserInfo fetches user information from the OAuth provider
	GetUserInfo(ctx context.Context, token string) (*model.OAuthUserInfo, error)
	// GetProviderName returns the name of the provider
	GetProviderName() string
}

// OAuthService interface defines OAuth authentication operations
type OAuthService interface {
	// Login authenticates a user using an OAuth provider
	Login(ctx context.Context, provider string, token string) (*model.AuthResponse, error)
	// GetProvider returns the appropriate provider implementation
	GetProvider(provider string) (OAuthProvider, error)
}

// oauthService implements OAuthService
type oauthService struct {
	userRepo       repository.UserRepository
	jwtSvc         *util.JwtService
	googleConfig   *GoogleOAuthConfig
	facebookConfig *FacebookOAuthConfig
	logger         *util.Logger
	authService    AuthService
}

// GoogleOAuthConfig holds configuration for Google OAuth
type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
}

// FacebookOAuthConfig holds configuration for Facebook OAuth
type FacebookOAuthConfig struct {
	AppID     string
	AppSecret string
}

// NewOAuthService creates a new OAuthService instance
func NewOAuthService(
	userRepo repository.UserRepository,
	jwtSvc *util.JwtService,
	googleConfig *GoogleOAuthConfig,
	facebookConfig *FacebookOAuthConfig,
	logger *util.Logger,
	authService AuthService,
) OAuthService {
	return &oauthService{
		userRepo:       userRepo,
		jwtSvc:         jwtSvc,
		googleConfig:   googleConfig,
		facebookConfig: facebookConfig,
		logger:         logger,
		authService:    authService,
	}
}

// Login authenticates a user using an OAuth provider
func (s *oauthService) Login(ctx context.Context, providerName string, token string) (*model.AuthResponse, error) {
	provider, err := s.GetProvider(providerName)
	if err != nil {
		return nil, err
	}

	userInfo, err := provider.GetUserInfo(ctx, token)
	if err != nil {
		s.logger.Error("Error fetching user info from OAuth provider",
			zap.String("provider", providerName),
			util.WithError(err))
		return nil, err
	}

	if userInfo.Email == "" {
		s.logger.Error("OAuth provider did not return an email",
			zap.String("provider", providerName))
		return nil, errors.New("oauth provider did not return an email")
	}

	// Check if user exists by email
	user, err := s.userRepo.FindByEmail(ctx, userInfo.Email)
	if err != nil {
		s.logger.Error("Error finding user by email",
			zap.String("email", userInfo.Email),
			util.WithError(err))
		return nil, err
	}

	// Get provider enum from string
	var authProvider model.AuthProvider
	switch strings.ToLower(providerName) {
	case "google":
		authProvider = model.GoogleAuth
	case "facebook":
		authProvider = model.FacebookAuth
	default:
		authProvider = model.LocalAuth
	}

	if user == nil {
		// Create a new user
		username := generateUsername(userInfo.Email)

		user = &model.User{
			Email:        userInfo.Email,
			Username:     username,
			FirstName:    userInfo.FirstName,
			LastName:     userInfo.LastName,
			Password:     uuid.NewString(), // Random password for OAuth users
			Provider:     authProvider,
			ProviderId:   userInfo.ID,
			StorageUsed:  0,
			StorageLimit: 15000,
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			s.logger.Error("Error creating user from OAuth",
				zap.String("email", userInfo.Email),
				util.WithError(err))
			return nil, err
		}
	} else if user.Provider == model.LocalAuth {
		// Update existing user with OAuth info if they were using local auth
		user.Provider = authProvider
		user.ProviderId = userInfo.ID

		if err := s.userRepo.Update(ctx, user); err != nil {
			s.logger.Error("Error updating user with OAuth info",
				util.WithUserID(user.ID),
				util.WithError(err))
			return nil, err
		}
	}

	// Generate tokens
	accessToken, err := s.jwtSvc.GenerateAccessToken(user.ID)
	if err != nil {
		s.logger.Error("Error generating access token",
			util.WithUserID(user.ID),
			util.WithError(err))
		return nil, err
	}

	refreshToken, err := s.jwtSvc.GenerateRefreshToken(user.ID)
	if err != nil {
		s.logger.Error("Error generating refresh token",
			util.WithUserID(user.ID),
			util.WithError(err))
		return nil, err
	}

	s.logger.Info("User logged in with OAuth successfully",
		util.WithUserID(user.ID),
		zap.String("provider", providerName))

	return &model.AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GetProvider returns the appropriate provider implementation
func (s *oauthService) GetProvider(provider string) (OAuthProvider, error) {
	switch strings.ToLower(provider) {
	case "google":
		return NewGoogleOAuthProvider(s.googleConfig), nil
	case "facebook":
		return NewFacebookOAuthProvider(s.facebookConfig), nil
	default:
		return nil, ErrUnsupportedProvider
	}
}

// generateUsername creates a username from email
func generateUsername(email string) string {
	parts := strings.Split(email, "@")
	base := parts[0]
	// Add random suffix to ensure uniqueness
	suffix := fmt.Sprintf("%d", time.Now().UnixNano()%10000)
	return base + suffix
}

// GoogleOAuthProvider implements OAuthProvider for Google
type GoogleOAuthProvider struct {
	config *GoogleOAuthConfig
}

// NewGoogleOAuthProvider creates a new Google OAuth provider
func NewGoogleOAuthProvider(config *GoogleOAuthConfig) OAuthProvider {
	return &GoogleOAuthProvider{
		config: config,
	}
}

// GetProviderName returns the provider name
func (p *GoogleOAuthProvider) GetProviderName() string {
	return "google"
}

// GetUserInfo fetches user information from Google
func (p *GoogleOAuthProvider) GetUserInfo(ctx context.Context, token string) (*model.OAuthUserInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("invalid response from Google: %s - %s",
			resp.Status, string(body))
	}

	var result struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &model.OAuthUserInfo{
		ID:        result.Sub,
		Email:     result.Email,
		FirstName: result.GivenName,
		LastName:  result.FamilyName,
		Picture:   result.Picture,
	}, nil
}

// FacebookOAuthProvider implements OAuthProvider for Facebook
type FacebookOAuthProvider struct {
	config *FacebookOAuthConfig
}

// NewFacebookOAuthProvider creates a new Facebook OAuth provider
func NewFacebookOAuthProvider(config *FacebookOAuthConfig) OAuthProvider {
	return &FacebookOAuthProvider{
		config: config,
	}
}

// GetProviderName returns the provider name
func (p *FacebookOAuthProvider) GetProviderName() string {
	return "facebook"
}

// GetUserInfo fetches user information from Facebook
func (p *FacebookOAuthProvider) GetUserInfo(ctx context.Context, token string) (*model.OAuthUserInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	// Create URL with fields we need
	u, err := url.Parse("https://graph.facebook.com/v18.0/me")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("fields", "id,email,first_name,last_name,picture")
	q.Add("access_token", token)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("invalid response from Facebook: %s - %s",
			resp.Status, string(body))
	}

	var result struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Picture   struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &model.OAuthUserInfo{
		ID:        result.ID,
		Email:     result.Email,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Picture:   result.Picture.Data.URL,
	}, nil
}
