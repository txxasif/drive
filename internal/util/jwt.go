package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken is returned when a token is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenExpired is returned when a token has expired
	ErrTokenExpired = errors.New("token expired")
	// ErrInvalidSigningMethod is returned when token signing method is invalid
	ErrInvalidSigningMethod = errors.New("invalid signing method")
)

// TokenType represents the type of JWT token
type TokenType string

const (
	// AccessToken is a short-lived token for API access
	AccessToken TokenType = "access"
	// RefreshToken is a long-lived token for refreshing access tokens
	RefreshToken TokenType = "refresh"
)

// Claims represents the claims in a JWT token
type Claims struct {
	UserID uint      `json:"user_id"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

// TokenService defines the interface for JWT operations
type TokenService interface {
	// GenerateAccessToken generates a new access token
	GenerateAccessToken(userID uint) (string, error)
	// GenerateRefreshToken generates a new refresh token
	GenerateRefreshToken(userID uint) (string, error)
	// ValidateToken validates a token and returns the user ID
	ValidateToken(tokenString string) (uint, TokenType, error)
	// ParseToken parses a token without validation
	ParseToken(tokenString string) (*Claims, error)
}

// Service implements the TokenService interface
type JwtService struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// ServiceConfig holds configuration for the JWT service
type ServiceConfig struct {
	SecretKey     string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

// NewServiceWithConfig creates a new TokenService with custom configuration
func NewJwtService(config ServiceConfig) *JwtService {
	return &JwtService{
		secretKey:     []byte(config.SecretKey),
		accessExpiry:  config.AccessExpiry,
		refreshExpiry: config.RefreshExpiry,
	}
}

// GenerateAccessToken generates a new access token for a user
func (s *JwtService) GenerateAccessToken(userID uint) (string, error) {
	return s.generateToken(userID, AccessToken, s.accessExpiry)
}

// GenerateRefreshToken generates a new refresh token for a user
func (s *JwtService) GenerateRefreshToken(userID uint) (string, error) {
	return s.generateToken(userID, RefreshToken, s.refreshExpiry)
}

// generateToken generates a token with the given parameters
func (s *JwtService) generateToken(userID uint, tokenType TokenType, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

// ValidateToken validates a token and returns the user ID and token type
func (s *JwtService) ValidateToken(tokenString string) (uint, TokenType, error) {
	claims, err := s.ParseToken(tokenString)
	if err != nil {
		return 0, "", err
	}

	return claims.UserID, claims.Type, nil
}

// ParseToken parses a token without validation
func (s *JwtService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
