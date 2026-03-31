package domain

import "time"

type Sale struct {
	ID          int64
	ProductID   int64
	ProductName string
	UnitPrice   int64
	MaxPerUser  int
	StartAt     time.Time
	EndAt       time.Time
}

type SaleOrderItem struct {
	SaleID   int64
	Quantity int
}

func (s *Sale) IsActive(now time.Time) bool {
	return now.After(s.StartAt) && now.Before(s.EndAt)
}

func (s *Sale) CheckPurchaseLimit(currentCount int, newQuantity int) error {
	if currentCount+newQuantity > s.MaxPerUser {
		return ErrExceedPurchaseLimit
	}
	return nil
}
