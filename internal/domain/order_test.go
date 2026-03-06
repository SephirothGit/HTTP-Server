package domain

import (
	"testing"
)

func TestChangeStatus(t *testing.T) {

	tests := []struct {
		name        string
		from        string
		to          string
		expectError error
		expectVer   int
	}{
		{
			name:        "valid transition created -> paid",
			from:        StatusCreated,
			to:          StatusPaid,
			expectError: nil,
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

		t.Run(tt.name, func(t *testing.T) {

			order := &Order{
				ID:      "123",
				Status:  tt.from,
				Version: 0,
			}

			err := order.ChangeStatus(tt.to)

			if err != tt.expectError {
				t.Fatalf("expected error %v got %v", tt.expectError, err)
			}

			if order.Version != tt.expectVer {
				t.Fatalf("expected version %v got %v", tt.expectVer, order.Version)
			}
		})
	}
}

func TestChangeStatus_CreatesEvent(t *testing.T) {

	order := &Order{
		ID:     "123",
		Status: StatusCreated,
	}

	err := order.ChangeStatus(StatusPaid)

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	events := order.PullEvents()

	if len(events) != 1 {
		t.Fatalf("expected 1 event got %v", len(events))
	}

	event := events[0].(OrderStatusChanged)

	if event.From != StatusCreated {
		t.Fatalf("expected from created got %s", event.From)
	}

	if event.To != StatusPaid {
		t.Fatalf("expected to paid got %s", event.To)
	}
}

func TestPullEvents_ClearEvents(t *testing.T) {
	order := &Order{
		ID:     "123",
		Status: StatusCreated,
	}

	_ = order.ChangeStatus(StatusPaid)

	events := order.PullEvents()

	if len(events) != 1 {
		t.Fatalf("expected 1 event")
	}

	events = order.PullEvents()

	if len(events) != 0 {
		t.Fatalf("expected events to be cleared")
	}
}

func TestCanTransition(t *testing.T) {
	if !CanTransition(StatusCreated, StatusPaid) {
		t.Fatalf("expected transition to be valid")
	}

	if CanTransition(StatusCreated, StatusShipped) {
		t.Fatalf("expected transition to be invalid")
	}
}
