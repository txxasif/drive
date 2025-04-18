package service

import (
	"drive/internal/repository"
	"drive/internal/util"
)

type Services struct {
	Auth AuthService
}

func NewServices(repos repository.Repositories, jwtSvc *util.JwtService, logger *util.Logger) *Services {
	return &Services{
		Auth: NewAuthService(repos.User, jwtSvc, logger),
	}
}
