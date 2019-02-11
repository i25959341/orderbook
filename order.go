package orderbook

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// Order strores information about request
type Order struct {
	timestamp time.Time
	quantity  decimal.Decimal
	price     decimal.Decimal
	orderID   string

	queue *OrderQueue
	next  *Order
	prev  *Order
}

// NewOrder creates new constant object Order
func NewOrder(orderID string, quantity, price decimal.Decimal, timestamp time.Time) *Order {
	return &Order{
		timestamp: timestamp,
		orderID:   orderID,
		quantity:  quantity,
		price:     price,
	}
}

// ID returns orderID field copy
func (o *Order) ID() string {
	return o.orderID
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

// Queue returns pointer to queue field
func (o *Order) Queue() *OrderQueue {
	return o.queue
}

// Next returns pointer to next Order in queue
func (o *Order) Next() *Order {
	return o.next
}

// Next returns pointer to previous Order in queue
func (o *Order) Prev() *Order {
	return o.prev
}

// // Update updates order quantity and moves it to the end of chain
// func (o *Order) Update(quantity decimal.Decimal, timestamp time.Time) {
// 	if quantity.GreaterThan(o.quantity) && o.queue.tail != o {
// 		o.queue.MoveToTail(o)
// 	}
// 	o.queue.volume = o.queue.volume.Sub(o.quantity.Sub(quantity))
// 	o.timestamp = timestamp
// 	o.quantity = quantity
// }

func (o *Order) String() string {
	return fmt.Sprintf("\n\"%s\":\n\tquantity: %s\n\tprice: %s\n\ttime: %s\n\tnext: %v\n\tprev: %v\n\tqueue: %v\n", o.ID(), o.Quantity(), o.Price(), o.Time(), o.Next(), o.Prev(), o.Queue())
}
