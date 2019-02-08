package orderbook

import (
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
	lane "gopkg.in/oleiade/lane.v1"
)

type OrderBook struct {
	deque       *lane.Deque
	bids        *OrderTree
	asks        *OrderTree
	time        int
	nextOrderID int
}

func NewOrderBook() *OrderBook {
	deque := lane.NewDeque()
	bids := NewOrderTree()
	asks := NewOrderTree()
	return &OrderBook{deque, bids, asks, 0, 0}
}

func (ob *OrderBook) UpdateTime() {
	ob.time++
}

func (ob *OrderBook) BestBid() (value decimal.Decimal) {
	value = ob.bids.MaxPrice()
	return
}

func (ob *OrderBook) BestAsk() (value decimal.Decimal) {
	value = ob.asks.MinPrice()
	return
}

func (ob *OrderBook) WorstBid() (value decimal.Decimal) {
	value = ob.bids.MinPrice()
	return
}

func (ob *OrderBook) WorstAsk() (value decimal.Decimal) {
	value = ob.asks.MaxPrice()
	return
}

func (ob *OrderBook) ProcessMarketOrder(quote map[string]string, verbose bool) []map[string]string {
	var trades []map[string]string
	quantityToTrade, _ := decimal.NewFromString(quote["quantity"])
	side := quote["side"]
	var newTrades []map[string]string

	if side == "bid" {
		for quantityToTrade.GreaterThan(decimal.Zero) && ob.asks.Length() > 0 {
			bestPriceAsks := ob.asks.MinPriceQueue()
			quantityToTrade, newTrades = ob.ProcessOrderList("ask", bestPriceAsks, quantityToTrade, quote, verbose)
			trades = append(trades, newTrades...)
		}
	} else if side == "ask" {
		for quantityToTrade.GreaterThan(decimal.Zero) && ob.bids.Length() > 0 {
			bestPriceBids := ob.bids.MaxPriceQueue()
			quantityToTrade, newTrades = ob.ProcessOrderList("bid", bestPriceBids, quantityToTrade, quote, verbose)
			trades = append(trades, newTrades...)
		}
	}
	return trades
}

func (ob *OrderBook) ProcessLimitOrder(quote map[string]string, verbose bool) ([]map[string]string, map[string]string) {
	var trades []map[string]string
	quantityToTrade, _ := decimal.NewFromString(quote["quantity"])
	side := quote["side"]
	price, _ := decimal.NewFromString(quote["price"])
	var newTrades []map[string]string

	var orderInBook map[string]string

	if side == "bid" {
		minPrice := ob.asks.MinPrice()
		for quantityToTrade.GreaterThan(decimal.Zero) && ob.asks.Length() > 0 && price.GreaterThanOrEqual(minPrice) {
			bestPriceAsks := ob.asks.MinPriceQueue()
			quantityToTrade, newTrades = ob.ProcessOrderList("ask", bestPriceAsks, quantityToTrade, quote, verbose)
			trades = append(trades, newTrades...)
			minPrice = ob.asks.MinPrice()
		}

		if quantityToTrade.GreaterThan(decimal.Zero) {
			quote["order_id"] = strconv.Itoa(ob.nextOrderID)
			quote["quantity"] = quantityToTrade.String()
			ob.bids.InsertOrderFromMap(quote)
			orderInBook = quote
		}

	} else if side == "ask" {
		maxPrice := ob.bids.MaxPrice()
		for quantityToTrade.GreaterThan(decimal.Zero) && ob.bids.Length() > 0 && price.LessThanOrEqual(maxPrice) {
			bestPriceBids := ob.bids.MaxPriceQueue()
			quantityToTrade, newTrades = ob.ProcessOrderList("bid", bestPriceBids, quantityToTrade, quote, verbose)
			trades = append(trades, newTrades...)
			maxPrice = ob.bids.MaxPrice()
		}

		if quantityToTrade.GreaterThan(decimal.Zero) {
			quote["order_id"] = strconv.Itoa(ob.nextOrderID)
			quote["quantity"] = quantityToTrade.String()
			ob.asks.InsertOrderFromMap(quote)
			orderInBook = quote
		}
	}
	return trades, orderInBook
}

