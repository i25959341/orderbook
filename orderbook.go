package orderbook

import (
	"github.com/shopspring/decimal"
)

type OrderBook struct {
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

func (orderBook *OrderBook) ProcessMarketOrder(quote map[string]string, verbose bool) {

	quantity_to_trade, _ := decimal.NewFromString(quote["quantity"])
	side := quote["side"]

	if side == "bid" {
		for (quantity_to_trade.GreaterThan(decimal.Zero)) && (orderBook.asks.Length() > 0) {

		}

	} else if side == "ask" {

	} else {

	}
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
