package orderbook

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func BenchmarkOrderQueue(b *testing.B) {
	price := decimal.New(100, 0)
	orderQueue := NewOrderQueue(price)
	stopwatch := time.Now()

	var o *Order
	for i := 0; i < b.N; i++ {
		o = NewOrder(
			fmt.Sprintf("order-%d", i),
			Buy,
			decimal.New(100, 0),
			decimal.New(int64(i), 0),
			stopwatch,
		)
		orderQueue.Append(o)
	}
	elapsed := time.Since(stopwatch)
	fmt.Printf("\n\nElapsed: %s\nTransactions per second: %f\n", elapsed, float64(b.N)/elapsed.Seconds())
}
