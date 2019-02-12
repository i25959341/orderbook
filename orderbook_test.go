package orderbook

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

func TestLimitOrder(t *testing.T) {
	ob := NewOrderbook()

	for i := 50; i < 100; i = i + 10 {
		ob.ProcessLimitOrder(Buy, fmt.Sprintf("o-%d", i), decimal.New(1, 0), decimal.New(int64(i), 0))
	}

	for i := 100; i < 150; i = i + 10 {
		ob.ProcessLimitOrder(Sell, fmt.Sprintf("o-%d", i), decimal.New(1, 0), decimal.New(int64(i), 0))
	}

	t.Log(ob.ProcessMarketOrder(Buy, decimal.New(3, 0)))

	t.Log(ob.ProcessMarketOrder(Sell, decimal.New(3, 0)))

	// t.Log(ob.ProcessLimitOrder(Buy, "order-b100", decimal.New(20, 0), decimal.New(100, 0)))

	// t.Log(ob.ProcessLimitOrder(Sell, "order-s100", decimal.New(30, 0), decimal.New(10, 0)))

	// t.Log(ob.ProcessLimitOrder(Buy, "order-b1000", decimal.New(30, 0), decimal.New(1000, 0)))

	t.Log(ob)
}
