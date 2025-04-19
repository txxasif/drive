package handler

import (
	"drive/internal/model"
	"drive/internal/response"
	"drive/internal/service"
	"drive/internal/util"
	"errors"
	"net/http"
)

// OAuthHandler handles OAuth-related requests
type OAuthHandler struct {
	oauthService service.OAuthService
}

// NewOAuthHandler creates a new OAuth handler
func NewOAuthHandler(oauthService service.OAuthService) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
	}
}

// Login handles OAuth login requests
func (h *OAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.OAuthRequest

	// Validate request
	if fieldErrors := util.ValidateRequestWithFields(r, &req); fieldErrors != nil {
		response.ValidationErrorWithFields(w, fieldErrors)
		return
	}

	// Login with provider
	resp, err := h.oauthService.Login(r.Context(), req.Provider, req.Token)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOAuthToken) {
			response.Error(w, http.StatusUnauthorized, "Invalid OAuth token", err.Error())
			return
		}
		if errors.Is(err, service.ErrUnsupportedProvider) {
			response.Error(w, http.StatusBadRequest, "Unsupported OAuth provider", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to authenticate with OAuth provider", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, resp)
}
