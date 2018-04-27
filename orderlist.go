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
		order.nextOrder = nil
		order.prevOrder = nil
		orderlist.head_order = order
		orderlist.tail_order = order
	} else {
		order.prevOrder = orderlist.tail_order
		order.nextOrder = nil
		orderlist.tail_order.nextOrder = order
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

	nextOrder := order.nextOrder
	prevOrder := order.prevOrder

	if nextOrder != nil && prevOrder != nil {
		nextOrder.prevOrder = prevOrder
		prevOrder.nextOrder = nextOrder
	} else if nextOrder != nil {
		nextOrder.prevOrder = nil
		orderlist.head_order = nextOrder
	} else if prevOrder != nil {
		prevOrder.nextOrder = nil
		orderlist.tail_order = prevOrder
	}
}

func (orderlist *OrderList) MoveToTail(order *Order) {
	if order.prevOrder != nil { // This Order is not the first Order in the OrderList
		order.prevOrder.nextOrder = order.nextOrder // Link the previous Order to the next Order, then move the Order to tail
	} else { // This Order is the first Order in the OrderList
		orderlist.head_order = order.nextOrder // Make next order the first
	}
	order.nextOrder.prevOrder = order.prevOrder

	// Move Order to the last position. Link up the previous last position Order.
	orderlist.tail_order.nextOrder = order
	orderlist.tail_order = order
}
