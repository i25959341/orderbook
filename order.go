package orderbook

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// Order strores information about request
type Order struct {
	side      Side
	id        string
	timestamp time.Time
	quantity  decimal.Decimal
	price     decimal.Decimal
}

// NewOrder creates new constant object Order
func NewOrder(orderID string, side Side, quantity, price decimal.Decimal, timestamp time.Time) *Order {
	return &Order{
		id:        orderID,
		side:      side,
		quantity:  quantity,
		price:     price,
		timestamp: timestamp,
	}
}

// ID returns orderID field copy
func (o *Order) ID() string {
	return o.id
}

// Side returns side of the order
func (o *Order) Side() Side {
	return o.side
}

// Quantity returns quantity field copy
func (o *Order) Quantity() decimal.Decimal {
	return o.quantity
}

// Price returns price field copy
func (o *Order) Price() decimal.Decimal {
	return o.price
}

// Time returns timestamp field copy
func (o *Order) Time() time.Time {
	return o.timestamp
}

func (o *Order) String() string {
	return fmt.Sprintf("\n\"%s\":\n\tside: %s\n\tquantity: %s\n\tprice: %s\n\ttime: %s\n", o.ID(), o.Side(), o.Quantity(), o.Price(), o.Time())
}
