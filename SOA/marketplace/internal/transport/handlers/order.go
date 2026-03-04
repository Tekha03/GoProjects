package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"marketplace/internal/domain"
	"marketplace/internal/transport/middleware"
	"marketplace/internal/service"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

type createOrderRequest struct {
	Items     []domain.OrderItem `json:"items"`
	PromoCode *string            `json:"promo_code"`
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, domain.ErrInvalidInput)
		return
	}

	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, domain.ErrUnauthorized)
		return
	}

	order, err := h.service.CreateOrder(
		r.Context(),
		claims.UserID,
		req.Items,
		req.PromoCode,
	)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	order, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Status domain.OrderStatus `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, domain.ErrInvalidInput)
		return
	}

	if err := h.service.ChangeStatus(r.Context(), id, req.Status); err != nil {
		writeDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
