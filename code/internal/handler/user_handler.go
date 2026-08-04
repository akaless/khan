package handler

import (
	"encoding/json"
	"net/http"

	"khan/internal/models"
	"khan/internal/repository"
	"khan/internal/service"
)

// UserHandler serves /api/users/*
type UserHandler struct {
	users *service.UserService
	repo  *repository.UserRepo
	auth  *service.AuthService
}

func NewUserHandler(users *service.UserService, repo *repository.UserRepo, auth *service.AuthService) *UserHandler {
	return &UserHandler{users: users, repo: repo, auth: auth}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.ListVisible()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "خطای داخلی")
		return
	}
	sanitized := make([]*models.User, 0, len(users))
	for i := range users {
		sanitized = append(sanitized, users[i].Sanitized())
	}
	writeJSON(w, http.StatusOK, sanitized)
}

func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		h.List(w, r)
		return
	}
	users, err := h.repo.Search(q)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "خطای داخلی")
		return
	}
	sanitized := make([]*models.User, 0, len(users))
	for i := range users {
		sanitized = append(sanitized, users[i].Sanitized())
	}
	writeJSON(w, http.StatusOK, sanitized)
}

type createUserReq struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
	Role        string `json:"role"`
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	actor := userFrom(r)
	var req createUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	u, err := h.users.CreateUser(actor, req.Username, req.DisplayName, req.Password, req.Role)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, u)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	actor := userFrom(r)
	id, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	if err := h.users.DeleteUser(actor, id); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

type resetPwdReq struct {
	NewPassword string `json:"new_password"`
}

func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	actor := userFrom(r)
	id, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	var req resetPwdReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	if err := h.users.ResetPassword(actor, id, req.NewPassword); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *UserHandler) ToggleActive(w http.ResponseWriter, r *http.Request) {
	actor := userFrom(r)
	id, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	if err := h.users.ToggleActive(actor, id); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

type setRoleReq struct {
	Role string `json:"role"`
}

func (h *UserHandler) SetRole(w http.ResponseWriter, r *http.Request) {
	actor := userFrom(r)
	id, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	var req setRoleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	if err := h.users.SetRole(actor, id, req.Role); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
