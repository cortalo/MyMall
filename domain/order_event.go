package domain

import "time"

type OrderCreatedEvent struct {
	OrderID   int64
	CreatedAt time.Time
}

func (e OrderCreatedEvent) EventName() string {
	return "order.created"
}
