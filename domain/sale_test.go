// domain/sale_test.go
package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSale_IsActive(t *testing.T) {
	sale := &Sale{
		StartAt: time.Now().Add(-1 * time.Hour),
		EndAt:   time.Now().Add(1 * time.Hour),
	}

	require.True(t, sale.IsActive(time.Now()))
}

func TestSale_IsActive_BeforeStart(t *testing.T) {
	sale := &Sale{
		StartAt: time.Now().Add(1 * time.Hour),
		EndAt:   time.Now().Add(2 * time.Hour),
	}

	require.False(t, sale.IsActive(time.Now()))
}

func TestSale_IsActive_AfterEnd(t *testing.T) {
	sale := &Sale{
		StartAt: time.Now().Add(-2 * time.Hour),
		EndAt:   time.Now().Add(-1 * time.Hour),
	}

	require.False(t, sale.IsActive(time.Now()))
}

func TestSale_CheckPurchaseLimit(t *testing.T) {
	sale := &Sale{MaxPerUser: 2}

	require.NoError(t, sale.CheckPurchaseLimit(0, 1))
	require.NoError(t, sale.CheckPurchaseLimit(1, 1))
	require.NoError(t, sale.CheckPurchaseLimit(0, 2))
}

func TestSale_CheckPurchaseLimit_Exceed(t *testing.T) {
	sale := &Sale{MaxPerUser: 2}

	require.ErrorIs(t, sale.CheckPurchaseLimit(1, 2), ErrExceedPurchaseLimit)
	require.ErrorIs(t, sale.CheckPurchaseLimit(2, 1), ErrExceedPurchaseLimit)
}
