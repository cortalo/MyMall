package application

import (
	"MyMall/domain"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockInventoryRepository struct {
	findByProductID   func(ctx context.Context, productID int64) (*domain.Inventory, error)
	save              func(ctx context.Context, inventory *domain.Inventory) error
	deductByProductID func(ctx context.Context, productID int64, amount int) error
	withTx            func(ctx context.Context, fn func(ctx context.Context) error) error
}

func (m *mockInventoryRepository) FindByProductID(ctx context.Context, productID int64) (*domain.Inventory, error) {
	return m.findByProductID(ctx, productID)
}
func (m *mockInventoryRepository) Save(ctx context.Context, inventory *domain.Inventory) error {
	return m.save(ctx, inventory)
}
func (m *mockInventoryRepository) DeductByProductID(ctx context.Context, productID int64, amount int) error {
	return m.deductByProductID(ctx, productID, amount)
}
func (m *mockInventoryRepository) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.withTx(ctx, fn)
}

func newMockInventoryRepository(t *testing.T) *mockInventoryRepository {
	return &mockInventoryRepository{
		withTx: func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		},
		findByProductID: func(ctx context.Context, productID int64) (*domain.Inventory, error) {
			t.Fatal("findByProductID should not be called")
			return nil, nil
		},
		save: func(ctx context.Context, inventory *domain.Inventory) error {
			t.Fatal("save should not be called")
			return nil
		},
		deductByProductID: func(ctx context.Context, productID int64, amount int) error {
			t.Fatal("deductByProductID should not be called")
			return nil
		},
	}
}

type mockProcessedEventRepository struct {
	isProcessed   func(ctx context.Context, eventName string, uniqueKey string) (bool, error)
	markProcessed func(ctx context.Context, eventName string, uniqueKey string) error
}

func (m *mockProcessedEventRepository) IsProcessed(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
	return m.isProcessed(ctx, eventName, uniqueKey)
}

func (m *mockProcessedEventRepository) MarkProcessed(ctx context.Context, eventName string, uniqueKey string) error {
	return m.markProcessed(ctx, eventName, uniqueKey)
}

func newMockProcessedEventRepository(t *testing.T) *mockProcessedEventRepository {
	return &mockProcessedEventRepository{
		isProcessed: func(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
			t.Fatal("isProcessed should not be called")
			return false, nil
		},
		markProcessed: func(ctx context.Context, eventName string, uniqueKey string) error {
			t.Fatal("markProcessed should not be called")
			return nil
		},
	}
}

func TestInventoryService_HandleOrderCreated(t *testing.T) {
	// 用 map 模拟数据库中的库存数据
	inventoryStore := map[int64]*domain.Inventory{
		1: {ID: 1, ProductID: 1, Stock: 100},
		2: {ID: 2, ProductID: 2, Stock: 50},
	}

	repo := newMockInventoryRepository(t)
	repo.deductByProductID = func(ctx context.Context, productID int64, amount int) error {
		inv, ok := inventoryStore[productID]
		require.True(t, ok)
		inv.Stock -= amount
		return nil
	}
	processedEventRepo := newMockProcessedEventRepository(t)
	processedEventRepo.isProcessed = func(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
		return false, nil
	}
	processedEventRepo.markProcessed = func(ctx context.Context, eventName string, uniqueKey string) error {
		return nil
	}
	service := NewInventoryService(repo, processedEventRepo)

	event := domain.OrderCreatedEvent{
		OrderID: 1,
		Items: []*domain.OrderItem{
			{ProductID: 1, Quantity: 3},
			{ProductID: 2, Quantity: 5},
		},
	}

	err := service.HandleOrderCreated(context.Background(), event)

	require.NoError(t, err)
	// 验证库存是否被正确扣减
	require.Equal(t, 97, inventoryStore[1].Stock, "product 1 stock should be decremented by 3")
	require.Equal(t, 45, inventoryStore[2].Stock, "product 2 stock should be decremented by 5")
}

func TestInventoryService_HandleOrderCreated_EmptyItems(t *testing.T) {
	repo := newMockInventoryRepository(t)
	processedEventRepo := newMockProcessedEventRepository(t)
	processedEventRepo.isProcessed = func(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
		return false, nil
	}
	processedEventRepo.markProcessed = func(ctx context.Context, eventName string, uniqueKey string) error {
		return nil
	}
	service := NewInventoryService(repo, processedEventRepo)

	event := domain.OrderCreatedEvent{
		OrderID: 1,
		Items:   []*domain.OrderItem{},
	}

	err := service.HandleOrderCreated(context.Background(), event)

	require.NoError(t, err)
}

