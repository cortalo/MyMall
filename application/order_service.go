package application

import (
	"MyMall/domain"
	"MyMall/domain/shared"
	"context"
)

type OrderService interface {
	CreateOrder(ctx context.Context, customerID int64, inputs []domain.OrderItemInput, operator shared.Operator,
		idempotencyKey string) (*domain.Order, error)
}

type OrderRepository interface {
	Save(ctx context.Context, order *domain.Order) error
	FindByID(ctx context.Context, id int64) (*domain.Order, error)
}

type IdempotencyRepository interface {
	FindOrderIDByKey(ctx context.Context, key string) (int64, bool, error)
	SaveKey(ctx context.Context, key string, orderID int64) error
}

type orderService struct {
	repo            OrderRepository
	idempotencyRepo IdempotencyRepository
}

func NewOrderService(repo OrderRepository, idempotencyRepo IdempotencyRepository) OrderService {
	return &orderService{repo: repo, idempotencyRepo: idempotencyRepo}
}

func (s *orderService) CreateOrder(ctx context.Context, customerID int64, inputs []domain.OrderItemInput,
	operator shared.Operator, idempotencyKey string) (*domain.Order, error) {
	orderID, exists, err := s.idempotencyRepo.FindOrderIDByKey(ctx, idempotencyKey)
	if err != nil {
		return nil, err
	}
	if exists {
		return s.repo.FindByID(ctx, orderID)
	}
	order, _, err := domain.CreateOrder(customerID, inputs, operator)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, order); err != nil {
		return nil, err
	}
	if err := s.idempotencyRepo.SaveKey(ctx, idempotencyKey, order.ID); err != nil {
		return nil, err
	}
	return order, nil
}
