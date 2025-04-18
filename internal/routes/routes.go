package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"drive/internal/handler"
	"drive/internal/service"
)

func SetupRoutes(h *handler.Handler, authService service.AuthService) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)

	r.Route("/api", func(r chi.Router) {
		// Health check route
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		r.Group(func(r chi.Router) {
			AuthRoutes(r, h)
		})

	})

	return r
}
