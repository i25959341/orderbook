package orderbook

import (
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
	lane "gopkg.in/oleiade/lane.v1"
)

type OrderBook struct {
	deque         *lane.Deque
	bids          OrderTree
	asks          OrderTree
	time          int
	next_order_id int
}

func (orderBook *OrderBook) UpdateTime() {
	orderBook.time++
}

func (orderBook *OrderBook) BestBid() (value interface{}, found bool) {
	value, found = orderBook.bids.MaxPrice()
	return value, found
}

func (orderBook *OrderBook) BestAsk() (value interface{}, found bool) {
	value, found = orderBook.asks.MinPrice()
	return value, found
}

func (orderBook *OrderBook) WorstBid() (value interface{}, found bool) {
	value, found = orderBook.bids.MinPrice()
	return value, found
}

func (orderBook *OrderBook) WorstAsk() (value interface{}, found bool) {
	value, found = orderBook.asks.MaxPrice()
	return value, found
}

func (orderBook *OrderBook) ProcessOrderList(side string, orderList OrderList, quantityStillToTrade decimal.Decimal, quote map[string]string, verbose bool) (decimal.Decimal, []map[string]string) {
	quantityToTrade := quantityStillToTrade
	var trades []map[string]string

	for orderList.Length() > 0 && quantityToTrade.GreaterThan(decimal.Zero) {
		headOrder := orderList.HeadOrder()
		tradedPrice := headOrder.price
		counterParty := headOrder.trade_id
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
				orderBook.bids.RemoveOrderById(headOrder.order_id)
			} else {
				orderBook.asks.RemoveOrderById(headOrder.order_id)
			}
			quantityToTrade = decimal.Zero

		} else {
			tradedQuantity = headOrder.quantity
			if side == "bid" {
				orderBook.bids.RemoveOrderById(headOrder.order_id)
			} else {
				orderBook.asks.RemoveOrderById(headOrder.order_id)
			}
		}

		if verbose {
			fmt.Println("TRADE: Time - %v, Price - %v, Quantity - %v, TradeID - %v, Matching TradeID - %v", orderBook.time, tradedPrice.String(), tradedQuantity.String(), counterParty, quote["trade_id"])
		}

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

// def process_market_order(self, quote, verbose):
// trades = []
// quantity_to_trade = quote['quantity']
// side = quote['side']
// if side == 'bid':
// 	while quantity_to_trade > 0 and self.asks:
// 		best_price_asks = self.asks.min_price_list()
// 		quantity_to_trade, new_trades = self.process_order_list('ask', best_price_asks, quantity_to_trade, quote, verbose)
// 		trades += new_trades
// elif side == 'ask':
// 	while quantity_to_trade > 0 and self.bids:
// 		best_price_bids = self.bids.max_price_list()
// 		quantity_to_trade, new_trades = self.process_order_list('bid', best_price_bids, quantity_to_trade, quote, verbose)
// 		trades += new_trades
// else:
// 	sys.exit('process_market_order() recieved neither "bid" nor "ask"')
// return trades

// def process_limit_order(self, quote, from_data, verbose):
// order_in_book = None
// trades = []
// quantity_to_trade = quote['quantity']
// side = quote['side']
// price = quote['price']
// if side == 'bid':
// 	while (self.asks and price >= self.asks.min_price() and quantity_to_trade > 0):
// 		best_price_asks = self.asks.min_price_list()
// 		quantity_to_trade, new_trades = self.process_order_list('ask', best_price_asks, quantity_to_trade, quote, verbose)
// 		trades += new_trades
// 	# If volume remains, need to update the book with new quantity
// 	if quantity_to_trade > 0:
// 		if not from_data:
// 			quote['order_id'] = self.next_order_id
// 		quote['quantity'] = quantity_to_trade
// 		self.bids.insert_order(quote)
// 		order_in_book = quote
// elif side == 'ask':
// 	while (self.bids and price <= self.bids.max_price() and quantity_to_trade > 0):
// 		best_price_bids = self.bids.max_price_list()
// 		quantity_to_trade, new_trades = self.process_order_list('bid', best_price_bids, quantity_to_trade, quote, verbose)
// 		trades += new_trades
// 	# If volume remains, need to update the book with new quantity
// 	if quantity_to_trade > 0:
// 		if not from_data:
// 			quote['order_id'] = self.next_order_id
// 		quote['quantity'] = quantity_to_trade
// 		self.asks.insert_order(quote)
// 		order_in_book = quote
// else:
// 	sys.exit('process_limit_order() given neither "bid" nor "ask"')
// return trades, order_in_book

// def process_order(self, quote, from_data, verbose):
//         order_type = quote['type']
//         order_in_book = None
//         if from_data:
//             self.time = quote['timestamp']
//         else:
//             self.update_time()
//             quote['timestamp'] = self.time
//         if quote['quantity'] <= 0:
//             sys.exit('process_order() given order of quantity <= 0')
//         if not from_data:
//             self.next_order_id += 1
//         if order_type == 'market':
//             trades = self.process_market_order(quote, verbose)
//         elif order_type == 'limit':
//             quote['price'] = Decimal(quote['price'])
//             trades, order_in_book = self.process_limit_order(quote, from_data, verbose)
//         else:
//             sys.exit("order_type for process_order() is neither 'market' or 'limit'")
//         return trades, order_in_book
