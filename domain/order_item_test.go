package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrderItem_NewOrderItem(t *testing.T) {
	orderItem, err := newOrderItem(1, "product", 199, 3)
	require.NoError(t, err)
	require.NotNil(t, orderItem)
	require.Equal(t, int64(0), orderItem.ID)
	require.Equal(t, int64(1), orderItem.ProductID)
	require.Equal(t, "product", orderItem.ProductName)
	require.Equal(t, int64(199), orderItem.UnitPrice)
	require.Equal(t, 3, orderItem.Quantity)
	require.Equal(t, int64(199*3), orderItem.Subtotal)
}

func TestOrderItem_NewOrderItemErrInvalidQuantity(t *testing.T) {
	orderItem, err := newOrderItem(1, "product", 199, -1)
	require.Nil(t, orderItem)
	require.ErrorIs(t, err, ErrOrderItemInvalidQuantity)
}

func TestOrderItem_NewOrderItemErrInvalidPrice(t *testing.T) {
	orderItem, err := newOrderItem(1, "product", -1, 3)
	require.Nil(t, orderItem)
	require.ErrorIs(t, err, ErrOrderItemInvalidPrice)
}
