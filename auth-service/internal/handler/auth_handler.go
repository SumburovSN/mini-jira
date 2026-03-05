package handler

import (
	"encoding/json"
	"net/http"

	"mini-jira/auth-service/internal/service"
)

type AuthHandler struct {
	s *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{s: s}
}

type creds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var c creds
	json.NewDecoder(r.Body).Decode(&c)

	err := h.s.Register(r.Context(), c.Email, c.Password)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Write([]byte("ok"))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var c creds
	json.NewDecoder(r.Body).Decode(&c)

	token, err := h.s.Login(r.Context(), c.Email, c.Password)
	if err != nil {
		http.Error(w, "invalid creds", 401)
		return
	}
	w.Write([]byte(token))
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)

	u, err := h.s.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}

	w.Write([]byte(u.Email))
}
