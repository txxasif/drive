package model

import "time"

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,strong_password"`
	Username  string `json:"username" validate:"required,min=3,max=20,alphanum"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// OAuthRequest contains the OAuth provider token
type OAuthRequest struct {
	Token    string `json:"token" validate:"required"`
	Provider string `json:"provider" validate:"required,oneof=google facebook"`
}

// OAuthUserInfo represents user information from OAuth providers
type OAuthUserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Picture   string `json:"picture,omitempty"`
}

type UserResponse struct {
	ID           uint         `json:"id"`
	Email        string       `json:"email"`
	Username     string       `json:"username"`
	FirstName    string       `json:"first_name"`
	LastName     string       `json:"last_name"`
	StorageUsed  float64      `json:"storage_used"`
	StorageLimit float64      `json:"storage_limit"`
	Provider     AuthProvider `json:"provider"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:           u.ID,
		Email:        u.Email,
		Username:     u.Username,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		StorageUsed:  u.StorageUsed,
		StorageLimit: u.StorageLimit,
		Provider:     u.Provider,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
}
