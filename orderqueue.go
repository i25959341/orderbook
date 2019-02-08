package orderbook

import (
	"github.com/shopspring/decimal"
)

type OrderQueue struct {
	volume decimal.Decimal
	price  decimal.Decimal

	length int
	head   *Order
	tail   *Order
	last   *Order
}

func NewOrderQueue(price decimal.Decimal) *OrderQueue {
	return &OrderQueue{
		price:  price,
		volume: decimal.Zero,
	}
}

func (oq *OrderQueue) Length() int {
	return oq.length
}

func (oq *OrderQueue) Head() *Order {
	return oq.head
}

func (oq *OrderQueue) Append(order *Order) {
	if oq.Length() == 0 {
		order.next = nil
		order.prev = nil
		oq.head = order
		oq.tail = order
	} else {
		order.prev = oq.tail
		order.next = nil
		oq.tail.next = order
		oq.tail = order
	}
	oq.length = oq.length + 1
	oq.volume = oq.volume.Add(order.quantity)
}

func (oq *OrderQueue) Remove(order *Order) {
	oq.volume = oq.volume.Sub(order.quantity)
	oq.length = oq.length - 1
	if oq.length == 0 {
		return
	}

	nextOrder := order.next
	prevOrder := order.prev

	if nextOrder != nil && prevOrder != nil {
		nextOrder.prev = prevOrder
		prevOrder.next = nextOrder
	} else if nextOrder != nil {
		nextOrder.prev = nil
		oq.head = nextOrder
	} else if prevOrder != nil {
		prevOrder.next = nil
		oq.tail = prevOrder
	}
}

func (oq *OrderQueue) MoveToTail(order *Order) {
	if order.prev != nil {
		order.prev.next = order.next
	} else {
		oq.head = order.next
	}
	order.next.prev = order.prev
	oq.tail.next = order
	oq.tail = order
}
