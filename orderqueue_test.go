package orderbook

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestNewOrderQueue(t *testing.T) {
	t.Log(NewOrderQueue(decimal.New(100, 0)))
}

func TestOrderAppendRemove(t *testing.T) {
	price := decimal.New(100, 0)
	volume := decimal.New(10, 0)
	length := 10

	oq := NewOrderQueue(price)
	orders := make([]*Order, length)
	for i := 0; i < length; i++ {
		orders[i] = NewOrder(fmt.Sprintf("order-%d", i), volume, price, time.Now().UTC())
		if err := oq.Append(orders[i]); err != nil {
			t.Fatal(err)
		}
	}
	t.Log(oq)

	if oq.Head() != orders[0] {
		t.Fatal("invalid head order")
	}

	if oq.Len() != length {
		t.Fatalf("wrong length: got %d, want %d", oq.Len(), length)
	}

	if !oq.Volume().Equal(volume.Mul(decimal.New(int64(length), 0))) {
		t.Fatalf("wrong length: got %s, want %s", oq.Volume(), volume.Mul(decimal.New(int64(length), 0)))
	}

	if err := NewOrderQueue(price).Append(orders[1]); err == nil {
		t.Fatal("it is possible to append already linked order")
	}

	if err := oq.Append(NewOrder("fakeOrder", decimal.Zero, price, time.Now().UTC())); err == nil {
		t.Fatal("it is possible to append zero volume order")
	}

	if err := oq.Append(NewOrder("fakeOrder", decimal.New(10, 0), decimal.Zero, time.Now().UTC())); err == nil {
		t.Fatal("it is possible to append zero price order")
	}

	if err := oq.Append(NewOrder("fakeOrder", decimal.New(10, 0), price.Add(decimal.New(10, 0)), time.Now().UTC())); err == nil {
		t.Fatal("it is possible to append different price order")
	}

	if err := oq.Remove(NewOrder("fakeOrder", decimal.New(10, 0), decimal.Zero, time.Now().UTC())); err == nil {
		t.Fatal("it is possible to remove invalid order")
	}

	if err := oq.Remove(orders[2]); err != nil {
		t.Fatal(err)
	}

	t.Log(oq)

	if err := oq.Remove(orders[0]); err != nil {
		t.Fatal(err)
	}

	if oq.Head() != orders[1] {
		t.Fatalf("invalid head: got %v, want %v", oq.Head(), orders[1])
	}
	t.Log(oq)
}

func BenchmarkOrderQueue(b *testing.B) {
	price := decimal.New(100, 0)
	orderQueue := NewOrderQueue(price)
	stopwatch := time.Now()
	for i := 0; i < b.N; i++ {
		order := NewOrder(fmt.Sprintf("order-%d", i), decimal.New(100, 0), decimal.New(int64(i), 0), time.Now().UTC())
		orderQueue.Append(order)
	}
	elapsed := time.Since(stopwatch)
	fmt.Printf("\n\nElapsed: %s\nTransactions per second: %f\n", elapsed, float64(b.N)/elapsed.Seconds())
}
