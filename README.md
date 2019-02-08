# OrderBook

Matching engine based on a limit order book written in Go.

Features:

- Standard price-time priority
- Supports both market and limit orders
- Add, cancel, update orders

## Data Structure

Orders are sent to the order book using the ProcessOrder function. The Order is created using a quote.

```Go
// For a limit order
quote := map[string]string {
        "type"     : "limit",
        "side"     : "bid", 
        "quantity" : "6", 
        "price"    : "108.2", 
        "trade_id" : "001",
    }
         
// For a market order
quote := map[string]string {
        "type"     : "market",
        "side"     : "ask", 
        "quantity" : "6", 
        "trade_id" : "002",
    }
```