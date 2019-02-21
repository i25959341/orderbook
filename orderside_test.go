package orderbook

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestOrderSide(t *testing.T) {
	ot := NewOrderSide()

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

func TestOrderSideJSON(t *testing.T) {
	data := NewOrderSide()

	data.Append(NewOrder("one", Buy, decimal.New(11, -1), decimal.New(11, 1), time.Now().UTC()))
	data.Append(NewOrder("two", Buy, decimal.New(22, -1), decimal.New(22, 1), time.Now().UTC()))
	data.Append(NewOrder("three", Sell, decimal.New(33, -1), decimal.New(33, 1), time.Now().UTC()))
	data.Append(NewOrder("four", Sell, decimal.New(44, -1), decimal.New(44, 1), time.Now().UTC()))

	data.Append(NewOrder("five", Buy, decimal.New(11, -1), decimal.New(11, 1), time.Now().UTC()))
	data.Append(NewOrder("six", Buy, decimal.New(22, -1), decimal.New(22, 1), time.Now().UTC()))
	data.Append(NewOrder("seven", Sell, decimal.New(33, -1), decimal.New(33, 1), time.Now().UTC()))
	data.Append(NewOrder("eight", Sell, decimal.New(44, -1), decimal.New(44, 1), time.Now().UTC()))

	result, _ := json.Marshal(data)
	t.Log(string(result))

	data = NewOrderSide()
	if err := json.Unmarshal(result, data); err != nil {
		t.Fatal(err)
	}

	t.Log(data)

	err := json.Unmarshal([]byte(`[{"side":"fake"}]`), &data)
	if err == nil {
		t.Fatal("can unmarshal unsupported value")
	}
}

func TestPriceFinding(t *testing.T) {
	os := NewOrderSide()

	os.Append(NewOrder("five", Sell, decimal.New(5, 0), decimal.New(130, 0), time.Now().UTC()))
	os.Append(NewOrder("one", Sell, decimal.New(5, 0), decimal.New(170, 0), time.Now().UTC()))
	os.Append(NewOrder("eight", Sell, decimal.New(5, 0), decimal.New(100, 0), time.Now().UTC()))
	os.Append(NewOrder("two", Sell, decimal.New(5, 0), decimal.New(160, 0), time.Now().UTC()))
	os.Append(NewOrder("four", Sell, decimal.New(5, 0), decimal.New(140, 0), time.Now().UTC()))
	os.Append(NewOrder("six", Sell, decimal.New(5, 0), decimal.New(120, 0), time.Now().UTC()))
	os.Append(NewOrder("three", Sell, decimal.New(5, 0), decimal.New(150, 0), time.Now().UTC()))
	os.Append(NewOrder("seven", Sell, decimal.New(5, 0), decimal.New(110, 0), time.Now().UTC()))

	if !os.Volume().Equals(decimal.New(40, 0)) {
		t.Fatal("invalid volume")
	}

	if !os.LessThan(decimal.New(101, 0)).Price().Equals(decimal.New(100, 0)) ||
		!os.LessThan(decimal.New(150, 0)).Price().Equals(decimal.New(140, 0)) ||
		os.LessThan(decimal.New(100, 0)) != nil {
		t.Fatal("LessThan return invalid price")
	}

	if !os.GreaterThan(decimal.New(169, 0)).Price().Equals(decimal.New(170, 0)) ||
		!os.GreaterThan(decimal.New(150, 0)).Price().Equals(decimal.New(160, 0)) ||
		os.GreaterThan(decimal.New(170, 0)) != nil {
		t.Fatal("GreaterThan return invalid price")
	}

	t.Log(os.LessThan(decimal.New(101, 0)))
	t.Log(os.GreaterThan(decimal.New(169, 0)))
}

func BenchmarkOrderSide(b *testing.B) {
	ot := NewOrderSide()
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
