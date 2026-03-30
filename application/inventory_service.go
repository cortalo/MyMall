package application

import (
	"MyMall/domain"
	"context"
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

type inventoryService struct {
	repo InventoryRepository
}

func NewInventoryService(repo InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) HandleOrderCreated(ctx context.Context, event domain.OrderCreatedEvent) error {
	return s.repo.WithTx(ctx, func(ctx context.Context) error {
		for _, item := range event.Items {
			if err := s.repo.DeductByProductID(ctx, item.ProductID, item.Quantity); err != nil {
				return err
			}
		}
		return nil
	})
}
