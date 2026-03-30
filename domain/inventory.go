package domain

type Inventory struct {
	ID        int64
	ProductID int64
	Stock     int
}

func (inv *Inventory) Deduct(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidDeductQuantity
	}
	if inv.Stock < quantity {
		return ErrInsufficientStock
	}
	inv.Stock -= quantity
	return nil
}
