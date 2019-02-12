# Go orderbook

Improved matching engine written in Go

## Features

- Standard price-time priority
- Supports both market and limit orders
- Supports order cancelling
- High performance (above 100k trades per second)
- Optimal memory usage

## Usage

```go
package main

import (
    "fmt"
    ob "github.com/muzykantov/orderbook"
    "github.com/shopspring/decimal"
)

func main() {
    orderBook := ob.NewOrderBook(actions)

    for i := 50; i < 100; i = i + 10 {
		orderBook.ProcessLimitOrder(ob.Buy, fmt.Sprintf("b-%d", i), decimal.New(2, 0), decimal.New(int64(i), 0))
	}

	for i := 100; i < 150; i = i + 10 {
		orderBook.ProcessLimitOrder(ob.Sell, fmt.Sprintf("s-%d", i), decimal.New(2, 0), decimal.New(int64(i), 0))
    }
    fmt.Println(orderBook)

    ordersDone, partialDone, _ := orderBook.ProcessMarketOrder(Buy, decimal.New(3, 0))
    fmt.Println(orderBook)

    fmt.Println("Done:", ordersDone)
    fmt.Println("Partial:", partialDone)
}
```