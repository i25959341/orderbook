package orderbook

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestOrderTree(t *testing.T) {
	ot := NewOrderTree()

	o1 := NewOrder(
		"order-1",
		Buy,
		decimal.New(10, 0),
		decimal.New(10, 0),
		time.Now().UTC(),
	)

	o2 := NewOrder(
		"order-2",
		Buy,
		decimal.New(10, 0),
		decimal.New(20, 0),
		time.Now().UTC(),
	)

	if ot.MinPriceQueue() != nil || ot.MaxPriceQueue() != nil {
		t.Fatal("invalid price levels")
	}

	el1 := ot.Append(o1)

	if ot.MinPriceQueue() != ot.MaxPriceQueue() {
		t.Fatal("invalid price levels")
	}

	el2 := ot.Append(o2)

	if ot.Depth() != 2 {
		t.Fatal("invalid depth")
	}

	if ot.Len() != 2 {
		t.Fatal("invalid orders count")
	}

	t.Log(ot)

	if ot.MinPriceQueue().Head() != el1 || ot.MinPriceQueue().Tail() != el1 ||
		ot.MaxPriceQueue().Head() != el2 || ot.MaxPriceQueue().Tail() != el2 {
		t.Fatal("invalid price levels")
	}

	if o := ot.Remove(el1); o != o1 {
		t.Fatal("invalid order")
	}

	if ot.MinPriceQueue() != ot.MaxPriceQueue() {
		t.Fatal("invalid price levels")
	}

	t.Log(ot)
}

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
