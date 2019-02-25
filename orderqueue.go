package orderbook

import (
	"container/list"
	"encoding/json"
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
	oq.volume = oq.volume.Add(o.Quantity())
	return oq.orders.PushBack(o)
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
	oq.volume = oq.volume.Sub(e.Value.(*Order).Quantity())
	return oq.orders.Remove(e).(*Order)
}

// String implements fmt.Stringer interface
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

// MarshalJSON implements json.Marshaler interface
func (oq *OrderQueue) MarshalJSON() ([]byte, error) {
	iter := oq.Head()

	var orders []*Order
	for iter != nil {
		orders = append(orders, iter.Value.(*Order))
		iter = iter.Next()
	}

	return json.Marshal(
		&struct {
			Volume decimal.Decimal `json:"volume"`
			Price  decimal.Decimal `json:"price"`
			Orders []*Order        `json:"orders"`
		}{
			Volume: oq.Volume(),
			Price:  oq.Price(),
			Orders: orders,
		},
	)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (oq *OrderQueue) UnmarshalJSON(data []byte) error {
	obj := struct {
		Volume decimal.Decimal `json:"volume"`
		Price  decimal.Decimal `json:"price"`
		Orders []*Order        `json:"orders"`
	}{}

	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	oq.volume = obj.Volume
	oq.price = obj.Price
	oq.orders = list.New()
	for _, order := range obj.Orders {
		oq.orders.PushBack(order)
	}
	return nil
}
