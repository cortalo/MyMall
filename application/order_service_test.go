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
	save                 func(ctx context.Context, order *domain.Order, idempotencyKey string) error
	findByID             func(ctx context.Context, id int64) (*domain.Order, error)
	findByIdempotencyKey func(ctx context.Context, idempotencyKey string) (*domain.Order, error)
}

func (m *mockOrderRepository) Save(ctx context.Context, order *domain.Order, idempotencyKey string) error {
	return m.save(ctx, order, idempotencyKey)
}

func (m *mockOrderRepository) FindByID(ctx context.Context, id int64) (*domain.Order, error) {
	return m.findByID(ctx, id)
}

func (m *mockOrderRepository) FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Order, error) {
	return m.findByIdempotencyKey(ctx, idempotencyKey)
}

// mock bean
type mockEventPublisher struct {
	publish func(ctx context.Context, event domain.Event) error
}

func (m *mockEventPublisher) Publish(ctx context.Context, event domain.Event) error {
	return m.publish(ctx, event)
}

func TestOrderService_CreateOrder(t *testing.T) {
	repo := &mockOrderRepository{
		save: func(ctx context.Context, order *domain.Order, idempotencyKey string) error {
			require.Equal(t, "random-key-72", idempotencyKey)
			return nil
		},
	}
	publisher := &mockEventPublisher{
		publish: func(ctx context.Context, event domain.Event) error {
			return nil
		},
	}
	service := NewOrderService(repo, publisher)

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
		save: func(ctx context.Context, order *domain.Order, idempotencyKey string) error {
			return domain.ErrDuplicateIdempotencyKey
		},
		findByID: func(ctx context.Context, id int64) (*domain.Order, error) {
			t.Fatal("FindByID should not be called")
			return nil, nil
		},
		findByIdempotencyKey: func(ctx context.Context, idempotencyKey string) (*domain.Order, error) {
			require.Equal(t, "idempotency-key-1", idempotencyKey)
			return existingOrder, nil
		},
	}
	publisher := &mockEventPublisher{
		publish: func(ctx context.Context, event domain.Event) error {
			return nil
		},
	}
	service := NewOrderService(repo, publisher)

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
