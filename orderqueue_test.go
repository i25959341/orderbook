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

	if oq.Length() != length {
		t.Fatalf("wrong length: got %d, want %d", oq.Length(), length)
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
	if orders[1].next != orders[3] || orders[3].prev != orders[1] {
		t.Fatal("invalid order chain")
	}

	t.Log(oq)

	if err := oq.Remove(orders[0]); err != nil {
		t.Fatal(err)
	}

	if oq.Head() != orders[1] {
		t.Fatalf("invalid head: got %v, want %v", oq.Head(), orders[1])
	}
	t.Log(oq)

	if err := oq.Remove(orders[length-1]); err != nil {
		t.Fatal(err)
	}

	if oq.Tail() != orders[length-2] {
		t.Fatalf("invalid tail: got %v, want %v", oq.Tail(), orders[length-2])
	}
	t.Log(oq)

	if err := oq.MoveToTail(oq.Head()); err != nil {
		t.Fatal(err)
	}

	if oq.Head() != orders[3] {
		t.Fatalf("invalid head: got %v, want %v", oq.Head(), orders[3])
	}

	t.Log(oq)
}

func TestOrderMoveToTail(t *testing.T) {
	price := decimal.New(100, 0)
	volume := decimal.New(10, 0)
	length := 2

	oq := NewOrderQueue(price)
	orders := make([]*Order, length)
	for i := 0; i < length; i++ {
		orders[i] = NewOrder(fmt.Sprintf("order-%d", i), volume, price, time.Now().UTC())
		if err := oq.Append(orders[i]); err != nil {
			t.Fatal(err)
		}
	}
	t.Log(oq)

	if err := oq.MoveToTail(NewOrder("fakeOrder", decimal.New(10, 0), decimal.Zero, time.Now().UTC())); err == nil {
		t.Fatal("it is possible to move order from another queue")
	}

	oq.MoveToTail(orders[0])
	if oq.Head() != orders[1] || oq.Tail() != orders[0] {
		t.Fatal("invalid order head or tail")
	}

	if orders[1].Prev() != nil || oq.Tail().Next() != nil {
		t.Fatal("invalid order connection")
	}
	t.Log(oq)

	orders = append(orders, NewOrder("Order-2", volume, price, time.Now().UTC()))
	oq.Append(orders[2])
	t.Log(oq)

	oq.MoveToTail(orders[0])
	if oq.Head() != orders[1] || oq.Tail() != orders[0] {
		t.Fatal("invalid order connection")
	}
	t.Log(oq)

	oq.MoveToTail(orders[0])

	oq.Remove(orders[0])
	t.Log(oq)

	oq.Remove(orders[1])
	t.Log(oq)

	oq.Remove(orders[2])
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
