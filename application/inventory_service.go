package application

import (
	"MyMall/domain"
	"context"
	"errors"
	"strconv"
)

type InventoryService interface {
	HandleOrderCreated(ctx context.Context, event domain.OrderCreatedEvent) error
}

type Logger interface {
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
}
type Field struct {
	Key   string
	Value any
}

func F(key string, value any) Field {
	return Field{Key: key, Value: value}
}

type InventoryRepository interface {
	FindByProductID(ctx context.Context, productID int64) (*domain.Inventory, error)
	Save(ctx context.Context, inventory *domain.Inventory) error
	DeductByProductID(ctx context.Context, productID int64, amount int) error
}

type ProcessedEventRepository interface {
	IsProcessed(ctx context.Context, eventName string, uniqueKey string) (bool, error)
	MarkProcessed(ctx context.Context, eventName string, uniqueKey string) error
}

type UnitOfWork interface {
	InventoryRepo() InventoryRepository
	ProcessedEventRepo() ProcessedEventRepository
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type UnitOfWorkFactory interface {
	New(ctx context.Context) (UnitOfWork, error)
}

type inventoryService struct {
	uowFactory UnitOfWorkFactory
	logger     Logger
}

func NewInventoryService(uowFactory UnitOfWorkFactory, logger Logger) InventoryService {
	return &inventoryService{uowFactory: uowFactory, logger: logger}
}

func (s *inventoryService) HandleOrderCreated(ctx context.Context, event domain.OrderCreatedEvent) error {
	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := uow.Rollback(ctx); err != nil {
			s.logger.Warn(ctx, "rollback failed",
				F("event_name", event.EventName()),
				F("order_id", event.OrderID),
				F("error", err),
			)
		}
	}()

	uniqueKey := strconv.FormatInt(event.OrderID, 10)
	processed, err := uow.ProcessedEventRepo().IsProcessed(ctx, event.EventName(), uniqueKey)
	if err != nil {
		return err
	}
	if processed {
		return nil
	}
	for _, item := range event.Items {
		if err := uow.InventoryRepo().DeductByProductID(ctx, item.ProductID, item.Quantity); err != nil {
			return err
		}
	}
	if err := uow.ProcessedEventRepo().MarkProcessed(ctx, event.EventName(), uniqueKey); err != nil {
		if errors.Is(err, domain.ErrDuplicateProcessedEvent) {
			return nil
		}
		return err
	}
	return uow.Commit(ctx)
}
