package orderbook

import (
	"container/list"
	"time"

	"github.com/shopspring/decimal"
)

// Orderbook implements standard matching algorithm
type Orderbook struct {
	orders map[string]*list.Element // orderID -> *Order (*list.Element.Value.(*Order))

	asks *OrderTree
	bids *OrderTree
}

// NewOrder creates Orderbook object
func NewOrderbook() *Orderbook {
	return &Orderbook{
		orders: map[string]*list.Element{},
		bids:   NewOrderTree(),
		asks:   NewOrderTree(),
	}
}

// ProcessLimitOrder places limit order to the orderbook
func (ob *Orderbook) ProcessLimitOrder(side Side, orderID string, quantity, price decimal.Decimal) (done []*Order, partial *Order, err error) {
	if _, ok := ob.orders[orderID]; ok {
		return nil, nil, ErrOrderExists
	}

	if quantity.Sign() <= 0 {
		return nil, nil, ErrInvalidQuantity
	}

	if price.Sign() <= 0 {
		return nil, nil, ErrInvalidPrice
	}

	if side == Buy {
		minPrice := ob.asks.MinPriceQueue()
		for quantity.Sign() > 0 && ob.asks.Len() > 0 && price.GreaterThanOrEqual(minPrice.Price()) {
			ordersDone, partialDone, quantityLeft := ob.processQueue(minPrice, quantity)
			done = append(done, ordersDone...)
			partial = partialDone
			quantity = quantityLeft
			minPrice = ob.asks.MinPriceQueue()
		}

		o := NewOrder(orderID, side, quantity, price, time.Now().UTC())
		if quantity.Sign() > 0 {
			partial = o
			ob.orders[orderID] = ob.bids.Append(o)
		} else {
			done = append(done, o)
		}
	} else {
		maxPrice := ob.bids.MaxPriceQueue()
		for quantity.Sign() > 0 && ob.bids.Len() > 0 && price.LessThanOrEqual(maxPrice.Price()) {
			ordersDone, partialDone, quantityLeft := ob.processQueue(maxPrice, quantity)
			done = append(done, ordersDone...)
			partial = partialDone
			quantity = quantityLeft
			maxPrice = ob.bids.MaxPriceQueue()
		}

		o := NewOrder(orderID, side, quantity, price, time.Now().UTC())
		if quantity.Sign() > 0 {
			partial = o
			ob.orders[orderID] = ob.asks.Append(o)
		} else {
			done = append(done, o)
		}
	}
	return
}

func (ob *Orderbook) processQueue(orderQueue *OrderQueue, quantityToTrade decimal.Decimal) (done []*Order, partial *Order, quantityLeft decimal.Decimal) {
	quantityLeft = quantityToTrade

	for orderQueue.Len() > 0 && quantityLeft.Sign() > 0 {
		headOrderEl := orderQueue.Head()
		headOrder := headOrderEl.Value.(*Order)

		if quantityLeft.LessThan(headOrder.Quantity()) {
			partial = NewOrder(headOrder.ID(), headOrder.Side(), headOrder.Quantity().Sub(quantityLeft), headOrder.Price(), headOrder.Time())
			orderQueue.Update(headOrderEl, partial)
			quantityLeft = decimal.Zero
		} else {
			quantityLeft = quantityLeft.Sub(headOrder.Quantity())
			done = append(done, ob.CancelOrder(headOrder.ID()))
		}
	}

	return
}

// CancelOrder removes order from orderbook
func (ob *Orderbook) CancelOrder(orderID string) *Order {
	e, ok := ob.orders[orderID]
	if !ok {
		return nil
	}

	delete(ob.orders, orderID)

	if e.Value.(*Order).Side() == Buy {
		return ob.bids.Remove(e)
	}

	return ob.asks.Remove(e)
}

func (ob *Orderbook) String() string {
	return ob.asks.String() + "\r\n------------------------------------" + ob.bids.String()
}
