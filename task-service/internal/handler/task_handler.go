package handler

import (
	"encoding/json"
	"errors"
	"mini-jira/task-service/internal/middleware"
	"mini-jira/task-service/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: s}
}

type createReq struct {
	ProjectID   int    `json:"projectID"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	AssigneeID  int    `json:"assigneeID"`
}

type updateReq struct {
	ProjectID   int    `json:"projectID"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	AssigneeID  int    `json:"assigneeID"`
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста
	userID := middleware.UserIDFromContext(r.Context())

	t, err := h.svc.Create(r.Context(), req.ProjectID, req.Title, req.Description, req.Status, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из контекста
	userID := middleware.UserIDFromContext(r.Context())

	tasks, err := h.svc.List(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *TaskHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	userID := middleware.UserIDFromContext(r.Context())

	t, err := h.svc.GetById(r.Context(), id, userID)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req updateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := middleware.UserIDFromContext(r.Context())

	err = h.svc.Update(r.Context(), req.Title, req.Description, req.Status, id, userID)
	if err != nil {
		if errors.Is(err, service.ErrPTaskNotFoundOrForbidden) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		if errors.Is(err, service.ErrTaskTitleEmpty) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
