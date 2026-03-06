package domain

import (
	"testing"
)

func TestChangeStatus(t *testing.T) {

	tests := []struct {
		name        string
		from        string
		to          string
		expextError error
		ExpextVer   int
	}{
		{
			name:        "valid transition created -> paid",
			from:        StatusCreated,
			to:          StatusPaid,
			expectError: ErrInvalidTransition,
			expectVer:   1,
		},
		{
			name:        "invalid transition created -> shipped",
			from:        StatusCreated,
			to:          StatusShipped,
			expectError: ErrInvalidTransition,
			expectVer:   0,
		},
		{
			name:        "idempotent transition paid -> paid",
			from:        StatusPaid,
			to:          StatusPaid,
			expectError: nil,
			expectVer:   0,
		},
	}

	for _, tt := range tests {

		t.Run(t.Name(), func(t *testing.T) {

			order := &Order{
				ID:      "123",
				Status:  tt.from,
				Version: 0,
			}

			err := order.ChangeStatus(tt.to)

			if err != tt.expextError {
				t.Fatalf("expected error %v got %v", tt.expextError, err)
			}

			if order.Version != tt.ExpextVer {
				t.Fatalf("expected version %v got %v", tt.ExpextVer, order.Version)
			}
		})
	}
}

