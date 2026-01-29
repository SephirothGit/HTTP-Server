package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func updateOrderStatusHandler(svc OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//Method check
		if r.Method != http.MethodPut {
			http.Error(w, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		//Extract id
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 4 || parts[2] == "" {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}

		id := parts[2]

		//Parse body
		type req struct {
			Status string `json:"status"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if req.Status == "" {
			http.Error(w, "status is required", http.StatusBadRequest)
			return
		}

		//Call use-case

		//Map errors
		if err := svc.UpdateStatus(r.Context(), id, req.Status); err != nil {
			switch {
			case errors.Is(err, domain.ErrOrderNotFound):
				http.Error(w, "order not found", http.StatusNotFound)
			case errors.Is(err, domain.ErrInvalidTransition):
				http.Error(w, "invalid status transition", http.StatusConflict)
			default:
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}
		//Write response
		w.WriteHeader(http.StatusNoContent)
	}
}
