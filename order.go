package orderbook

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	timestamp time.Time
	quantity  decimal.Decimal
	price     decimal.Decimal
	orderID   string

	orderList *OrderQueue
	next      *Order
	prev      *Order
}

func NewOrderFromMap(quote map[string]string, orderQueue *OrderQueue) *Order {
	timestamp, _ := time.Parse(time.RFC3339Nano, quote["timestamp"])
	quantity, _ := decimal.NewFromString(quote["quantity"])
	price, _ := decimal.NewFromString(quote["price"])
	orderID := quote["order_id"]
	return NewOrder(orderQueue, orderID, quantity, price, timestamp)
}

func NewOrder(orderQueue *OrderQueue, orderID string, quantity, price decimal.Decimal, timestamp time.Time) *Order {
	return &Order{
		timestamp: timestamp,
		orderID:   orderID,
		quantity:  quantity,
		price:     price,
		orderList: orderQueue,
	}
}

func (o *Order) Next() *Order {
	return o.next
}

func (o *Order) Prev() *Order {
	return o.prev
}

func (o *Order) Update(quantity decimal.Decimal, timestamp time.Time) {
	if quantity.GreaterThan(o.quantity) && o.orderList.tail != o {
		o.orderList.MoveToTail(o)
	}
	o.orderList.volume = o.orderList.volume.Sub(o.quantity.Sub(quantity))
	o.timestamp = timestamp
	o.quantity = quantity
}
