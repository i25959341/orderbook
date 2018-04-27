package orderbook

import (
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

func (orderBook *OrderBook) UpdateTime() {
	orderBook.time++
}

func (orderBook *OrderBook) BestBid() (value decimal.Decimal) {
	value = orderBook.bids.MaxPrice()
	return
}

func (orderBook *OrderBook) BestAsk() (value decimal.Decimal) {
	value = orderBook.asks.MinPrice()
	return
}

func (orderBook *OrderBook) WorstBid() (value decimal.Decimal) {
	value = orderBook.bids.MinPrice()
	return
}

func (orderBook *OrderBook) WorstAsk() (value decimal.Decimal) {
	value = orderBook.asks.MaxPrice()
	return
}

func (orderBook *OrderBook) ProcessMarketOrder(quote map[string]string, verbose bool) []map[string]string {
	var trades []map[string]string
	quantity_to_trade, _ := decimal.NewFromString(quote["quantity"])
	side := quote["side"]
	var new_trades []map[string]string

	if side == "bid" {
		for quantity_to_trade.GreaterThan(decimal.Zero) && orderBook.asks.Length() > 0 {
			best_price_asks := orderBook.asks.MinPriceList()
			quantity_to_trade, new_trades = orderBook.ProcessOrderList("ask", best_price_asks, quantity_to_trade, quote, verbose)
			trades = append(trades, new_trades...)
		}
	} else if side == "ask" {
		for quantity_to_trade.GreaterThan(decimal.Zero) && orderBook.bids.Length() > 0 {
			best_price_bids := orderBook.bids.MaxPriceList()
			quantity_to_trade, new_trades = orderBook.ProcessOrderList("bid", best_price_bids, quantity_to_trade, quote, verbose)
			trades = append(trades, new_trades...)
		}
	}
	return trades
}

func (orderBook *OrderBook) ProcessLimitOrder(quote map[string]string, verbose bool) ([]map[string]string, map[string]string) {
	var trades []map[string]string
	quantity_to_trade, _ := decimal.NewFromString(quote["quantity"])
	side := quote["side"]
	price, _ := decimal.NewFromString(quote["price"])
	var new_trades []map[string]string

	var order_in_book map[string]string

	if side == "bid" {
		minPrice := orderBook.asks.MinPrice()
		for quantity_to_trade.GreaterThan(decimal.Zero) && orderBook.asks.Length() > 0 && price.GreaterThanOrEqual(minPrice) {
			best_price_asks := orderBook.asks.MinPriceList()
			quantity_to_trade, new_trades = orderBook.ProcessOrderList("ask", best_price_asks, quantity_to_trade, quote, verbose)
			trades = append(trades, new_trades...)
			minPrice = orderBook.asks.MinPrice()
		}

		if quantity_to_trade.GreaterThan(decimal.Zero) {
			quote["order_id"] = strconv.Itoa(orderBook.nextOrderID)
			quote["quantity"] = quantity_to_trade.String()
			orderBook.bids.InsertOrder(quote)
			order_in_book = quote
		}

	} else if side == "ask" {
		maxPrice := orderBook.bids.MaxPrice()
		for quantity_to_trade.GreaterThan(decimal.Zero) && orderBook.bids.Length() > 0 && price.LessThanOrEqual(maxPrice) {
			best_price_bids := orderBook.bids.MaxPriceList()
			quantity_to_trade, new_trades = orderBook.ProcessOrderList("bid", best_price_bids, quantity_to_trade, quote, verbose)
			trades = append(trades, new_trades...)
			maxPrice = orderBook.bids.MaxPrice()
		}

		if quantity_to_trade.GreaterThan(decimal.Zero) {
			quote["order_id"] = strconv.Itoa(orderBook.nextOrderID)
			quote["quantity"] = quantity_to_trade.String()
			orderBook.asks.InsertOrder(quote)
			order_in_book = quote
		}
	}
	return trades, order_in_book
}

func (orderBook *OrderBook) ProcessOrder(quote map[string]string, verbose bool) ([]map[string]string, map[string]string) {
	order_type := quote["type"]
	var order_in_book map[string]string
	var trades []map[string]string

	orderBook.UpdateTime()
	quote["timestamp"] = strconv.Itoa(orderBook.time)
	orderBook.nextOrderID++

	if order_type == "market" {
		trades = orderBook.ProcessMarketOrder(quote, verbose)
	} else {
		trades, order_in_book = orderBook.ProcessLimitOrder(quote, verbose)
	}
	return trades, order_in_book
}

func (orderBook *OrderBook) ProcessOrderList(side string, orderList *OrderList, quantityStillToTrade decimal.Decimal, quote map[string]string, verbose bool) (decimal.Decimal, []map[string]string) {
	quantityToTrade := quantityStillToTrade
	var trades []map[string]string

	for orderList.Length() > 0 && quantityToTrade.GreaterThan(decimal.Zero) {
		headOrder := orderList.HeadOrder()
		tradedPrice := headOrder.price
		// counterParty := headOrder.trade_id
		var newBookQuantity decimal.Decimal
		var tradedQuantity decimal.Decimal

		if quantityToTrade.LessThan(headOrder.quantity) {
			tradedQuantity = quantityToTrade
			// Do the transaction
			newBookQuantity = headOrder.quantity.Sub(quantityToTrade)
			headOrder.UpdateQuantity(newBookQuantity, headOrder.timestamp)
			quantityToTrade = decimal.Zero

		} else if quantityToTrade.Equal(headOrder.quantity) {
			tradedQuantity = quantityToTrade
			if side == "bid" {
				orderBook.bids.RemoveOrderById(headOrder.orderID)
			} else {
				orderBook.asks.RemoveOrderById(headOrder.orderID)
			}
			quantityToTrade = decimal.Zero

		} else {
			tradedQuantity = headOrder.quantity
			if side == "bid" {
				orderBook.bids.RemoveOrderById(headOrder.orderID)
			} else {
				orderBook.asks.RemoveOrderById(headOrder.orderID)
			}
		}

		// if verbose {
		// 	fmt.Println("TRADE: Time - %v, Price - %v, Quantity - %v, TradeID - %v, Matching TradeID - %v", orderBook.time, tradedPrice.String(), tradedQuantity.String(), counterParty, quote["trade_id"])
		// }

		transactionRecord := make(map[string]string)
		transactionRecord["timestamp"] = strconv.Itoa(orderBook.time)
		transactionRecord["price"] = tradedPrice.String()
		transactionRecord["quantity"] = tradedQuantity.String()
		transactionRecord["time"] = strconv.Itoa(orderBook.time)

		orderBook.deque.Append(transactionRecord)
		trades = append(trades, transactionRecord)
	}
	return quantityToTrade, trades
}

func (orderBook *OrderBook) CancelOrder(side string, order_id int) {
	orderBook.UpdateTime()

	if side == "bid" {
		if orderBook.bids.OrderExist(strconv.Itoa(order_id)) {
			orderBook.bids.RemoveOrderById(strconv.Itoa(order_id))
		}
	} else {
		if orderBook.asks.OrderExist(strconv.Itoa(order_id)) {
			orderBook.asks.RemoveOrderById(strconv.Itoa(order_id))
		}
	}
}

func (orderBook *OrderBook) ModifyOrder(quoteUpdate map[string]string, order_id int) {
	orderBook.UpdateTime()

	side := quoteUpdate["side"]
	quoteUpdate["order_id"] = strconv.Itoa(order_id)
	quoteUpdate["timestamp"] = strconv.Itoa(orderBook.time)

	if side == "bid" {
		if orderBook.bids.OrderExist(strconv.Itoa(order_id)) {
			orderBook.bids.UpdateOrder(quoteUpdate)
		}
	} else {
		if orderBook.asks.OrderExist(strconv.Itoa(order_id)) {
			orderBook.asks.UpdateOrder(quoteUpdate)
		}
	}
}

func (orderBook *OrderBook) VolumeAtPrice(side string, price decimal.Decimal) decimal.Decimal {
	if side == "bid" {
		volume := decimal.Zero
		if orderBook.bids.PriceExist(price) {
			volume = orderBook.bids.PriceList(price).volume
		}

		return volume

	} else {
		volume := decimal.Zero
		if orderBook.asks.PriceExist(price) {
			volume = orderBook.asks.PriceList(price).volume
		}
		return volume
	}
}
