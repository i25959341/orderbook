package orderbook

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

func createDepth(quantity decimal.Decimal) (ob *OrderBook) {
	ob = NewOrderBook()

	for i := 50; i < 100; i = i + 10 {
		ob.ProcessLimitOrder(Buy, fmt.Sprintf("buy-%d", i), quantity, decimal.New(int64(i), 0))
	}

	for i := 100; i < 150; i = i + 10 {
		ob.ProcessLimitOrder(Sell, fmt.Sprintf("sell-%d", i), quantity, decimal.New(int64(i), 0))
	}

	return
}

func TestLimitPlace(t *testing.T) {
	ob := NewOrderBook()
	quantity := decimal.New(2, 0)
	for i := 50; i < 100; i = i + 10 {
		done, partial, err := ob.ProcessLimitOrder(Buy, fmt.Sprintf("buy-%d", i), quantity, decimal.New(int64(i), 0))
		if len(done) != 0 {
			t.Fatal("OrderBook failed to process limit order (done is not empty)")
		}
		if partial != nil {
			t.Fatal("OrderBook failed to process limit order (partial is not empty)")
		}
		if partial != nil {
			t.Fatal(err)
		}
	}

	for i := 100; i < 150; i = i + 10 {
		done, partial, err := ob.ProcessLimitOrder(Sell, fmt.Sprintf("sell-%d", i), quantity, decimal.New(int64(i), 0))
		if len(done) != 0 {
			t.Fatal("OrderBook failed to process limit order (done is not empty)")
		}
		if partial != nil {
			t.Fatal("OrderBook failed to process limit order (partial is not empty)")
		}
		if partial != nil {
			t.Fatal(err)
		}
	}

	t.Log(ob)
	return
}

func TestLimitDone(t *testing.T) {
	ob := createDepth(decimal.New(2, 0))

	done, partial, err := ob.ProcessLimitOrder(Buy, "order-b100", decimal.New(1, 0), decimal.New(100, 0))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Done:", done)
	if done[0].ID() != "order-b100" {
		t.Fatal("Wrong done id")
	}

	t.Log("Partial:", partial)
	if partial.ID() != "sell-100" {
		t.Fatal("Wrong partial id")
	}

	t.Log(ob)

	done, partial, err = ob.ProcessLimitOrder(Buy, "order-b150", decimal.New(10, 0), decimal.New(150, 0))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Done:", done)
	if len(done) != 5 {
		t.Fatal("Wrong done quantity")
	}

	t.Log("Partial:", partial)
	if partial.ID() != "order-b150" {
		t.Fatal("Wrong partial id")
	}

	t.Log(ob)

	done, partial, err = ob.ProcessLimitOrder(Sell, "order-s40", decimal.New(11, 0), decimal.New(40, 0))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Done:", done)
	if len(done) != 7 {
		t.Fatal("Wrong done quantity")
	}

	if partial != nil {
		t.Fatal("Wrong partial")
	}

	t.Log(ob)
}

// ob := NewOrderBook()
// for i := 50; i < 100; i = i + 10 {
// 	ob.ProcessLimitOrder(Buy, fmt.Sprintf("o-%d", i), decimal.New(1, 0), decimal.New(int64(i), 0))
// }
// for i := 100; i < 150; i = i + 10 {
// 	ob.ProcessLimitOrder(Sell, fmt.Sprintf("o-%d", i), decimal.New(1, 0), decimal.New(int64(i), 0))
// }
// t.Log(ob.ProcessMarketOrder(Buy, decimal.New(3, 0)))
// t.Log(ob.ProcessMarketOrder(Sell, decimal.New(3, 0)))
// t.Log(ob.ProcessLimitOrder(Buy, "order-b100", decimal.New(20, 0), decimal.New(100, 0)))
// t.Log(ob.ProcessLimitOrder(Sell, "order-s100", decimal.New(30, 0), decimal.New(10, 0)))
// t.Log(ob.ProcessLimitOrder(Buy, "order-b1000", decimal.New(30, 0), decimal.New(1000, 0)))
//t.Log(createDepth(decimal.New(1, 0)))
