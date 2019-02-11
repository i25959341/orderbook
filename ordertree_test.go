package orderbook

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestOrderTree(t *testing.T) {
	ot := NewOrderTree()
	length := 15

	for i := 0; i < length; i++ {
		ot.CreateOrder(fmt.Sprintf("order-%d", i), decimal.New(10, 0), decimal.New(int64(i), 0), time.Now().UTC())
	}
	for i := 0; i < length; i++ {
		ot.CreateOrder(fmt.Sprintf("order2-%d", i), decimal.New(10, 0), decimal.New(int64(i), 0), time.Now().UTC())
	}

	t.Log(ot)
}

func BenchmarkOrderTree(b *testing.B) {
	ot := NewOrderTree()
	stopwatch := time.Now()
	for i := 0; i < b.N; i++ {
		ot.CreateOrder(fmt.Sprintf("order-%d", i), decimal.New(10, 0), decimal.New(int64(i), 0), stopwatch)
	}
	elapsed := time.Since(stopwatch)
	fmt.Printf("\n\nElapsed: %s\nTransactions per second: %f\n", elapsed, float64(b.N)/elapsed.Seconds())
}
