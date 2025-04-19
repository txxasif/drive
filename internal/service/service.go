package service

import (
	"drive/internal/config"
	"drive/internal/repository"
	"drive/internal/util"
)

type Services struct {
	Auth  AuthService
	OAuth OAuthService
}

func NewServices(repos repository.Repositories, jwtSvc *util.JwtService, logger *util.Logger, cfg *config.Config) *Services {
	authService := NewAuthService(repos.User, jwtSvc, logger)

	// Create OAuth configs
	googleConfig := &GoogleOAuthConfig{
		ClientID:     cfg.OAuth.GoogleClientID,
		ClientSecret: cfg.OAuth.GoogleClientSecret,
	}

	facebookConfig := &FacebookOAuthConfig{
		AppID:     cfg.OAuth.FacebookAppID,
		AppSecret: cfg.OAuth.FacebookAppSecret,
	}

	return &Services{
		Auth:  authService,
		OAuth: NewOAuthService(repos.User, jwtSvc, googleConfig, facebookConfig, logger, authService),
	}
}
