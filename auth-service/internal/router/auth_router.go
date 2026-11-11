package router

import (
	"mini-jira/auth-service/internal/handler"

	"github.com/go-chi/chi/v5"
)

func RegisterPublicAuthRoutes(r chi.Router, h *handler.AuthHandler) {
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
}

func RegisterPrivateAuthRoutes(r chi.Router, h *handler.AuthHandler) {
	r.Get("/me", h.Me)
}
