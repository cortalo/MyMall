// application/sale_service_test.go
package application

import (
	"MyMall/domain"
	"MyMall/domain/shared"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type mockSaleRepository struct {
	findByID func(ctx context.Context, id int64) (*domain.Sale, error)
}

func (m *mockSaleRepository) FindByID(ctx context.Context, id int64) (*domain.Sale, error) {
	return m.findByID(ctx, id)
}

type mockOrderReadRepository struct {
	countByUserAndProduct func(ctx context.Context, userID int64, productID int64) (int, error)
	findByIdempotencyKey  func(ctx context.Context, idempotencyKey string) (*domain.Order, error)
}

func (m *mockOrderReadRepository) CountByUserAndProduct(ctx context.Context, userID int64, productID int64) (int, error) {
	return m.countByUserAndProduct(ctx, userID, productID)
}

func (m *mockOrderReadRepository) FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Order, error) {
	return m.findByIdempotencyKey(ctx, idempotencyKey)
}

func newMockSaleRepository(t *testing.T) *mockSaleRepository {
	return &mockSaleRepository{
		findByID: func(ctx context.Context, id int64) (*domain.Sale, error) {
			t.Fatal("findByID should not be called")
			return nil, nil
		},
	}
}

func newMockOrderReadRepository(t *testing.T) *mockOrderReadRepository {
	return &mockOrderReadRepository{
		countByUserAndProduct: func(ctx context.Context, userID int64, productID int64) (int, error) {
			t.Fatal("countByUserAndProduct should not be called")
			return 0, nil
		},
		findByIdempotencyKey: func(ctx context.Context, idempotencyKey string) (*domain.Order, error) {
			return nil, domain.ErrOrderEmptyItems
		},
	}
}

type mockOrderCreator struct {
	createOrder func(ctx context.Context, customerID int64, inputs []domain.OrderItemInput, operator shared.Operator, idempotencyKey string) (*domain.Order, error)
}

func (m *mockOrderCreator) CreateOrder(ctx context.Context, customerID int64, inputs []domain.OrderItemInput, operator shared.Operator, idempotencyKey string) (*domain.Order, error) {
	return m.createOrder(ctx, customerID, inputs, operator, idempotencyKey)
}

func newMockOrderCreator(t *testing.T) *mockOrderCreator {
	return &mockOrderCreator{
		createOrder: func(ctx context.Context, customerID int64, inputs []domain.OrderItemInput, operator shared.Operator, idempotencyKey string) (*domain.Order, error) {
			t.Fatal("createOrder should not be called")
			return nil, nil
		},
	}
}

func TestSaleService_CreateOrder_Success(t *testing.T) {
	saleRepo := newMockSaleRepository(t)
	saleRepo.findByID = func(ctx context.Context, id int64) (*domain.Sale, error) {
		return &domain.Sale{
			ID:          1,
			ProductID:   1,
			ProductName: "Apple",
			UnitPrice:   500,
			MaxPerUser:  2,
			StartAt:     time.Now().Add(-1 * time.Hour),
			EndAt:       time.Now().Add(1 * time.Hour),
		}, nil
	}

	orderReadRepo := newMockOrderReadRepository(t)
	orderReadRepo.countByUserAndProduct = func(ctx context.Context, userID int64, productID int64) (int, error) {
		return 0, nil
	}

	orderRepo := &mockOrderRepository{
		save: func(ctx context.Context, order *domain.Order, idempotencyKey string) error {
			order.ID = 1
			return nil
		},
	}
	publisher := &mockEventPublisher{
		publish: func(ctx context.Context, event domain.Event) error {
			return nil
		},
	}
	orderCreator := NewOrderService(orderRepo, publisher)

	service := NewSaleService(saleRepo, orderReadRepo, orderCreator)

	operator := shared.Operator{ID: 1, Username: "test_user"}
	items := []domain.SaleOrderItem{
		{SaleID: 1, Quantity: 1},
	}

	order, err := service.CreateOrder(context.Background(), 42, items, operator, "random-key-72")

	require.NoError(t, err)
	require.NotNil(t, order)
	require.Equal(t, int64(500), order.TotalAmount)
}

