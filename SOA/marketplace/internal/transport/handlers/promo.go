package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"marketplace/internal/service"
)

type PromoHandler struct {
	service service.PromoService
}

func NewPromoHandler(service service.PromoService) *PromoHandler {
	return &PromoHandler{service: service}
}

func (h *PromoHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	promo, err := h.service.GetByCode(r.Context(), code)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, promo)
}

func (h *PromoHandler) Validate(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	promo, err := h.service.Validate(r.Context(), code)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, promo)
}
