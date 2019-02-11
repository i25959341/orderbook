package orderbook

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidQuantity = errors.New("orderqueue: invalid order quantity")
	ErrInvalidPrice    = errors.New("orderqueue: invalid order price")
	ErrInvalidQueue    = errors.New("orderqueue: invalid order queue")
	ErrAlreadyLinked   = errors.New("orderqueue: order links are not empty")
)

// OrderQueue stores and manage chain of orders
type OrderQueue struct {
	volume decimal.Decimal
	price  decimal.Decimal
	length int

	head *Order
	tail *Order
	last *Order
}

// NewOrderQueue creates and initialize OrderQueue object
func NewOrderQueue(price decimal.Decimal) *OrderQueue {
	return &OrderQueue{
		price:  price,
		volume: decimal.Zero,
	}
}

// Length returns amount of orders in queue
func (oq *OrderQueue) Length() int {
	return oq.length
}

// Price returns price level of the queue
func (oq *OrderQueue) Price() decimal.Decimal {
	return oq.price
}

// Volume returns total orders volume
func (oq *OrderQueue) Volume() decimal.Decimal {
	return oq.volume
}

// Head returns top order in queue
func (oq *OrderQueue) Head() *Order {
	return oq.head
}

// Tail returns bottom order in queue
func (oq *OrderQueue) Tail() *Order {
	return oq.tail
}

// Append adds order to tail of the queue
func (oq *OrderQueue) Append(order *Order) error {
	if order.quantity.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidQuantity
	}

	if order.price.LessThanOrEqual(decimal.Zero) || !order.price.Equal(oq.price) {
		return ErrInvalidPrice
	}

	if order.queue != nil || order.next != nil || order.prev != nil {
		return ErrAlreadyLinked
	}

	if oq.Length() == 0 {
		oq.head = order
		oq.tail = order
	} else {
		order.prev = oq.tail
		oq.tail.next = order
		oq.tail = order
	}

	oq.length = oq.length + 1
	oq.volume = oq.volume.Add(order.quantity)

	order.queue = oq
	return nil
}

// Remove removes order from the queue and link order chain
func (oq *OrderQueue) Remove(order *Order) error {
	if order.queue != oq {
		return ErrInvalidQueue
	}

	nextOrder := order.next
	prevOrder := order.prev

	order.next = nil
	order.prev = nil
	order.queue = nil

	oq.volume = oq.volume.Sub(order.quantity)
	oq.length = oq.length - 1
	if oq.length == 0 {
		oq.head = nil
		oq.tail = nil
		return nil
	}

	if nextOrder != nil && prevOrder != nil {
		nextOrder.prev = prevOrder
		prevOrder.next = nextOrder
		return nil
	}

	if nextOrder != nil {
		nextOrder.prev = nil
		oq.head = nextOrder
		return nil
	}

	prevOrder.next = nil
	oq.tail = prevOrder
	return nil
}

// MoveToTail moves order to end of the chain
func (oq *OrderQueue) MoveToTail(order *Order) error {
	if order.queue != oq {
		return ErrInvalidQueue
	}

	if order == oq.tail {
		return nil
	}

	if order.prev == nil {
		oq.head = order.next
	}

	if order.prev != nil {
		order.prev.next = order.next
	}

	order.next.prev = order.prev
	oq.tail.next = order
	order.prev = oq.tail
	order.next = nil
	oq.tail = order

	return nil
}

func (oq *OrderQueue) String() string {
	sb := strings.Builder{}
	iter := oq.head
	sb.WriteString(fmt.Sprintf("\nqueue length: %d, price: %s, volume: %s, orders:", oq.Length(), oq.Price(), oq.Volume()))
	for iter != nil {
		str := fmt.Sprintf("\n\tid: %s, volume: %s, time: %s", iter.orderID, iter.quantity, iter.timestamp)
		sb.WriteString(str)
		iter = iter.next
	}
	return sb.String()
}