func TestSaleService_CreateOrder_SaleNotFound(t *testing.T) {
	saleRepo := newMockSaleRepository(t)
	saleRepo.findByID = func(ctx context.Context, id int64) (*domain.Sale, error) {
		return nil, domain.ErrSaleNotFound
	}

	service := NewSaleService(saleRepo, newMockOrderReadRepository(t), newMockOrderCreator(t))

	operator := shared.Operator{ID: 1, Username: "test_user"}
	items := []domain.SaleOrderItem{
		{SaleID: 1, Quantity: 1},
	}

	_, err := service.CreateOrder(context.Background(), 42, items, operator, "random-key-72")

	require.ErrorIs(t, err, domain.ErrSaleNotFound)
}

func TestSaleService_CreateOrder_SaleNotActive(t *testing.T) {
	saleRepo := newMockSaleRepository(t)
	saleRepo.findByID = func(ctx context.Context, id int64) (*domain.Sale, error) {
		return &domain.Sale{
			ID:        1,
			ProductID: 1,
			StartAt:   time.Now().Add(-2 * time.Hour),
			EndAt:     time.Now().Add(-1 * time.Hour),
		}, nil
	}

	service := NewSaleService(saleRepo, newMockOrderReadRepository(t), newMockOrderCreator(t))

	operator := shared.Operator{ID: 1, Username: "test_user"}
	items := []domain.SaleOrderItem{
		{SaleID: 1, Quantity: 1},
	}

	_, err := service.CreateOrder(context.Background(), 42, items, operator, "random-key-72")

	require.ErrorIs(t, err, domain.ErrSaleNotActive)
}

func TestSaleService_CreateOrder_ExceedPurchaseLimit(t *testing.T) {
	saleRepo := newMockSaleRepository(t)
	saleRepo.findByID = func(ctx context.Context, id int64) (*domain.Sale, error) {
		return &domain.Sale{
			ID:         1,
			ProductID:  1,
			MaxPerUser: 2,
			StartAt:    time.Now().Add(-1 * time.Hour),
			EndAt:      time.Now().Add(1 * time.Hour),
		}, nil
	}

	orderCountRepo := newMockOrderReadRepository(t)
	orderCountRepo.countByUserAndProduct = func(ctx context.Context, userID int64, productID int64) (int, error) {
		return 2, nil
	}

	service := NewSaleService(saleRepo, orderCountRepo, newMockOrderCreator(t))

	operator := shared.Operator{ID: 1, Username: "test_user"}
	items := []domain.SaleOrderItem{
		{SaleID: 1, Quantity: 1},
	}

	_, err := service.CreateOrder(context.Background(), 42, items, operator, "random-key-72")

	require.ErrorIs(t, err, domain.ErrExceedPurchaseLimit)
}

func TestSaleService_CreateOrder_IdempotentReturn(t *testing.T) {
	existingOrder := &domain.Order{
		ID:          99,
		TotalAmount: 500,
	}

	orderReadRepo := newMockOrderReadRepository(t)
	orderReadRepo.findByIdempotencyKey = func(ctx context.Context, idempotencyKey string) (*domain.Order, error) {
		require.Equal(t, "random-key-72", idempotencyKey)
		return existingOrder, nil
	}

	// saleRepo 和 orderCreator 用默认的 mock——命中幂等时两者都不应该被调用
	service := NewSaleService(newMockSaleRepository(t), orderReadRepo, newMockOrderCreator(t))

	operator := shared.Operator{ID: 1, Username: "test_user"}
	items := []domain.SaleOrderItem{
		{SaleID: 1, Quantity: 1},
	}

	order, err := service.CreateOrder(context.Background(), 42, items, operator, "random-key-72")

	require.NoError(t, err)
	require.Equal(t, existingOrder, order)
}
