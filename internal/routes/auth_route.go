package routes

import (
	"drive/internal/handler"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router, handler *handler.Handler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", handler.UserHandler.Register)
		r.Post("/login", handler.UserHandler.Login)
	})
}
