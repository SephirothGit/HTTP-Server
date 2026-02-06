package service

import (
	domain "_/C_/GOPATH/Http/echo-server"
	"context"
)

type OrderRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	Save(ctx context.Context, order *domain.Order) error
}

//Use-case API
type OrderService interface {
	UpdateStatus(ctx context.Context, id string, status string) error
}

type orderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) OrderService {
	return &orderService{
		repo: repo,
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

	if !domain.CanTransition(order.Status, status) {
		return domain.ErrInvalidTransition
	}

	order.Status = status

	if err := s.repo.Save(ctx, order); err != nil {
		return err
	}
	
	return nil
}