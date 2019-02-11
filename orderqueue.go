package orderbook

import (
	"container/list"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// OrderQueue stores and manage chain of orders
type OrderQueue struct {
	volume decimal.Decimal
	price  decimal.Decimal
	orders *list.List
}

// NewOrderQueue creates and initialize OrderQueue object
func NewOrderQueue(price decimal.Decimal) *OrderQueue {
	return &OrderQueue{
		price:  price,
		volume: decimal.Zero,
		orders: list.New(),
	}
}

// Len returns amount of orders in queue
func (oq *OrderQueue) Len() int {
	return oq.orders.Len()
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
	return oq.orders.Front().Value.(*Order)
}

// Tail returns bottom order in queue
func (oq *OrderQueue) Tail() *Order {
	return oq.orders.Back().Value.(*Order)
}

// Append adds order to tail of the queue
func (oq *OrderQueue) Append(order *Order) error {
	if order.owner != nil {
		return ErrAlreadyLinked
	}
	if order.quantity.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidQuantity
	}
	if order.price.LessThanOrEqual(decimal.Zero) || !order.price.Equal(oq.price) {
		return ErrInvalidPrice
	}

	element := oq.orders.PushBack(order)
	if element == nil {
		return ErrInvalidOrder
	}
	order.owner = element
	oq.volume = oq.volume.Add(order.quantity)

	return nil
}

// Remove removes order from the queue and link order chain
func (oq *OrderQueue) Remove(order *Order) error {
	if order.owner == nil {
		return ErrOrderNotExists
	}

	result := oq.orders.Remove(order.owner)
	if result == nil {
		return ErrOrderNotExists
	}

	order.owner = nil
	oq.volume = oq.volume.Sub(order.quantity)
	return nil
}

func (oq *OrderQueue) String() string {
	sb := strings.Builder{}
	iter := oq.orders.Front()
	sb.WriteString(fmt.Sprintf("\nqueue length: %d, price: %s, volume: %s, orders:", oq.Len(), oq.Price(), oq.Volume()))
	for iter != nil {
		order := iter.Value.(*Order)
		str := fmt.Sprintf("\n\tid: %s, volume: %s, time: %s", order.ID(), order.Quantity(), order.Price())
		sb.WriteString(str)
		iter = iter.Next()
	}
	return sb.String()
}
