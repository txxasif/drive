package service

import (
	"context"
	"drive/internal/model"
	"drive/internal/repository"
	"drive/internal/util"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

var (
	ErrEmailAlreadyExists = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUnauthorized       = errors.New("unauthorized")
)

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error)
	GetUserByToken(ctx context.Context, token string) (*model.User, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*RefreshResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwtSvc   *util.JwtService
	logger   *util.Logger
}

func NewAuthService(userRepo repository.UserRepository, jwtSvc *util.JwtService, logger *util.Logger) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtSvc:   jwtSvc,
		logger:   logger,
	}
}

func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	logger := s.logger.WithEmail(req.Email)

	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)

	if err != nil {
		logger.Error("Error checking existing user", zap.Error(err))
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}

	if existingUser != nil {
		logger.Warn("User already exists")
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		logger.Error("Error hashing password", util.WithError(err))
		return nil, err
	}

	user := &model.User{
		Email:        req.Email,
		Password:     hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		StorageUsed:  0,
		StorageLimit: 15000,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.Error("Error creating user", util.WithError(err))
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	token, err := s.jwtSvc.GenerateAccessToken(user.ID)
	if err != nil {
		logger.Error("Error generating access token", util.WithError(err))
		return nil, fmt.Errorf("error generating token: %w", err)
	}
	refreshToken, err := s.jwtSvc.GenerateRefreshToken(user.ID)
	if err != nil {
		logger.Error("Error generating refresh token", util.WithError(err))
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	logger.Info("User registered successfully")

	user.Password = ""

	return &model.AuthResponse{
		User:         user.ToResponse(),
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	logger := s.logger.WithEmail(req.Email)

	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		logger.Error("Error finding user", util.WithError(err))
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	if user == nil {
		logger.Warn("User not found")
		return nil, ErrInvalidCredentials
	}

	if err := util.CheckPassword(user.Password, req.Password); err != nil {
		logger.Warn("Invalid password")
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.jwtSvc.GenerateAccessToken(user.ID)
	if err != nil {
		logger.Error("Error generating access token", util.WithError(err))
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	refreshToken, err := s.jwtSvc.GenerateRefreshToken(user.ID)
	if err != nil {
		logger.Error("Error generating refresh token", util.WithError(err))
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	logger.Info("User logged in successfully")

	// Clear sensitive data
	user.Password = ""

	return &model.AuthResponse{
		User:         user.ToResponse(),
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) GetUserByToken(ctx context.Context, token string) (*model.User, error) {
	userID, tokenType, err := s.jwtSvc.ValidateToken(token)
	if err != nil {
		s.logger.Error("Error validating token", util.WithError(err))
		return nil, err
	}

	if tokenType != util.AccessToken {
		s.logger.Warn("Invalid token type", util.WithUserID(userID))
		return nil, ErrUnauthorized
	}

	user, err := s.userRepo.GetById(ctx, userID)
	if err != nil {
		s.logger.Error("Error getting user by ID", util.WithUserID(userID), util.WithError(err))
		return nil, err
	}
	if user == nil {
		s.logger.Warn("User not found", util.WithUserID(userID))
		return nil, errors.New("user not found")
	}

	return user, nil
}

// RefreshTokens refreshes the access and refresh tokens
func (s *authService) RefreshTokens(ctx context.Context, refreshToken string) (*RefreshResponse, error) {
	userID, tokenType, err := s.jwtSvc.ValidateToken(refreshToken)
	if err != nil {
		s.logger.Error("Error validating refresh token", util.WithError(err))
		return nil, err
	}

	// Verify it's a refresh token
	if tokenType != util.RefreshToken {
		s.logger.Warn("Invalid token type for refresh", util.WithUserID(userID))
		return nil, ErrUnauthorized
	}

	// Verify user exists
	user, err := s.userRepo.GetById(ctx, userID)
	if err != nil || user == nil {
		s.logger.Error("Error getting user for refresh", util.WithUserID(userID), util.WithError(err))
		return nil, ErrUnauthorized
	}

	// Generate new tokens
	newAccessToken, err := s.jwtSvc.GenerateAccessToken(user.ID)
	if err != nil {
		s.logger.Error("Error generating new access token", util.WithUserID(userID), util.WithError(err))
		return nil, err
	}

	newRefreshToken, err := s.jwtSvc.GenerateRefreshToken(user.ID)
	if err != nil {
		s.logger.Error("Error generating new refresh token", util.WithUserID(userID), util.WithError(err))
		return nil, err
	}

	s.logger.Info("Tokens refreshed successfully", util.WithUserID(userID))

	return &RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
