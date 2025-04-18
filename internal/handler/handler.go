package handler

import "drive/internal/service"

type Handler struct {
	UserHandler *UserHandler
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		UserHandler: NewUserHandler(services.Auth),
	}
}
