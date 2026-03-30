package application

import (
	"MyMall/domain"
	"MyMall/domain/shared"
	"context"
	"errors"
)

type OrderService interface {
	CreateOrder(ctx context.Context, customerID int64, inputs []domain.OrderItemInput, operator shared.Operator,
		idempotencyKey string) (*domain.Order, error)
}

type OrderRepository interface {
	Save(ctx context.Context, order *domain.Order, idempotencyKey string) error
	FindByID(ctx context.Context, id int64) (*domain.Order, error)
	FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Order, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, event domain.Event) error
}

type orderService struct {
	repo      OrderRepository
	publisher EventPublisher
}

func NewOrderService(repo OrderRepository, publisher EventPublisher) OrderService {
	return &orderService{repo: repo, publisher: publisher}
}

func (s *orderService) CreateOrder(
	ctx context.Context,
	customerID int64,
	inputs []domain.OrderItemInput,
	operator shared.Operator,
	idempotencyKey string,
) (*domain.Order, error) {
	order, events, err := domain.CreateOrder(customerID, inputs, operator)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, order, idempotencyKey); err != nil {
		if errors.Is(err, domain.ErrDuplicateIdempotencyKey) {
			return s.repo.FindByIdempotencyKey(ctx, idempotencyKey)
		}
		return nil, err
	}

	for _, event := range events {
		if err := s.publisher.Publish(ctx, event); err != nil {
			return nil, err
		}
	}
	return order, nil
}
