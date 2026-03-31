package domain

import "errors"

var ErrRegionAbandoned = errors.New("region abandoned")
var ErrRegionNotAllowedStatus = errors.New("region not allowed status")
var ErrOrderItemInvalidQuantity = errors.New("order item: quantity must be greater than zero")
var ErrOrderItemInvalidPrice = errors.New("order item: unit price must be greater than zero")
var ErrOrderEmptyItems = errors.New("order: must contain at least one item")
var ErrDuplicateIdempotencyKey = errors.New("idempotency key already exists")
var ErrInsufficientStock = errors.New("inventory: insufficient stock")
var ErrInvalidDeductQuantity = errors.New("inventory: deduct quantity must be greater than zero")
var ErrDuplicateProcessedEvent = errors.New("processed event: duplicate event")
var ErrExceedPurchaseLimit = errors.New("sale: exceed purchase limit")
var ErrSaleNotFound = errors.New("sale: not found")
var ErrSaleNotActive = errors.New("sale: not active")
var ErrOrderNotFound = errors.New("order: not found")