func TestInventoryService_HandleOrderCreated_InsufficientStock(t *testing.T) {
	repo := newMockInventoryRepository(t)
	repo.deductByProductID = func(ctx context.Context, productID int64, amount int) error {
		return domain.ErrInsufficientStock
	}

	processedEventRepo := newMockProcessedEventRepository(t)
	processedEventRepo.isProcessed = func(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
		return false, nil
	}
	processedEventRepo.markProcessed = func(ctx context.Context, eventName string, uniqueKey string) error {
		return nil
	}
	service := NewInventoryService(repo, processedEventRepo)

	event := domain.OrderCreatedEvent{
		OrderID: 1,
		Items: []*domain.OrderItem{
			{ProductID: 1, Quantity: 200},
		},
	}

	err := service.HandleOrderCreated(context.Background(), event)

	require.ErrorIs(t, err, domain.ErrInsufficientStock)
}

func TestInventoryService_HandleOrderCreated_PartialFailure(t *testing.T) {
	// 注意：mock 的 WithTx 是透传的，不会真正回滚
	// 这个测试只验证 error 能被正确传播，回滚逻辑由集成测试覆盖
	inventoryStore := map[int64]*domain.Inventory{
		1: {ID: 1, ProductID: 1, Stock: 100},
		2: {ID: 2, ProductID: 2, Stock: 5},
	}

	repo := newMockInventoryRepository(t)
	repo.deductByProductID = func(ctx context.Context, productID int64, amount int) error {
		inv, ok := inventoryStore[productID]
		require.True(t, ok)
		if inv.Stock < amount {
			return domain.ErrInsufficientStock
		}
		inv.Stock -= amount
		return nil
	}

	processedEventRepo := newMockProcessedEventRepository(t)
	processedEventRepo.isProcessed = func(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
		return false, nil
	}
	processedEventRepo.markProcessed = func(ctx context.Context, eventName string, uniqueKey string) error {
		return nil
	}
	service := NewInventoryService(repo, processedEventRepo)

	event := domain.OrderCreatedEvent{
		OrderID: 1,
		Items: []*domain.OrderItem{
			{ProductID: 1, Quantity: 3},
			{ProductID: 2, Quantity: 10},
		},
	}

	err := service.HandleOrderCreated(context.Background(), event)

	require.ErrorIs(t, err, domain.ErrInsufficientStock)
}

func TestInventoryService_HandleOrderCreated_Idempotent(t *testing.T) {
	repo := newMockInventoryRepository(t)
	repo.deductByProductID = func(ctx context.Context, productID int64, amount int) error {
		t.Fatal("deductByProductID should not be called for duplicate event")
		return nil
	}

	processedEventRepo := &mockProcessedEventRepository{
		isProcessed: func(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
			return true, nil // 已经处理过了
		},
		markProcessed: func(ctx context.Context, eventName string, uniqueKey string) error {
			t.Fatal("markProcessed should not be called for duplicate event")
			return nil
		},
	}

	service := NewInventoryService(repo, processedEventRepo)

	event := domain.OrderCreatedEvent{
		OrderID: 1,
		Items:   []*domain.OrderItem{{ProductID: 1, Quantity: 3}},
	}

	err := service.HandleOrderCreated(context.Background(), event)

	require.NoError(t, err)
}

func TestInventoryService_HandleOrderCreated_MarkProcessedWithinTx(t *testing.T) {
	var markCalledInTx bool

	repo := newMockInventoryRepository(t)
	repo.withTx = func(ctx context.Context, fn func(ctx context.Context) error) error {
		// 执行 fn 之后检查 markProcessed 是否在事务内被调用
		err := fn(ctx)
		require.True(t, markCalledInTx, "markProcessed should be called within transaction")
		return err
	}
	repo.deductByProductID = func(ctx context.Context, productID int64, amount int) error {
		return nil
	}

	processedEventRepo := &mockProcessedEventRepository{
		isProcessed: func(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
			return false, nil
		},
		markProcessed: func(ctx context.Context, eventName string, uniqueKey string) error {
			markCalledInTx = true
			return nil
		},
	}

	service := NewInventoryService(repo, processedEventRepo)

	event := domain.OrderCreatedEvent{
		OrderID: 1,
		Items:   []*domain.OrderItem{{ProductID: 1, Quantity: 3}},
	}

	err := service.HandleOrderCreated(context.Background(), event)

	require.NoError(t, err)
}

func TestInventoryService_HandleOrderCreated_DuplicateProcessedEvent(t *testing.T) {
	// 模拟并发下 MarkProcessed 唯一键冲突
	repo := newMockInventoryRepository(t)
	repo.deductByProductID = func(ctx context.Context, productID int64, amount int) error {
		return nil
	}

	processedEventRepo := newMockProcessedEventRepository(t)
	processedEventRepo.isProcessed = func(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
		return false, nil
	}
	processedEventRepo.markProcessed = func(ctx context.Context, eventName string, uniqueKey string) error {
		return domain.ErrDuplicateProcessedEvent
	}

	service := NewInventoryService(repo, processedEventRepo)

	event := domain.OrderCreatedEvent{
		OrderID: 1,
		Items:   []*domain.OrderItem{{ProductID: 1, Quantity: 3}},
	}

	err := service.HandleOrderCreated(context.Background(), event)

	require.ErrorIs(t, err, domain.ErrDuplicateProcessedEvent)
}