func (ob *OrderBook) ProcessOrder(quote map[string]string, verbose bool) ([]map[string]string, map[string]string) {
	orderType := quote["type"]
	var orderInBook map[string]string
	var trades []map[string]string

	ob.UpdateTime()
	quote["timestamp"] = strconv.Itoa(ob.time)
	ob.nextOrderID++

	if orderType == "market" {
		trades = ob.ProcessMarketOrder(quote, verbose)
	} else {
		trades, orderInBook = ob.ProcessLimitOrder(quote, verbose)
	}
	return trades, orderInBook
}

func (ob *OrderBook) ProcessOrderList(side string, orderList *OrderQueue, quantityStillToTrade decimal.Decimal, quote map[string]string, verbose bool) (decimal.Decimal, []map[string]string) {
	quantityToTrade := quantityStillToTrade
	var trades []map[string]string

	for orderList.Length() > 0 && quantityToTrade.GreaterThan(decimal.Zero) {
		headOrder := orderList.Head()
		tradedPrice := headOrder.price
		// counterParty := headOrder.trade_id
		var newBookQuantity decimal.Decimal
		var tradedQuantity decimal.Decimal

		if quantityToTrade.LessThan(headOrder.quantity) {
			tradedQuantity = quantityToTrade
			// Do the transaction
			newBookQuantity = headOrder.quantity.Sub(quantityToTrade)
			headOrder.Update(newBookQuantity, headOrder.timestamp)
			quantityToTrade = decimal.Zero

		} else if quantityToTrade.Equal(headOrder.quantity) {
			tradedQuantity = quantityToTrade
			if side == "bid" {
				ob.bids.RemoveOrder(headOrder.orderID)
			} else {
				ob.asks.RemoveOrder(headOrder.orderID)
			}
			quantityToTrade = decimal.Zero

		} else {
			tradedQuantity = headOrder.quantity
			if side == "bid" {
				ob.bids.RemoveOrder(headOrder.orderID)
			} else {
				ob.asks.RemoveOrder(headOrder.orderID)
			}
		}

		if verbose {
			fmt.Printf("TRADE: Time - %v, Price - %v, Quantity - %v, Matching TradeID - %v\n", ob.time, tradedPrice.String(), tradedQuantity.String(), quote["trade_id"])
		}

		transactionRecord := make(map[string]string)
		transactionRecord["timestamp"] = strconv.Itoa(ob.time)
		transactionRecord["price"] = tradedPrice.String()
		transactionRecord["quantity"] = tradedQuantity.String()
		transactionRecord["time"] = strconv.Itoa(ob.time)

		ob.deque.Append(transactionRecord)
		trades = append(trades, transactionRecord)
	}
	return quantityToTrade, trades
}

func (ob *OrderBook) CancelOrder(side string, order_id int) {
	ob.UpdateTime()

	if side == "bid" {
		if ob.bids.OrderExist(strconv.Itoa(order_id)) {
			ob.bids.RemoveOrder(strconv.Itoa(order_id))
		}
	} else {
		if ob.asks.OrderExist(strconv.Itoa(order_id)) {
			ob.asks.RemoveOrder(strconv.Itoa(order_id))
		}
	}
}

func (ob *OrderBook) ModifyOrder(quoteUpdate map[string]string, order_id int) {
	ob.UpdateTime()

	side := quoteUpdate["side"]
	quoteUpdate["order_id"] = strconv.Itoa(order_id)
	quoteUpdate["timestamp"] = strconv.Itoa(ob.time)

	if side == "bid" {
		if ob.bids.OrderExist(strconv.Itoa(order_id)) {
			ob.bids.UpdateOrderFromMap(quoteUpdate)
		}
	} else {
		if ob.asks.OrderExist(strconv.Itoa(order_id)) {
			ob.asks.UpdateOrderFromMap(quoteUpdate)
		}
	}
}

func (ob *OrderBook) VolumeAtPrice(side string, price decimal.Decimal) decimal.Decimal {
	if side == "bid" {
		volume := decimal.Zero
		if ob.bids.PriceExist(price) {
			volume = ob.bids.PriceQueue(price).volume
		}

		return volume

	} else {
		volume := decimal.Zero
		if ob.asks.PriceExist(price) {
			volume = ob.asks.PriceQueue(price).volume
		}
		return volume
	}
}

func (ob *OrderBook) String() string {
	panic("not implemented")
}
