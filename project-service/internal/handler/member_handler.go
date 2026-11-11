package handler

import (
	"encoding/json"
	"errors"
	"mini-jira/project-service/internal/middleware"
	"mini-jira/project-service/internal/model"
	"mini-jira/project-service/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MemberHandler struct {
	svc *service.MemberService
}

func NewMemberHandler(s *service.MemberService) *MemberHandler {
	return &MemberHandler{svc: s}
}

type createMemReq struct {
	ProjectID int    `json:"project_id"`
	UserID    int    `json:"user_id"`
	Role      string `json:"role"`
}

type updateMemReq struct {
	ProjectID int    `json:"project_id"`
	UserID    int    `json:"user_id"`
	Role      string `json:"role"`
}

func (h *MemberHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createMemReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем userID из контекста
	userID := middleware.UserIDFromContext(r.Context())
	role, err := h.svc.GetUserRole(r.Context(), req.ProjectID, userID)
	if role != model.RoleOwner {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Создаем нового участника проекта
	m, err := h.svc.Create(r.Context(), req.ProjectID, req.UserID, model.Role(req.Role))
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Возвращаем участника проекта в ответе
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func (h *MemberHandler) List(w http.ResponseWriter, r *http.Request) {
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

func (h *MemberHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// userID := middleware.UserIDFromContext(r.Context())

	err = h.svc.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MemberHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// userID := middleware.UserIDFromContext(r.Context())

	p, err := h.svc.GetById(r.Context(), id)
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

func (h *MemberHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req updateMemReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// userID := middleware.UserIDFromContext(r.Context())

	err = h.svc.UpdateRole(r.Context(), req.UserID, model.Role(req.Role))
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFoundOrForbidden) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		if errors.Is(err, service.ErrProjectNameEmpty) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
