package application

import (
	"MyMall/domain"
	"MyMall/domain/shared"
	"context"
	"errors"
)
import appshared "MyMall/application/shared"

type OrderService interface {
	CreateOrder(ctx context.Context, customerID int64, inputs []domain.OrderItemInput, operator shared.Operator,
		idempotencyKey string) (*domain.Order, error)
}

type OrderRepository interface {
	Save(ctx context.Context, order *domain.Order, idempotencyKey string) error
	FindByID(ctx context.Context, id int64) (*domain.Order, error)
	FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Order, error)
}

type OrderUnitOfWork interface {
	OrderRepo() OrderRepository
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type OrderUnitOfWorkFactory interface {
	New(ctx context.Context) (OrderUnitOfWork, error)
}

type orderService struct {
	uowFactory OrderUnitOfWorkFactory
	publisher  appshared.EventPublisher
	logger     Logger
}

func NewOrderService(uowFactory OrderUnitOfWorkFactory, publisher appshared.EventPublisher, logger Logger) OrderService {
	return &orderService{uowFactory: uowFactory, publisher: publisher, logger: logger}
}

func (s *orderService) CreateOrder(
	ctx context.Context,
	customerID int64,
	inputs []domain.OrderItemInput,
	operator shared.Operator,
	idempotencyKey string,
) (*domain.Order, error) {
	order, err := domain.CreateOrder(customerID, inputs, operator)
	if err != nil {
		return nil, err
	}
	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := uow.Rollback(ctx); err != nil {
			s.logger.Warn(ctx, "rollback failed",
				F("error", err),
			)
		}
	}()

	if err := uow.OrderRepo().Save(ctx, order, idempotencyKey); err != nil {
		if errors.Is(err, domain.ErrDuplicateIdempotencyKey) {
			return uow.OrderRepo().FindByIdempotencyKey(ctx, idempotencyKey)
		}
		return nil, err
	}

	for _, event := range order.BuildEvents() {
		if err := s.publisher.Publish(ctx, event); err != nil {
			return nil, err
		}
	}
	return order, uow.Commit(ctx)
}
