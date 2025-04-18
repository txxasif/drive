package handler

import (
	"drive/internal/model"
	"drive/internal/response"
	"drive/internal/service"
	"drive/internal/util"

	"errors"
	"net/http"
)

type UserHandler struct {
	authService service.AuthService
}

func NewUserHandler(authService service.AuthService) *UserHandler {
	return &UserHandler{
		authService: authService,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var userDTO model.RegisterRequest

	// Validate request with field errors
	if fieldErrors := util.ValidateRequestWithFields(r, &userDTO); fieldErrors != nil {
		response.ValidationErrorWithFields(w, fieldErrors)
		return
	}

	user, err := h.authService.Register(r.Context(), &userDTO)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			response.Error(w, http.StatusConflict, "User with this email already exists", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to register user", err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var userDTO model.LoginRequest

	// Validate request with field errors
	if fieldErrors := util.ValidateRequestWithFields(r, &userDTO); fieldErrors != nil {
		response.ValidationErrorWithFields(w, fieldErrors)
		return
	}

	token, err := h.authService.Login(r.Context(), &userDTO)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, "Invalid email or password", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to login", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"user":          token.User,
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}
