package domain

import "errors"

var ErrRegionAbandoned = errors.New("region abandoned")
var ErrRegionNotAllowedStatus = errors.New("region not allowed status")
var ErrOrderItemInvalidQuantity = errors.New("order item: quantity must be greater than zero")
var ErrOrderItemInvalidPrice = errors.New("order item: unit price must be greater than zero")
var ErrOrderEmptyItems = errors.New("order: must contain at least one item")
var ErrDuplicateIdempotencyKey = errors.New("idempotency key already exists")
