package constant

import "errors"

var (
	ErrMerchantNotFound       = errors.New("merchant not found")
	ErrInvalidMerchant        = errors.New("invalid merchant")
	ErrInvalidItem            = errors.New("invalid item")
	ErrInvalidStartingPoint   = errors.New("invalid starting point")
	ErrOrderNotFound          = errors.New("order not found")
	ErrInvalidOrderStatus     = errors.New("invalid order status")
	ErrCannotUpdateOrder      = errors.New("cannot update order")
	ErrInsufficientOrderItems = errors.New("insufficient order items")
	ErrOrderNotAuthorized     = errors.New("order not authorized")
)
