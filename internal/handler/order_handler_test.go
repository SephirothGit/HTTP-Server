package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SephirothGit/Backend-service/internal/domain"
)

type mockService struct {
	err error
}

func (m *mockService) UpdateStatus(ctx context.Context, id string, status string) error {
	return m.err
}

func TestUpdateOrderStatus(t *testing.T) {

	tests := []struct {
		name string
		body interface{}
		serviceError error
		expectedCode int
	}{
		{
			name: "invalid json",
			body: "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "order not found",
			body: map[string]string{
				"status": "paid",
			},
			serviceError: domain.ErrOrderNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name: "invalid transition",
			body: map[string]string{
				"status": "shipped",
			},
			serviceError: domain.ErrInvalidTransition,
			expectedCode: http.StatusConflict,
		},
		{
			name: "success",
			body: map[string]string{
				"status": "paid",
			},
			serviceError: nil,
			expectedCode: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc := &mockService{
				err: tt.serviceError,
			}

			h := NewOrderHandler(svc)

			var body []byte

			switch v := tt.body.(type) {

			case string:
				body = []byte(v)

			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest(
				http.MethodPut,
				"/api/v1/orders/123",
				bytes.NewBuffer(body),
			)

			rec := httptest.NewRecorder()

			h(rec, req)

			if rec.Code != tt.expectedCode {
				t.Fatalf("expected status %d got %d", tt.expectedCode, rec.Code)
			}
		})
	}
}