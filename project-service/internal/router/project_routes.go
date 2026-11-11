package router

import (
	"mini-jira/project-service/internal/handler"

	"github.com/go-chi/chi/v5"
)

func RegisterProjectRoutes(r chi.Router, h *handler.ProjectHandler) {
	r.Post("/projects", h.Create)
	r.Get("/projects", h.List)
	r.Delete("/projects/{id}", h.Delete)
	r.Get("/projects/{id}", h.GetById)
	r.Patch("/projects/{id}", h.Update)
}
