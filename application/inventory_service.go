package application

import (
	"MyMall/domain"
	"context"
	"strconv"
)

type InventoryService interface {
	HandleOrderCreated(ctx context.Context, event domain.OrderCreatedEvent) error
}

type InventoryRepository interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
	FindByProductID(ctx context.Context, productID int64) (*domain.Inventory, error)
	Save(ctx context.Context, inventory *domain.Inventory) error
	DeductByProductID(ctx context.Context, productID int64, amount int) error
}

type ProcessedEventRepository interface {
	IsProcessed(ctx context.Context, eventName string, uniqueKey string) (bool, error)
	MarkProcessed(ctx context.Context, eventName string, uniqueKey string) error
}

type inventoryService struct {
	repo               InventoryRepository
	processedEventRepo ProcessedEventRepository
}

func NewInventoryService(repo InventoryRepository, processedEventRepo ProcessedEventRepository) InventoryService {
	return &inventoryService{repo: repo, processedEventRepo: processedEventRepo}
}

func (s *inventoryService) HandleOrderCreated(ctx context.Context, event domain.OrderCreatedEvent) error {
	processed, err := s.processedEventRepo.IsProcessed(ctx, event.EventName(), strconv.FormatInt(event.OrderID, 10))
	if err != nil {
		return err
	}
	if processed {
		return nil
	}
	return s.repo.WithTx(ctx, func(ctx context.Context) error {
		for _, item := range event.Items {
			if err := s.repo.DeductByProductID(ctx, item.ProductID, item.Quantity); err != nil {
				return err
			}
		}
		return s.processedEventRepo.MarkProcessed(ctx, event.EventName(), strconv.FormatInt(event.OrderID, 10))
	})
}
