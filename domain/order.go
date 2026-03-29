package domain

import (
	"MyMall/domain/shared"
	"time"

	"github.com/samber/lo"
)

const (
	OrderStatusPending int8 = iota
	OrderStatusPaid
	OrderStatusShipped
	OrderStatusDelivered
	OrderStatusCancelled
)

var orderAllowedTransitions = map[int8][]int8{
	OrderStatusPending: {OrderStatusPaid, OrderStatusCancelled},
	OrderStatusPaid:    {OrderStatusShipped, OrderStatusCancelled},
	OrderStatusShipped: {OrderStatusDelivered},
}

type OrderItemInput struct {
	ProductID   int64
	ProductName string
	UnitPrice   int64
	Quantity    int
}

type Order struct {
	ID           int64
	CustomerID   int64
	CreatorID    int64
	CreatorName  string
	ModifierID   int64
	ModifierName string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Status       int8
	TotalAmount  int64 // int cents
	Items        []*OrderItem
}

func CreateOrder(customerID int64, inputs []OrderItemInput, operator shared.Operator) (*Order, []Event, error) {
	if len(inputs) == 0 {
		return nil, nil, ErrOrderEmptyItems
	}
	items, err := lo.MapErr(inputs, func(input OrderItemInput, _ int) (*OrderItem, error) {
		return newOrderItem(input.ProductID, input.ProductName, input.UnitPrice, input.Quantity)
	})
	if err != nil {
		return nil, nil, err
	}
	totalAmount := lo.SumBy(items, func(item *OrderItem) int64 {
		return item.Subtotal
	})
	order := &Order{
		CustomerID:  customerID,
		CreatorID:   operator.ID,
		CreatorName: operator.Username,
		Status:      OrderStatusPending,
		TotalAmount: totalAmount,
		Items:       items,
	}

	events := []Event{
		OrderCreatedEvent{OrderID: order.ID, CreatedAt: time.Now()},
	}
	return order, events, nil
}
