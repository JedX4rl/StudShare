package handler

import (
	"StudShare/internal/domain"
	"StudShare/internal/my_errors"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

// @Summary Создание объявления
// @Tags Listings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body domain.CreateListingRequest true "Данные объявления"
// @Success 201 {string} string "Создано"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Неавторизован"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/listings [post]
func (h *Handler) CreateListing(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("userID").(string)
	if !ok || id == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req domain.CreateListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if err := h.services.Listings.CreateListing(r.Context(), &domain.Listing{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		PreviewURL:  req.PreviewURL,
		City:        req.City,
		Street:      req.Street,
		Images:      req.Images,
		Owner: domain.Owner{
			ID: id,
		},
		Status: req.Status,
	}); err != nil {
		http.Error(w, "failed to create listing", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Получить объявление по ID
// @Tags Listings
// @Produce json
// @Param id query string true "ID объявления"
// @Success 200 {object} domain.Listing
// @Failure 404 {string} string "Объявление не найдено"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/listings/ [get]
func (h *Handler) GetListingByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	listing, err := h.services.Listings.GetListingByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, my_errors.ErrNotFound) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to fetch listing", http.StatusInternalServerError)
	}
	if listing == nil {
		http.Error(w, "listing not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listing)

}

// @Summary Получить все объявления (по статусу)
// @Tags Listings
// @Produce json
// @Param status query string false "Статус объявления"
// @Success 200 {array} domain.Listing
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/listings [get]
func (h *Handler) GetAllListings(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	listings, err := h.services.GetAllListings(r.Context(), status)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get listings: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listings)
}

// @Summary Получить объявления поблизости
// @Tags Listings
// @Produce json
// @Param lat query number true "Широта"
// @Param lon query number true "Долгота"
// @Param radius query number true "Радиус поиска (км)"
// @Param status query string false "Статус"
// @Success 200 {array} domain.Listing
// @Failure 400 {string} string "Некорректные координаты"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/listings/near [get]
func (h *Handler) GetListingsNear(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	latStr := q.Get("lat")
	lonStr := q.Get("lon")
	radiusStr := q.Get("radius")
	status := q.Get("status")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "invalid latitude", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "invalid longitude", http.StatusBadRequest)
		return
	}
	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		http.Error(w, "invalid radius", http.StatusBadRequest)
		return
	}

	listings, err := h.services.GetListingsNear(r.Context(), lat, lon, radius, status)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get nearby listings: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listings)
}

// @Summary Обновить объявление
// @Tags Listings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "ID объявления"
// @Param input body domain.Listing true "Объявление"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Неавторизован"
// @Failure 403 {string} string "Нет доступа"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/listings [put]
func (h *Handler) UpdateListing(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userID").(string)
	if !ok || userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var listing domain.Listing
	if err := json.NewDecoder(r.Body).Decode(&listing); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if listing.Owner.ID != userId {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	listingID := r.URL.Query().Get("id")
	if listingID == "" {
		http.Error(w, "missing listing ID", http.StatusBadRequest)
		return
	}

	listing.ID = listingID

	err := h.services.Listings.UpdateListing(r.Context(), &listing)
	if err != nil {
		http.Error(w, "failed to update listing: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Удалить объявление
// @Tags Listings
// @Produce json
// @Param id query string true "ID объявления"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 403 {string} string "Нет доступа"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/listings [delete]
func (h *Handler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	listingID := r.URL.Query().Get("id")
	if listingID == "" {
		http.Error(w, "missing listing ID", http.StatusBadRequest)
		return
	}

	if err := h.services.Listings.DeleteListing(r.Context(), listingID); err != nil {
		if errors.Is(err, my_errors.ErrPermissionDenied) {
			http.Error(w, "permission denied", http.StatusForbidden)
			return
		}
		http.Error(w, "failed to delete listing: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
