package router

import (
	"mini-jira/project-service/internal/handler"

	"github.com/go-chi/chi/v5"
)

func RegisterMemberRoutes(r chi.Router, h *handler.MemberHandler) {

	r.Post("/projects/{id}/members", h.Create)
	r.Get("/projects/{id}/members", h.List)

	r.Get("/members/{id}", h.GetById)
	r.Patch("/members/{id}/role", h.UpdateRole)
	r.Delete("/members/{id}", h.Delete)
}
