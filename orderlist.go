package orderbook

import (
	"github.com/HuKeping/rbtree"
	"github.com/shopspring/decimal"
)

type Item interface {
	Less(than Item) bool
}

type OrderList struct {
	headOrder *Order
	tailOrder *Order
	length    int
	volume    decimal.Decimal
	lastOrder *Order
	price     decimal.Decimal
}

func NewOrderList(price decimal.Decimal) *OrderList {
	return &OrderList{headOrder: nil, tailOrder: nil, length: 0, volume: decimal.Zero,
		lastOrder: nil, price: price}
}

func (orderlist *OrderList) Less(than rbtree.Item) bool {
	return orderlist.price.LessThan(than.(*OrderList).price)
}

func (orderlist *OrderList) Length() int {
	return orderlist.length
}

func (orderlist *OrderList) HeadOrder() *Order {
	return orderlist.headOrder
}

func (orderlist *OrderList) AppendOrder(order *Order) {
	if orderlist.Length() == 0 {
		order.next = nil
		order.prev = nil
		orderlist.headOrder = order
		orderlist.tailOrder = order
	} else {
		order.prev = orderlist.tailOrder
		order.next = nil
		orderlist.tailOrder.next = order
		orderlist.tailOrder = order
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

	nextOrder := order.next
	prevOrder := order.prev

	if nextOrder != nil && prevOrder != nil {
		nextOrder.prev = prevOrder
		prevOrder.next = nextOrder
	} else if nextOrder != nil {
		nextOrder.prev = nil
		orderlist.headOrder = nextOrder
	} else if prevOrder != nil {
		prevOrder.next = nil
		orderlist.tailOrder = prevOrder
	}
}

func (orderlist *OrderList) MoveToTail(order *Order) {
	if order.prev != nil { // This Order is not the first Order in the OrderList
		order.prev.next = order.next // Link the previous Order to the next Order, then move the Order to tail
	} else { // This Order is the first Order in the OrderList
		orderlist.headOrder = order.next // Make next order the first
	}
	order.next.prev = order.prev

	// Move Order to the last position. Link up the previous last position Order.
	orderlist.tailOrder.next = order
	orderlist.tailOrder = order
}
