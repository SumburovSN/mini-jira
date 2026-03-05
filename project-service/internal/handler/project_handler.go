package handler

import (
	"encoding/json"
	"errors"
	"mini-jira/project-service/internal/middleware"
	"mini-jira/project-service/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProjectHandler struct {
	svc *service.ProjectService
}

func NewProjectHandler(s *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: s}
}

type createReq struct {
	Name string `json:"name"`
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста
	userID := middleware.UserIDFromContext(r.Context())

	// Создаем проект для этого пользователя
	p, err := h.svc.Create(r.Context(), req.Name, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем проект в ответе
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из контекста
	userID := middleware.UserIDFromContext(r.Context())

	// Получаем проекты для пользователя
	projects, err := h.svc.List(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем список проектов в ответе
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	userID := middleware.UserIDFromContext(r.Context())

	err = h.svc.Delete(r.Context(), id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	userID := middleware.UserIDFromContext(r.Context())

	p, err := h.svc.GetById(r.Context(), id, userID)
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
