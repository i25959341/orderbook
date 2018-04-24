package orderbook

import (
	"strconv"

	"github.com/shopspring/decimal"
)

type Order struct {
	timestamp  int
	quantity   decimal.Decimal
	price      decimal.Decimal
	order_id   string
	trade_id   string
	next_order *Order
	prev_order *Order
	order_list *OrderList
}

func NewOrder(quote map[string]string, order_list *OrderList) *Order {
	timestamp, _ := strconv.Atoi(quote["timestamp"])
	quantity, _ := decimal.NewFromString(quote["quantity"])
	price, _ := decimal.NewFromString(quote["price"])
	order_id := quote["order_id"]
	trade_id := quote["trade_id"]
	return &Order{timestamp: timestamp, quantity: quantity, price: price, order_id: order_id,
		trade_id: trade_id, next_order: nil, prev_order: nil, order_list: order_list}
}

func (o *Order) NextOrder() *Order {
	return o.next_order
}

func (o *Order) PrevOrder() *Order {
	return o.prev_order
}

func (o *Order) UpdateQuantity(new_quantity decimal.Decimal, new_timestamp int) {
	if new_quantity.GreaterThan(o.quantity) && o.order_list.tail_order != o {
		o.order_list.MoveToTail(o)
	}
	o.order_list.volume = o.order_list.volume.Sub(o.quantity.Sub(new_quantity))
	o.timestamp = new_timestamp
	o.quantity = new_quantity
}
