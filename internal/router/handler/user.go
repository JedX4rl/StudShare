package handler

import (
	"StudShare/internal/domain"
	"StudShare/internal/my_errors"
	"StudShare/pkg"
	"encoding/json"
	"errors"
	"net/http"
)

// GetProfile получает профиль текущего пользователя
// @Summary Получить свой профиль
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.User
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Router /users [get]
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("userID").(string)
	if !ok || id == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.services.Users.GetProfile(r.Context(), id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetProfileByID получает профиль пользователя по ID
// @Summary Получить профиль пользователя по ID
// @Tags users
// @Produce json
// @Param id query string true "User ID"
// @Success 200 {object} domain.User
// @Failure 404 {string} string "User not found"
// @Router /users/profile [get]
func (h *Handler) GetProfileByID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")

	user, err := h.services.Users.GetProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateProfile обновляет профиль текущего пользователя
// @Summary Обновить свой профиль
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body domain.UpdateUserRequest true "Данные для обновления"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Invalid request or no changes"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Error updating user"
// @Router /users/update [put]
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("userID").(string)
	if !ok || id == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req domain.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, pkg.FormatValidationError(err), http.StatusBadRequest)
		return
	}

	user := &domain.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Surname:  req.Surname,
		Phone:    req.Phone,
		ID:       id,
	}

	if err := h.services.Users.UpdateProfile(r.Context(), user); err != nil {
		if errors.Is(err, my_errors.ErrNoChanges) {
			http.Error(w, "No changes made to profile", http.StatusBadRequest)
			return
		}
		if errors.Is(err, my_errors.ErrNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
