package orderbook

import (
	"container/list"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// Order strores information about request
type Order struct {
	id        string
	timestamp time.Time
	quantity  decimal.Decimal
	price     decimal.Decimal

	container *list.Element
	owner     *list.List
}

// NewOrder creates new constant object Order
func NewOrder(orderID string, quantity, price decimal.Decimal, timestamp time.Time) *Order {
	return &Order{
		timestamp: timestamp,
		id:        orderID,
		quantity:  quantity,
		price:     price,
	}
}

// ID returns orderID field copy
func (o *Order) ID() string {
	return o.id
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
	return fmt.Sprintf("\n\"%s\":\n\tquantity: %s\n\tprice: %s\n\ttime: %s\n", o.ID(), o.Quantity(), o.Price(), o.Time())
}
