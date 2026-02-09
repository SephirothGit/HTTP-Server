package repository

import (
	"context"
	"sync"

	"github.com/SephirothGit/Backend-service/internal/domain"
)

type InMemoryOrderRepository struct {
	mu sync.RWMutex
	orders map[string]*domain.Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]*domain.Order),
	}
}

func (r *InMemoryOrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	o, ok := r.orders[id]
	if !ok {
		return nil, domain.ErrOrderNotFound
	}
	return o, nil
}

func (r *InMemoryOrderRepository) Save(ctx context.Context, order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.ID] = order
	return nil
}