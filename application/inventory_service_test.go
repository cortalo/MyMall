package application

import (
	"MyMall/domain"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// ---- Mock Logger ----

type mockLogger struct{}

func (m *mockLogger) Info(ctx context.Context, msg string, fields ...Field)  {}
func (m *mockLogger) Warn(ctx context.Context, msg string, fields ...Field)  {}
func (m *mockLogger) Error(ctx context.Context, msg string, fields ...Field) {}

// ---- Mock Repos ----

type mockInventoryRepository struct {
	findByProductID   func(ctx context.Context, productID int64) (*domain.Inventory, error)
	save              func(ctx context.Context, inventory *domain.Inventory) error
	deductByProductID func(ctx context.Context, productID int64, amount int) error
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

func newMockInventoryRepository(t *testing.T) *mockInventoryRepository {
	return &mockInventoryRepository{
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

// ---- Mock UnitOfWork ----

type mockUnitOfWork struct {
	inventoryRepo      *mockInventoryRepository
	processedEventRepo *mockProcessedEventRepository
	commit             func(ctx context.Context) error
	rollback           func(ctx context.Context) error
}

func (m *mockUnitOfWork) InventoryRepo() InventoryRepository           { return m.inventoryRepo }
func (m *mockUnitOfWork) ProcessedEventRepo() ProcessedEventRepository { return m.processedEventRepo }
func (m *mockUnitOfWork) Commit(ctx context.Context) error             { return m.commit(ctx) }
func (m *mockUnitOfWork) Rollback(ctx context.Context) error           { return m.rollback(ctx) }

func newMockUnitOfWork(t *testing.T) *mockUnitOfWork {
	return &mockUnitOfWork{
		inventoryRepo:      newMockInventoryRepository(t),
		processedEventRepo: newMockProcessedEventRepository(t),
		commit: func(ctx context.Context) error {
			return nil
		},
		rollback: func(ctx context.Context) error {
			// Rollback 在 committed 后被 defer 调用是正常的，不报错
			return nil
		},
	}
}

type mockUnitOfWorkFactory struct {
	new func(ctx context.Context) (UnitOfWork, error)
}

func (m *mockUnitOfWorkFactory) New(ctx context.Context) (UnitOfWork, error) {
	return m.new(ctx)
}

func newMockUnitOfWorkFactory(uow UnitOfWork) *mockUnitOfWorkFactory {
	return &mockUnitOfWorkFactory{
		new: func(ctx context.Context) (UnitOfWork, error) {
			return uow, nil
		},
	}
}

// ---- Tests ----

func TestInventoryService_HandleOrderCreated(t *testing.T) {
	inventoryStore := map[int64]*domain.Inventory{
		1: {ID: 1, ProductID: 1, Stock: 100},
		2: {ID: 2, ProductID: 2, Stock: 50},
	}

	uow := newMockUnitOfWork(t)
	uow.inventoryRepo.deductByProductID = func(ctx context.Context, productID int64, amount int) error {
		inv, ok := inventoryStore[productID]
		require.True(t, ok)
		inv.Stock -= amount
		return nil
	}
	uow.processedEventRepo.isProcessed = func(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
		return false, nil
	}
	uow.processedEventRepo.markProcessed = func(ctx context.Context, eventName string, uniqueKey string) error {
		return nil
	}

	service := NewInventoryService(newMockUnitOfWorkFactory(uow), &mockLogger{})

	event := domain.OrderCreatedEvent{
		OrderID: 1,
		Items: []*domain.OrderItem{
			{ProductID: 1, Quantity: 3},
			{ProductID: 2, Quantity: 5},
		},
	}

	err := service.HandleOrderCreated(context.Background(), event)

	require.NoError(t, err)
	require.Equal(t, 97, inventoryStore[1].Stock, "product 1 stock should be decremented by 3")
	require.Equal(t, 45, inventoryStore[2].Stock, "product 2 stock should be decremented by 5")
}

func TestInventoryService_HandleOrderCreated_EmptyItems(t *testing.T) {
	uow := newMockUnitOfWork(t)
	uow.processedEventRepo.isProcessed = func(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
		return false, nil
	}
	uow.processedEventRepo.markProcessed = func(ctx context.Context, eventName string, uniqueKey string) error {
		return nil
	}

	service := NewInventoryService(newMockUnitOfWorkFactory(uow), &mockLogger{})

	event := domain.OrderCreatedEvent{
		OrderID: 1,
		Items:   []*domain.OrderItem{},
	}

	err := service.HandleOrderCreated(context.Background(), event)

	require.NoError(t, err)
}

func TestInventoryService_HandleOrderCreated_InsufficientStock(t *testing.T) {
	uow := newMockUnitOfWork(t)
	uow.inventoryRepo.deductByProductID = func(ctx context.Context, productID int64, amount int) error {
		return domain.ErrInsufficientStock
	}
	uow.processedEventRepo.isProcessed = func(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
		return false, nil
	}

	service := NewInventoryService(newMockUnitOfWorkFactory(uow), &mockLogger{})

	event := domain.OrderCreatedEvent{
		OrderID: 1,
		Items:   []*domain.OrderItem{{ProductID: 1, Quantity: 200}},
	}

	err := service.HandleOrderCreated(context.Background(), event)

	require.ErrorIs(t, err, domain.ErrInsufficientStock)
}

func TestInventoryService_HandleOrderCreated_PartialFailure(t *testing.T) {
	// 注意：mock 的 Commit/Rollback 不会真正操作数据库
	// 这个测试只验证 error 能被正确传播，回滚逻辑由集成测试覆盖
	inventoryStore := map[int64]*domain.Inventory{
		1: {ID: 1, ProductID: 1, Stock: 100},
		2: {ID: 2, ProductID: 2, Stock: 5},
	}

	uow := newMockUnitOfWork(t)
	uow.inventoryRepo.deductByProductID = func(ctx context.Context, productID int64, amount int) error {
		inv, ok := inventoryStore[productID]
		require.True(t, ok)
		if inv.Stock < amount {
			return domain.ErrInsufficientStock
		}
		inv.Stock -= amount
		return nil
	}
	uow.processedEventRepo.isProcessed = func(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
		return false, nil
	}

	service := NewInventoryService(newMockUnitOfWorkFactory(uow), &mockLogger{})

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
	uow := newMockUnitOfWork(t)
	uow.processedEventRepo.isProcessed = func(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
		return true, nil
	}
	// markProcessed 和 deductByProductID 保持 t.Fatal 默认值，调用即失败

	service := NewInventoryService(newMockUnitOfWorkFactory(uow), &mockLogger{})

	event := domain.OrderCreatedEvent{
		OrderID: 1,
		Items:   []*domain.OrderItem{{ProductID: 1, Quantity: 3}},
	}

	err := service.HandleOrderCreated(context.Background(), event)

	require.NoError(t, err)
}

func TestInventoryService_HandleOrderCreated_DuplicateProcessedEvent(t *testing.T) {
	uow := newMockUnitOfWork(t)
	uow.inventoryRepo.deductByProductID = func(ctx context.Context, productID int64, amount int) error {
		return nil
	}
	uow.processedEventRepo.isProcessed = func(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
		return false, nil
	}
	uow.processedEventRepo.markProcessed = func(ctx context.Context, eventName string, uniqueKey string) error {
		return domain.ErrDuplicateProcessedEvent
	}

	service := NewInventoryService(newMockUnitOfWorkFactory(uow), &mockLogger{})

	event := domain.OrderCreatedEvent{
		OrderID: 1,
		Items:   []*domain.OrderItem{{ProductID: 1, Quantity: 3}},
	}

	err := service.HandleOrderCreated(context.Background(), event)

	require.NoError(t, err)
}

func TestInventoryService_HandleOrderCreated_MarkProcessedError_RollbackCalled(t *testing.T) {
	rollbackCalled := false
	commitCalled := false

	uow := newMockUnitOfWork(t)
	uow.inventoryRepo.deductByProductID = func(ctx context.Context, productID int64, amount int) error {
		return nil
	}
	uow.processedEventRepo.isProcessed = func(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
		return false, nil
	}
	uow.processedEventRepo.markProcessed = func(ctx context.Context, eventName string, uniqueKey string) error {
		return domain.ErrDuplicateProcessedEvent
	}
	uow.commit = func(ctx context.Context) error {
		commitCalled = true
		return nil
	}
	uow.rollback = func(ctx context.Context) error {
		rollbackCalled = true
		return nil
	}

	service := NewInventoryService(newMockUnitOfWorkFactory(uow), &mockLogger{})

	event := domain.OrderCreatedEvent{
		OrderID: 1,
		Items:   []*domain.OrderItem{{ProductID: 1, Quantity: 3}},
	}

	err := service.HandleOrderCreated(context.Background(), event)

	require.NoError(t, err)
	require.True(t, rollbackCalled, "rollback should be called when MarkProcessed fails")
	require.False(t, commitCalled, "commit should not be called when MarkProcessed fails")
}

func TestInventoryService_HandleOrderCreated_UoWFactoryError(t *testing.T) {
	factory := &mockUnitOfWorkFactory{
		new: func(ctx context.Context) (UnitOfWork, error) {
			return nil, errors.New("db connection failed")
		},
	}

	service := NewInventoryService(factory, &mockLogger{})

	err := service.HandleOrderCreated(context.Background(), domain.OrderCreatedEvent{OrderID: 1})

	require.Error(t, err)
}

func TestInventoryService_HandleOrderCreated_IsProcessedError(t *testing.T) {
	uow := newMockUnitOfWork(t)
	uow.processedEventRepo.isProcessed = func(ctx context.Context, eventName string, uniqueKey string) (bool, error) {
		return false, errors.New("db error")
	}

	service := NewInventoryService(newMockUnitOfWorkFactory(uow), &mockLogger{})

	err := service.HandleOrderCreated(context.Background(), domain.OrderCreatedEvent{OrderID: 1})

	require.Error(t, err)
}
