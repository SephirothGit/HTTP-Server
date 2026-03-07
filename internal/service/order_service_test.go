package service

import (
	"context"
	"errors"
	"testing"

	"github.com/SephirothGit/Backend-service/internal/domain"
)

type mockRepo struct {
	order   *domain.Order
	getErr  error
	saveErr error
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	return m.order, m.getErr
}

func (m *mockRepo) Save(ctx context.Context, order *domain.Order) error {
	return m.saveErr
}

type mockPublisher struct {
	events []domain.Event
}

func (m *mockPublisher) Publish(ctx context.Context, events ...domain.Event) {
	m.events = append(m.events, events...)
}

func TestUpdateStatus(t *testing.T) {

	tests := []struct {
		name        string
		order       *domain.Order
		getErr      error
		saveErr     error
		newStatus   string
		expectError error
		expectEvent bool
	}{
		{
			name:        "order not found",
			order:       nil,
			getErr:      errors.New("db error"),
			newStatus:   domain.StatusPaid,
			expectError: domain.ErrOrderNotFound,
		},
		{
			name: "invalid transition",
			order: &domain.Order{
				ID:     "123",
				Status: domain.StatusCreated,
			},
			newStatus:   domain.StatusShipped,
			expectError: domain.ErrInvalidTransition,
		},
		{
			name: "idempotent request",
			order: &domain.Order{
				ID:     "123",
				Status: domain.StatusPaid,
			},
			newStatus:   domain.StatusPaid,
			expectError: nil,
			expectEvent: false,
		},
		{
			name: "save error",
			order: &domain.Order{
				ID:     "123",
				Status: domain.StatusCreated,
			},
			newStatus:   domain.StatusPaid,
			saveErr:     errors.New("db error"),
			expectError: errors.New("db error"),
		},
		{
			name: "success",
			order: &domain.Order{
				ID:     "123",
				Status: domain.StatusCreated,
			},
			newStatus:   domain.StatusPaid,
			expectError: nil,
			expectEvent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mockRepo{
				order:   tt.order,
				getErr:  tt.getErr,
				saveErr: tt.saveErr,
			}

			pub := &mockPublisher{}

			svc := &orderService{
				repo:      &repo,
				publisher: pub,
			}

			err := svc.UpdateStatus(context.Background(), "123", tt.newStatus)

			if tt.expectError != nil && err == nil {
				t.Fatalf("expected error")
			}
			if tt.expectError == nil && err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if tt.expectEvent && len(pub.events) == 0 {
				t.Fatalf("expected event to be published")
			}
		})
	}
}
