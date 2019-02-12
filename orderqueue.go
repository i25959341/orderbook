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
func (oq *OrderQueue) Head() *list.Element {
	return oq.orders.Front()
}

// Tail returns bottom order in queue
func (oq *OrderQueue) Tail() *list.Element {
	return oq.orders.Back()
}

// Append adds order to tail of the queue
func (oq *OrderQueue) Append(o *Order) *list.Element {
	e := oq.orders.PushBack(o)
	if e != nil {
		oq.volume = oq.volume.Add(o.Quantity())
		return e
	}
	return nil
}

// Update sets up new order to list value
func (oq *OrderQueue) Update(e *list.Element, o *Order) *list.Element {
	oq.volume = oq.volume.Sub(e.Value.(*Order).Quantity())
	oq.volume = oq.volume.Add(o.Quantity())
	e.Value = o
	return e
}

// Remove removes order from the queue and link order chain
func (oq *OrderQueue) Remove(e *list.Element) *Order {
	o := oq.orders.Remove(e)
	if o != nil {
		oq.volume = oq.volume.Sub(o.(*Order).Quantity())
		return o.(*Order)
	}
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
