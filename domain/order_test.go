package domain

import (
	"MyMall/domain/shared"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	operator := shared.Operator{ID: 1, Username: "test_user"}
	inputs := []OrderItemInput{
		{ProductID: 1, ProductName: "Apple", UnitPrice: 500, Quantity: 2},
		{ProductID: 1, ProductName: "Orange", UnitPrice: 399, Quantity: 3},
	}
	order, err := CreateOrder(42, inputs, operator)
	require.NotNil(t, order)
	require.NoError(t, err)
	require.Equal(t, int64(42), order.CustomerID)
	require.Equal(t, operator.ID, order.CreatorID)
	require.Equal(t, operator.Username, order.CreatorName)
	require.Equal(t, int64(0), order.ModifierID)
	require.Equal(t, "", order.ModifierName)
	require.Equal(t, OrderStatusPending, order.Status)
	require.Equal(t, int64(500*2+399*3), order.TotalAmount)
	require.Equal(t, 2, len(order.Items))
	//require.Equal(t, 1, len(events))
	//require.Equal(t, "order.created", events[0].EventName())
	//createdEvent, ok := events[0].(OrderCreatedEvent)
	//require.True(t, ok, "event is not OrderCreatedEvent")
	//require.Equal(t, 2, len(createdEvent.Items))
}

func TestCreateOrder_ErrOrderEmptyItems(t *testing.T) {
	operator := shared.Operator{ID: 1, Username: "test_user"}
	var inputs []OrderItemInput
	order, err := CreateOrder(42, inputs, operator)
	require.Nil(t, order)
	require.ErrorIs(t, err, ErrOrderEmptyItems)
}

func TestCreateOrder_ErrOrderItemInvalidQuantity(t *testing.T) {
	operator := shared.Operator{ID: 1, Username: "test_user"}
	inputs := []OrderItemInput{
		{ProductID: 1, ProductName: "Apple", UnitPrice: 500, Quantity: -2},
		{ProductID: 1, ProductName: "Orange", UnitPrice: 399, Quantity: 3},
	}
	order, err := CreateOrder(42, inputs, operator)
	require.Nil(t, order)
	require.ErrorIs(t, err, ErrOrderItemInvalidQuantity)
}

func TestOrder_BuildEvents(t *testing.T) {
	order := &Order{
		ID: 42,
		Items: []*OrderItem{
			{ProductID: 1},
			{ProductID: 2},
		},
	}

	events := order.BuildEvents()

	require.NotNil(t, events)
	require.Equal(t, 1, len(events))
	require.Equal(t, int64(42), events[0].OrderID)
	require.Equal(t, 2, len(events[0].Items))
}
