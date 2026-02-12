package service

import (
	"github.com/SephirothGit/Backend-service/internal/domain"
	"context"
	"log"
)

type OrderRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	Save(ctx context.Context, order *domain.Order) error
}

// Event publisher
type EventPublisher interface {
	Publish(ctx context.Context, events ...domain.Event)
}

// Simple logger
type LogPublisher struct{}

func (LogPublisher) Publish(ctx context.Context, events ...domain.Event) {
	for _, e := range events {
		log.Printf("EVENT: %s %+v\n", e.EventName(), e)
	}
}

// Service
type OrderService interface {
	UpdateStatus(ctx context.Context, id string, status string) error
}

type orderService struct {
	repo OrderRepository
	publisher EventPublisher
}

func NewOrderService(repo OrderRepository) OrderService {
	return &orderService{
		repo: repo,
		publisher: LogPublisher{},
	}
}

func (s *orderService) UpdateStatus(
	ctx context.Context,
	id string,
	status string,
) error {

	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.ErrOrderNotFound
	}

	// Idempotency
	if domain.IsSameStatus(order.Status, status) {
		return nil
	}

	// Domain logic inside entity
	if err := order.ChangeStatus(status); err != nil {
		return err
	}

	if err := s.repo.Save(ctx, order); err != nil {
		return err
	}

	// Publish event
	events := order.PullEvents()
	s.publisher.Publish(ctx, events...)

	return nil
}