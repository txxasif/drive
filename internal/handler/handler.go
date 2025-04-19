package handler

import "drive/internal/service"

type Handler struct {
	UserHandler  *UserHandler
	OAuthHandler *OAuthHandler
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		UserHandler:  NewUserHandler(services.Auth),
		OAuthHandler: NewOAuthHandler(services.OAuth),
	}
}
