package handler

import (
	"StudShare/internal/domain"
	"encoding/json"
	"net/http"
)

// CreateCreateDraft godoc
// @Summary      Создать черновик объявления
// @Tags         drafts
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        draft  body      domain.Draft  true  "Данные черновика"
// @Success      201    {string}  string        "created"
// @Failure      400    {string}  string
// @Failure      500    {string}  string
// @Router       /drafts/ [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("userID").(string)
	if id == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var draft domain.Draft
	if err := json.NewDecoder(r.Body).Decode(&draft); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	draft.OwnerID = id

	if err := h.services.Drafts.Create(r.Context(), &draft); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetAllDrafts godoc
// @Summary      Получить все черновики пользователя
// @Tags         drafts
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   domain.Draft
// @Failure      400  {string}  string
// @Failure      500  {string}  string
// @Router       /drafts/all [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("userID").(string)
	if id == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	drafts, err := h.services.Drafts.GetAllDrafts(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drafts)
}

// GetDraftByID godoc
// @Summary      Получить черновик по ID
// @Tags         drafts
// @Security     BearerAuth
// @Produce      json
// @Param        id   query     string  true  "ID черновика"
// @Success      200  {object}  domain.Draft
// @Failure      404  {string}  string
// @Router       /drafts/ [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	draft, err := h.services.Drafts.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(draft)
}

// DeleteDraftByID godoc
// @Summary      Удалить черновик по ID
// @Tags         drafts
// @Security     BearerAuth
// @Produce      json
// @Param        id   query     string  true  "ID черновика"
// @Success      200  {string}  string  "deleted"
// @Failure      500  {string}  string
// @Router       /drafts/ [delete]
func (h *Handler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if err := h.services.Drafts.DeleteByID(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
