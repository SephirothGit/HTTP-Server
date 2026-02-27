package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SephirothGit/Backend-service/internal/domain"
	"github.com/SephirothGit/Backend-service/internal/service"
	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	svc service.OrderService
}

func NewOrderHandler(svc service.OrderService) http.HandlerFunc {
	h := OrderHandler{svc: svc}
	return h.updateOrderStatus
}

func (h *OrderHandler) updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateStatus(r.Context(), id, req.Status); err != nil {
		switch {
		case errors.Is(err, domain.ErrOrderNotFound):
			http.Error(w, "order not found", http.StatusNotFound)
		case errors.Is(err, domain.ErrInvalidTransition):
			http.Error(w, "invalid transition", http.StatusConflict)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
