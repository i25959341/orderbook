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

	orderList *OrderList
	next      *Order
	prev      *Order
}

func NewOrderFromMap(quote map[string]string, orderList *OrderList) *Order {
	timestamp, _ := strconv.Atoi(quote["timestamp"])
	quantity, _ := decimal.NewFromString(quote["quantity"])
	price, _ := decimal.NewFromString(quote["price"])
	orderID := quote["order_id"]
	tradeID := quote["trade_id"]
	return NewOrder(orderList, orderID, tradeID, quantity, price, timestamp)
}

func NewOrder(orderList *OrderList, orderID, tradeID string, quantity, price decimal.Decimal, timestamp int) *Order {
	return &Order{
		timestamp: timestamp,
		orderID:   orderID,
		tradeID:   tradeID,
		quantity:  quantity,
		price:     price,
		orderList: orderList,
		next:      nil,
		prev:      nil,
	}
}

func (o *Order) Next() *Order {
	return o.next
}

func (o *Order) Prev() *Order {
	return o.prev
}

func (o *Order) Update(quantity decimal.Decimal, timestamp int) {
	if quantity.GreaterThan(o.quantity) && o.orderList.tailOrder != o {
		o.orderList.MoveToTail(o)
	}
	o.orderList.volume = o.orderList.volume.Sub(o.quantity.Sub(quantity))
	o.timestamp = timestamp
	o.quantity = quantity
}
