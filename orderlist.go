package orderbook

import (
	"github.com/HuKeping/rbtree"
	"github.com/shopspring/decimal"
)

type Item interface {
	Less(than Item) bool
}

type OrderList struct {
	head_order *Order
	tail_order *Order
	length     int
	volume     decimal.Decimal
	last_order *Order
	price      decimal.Decimal
}

func NewOrderList(price decimal.Decimal) *OrderList {
	return &OrderList{head_order: nil, tail_order: nil, length: 0, volume: decimal.Zero,
		last_order: nil, price: price}
}

func (orderlist *OrderList) Less(than rbtree.Item) bool {
	return orderlist.price.LessThan(than.(*OrderList).price)
}

func (orderlist *OrderList) Length() int {
	return orderlist.length
}

func (orderlist *OrderList) HeadOrder() *Order {
	return orderlist.head_order
}

func (orderlist *OrderList) AppendOrder(order *Order) {
	if orderlist.Length() == 0 {
		order.next_order = nil
		order.prev_order = nil
		orderlist.head_order = order
		orderlist.tail_order = order
	} else {
		order.prev_order = orderlist.tail_order
		order.next_order = nil
		orderlist.tail_order.next_order = order
		orderlist.tail_order = order
	}
	orderlist.length = orderlist.length + 1
	orderlist.volume = orderlist.volume.Add(order.quantity)
}

func (orderlist *OrderList) RemoveOrder(order *Order) {
	orderlist.volume = orderlist.volume.Sub(order.quantity)
	orderlist.length = orderlist.length - 1
	if orderlist.length == 0 {
		return
	}

	next_order := order.next_order
	prev_order := order.prev_order

	if next_order != nil && prev_order != nil {
		next_order.prev_order = prev_order
		prev_order.next_order = next_order
	} else if next_order != nil {
		next_order.prev_order = nil
		orderlist.head_order = next_order
	} else if prev_order != nil {
		prev_order.next_order = nil
		orderlist.tail_order = prev_order
	}
}

func (orderlist *OrderList) MoveToTail(order *Order) {
	if order.prev_order != nil { // This Order is not the first Order in the OrderList
		order.prev_order.next_order = order.next_order // Link the previous Order to the next Order, then move the Order to tail
	} else { // This Order is the first Order in the OrderList
		orderlist.head_order = order.next_order // Make next order the first
	}
	order.next_order.prev_order = order.prev_order

	// Move Order to the last position. Link up the previous last position Order.
	orderlist.tail_order.next_order = order
	orderlist.tail_order = order
}
