package orderbook

import (
	"container/list"
	"time"

	"github.com/shopspring/decimal"
)

// OrderBook implements standard matching algorithm
type OrderBook struct {
	orders map[string]*list.Element // orderID -> *Order (*list.Element.Value.(*Order))

	asks *OrderSide
	bids *OrderSide
}

// NewOrderBook creates Orderbook object
func NewOrderBook() *OrderBook {
	return &OrderBook{
		orders: map[string]*list.Element{},
		bids:   NewOrderSide(),
		asks:   NewOrderSide(),
	}
}

// ProcessMarketOrder immediately gets definite quantity from the order book with market price
// Arguments:
//      side     - what do you want to do (ob.Sell or ob.Buy)
//      quantity - how much quantity you want to sell or buy
//      * to create new decimal number you should use decimal.New() func
//        read more at https://github.com/shopspring/decimal
// Return:
//      error        - not nil if price is less or equal 0
//      done         - not nil if your market order produse ends of anoter orders, this order will add to
//                     the "done" slice
//      partial      - not nil if your order has done but top order is not fully done
//      quantityLeft - more than zero if it is not enought orders to process all quantity
func (ob *OrderBook) ProcessMarketOrder(side Side, quantity decimal.Decimal) (done []*Order, partial *Order, quantityLeft decimal.Decimal, err error) {
	if quantity.Sign() <= 0 {
		return nil, nil, decimal.Zero, ErrInvalidQuantity
	}

	if side == Buy {
		for quantity.Sign() > 0 && ob.asks.Len() > 0 {
			bestPriceAsks := ob.asks.MinPriceQueue()
			ordersDone, partialDone, quantityLeft := ob.processQueue(bestPriceAsks, quantity)
			done = append(done, ordersDone...)
			partial = partialDone
			quantity = quantityLeft
		}
	} else {
		for quantity.Sign() > 0 && ob.bids.Len() > 0 {
			bestPriceBids := ob.bids.MaxPriceQueue()
			ordersDone, partialDone, quantityLeft := ob.processQueue(bestPriceBids, quantity)
			done = append(done, ordersDone...)
			partial = partialDone
			quantity = quantityLeft
		}
	}

	quantityLeft = quantity
	return
}

// ProcessLimitOrder places new order to the OrderBook
// Arguments:
//      side     - what do you want to do (ob.Sell or ob.Buy)
//      orderID  - unique order ID in depth
//      quantity - how much quantity you want to sell or buy
//      price    - no more expensive (or cheaper) this price
//      * to create new decimal number you should use decimal.New() func
//        read more at https://github.com/shopspring/decimal
// Return:
//      error   - not nil if quantity (or price) is less or equal 0. Or if order with given ID is exists
//      done    - not nil if your order produse ends of anoter order, this order will add to
//                the "done" slice. If your order have done too, it will be places to this array too
//      partial - not nil if your order has done but top order is not fully done. Or if your order is
//                partial done and placed to the orderbook without full quantity - partial will contain
//                your order with quantity to left
func (ob *OrderBook) ProcessLimitOrder(side Side, orderID string, quantity, price decimal.Decimal) (done []*Order, partial *Order, err error) {
	if _, ok := ob.orders[orderID]; ok {
		return nil, nil, ErrOrderExists
	}

	if quantity.Sign() <= 0 {
		return nil, nil, ErrInvalidQuantity
	}

	if price.Sign() <= 0 {
		return nil, nil, ErrInvalidPrice
	}

	quantityToTrade := quantity
	var sideToAdd *OrderSide
	if side == Buy {
		sideToAdd = ob.bids
		minPrice := ob.asks.MinPriceQueue()
		for quantityToTrade.Sign() > 0 && ob.asks.Len() > 0 && price.GreaterThanOrEqual(minPrice.Price()) {
			ordersDone, partialDone, quantityLeft := ob.processQueue(minPrice, quantityToTrade)
			done = append(done, ordersDone...)
			partial = partialDone
			quantityToTrade = quantityLeft
			minPrice = ob.asks.MinPriceQueue()
		}
	} else {
		sideToAdd = ob.asks
		maxPrice := ob.bids.MaxPriceQueue()
		for quantityToTrade.Sign() > 0 && ob.bids.Len() > 0 && price.LessThanOrEqual(maxPrice.Price()) {
			ordersDone, partialDone, quantityLeft := ob.processQueue(maxPrice, quantityToTrade)
			done = append(done, ordersDone...)
			partial = partialDone
			quantityToTrade = quantityLeft
			maxPrice = ob.bids.MaxPriceQueue()
		}
	}

	if quantityToTrade.Sign() > 0 {
		o := NewOrder(orderID, side, quantityToTrade, price, time.Now().UTC())
		if len(done) > 0 {
			partial = o
		}
		ob.orders[orderID] = sideToAdd.Append(o)
	} else {
		done = append(done, NewOrder(orderID, side, quantity, price, time.Now().UTC()))
	}
	return
}

func (ob *OrderBook) processQueue(orderQueue *OrderQueue, quantityToTrade decimal.Decimal) (done []*Order, partial *Order, quantityLeft decimal.Decimal) {
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

// CancelOrder removes order with given ID from the order book
func (ob *OrderBook) CancelOrder(orderID string) *Order {
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

func (ob *OrderBook) String() string {
	return ob.asks.String() + "\r\n------------------------------------" + ob.bids.String()
}
