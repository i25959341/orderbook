package orderbook

import (
	"strconv"

	"github.com/shopspring/decimal"
)

type Order struct {
	timestamp int
	quantity  decimal.Decimal
	price     decimal.Decimal
	orderID   string
	tradeID   string

	orderList *OrderQueue
	next      *Order
	prev      *Order
}

func NewOrderFromMap(quote map[string]string, orderQueue *OrderQueue) *Order {
	timestamp, _ := strconv.Atoi(quote["timestamp"])
	quantity, _ := decimal.NewFromString(quote["quantity"])
	price, _ := decimal.NewFromString(quote["price"])
	orderID := quote["order_id"]
	tradeID := quote["trade_id"]
	return NewOrder(orderQueue, orderID, tradeID, quantity, price, timestamp)
}

func NewOrder(orderQueue *OrderQueue, orderID, tradeID string, quantity, price decimal.Decimal, timestamp int) *Order {
	return &Order{
		timestamp: timestamp,
		orderID:   orderID,
		tradeID:   tradeID,
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

func (o *Order) Update(quantity decimal.Decimal, timestamp int) {
	if quantity.GreaterThan(o.quantity) && o.orderList.tail != o {
		o.orderList.MoveToTail(o)
	}
	o.orderList.volume = o.orderList.volume.Sub(o.quantity.Sub(quantity))
	o.timestamp = timestamp
	o.quantity = quantity
}
