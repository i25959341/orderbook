package orderbook

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func addDepth(ob *OrderBook, prefix string, quantity decimal.Decimal) {
	for i := 50; i < 100; i = i + 10 {
		ob.ProcessLimitOrder(Buy, fmt.Sprintf("%sbuy-%d", prefix, i), quantity, decimal.New(int64(i), 0))
	}

	for i := 100; i < 150; i = i + 10 {
		ob.ProcessLimitOrder(Sell, fmt.Sprintf("%ssell-%d", prefix, i), quantity, decimal.New(int64(i), 0))
	}

	return
}

func TestLimitPlace(t *testing.T) {
	ob := NewOrderBook()
	quantity := decimal.New(2, 0)
	for i := 50; i < 100; i = i + 10 {
		done, partial, partialQty, err := ob.ProcessLimitOrder(Buy, fmt.Sprintf("buy-%d", i), quantity, decimal.New(int64(i), 0))
		if len(done) != 0 {
			t.Fatal("OrderBook failed to process limit order (done is not empty)")
		}
		if partial != nil {
			t.Fatal("OrderBook failed to process limit order (partial is not empty)")
		}
		if partialQty.Sign() != 0 {
			t.Fatal("OrderBook failed to process limit order (partialQty is not zero)")
		}
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := 100; i < 150; i = i + 10 {
		done, partial, partialQty, err := ob.ProcessLimitOrder(Sell, fmt.Sprintf("sell-%d", i), quantity, decimal.New(int64(i), 0))
		if len(done) != 0 {
			t.Fatal("OrderBook failed to process limit order (done is not empty)")
		}
		if partial != nil {
			t.Fatal("OrderBook failed to process limit order (partial is not empty)")
		}
		if partialQty.Sign() != 0 {
			t.Fatal("OrderBook failed to process limit order (partialQty is not zero)")
		}
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log(ob)

	if ob.Order("fake") != nil {
		t.Fatal("can get fake order")
	}

	if ob.Order("sell-100") == nil {
		t.Fatal("can't get real order")
	}

	t.Log(ob.Depth())
	return
}

func TestLimitProcess(t *testing.T) {
	ob := NewOrderBook()
	addDepth(ob, "", decimal.New(2, 0))

	done, partial, partialQty, err := ob.ProcessLimitOrder(Buy, "order-b100", decimal.New(1, 0), decimal.New(100, 0))
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

	if !partialQty.Equal(decimal.New(1, 0)) {
		t.Fatal("Wrong partial quantity processed")
	}

	t.Log(ob)

	done, partial, partialQty, err = ob.ProcessLimitOrder(Buy, "order-b150", decimal.New(10, 0), decimal.New(150, 0))
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

	if !partialQty.Equal(decimal.New(9, 0)) {
		t.Fatal("Wrong partial quantity processed", partialQty)
	}

	t.Log(ob)

	if _, _, _, err := ob.ProcessLimitOrder(Sell, "buy-70", decimal.New(11, 0), decimal.New(40, 0)); err == nil {
		t.Fatal("Can add existing order")
	}

	if _, _, _, err := ob.ProcessLimitOrder(Sell, "fake-70", decimal.New(0, 0), decimal.New(40, 0)); err == nil {
		t.Fatal("Can add empty quantity order")
	}

	if _, _, _, err := ob.ProcessLimitOrder(Sell, "fake-70", decimal.New(10, 0), decimal.New(0, 0)); err == nil {
		t.Fatal("Can add zero price")
	}

	if o := ob.CancelOrder("order-b100"); o != nil {
		t.Fatal("Can cancel done order")
	}

	done, partial, partialQty, err = ob.ProcessLimitOrder(Sell, "order-s40", decimal.New(11, 0), decimal.New(40, 0))
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

	if partialQty.Sign() != 0 {
		t.Fatal("Wrong partialQty")
	}

	t.Log(ob)
}

func TestMarketProcess(t *testing.T) {
	ob := NewOrderBook()
	addDepth(ob, "", decimal.New(2, 0))

	done, partial, partialQty, left, err := ob.ProcessMarketOrder(Buy, decimal.New(3, 0))
	if err != nil {
		t.Fatal(err)
	}

	if left.Sign() > 0 {
		t.Fatal("Wrong quantity left")
	}

	if !partialQty.Equal(decimal.New(1, 0)) {
		t.Fatal("Wrong partial quantity left")
	}

	t.Log("Done", done)
	t.Log("Partial", partial)
	t.Log(ob)

	if _, _, _, _, err := ob.ProcessMarketOrder(Buy, decimal.New(0, 0)); err == nil {
		t.Fatal("Can add zero quantity order")
	}

	done, partial, partialQty, left, err = ob.ProcessMarketOrder(Sell, decimal.New(12, 0))
	if err != nil {
		t.Fatal(err)
	}

	if partial != nil {
		t.Fatal("Partial is not nil")
	}

	if partialQty.Sign() != 0 {
		t.Fatal("PartialQty is not nil")
	}

	if len(done) != 5 {
		t.Fatal("Invalid done amount")
	}

	if !left.Equal(decimal.New(2, 0)) {
		t.Fatal("Invalid left amount", left)
	}

	t.Log("Done", done)
	t.Log(ob)
}

func TestOrderBookJSON(t *testing.T) {
	data := NewOrderBook()

	result, _ := json.Marshal(data)
	t.Log(string(result))

	if err := json.Unmarshal(result, data); err != nil {
		t.Fatal(err)
	}

	addDepth(data, "01-", decimal.New(10, 0))
	addDepth(data, "02-", decimal.New(1, 0))
	addDepth(data, "03-", decimal.New(2, 0))

	result, _ = json.Marshal(data)
	t.Log(string(result))

	data = NewOrderBook()
	if err := json.Unmarshal(result, data); err != nil {
		t.Fatal(err)
	}

	t.Log(data)

	err := json.Unmarshal([]byte(`[{"side":"fake"}]`), &data)
	if err == nil {
		t.Fatal("can unmarshal unsupported value")
	}
}

func TestPriceCalculation(t *testing.T) {
	ob := NewOrderBook()
	addDepth(ob, "05-", decimal.New(10, 0))
	addDepth(ob, "10-", decimal.New(10, 0))
	addDepth(ob, "15-", decimal.New(10, 0))
	t.Log(ob)

	price, err := ob.CalculateMarketPrice(Buy, decimal.New(115, 0))
	if err != nil {
		t.Fatal(err)
	}

	if !price.Equal(decimal.New(13150, 0)) {
		t.Fatal("invalid price", price)
	}

	price, err = ob.CalculateMarketPrice(Buy, decimal.New(200, 0))
	if err == nil {
		t.Fatal("invalid quantity count")
	}

	if !price.Equal(decimal.New(18000, 0)) {
		t.Fatal("invalid price", price)
	}

	// -------

	price, err = ob.CalculateMarketPrice(Sell, decimal.New(115, 0))
	if err != nil {
		t.Fatal(err)
	}

	if !price.Equal(decimal.New(8700, 0)) {
		t.Fatal("invalid price", price)
	}

	price, err = ob.CalculateMarketPrice(Sell, decimal.New(200, 0))
	if err == nil {
		t.Fatal("invalid quantity count")
	}

	if !price.Equal(decimal.New(10500, 0)) {
		t.Fatal("invalid price", price)
	}
}

func BenchmarkLimitOrder(b *testing.B) {
	ob := NewOrderBook()
	stopwatch := time.Now()
	for i := 0; i < b.N; i++ {
		addDepth(ob, "05-", decimal.New(10, 0))                                           // 10 ts
		addDepth(ob, "10-", decimal.New(10, 0))                                           // 10 ts
		addDepth(ob, "15-", decimal.New(10, 0))                                           // 10 ts
		ob.ProcessLimitOrder(Buy, "order-b150", decimal.New(160, 0), decimal.New(150, 0)) // 1 ts
		ob.ProcessMarketOrder(Sell, decimal.New(200, 0))                                  // 1 ts = total 32
	}
	elapsed := time.Since(stopwatch)
	fmt.Printf("\n\nElapsed: %s\nTransactions per second (avg): %f\n", elapsed, float64(b.N*32)/elapsed.Seconds())
}
