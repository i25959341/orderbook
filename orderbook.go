package orderbook

import (
	"container/list"
	"encoding/json"
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

// PriceLevel contains price and volume in depth
type PriceLevel struct {
	Price    decimal.Decimal `json:"price"`
	Quantity decimal.Decimal `json:"quantity"`
}

// ProcessMarketOrder immediately gets definite quantity from the order book with market price
// Arguments:
//      side     - what do you want to do (ob.Sell or ob.Buy)
//      quantity - how much quantity you want to sell or buy
//      * to create new decimal number you should use decimal.New() func
//        read more at https://github.com/shopspring/decimal
// Return:
//      error        - not nil if price is less or equal 0
//      done         - not nil if your market order produces ends of anoter orders, this order will add to
//                     the "done" slice
//      partial      - not nil if your order has done but top order is not fully done
//      partialQuantityProcessed - if partial order is not nil this result contains processed quatity from partial order
//      quantityLeft - more than zero if it is not enought orders to process all quantity
func (ob *OrderBook) CalculatePriceAfterExecution(side Side,quantity decimal.Decimal) (price decimal.Decimal, err error) {
	price = decimal.Zero

	var (
		level *OrderQueue
		iter  func(decimal.Decimal) *OrderQueue
	)
	if side == Buy {
		level = ob.asks.MinPriceQueue()
		iter = ob.asks.GreaterThan
	}	else {
		level = ob.bids.MaxPriceQueue()
		iter = ob.bids.LessThan
	}
	for quantity.Sign() > 0 && level != nil {
		levelVolume := level.Volume()
		levelPrice := level.Price()
		if quantity.GreaterThanOrEqual(levelVolume) {
			price = levelPrice
			quantity = quantity.Sub(levelVolume)
			level = iter(levelPrice)
		} else {
			price = levelPrice
			quantity = decimal.Zero
		}
	}

	return


}





func (ob *OrderBook) ProcessMarketOrder(side Side, quantity decimal.Decimal) (done []*Order, partial *Order, partialQuantityProcessed, quantityLeft decimal.Decimal, err error) {
	if quantity.Sign() <= 0 {
		return nil, nil, decimal.Zero, decimal.Zero, ErrInvalidQuantity
	}

	var (
		iter          func() *OrderQueue
		sideToProcess *OrderSide
	)

	if side == Buy {
		iter = ob.asks.MinPriceQueue
		sideToProcess = ob.asks
	} else {
		iter = ob.bids.MaxPriceQueue
		sideToProcess = ob.bids
	}

	for quantity.Sign() > 0 && sideToProcess.Len() > 0 {
		bestPrice := iter()
		ordersDone, partialDone, partialProcessed, quantityLeft := ob.processQueue(bestPrice, quantity)
		done = append(done, ordersDone...)
		partial = partialDone
		partialQuantityProcessed = partialProcessed
		quantity = quantityLeft
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
//      done    - not nil if your order produces ends of anoter order, this order will add to
//                the "done" slice. If your order have done too, it will be places to this array too
//      partial - not nil if your order has done but top order is not fully done. Or if your order is
//                partial done and placed to the orderbook without full quantity - partial will contain
//                your order with quantity to left
//      partialQuantityProcessed - if partial order is not nil this result contains processed quatity from partial order
func (ob *OrderBook) ProcessLimitOrder(side Side, orderID string, quantity, price decimal.Decimal) (done []*Order, partial *Order, partialQuantityProcessed decimal.Decimal, err error) {
	if _, ok := ob.orders[orderID]; ok {
		return nil, nil, decimal.Zero, ErrOrderExists
	}

	if quantity.Sign() <= 0 {
		return nil, nil, decimal.Zero, ErrInvalidQuantity
	}

	if price.Sign() <= 0 {
		return nil, nil, decimal.Zero, ErrInvalidPrice
	}

	quantityToTrade := quantity
	var (
		sideToProcess *OrderSide
		sideToAdd     *OrderSide
		comparator    func(decimal.Decimal) bool
		iter          func() *OrderQueue
	)

	if side == Buy {
		sideToAdd = ob.bids
		sideToProcess = ob.asks
		comparator = price.GreaterThanOrEqual
		iter = ob.asks.MinPriceQueue
	} else {
		sideToAdd = ob.asks
		sideToProcess = ob.bids
		comparator = price.LessThanOrEqual
		iter = ob.bids.MaxPriceQueue
	}

	bestPrice := iter()
	for quantityToTrade.Sign() > 0 && sideToProcess.Len() > 0 && comparator(bestPrice.Price()) {
		ordersDone, partialDone, partialQty, quantityLeft := ob.processQueue(bestPrice, quantityToTrade)
		done = append(done, ordersDone...)
		partial = partialDone
		partialQuantityProcessed = partialQty
		quantityToTrade = quantityLeft
		bestPrice = iter()
	}

	if quantityToTrade.Sign() > 0 {
		o := NewOrder(orderID, side, quantityToTrade, price, time.Now().UTC())
		if len(done) > 0 {
			partialQuantityProcessed = quantity.Sub(quantityToTrade)
			partial = o
		}
		ob.orders[orderID] = sideToAdd.Append(o)
	} else {
		totalQuantity := decimal.Zero
		totalPrice := decimal.Zero

		for _, order := range done {
			totalQuantity = totalQuantity.Add(order.Quantity())
			totalPrice = totalPrice.Add(order.Price().Mul(order.Quantity()))
		}

		if partialQuantityProcessed.Sign() > 0 {
			totalQuantity = totalQuantity.Add(partialQuantityProcessed)
			totalPrice = totalPrice.Add(partial.Price().Mul(partialQuantityProcessed))
		}

		done = append(done, NewOrder(orderID, side, quantity, totalPrice.Div(totalQuantity), time.Now().UTC()))
	}
	return
}

func (ob *OrderBook) processQueue(orderQueue *OrderQueue, quantityToTrade decimal.Decimal) (done []*Order, partial *Order, partialQuantityProcessed, quantityLeft decimal.Decimal) {
	quantityLeft = quantityToTrade

	for orderQueue.Len() > 0 && quantityLeft.Sign() > 0 {
		headOrderEl := orderQueue.Head()
		headOrder := headOrderEl.Value.(*Order)

		if quantityLeft.LessThan(headOrder.Quantity()) {
			partial = NewOrder(headOrder.ID(), headOrder.Side(), headOrder.Quantity().Sub(quantityLeft), headOrder.Price(), headOrder.Time())
			partialQuantityProcessed = quantityLeft
			orderQueue.Update(headOrderEl, partial)
			quantityLeft = decimal.Zero
		} else {
			quantityLeft = quantityLeft.Sub(headOrder.Quantity())
			done = append(done, ob.CancelOrder(headOrder.ID()))
		}
	}

	return
}

// Order returns order by id
func (ob *OrderBook) Order(orderID string) *Order {
	e, ok := ob.orders[orderID]
	if !ok {
		return nil
	}

	return e.Value.(*Order)
}

// Depth returns price levels and volume at price level
func (ob *OrderBook) Depth() (asks, bids []*PriceLevel) {
	level := ob.asks.MaxPriceQueue()
	for level != nil {
		asks = append(asks, &PriceLevel{
			Price:    level.Price(),
			Quantity: level.Volume(),
		})
		level = ob.asks.LessThan(level.Price())
	}

	level = ob.bids.MaxPriceQueue()
	for level != nil {
		bids = append(bids, &PriceLevel{
			Price:    level.Price(),
			Quantity: level.Volume(),
		})
		level = ob.bids.LessThan(level.Price())
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

func (ob *OrderBook) GetOrders() *list {
	e := ob.orders
	return e
}

// CalculateMarketPrice returns total market price for requested quantity
// if err is not nil price returns total price of all levels in side
func (ob *OrderBook) CalculateMarketPrice(side Side, quantity decimal.Decimal) (price decimal.Decimal, err error) {
	price = decimal.Zero

	var (
		level *OrderQueue
		iter  func(decimal.Decimal) *OrderQueue
	)

	if side == Buy {
		level = ob.asks.MinPriceQueue()
		iter = ob.asks.GreaterThan
	} else {
		level = ob.bids.MaxPriceQueue()
		iter = ob.bids.LessThan
	}

	for quantity.Sign() > 0 && level != nil {
		levelVolume := level.Volume()
		levelPrice := level.Price()
		if quantity.GreaterThanOrEqual(levelVolume) {
			price = price.Add(levelPrice.Mul(levelVolume))
			quantity = quantity.Sub(levelVolume)
			level = iter(levelPrice)
		} else {
			price = price.Add(levelPrice.Mul(quantity))
			quantity = decimal.Zero
		}
	}

	if quantity.Sign() > 0 {
		err = ErrInsufficientQuantity
	}

	return
}

// String implements fmt.Stringer interface
func (ob *OrderBook) String() string {
	return ob.asks.String() + "\r\n------------------------------------" + ob.bids.String()
}

// MarshalJSON implements json.Marshaler interface
func (ob *OrderBook) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			Asks *OrderSide `json:"asks"`
			Bids *OrderSide `json:"bids"`
		}{
			Asks: ob.asks,
			Bids: ob.bids,
		},
	)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (ob *OrderBook) UnmarshalJSON(data []byte) error {
	obj := struct {
		Asks *OrderSide `json:"asks"`
		Bids *OrderSide `json:"bids"`
	}{}

	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	ob.asks = obj.Asks
	ob.bids = obj.Bids
	ob.orders = map[string]*list.Element{}

	for _, order := range ob.asks.Orders() {
		ob.orders[order.Value.(*Order).ID()] = order
	}

	for _, order := range ob.bids.Orders() {
		ob.orders[order.Value.(*Order).ID()] = order
	}

	return nil
}
