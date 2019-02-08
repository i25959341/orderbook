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
	return &OrderBook{
		deque: lane.NewDeque(),
		bids:  NewOrderTree(),
		asks:  NewOrderTree(),
	}
}

func (ob *OrderBook) UpdateTime() {
	ob.time++
}

func (ob *OrderBook) BestBid() decimal.Decimal {
	return ob.bids.MaxPrice()
}

func (ob *OrderBook) BestAsk() decimal.Decimal {
	return ob.asks.MinPrice()
}

func (ob *OrderBook) WorstBid() decimal.Decimal {
	return ob.bids.MinPrice()
}

func (ob *OrderBook) WorstAsk() decimal.Decimal {
	return ob.asks.MaxPrice()
}

func (ob *OrderBook) ProcessMarketOrderFromMap(quote map[string]string, verbose bool) []map[string]string {
	var trades []map[string]string
	quantityToTrade, _ := decimal.NewFromString(quote["quantity"])
	side := quote["side"]
	var newTrades []map[string]string

	if side == "bid" {
		for quantityToTrade.GreaterThan(decimal.Zero) && ob.asks.Length() > 0 {
			bestPriceAsks := ob.asks.MinPriceQueue()
			quantityToTrade, newTrades = ob.ProcessOrderListFromMap("ask", bestPriceAsks, quantityToTrade, quote, verbose)
			trades = append(trades, newTrades...)
		}
	} else if side == "ask" {
		for quantityToTrade.GreaterThan(decimal.Zero) && ob.bids.Length() > 0 {
			bestPriceBids := ob.bids.MaxPriceQueue()
			quantityToTrade, newTrades = ob.ProcessOrderListFromMap("bid", bestPriceBids, quantityToTrade, quote, verbose)
			trades = append(trades, newTrades...)
		}
	}
	return trades
}

func (ob *OrderBook) ProcessLimitOrderFromMap(quote map[string]string, verbose bool) ([]map[string]string, map[string]string) {
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
			quantityToTrade, newTrades = ob.ProcessOrderListFromMap("ask", bestPriceAsks, quantityToTrade, quote, verbose)
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
			quantityToTrade, newTrades = ob.ProcessOrderListFromMap("bid", bestPriceBids, quantityToTrade, quote, verbose)
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

func (ob *OrderBook) ProcessOrderFromMap(quote map[string]string, verbose bool) ([]map[string]string, map[string]string) {
	orderType := quote["type"]
	var orderInBook map[string]string
	var trades []map[string]string

	ob.UpdateTime()
	quote["timestamp"] = strconv.Itoa(ob.time)
	ob.nextOrderID++

	if orderType == "market" {
		trades = ob.ProcessMarketOrderFromMap(quote, verbose)
	} else {
		trades, orderInBook = ob.ProcessLimitOrderFromMap(quote, verbose)
	}
	return trades, orderInBook
}

func (ob *OrderBook) ProcessOrderListFromMap(side string, orderList *OrderQueue, quantityStillToTrade decimal.Decimal, quote map[string]string, verbose bool) (decimal.Decimal, []map[string]string) {
	quantityToTrade := quantityStillToTrade
	var trades []map[string]string

	for orderList.Length() > 0 && quantityToTrade.GreaterThan(decimal.Zero) {
		headOrder := orderList.Head()
		tradedPrice := headOrder.price
		var (
			newBookQuantity decimal.Decimal
			tradedQuantity  decimal.Decimal
		)

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

func (ob *OrderBook) CancelOrder(side string, orderID string) {
	ob.UpdateTime()

	if side == "bid" {
		if ob.bids.OrderExist(orderID) {
			ob.bids.RemoveOrder(orderID)
		}
	} else {
		if ob.asks.OrderExist(orderID) {
			ob.asks.RemoveOrder(orderID)
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
