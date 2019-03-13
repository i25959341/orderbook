# Go orderbook

Improved matching engine written in Go (Golang)

[![Go Report Card](https://goreportcard.com/badge/github.com/i25959341/orderbook)](https://goreportcard.com/report/github.com/i25959341/orderbook)
[![GoDoc](https://godoc.org/github.com/i25959341/orderbook?status.svg)](https://godoc.org/github.com/i25959341/orderbook)
[![gocover.run](https://gocover.run/github.com/i25959341/orderbook.svg?style=flat&tag=1.10)](https://gocover.run?tag=1.10&repo=github.com%2Fi25959341%2Forderbook)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)
[![Build Status](https://travis-ci.org/i25959341/orderbook.svg?branch=master)](https://travis-ci.org/i25959341/orderbook)

## Features

- Standard price-time priority
- Supports both market and limit orders
- Supports order cancelling
- High performance (above 300k trades per second)
- Optimal memory usage
- JSON Marshalling and Unmarsalling
- Calculating market price for definite quantity

## Usage

To start using order book you need to create object:

```go
import (
  "fmt" 
  ob "github.com/muzykantov/orderbook"
)

func main() {
  orderBook := ob.NewOrderBook()
  fmt.Println(orderBook)
}

```

Then you be able to use next primary functions:

```go

func (ob *OrderBook) ProcessLimitOrder(side Side, orderID string, quantity, price decimal.Decimal) (done []*Order, partial *Order, err error) { ... }

func (ob *OrderBook) ProcessMarketOrder(side Side, quantity decimal.Decimal) (done []*Order, partial *Order, quantityLeft decimal.Decimal, err error) { .. }

func (ob *OrderBook) CancelOrder(orderID string) *Order { ... }

```

## About primary functions

### ProcessLimitOrder

```go
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
func (ob *OrderBook) ProcessLimitOrder(side Side, orderID string, quantity, price decimal.Decimal) (done []*Order, partial *Order, err error) { ... }
```

For example:
```
ProcessLimitOrder(ob.Sell, "uinqueID", decimal.New(55, 0), decimal.New(100, 0))

asks: 110 -> 5      110 -> 5
      100 -> 1      100 -> 56
--------------  ->  --------------
bids: 90  -> 5      90  -> 5
      80  -> 1      80  -> 1

done    - nil
partial - nil

```

```
ProcessLimitOrder(ob.Buy, "uinqueID", decimal.New(7, 0), decimal.New(120, 0))

asks: 110 -> 5
      100 -> 1
--------------  ->  --------------
bids: 90  -> 5      120 -> 1
      80  -> 1      90  -> 5
                    80  -> 1

done    - 2 (or more orders)
partial - uinqueID order

```

```
ProcessLimitOrder(ob.Buy, "uinqueID", decimal.New(3, 0), decimal.New(120, 0))

asks: 110 -> 5
      100 -> 1      110 -> 3
--------------  ->  --------------
bids: 90  -> 5      90  -> 5
      80  -> 1      90  -> 5

done    - 1 order with 100 price, (may be also few orders with 110 price) + uinqueID order
partial - 1 order with price 110

```

### ProcessMarketOrder

```go
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
func (ob *OrderBook) ProcessMarketOrder(side Side, quantity decimal.Decimal) (done []*Order, partial *Order, quantityLeft decimal.Decimal, err error) { .. }
```

For example:
```
ProcessMarketOrder(ob.Sell, decimal.New(6, 0))

asks: 110 -> 5      110 -> 5
      100 -> 1      100 -> 1
--------------  ->  --------------
bids: 90  -> 5      80 -> 1
      80  -> 2

done         - 2 (or more orders)
partial      - 1 order with price 80
quantityLeft - 0

```

```
ProcessMarketOrder(ob.Buy, decimal.New(10, 0))

asks: 110 -> 5
      100 -> 1
--------------  ->  --------------
bids: 90  -> 5      90  -> 5
      80  -> 1      80  -> 1
                    
done         - 2 (or more orders)
partial      - nil
quantityLeft - 4

```

### CancelOrder

```go
// CancelOrder removes order with given ID from the order book
func (ob *OrderBook) CancelOrder(orderID string) *Order { ... }
```

```
CancelOrder("myUinqueID-Sell-1-with-100")

asks: 110 -> 5
      100 -> 1      110 -> 5
--------------  ->  --------------
bids: 90  -> 5      90  -> 5
      80  -> 1      80  -> 1
                    
done         - 2 (or more orders)
partial      - nil
quantityLeft - 4

```

## License

The MIT License (MIT)

See LICENSE and AUTHORS files
