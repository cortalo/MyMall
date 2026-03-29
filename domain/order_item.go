package domain

type OrderItem struct {
	ID          int64
	OrderID     int64
	ProductID   int64
	ProductName string
	UnitPrice   int64 // in cents
	Quantity    int
	Subtotal    int64 // in cents, UnitPrice * Quantity
}

func newOrderItem(productID int64, productName string, unitPrice int64, quantity int) (*OrderItem, error) {
	if quantity <= 0 {
		return nil, ErrOrderItemInvalidQuantity
	}
	if unitPrice <= 0 {
		return nil, ErrOrderItemInvalidPrice
	}
	return &OrderItem{
		ProductID:   productID,
		ProductName: productName,
		UnitPrice:   unitPrice,
		Quantity:    quantity,
		Subtotal:    unitPrice * int64(quantity),
	}, nil
}
