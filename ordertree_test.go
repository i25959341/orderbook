package orderbook

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func BenchmarkOrderTree(b *testing.B) {
	ot := NewOrderTree()
	stopwatch := time.Now()

	var o *Order
	for i := 0; i < b.N; i++ {
		o = NewOrder(
			fmt.Sprintf("order-%d", i),
			Buy,
			decimal.New(10, 0),
			decimal.New(int64(i), 0),
			stopwatch,
		)
		ot.Append(o)
	}
	elapsed := time.Since(stopwatch)
	fmt.Printf("\n\nElapsed: %s\nTransactions per second: %f\n", elapsed, float64(b.N)/elapsed.Seconds())
}
