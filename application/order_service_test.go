package application

import (
	"MyMall/domain"
	"MyMall/domain/shared"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// mock bean
type mockOrderRepository struct {
	save     func(ctx context.Context, order *domain.Order) error
	findByID func(ctx context.Context, id int64) (*domain.Order, error)
}

func (m *mockOrderRepository) Save(ctx context.Context, order *domain.Order) error {
	return m.save(ctx, order)
}

func (m *mockOrderRepository) FindByID(ctx context.Context, id int64) (*domain.Order, error) {
	return m.findByID(ctx, id)
}

// mock bean
type mockIdempotencyRepository struct {
	findOrderIDByKey func(ctx context.Context, key string) (int64, bool, error)
	saveKey          func(ctx context.Context, key string, orderID int64) error
}

func (m *mockIdempotencyRepository) FindOrderIDByKey(ctx context.Context, key string) (int64, bool, error) {
	return m.findOrderIDByKey(ctx, key)
}

func (m *mockIdempotencyRepository) SaveKey(ctx context.Context, key string, orderID int64) error {
	return m.saveKey(ctx, key, orderID)
}

func TestOrderService_CreateOrder(t *testing.T) {
	repo := &mockOrderRepository{
		save: func(ctx context.Context, order *domain.Order) error {
			return nil
		},
	}
	idempotencyRepo := &mockIdempotencyRepository{
		findOrderIDByKey: func(ctx context.Context, key string) (int64, bool, error) {
			require.Equal(t, "random-key-72", key)
			return 0, false, nil
		},
		saveKey: func(ctx context.Context, key string, orderID int64) error {
			require.Equal(t, "random-key-72", key)
			return nil
		},
	}
	service := NewOrderService(repo, idempotencyRepo)

	operator := shared.Operator{ID: 1, Username: "test_user"}
	inputs := []domain.OrderItemInput{
		{ProductID: 1, ProductName: "Apple", UnitPrice: 500, Quantity: 2},
	}

	order, err := service.CreateOrder(context.Background(), 42, inputs, operator, "random-key-72")

	require.NoError(t, err)
	require.NotNil(t, order)
	require.Equal(t, int64(42), order.CustomerID)
	require.Equal(t, domain.OrderStatusPending, order.Status)
	require.Equal(t, int64(1000), order.TotalAmount)
}

func TestOrderService_CreateOrder_IdempotentRetry(t *testing.T) {
	existingOrder := &domain.Order{
		ID:          99,
		CustomerID:  42,
		Status:      domain.OrderStatusPending,
		TotalAmount: 1000,
	}
	repo := &mockOrderRepository{
		save: func(ctx context.Context, order *domain.Order) error {
			t.Fatal("Save should not be called on idempotent retry")
			return nil
		},
		findByID: func(ctx context.Context, id int64) (*domain.Order, error) {
			require.Equal(t, int64(99), id)
			return existingOrder, nil
		},
	}
	idempotencyRepo := &mockIdempotencyRepository{
		findOrderIDByKey: func(ctx context.Context, key string) (int64, bool, error) {
			require.Equal(t, "idempotency-key-1", key)
			return 99, true, nil
		},
		saveKey: func(ctx context.Context, key string, orderID int64) error {
			t.Fatal("SaveKey should not be called on idempotent retry")
			return nil
		},
	}
	service := NewOrderService(repo, idempotencyRepo)

	operator := shared.Operator{ID: 1, Username: "test_user"}
	inputs := []domain.OrderItemInput{
		{ProductID: 1, ProductName: "Apple", UnitPrice: 500, Quantity: 2},
	}

	order, err := service.CreateOrder(context.Background(), 42, inputs, operator, "idempotency-key-1")

	require.NotNil(t, order)
	require.NoError(t, err)
	require.Equal(t, int64(99), order.ID)
	require.Equal(t, int64(42), order.CustomerID)
}
