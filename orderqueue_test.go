package orderbook

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestOrderQueue(t *testing.T) {
	price := decimal.New(100, 0)
	oq := NewOrderQueue(price)

	o1 := NewOrder(
		"order-1",
		Buy,
		decimal.New(100, 0),
		decimal.New(100, 0),
		time.Now().UTC(),
	)

	o2 := NewOrder(
		"order-2",
		Buy,
		decimal.New(100, 0),
		decimal.New(100, 0),
		time.Now().UTC(),
	)

	head := oq.Append(o1)
	tail := oq.Append(o2)

	if head == nil || tail == nil {
		t.Fatal("Could not append order to the OrderQueue")
	}

	if !oq.Volume().Equal(decimal.New(200, 0)) {
		t.Fatalf("Invalid order volume (have: %s, want: 200", oq.Volume())
	}

	if head.Value.(*Order) != o1 || tail.Value.(*Order) != o2 {
		t.Fatal("Invalid element value")
	}

	if oq.Head() != head || oq.Tail() != tail {
		t.Fatal("Invalid element position")
	}

	if oq.Head().Next() != oq.Tail() || oq.Tail().Prev() != head ||
		oq.Head().Prev() != nil || oq.Tail().Next() != nil {
		t.Fatal("Invalid element link")
	}

	o1 = NewOrder(
		"order-3",
		Buy,
		decimal.New(200, 0),
		decimal.New(200, 0),
		time.Now().UTC(),
	)

	oq.Update(head, o1)
	if !oq.Volume().Equal(decimal.New(300, 0)) {
		t.Fatalf("Invalid order volume (have: %s, want: 300", oq.Volume())
	}

	if o := oq.Remove(head); o != o1 {
		t.Fatal("Invalid element value")
	}

	if !oq.Volume().Equal(decimal.New(100, 0)) {
		t.Fatalf("Invalid order volume (have: %s, want: 100", oq.Volume())
	}

	t.Log(oq)
}

func TestOrderQueueJSON(t *testing.T) {
	data := NewOrderQueue(decimal.New(111, 0))

	data.Append(NewOrder("one", Buy, decimal.New(11, -1), decimal.New(11, 1), time.Now().UTC()))
	data.Append(NewOrder("two", Buy, decimal.New(22, -1), decimal.New(22, 1), time.Now().UTC()))
	data.Append(NewOrder("three", Sell, decimal.New(33, -1), decimal.New(33, 1), time.Now().UTC()))
	data.Append(NewOrder("four", Sell, decimal.New(44, -1), decimal.New(44, 1), time.Now().UTC()))

	result, _ := json.Marshal(data)
	t.Log(string(result))

	data = NewOrderQueue(decimal.Zero)
	if err := json.Unmarshal(result, data); err != nil {
		t.Fatal(err)
	}

	t.Log(data)

	err := json.Unmarshal([]byte(`[{"side":"fake"}]`), &data)
	if err == nil {
		t.Fatal("can unmarshal unsupported value")
	}
}

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
