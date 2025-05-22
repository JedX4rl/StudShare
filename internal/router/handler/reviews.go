package handler

import (
	"StudShare/internal/domain"
	"StudShare/internal/my_errors"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
)

// AddReview создает отзыв на пользователя
// @Summary Оставить отзыв
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body domain.CreateReviewRequest true "Отзыв"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "invalid request body or validation failed"
// @Failure 403 {string} string "cannot leave a review for yourself"
// @Failure 409 {string} string "review already exists"
// @Failure 500 {string} string "failed to add review"
// @Router /reviews [post]
func (h *Handler) AddReview(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		http.Error(w, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Comment) < 3 || len(req.Comment) > 100 {
		http.Error(w, "comment must be from 3 to 100 symbols", http.StatusBadRequest)
		return
	}

	authorID := r.Context().Value("userID").(string)
	if authorID == req.TargetID {
		http.Error(w, "cannot leave a review for yourself", http.StatusBadRequest)
		return
	}

	if authorID == req.TargetID {
		http.Error(w, "cannot leave a review for yourself", http.StatusForbidden)
		return
	}

	review := domain.Review{
		ID:       uuid.New().String(),
		AuthorID: authorID,
		TargetID: req.TargetID,
		Rating:   float64(req.Rating),
		Comment:  req.Comment,
	}

	if err := h.services.Reviews.AddReview(r.Context(), &review); err != nil {
		if errors.Is(err, my_errors.ErrNoChanges) {
			http.Error(w, "review already exists", http.StatusConflict)
			return
		}
		http.Error(w, "failed to add review: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

// GetReviewsForUser возвращает все отзывы о пользователе
// @Summary Получить отзывы о пользователе
// @Tags reviews
// @Produce json
// @Param id query string true "User ID"
// @Success 200 {array} domain.Review
// @Failure 400 {string} string "missing user ID"
// @Failure 500 {string} string "failed to get reviews"
// @Router /reviews/user [get]
func (h *Handler) GetReviewsForUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "missing user ID", http.StatusBadRequest)
		return
	}

	reviews, err := h.services.Reviews.GetReviewsForUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get reviews: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

// DeleteReview удаляет отзыв текущего пользователя
// @Summary Удалить свой отзыв
// @Tags reviews
// @Security BearerAuth
// @Produce json
// @Param id query string true "Review ID"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "missing review ID"
// @Failure 403 {string} string "permission denied"
// @Failure 404 {string} string "review not found"
// @Failure 500 {string} string "failed to delete review"
// @Router /reviews [delete]
func (h *Handler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	reviewID := r.URL.Query().Get("id")
	if reviewID == "" {
		http.Error(w, "missing review ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(string)

	err := h.services.Reviews.DeleteReview(r.Context(), reviewID, userID)
	if err != nil {
		switch {
		case errors.Is(err, my_errors.ErrPermissionDenied):
			http.Error(w, "permission denied", http.StatusForbidden)
		case errors.Is(err, my_errors.ErrNotFound):
			http.Error(w, "review not found", http.StatusNotFound)
		default:
			http.Error(w, "failed to delete review: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
