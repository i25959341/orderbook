package orderbook

import "errors"

var (
	ErrInvalidQuantity = errors.New("orderbook: invalid order quantity")
	ErrInvalidPrice    = errors.New("orderbook: invalid order price")
	ErrInvalidOrder    = errors.New("orderbook: invalid order ")
	ErrAlreadyLinked   = errors.New("orderbook: order links are not empty")
	ErrPriceNotExists  = errors.New("orderbook: price level does not exist")
	ErrOrderExists     = errors.New("orderbook: order already exists")
	ErrOrderNotExists  = errors.New("orderbook: order does not exist")
)
