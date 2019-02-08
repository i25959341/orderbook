package orderbook

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestNewOrder(t *testing.T) {
	var orderList OrderQueue
	dummyOrder := make(map[string]string)
	dummyOrder["timestamp"] = strconv.Itoa(testTimestamp)
	dummyOrder["quantity"] = testQuanity.String()
	dummyOrder["price"] = testPrice.String()
	dummyOrder["order_id"] = strconv.Itoa(testOrderId)
	dummyOrder["trade_id"] = strconv.Itoa(testTradeId)

	order := NewOrderFromMap(dummyOrder, &orderList)

	if !(order.timestamp == testTimestamp) {
		t.Errorf("Timesmape incorrect, got: %d, want: %d.", order.timestamp, testTimestamp)
	}

	if !(order.quantity.Equal(testQuanity)) {
		t.Errorf("quantity incorrect, got: %d, want: %d.", order.quantity, testQuanity)
	}

	if !(order.price.Equal(testPrice)) {
		t.Errorf("price incorrect, got: %d, want: %d.", order.price, testPrice)
	}

	if !(order.orderID == strconv.Itoa(testOrderId)) {
		t.Errorf("order id incorrect, got: %s, want: %d.", order.orderID, testOrderId)
	}

	if !(order.tradeID == strconv.Itoa(testTradeId)) {
		t.Errorf("trade id incorrect, got: %s, want: %d.", order.tradeID, testTradeId)
	}
}

func TestOrder(t *testing.T) {
	orderList := NewOrderQueue(testPrice)

	dummyOrder := make(map[string]string)
	dummyOrder["timestamp"] = strconv.Itoa(testTimestamp)
	dummyOrder["quantity"] = testQuanity.String()
	dummyOrder["price"] = testPrice.String()
	dummyOrder["order_id"] = strconv.Itoa(testOrderId)
	dummyOrder["trade_id"] = strconv.Itoa(testTradeId)

	order := NewOrderFromMap(dummyOrder, orderList)

	orderList.Append(order)

	order.Update(testQuanity1, testTimestamp1)

	if !(order.quantity.Equal(testQuanity1)) {
		t.Errorf("order id incorrect, got: %s, want: %d.", order.orderID, testOrderId)
	}

	if !(order.timestamp == testTimestamp1) {
		t.Errorf("trade id incorrect, got: %s, want: %d.", order.tradeID, testTradeId)
	}
}

func BenchmarkOrder(b *testing.B) {
	orderList := NewOrderQueue(testPrice)

	stopwatch := time.Now()
	for i := 0; i < b.N; i++ {
		order := NewOrderFromMap(map[string]string{
			"timestamp": strconv.Itoa(i),
			"quantity":  testQuanity.String(),
			"price":     decimal.New(int64(i), 0).String(),
			"order_id":  strconv.Itoa(i),
			"trade_id":  strconv.Itoa(i),
		}, orderList)
		orderList.Append(order)
	}
	elapsed := time.Since(stopwatch)
	fmt.Printf("\n\nElapsed: %s\nTransactions per second: %f\n", elapsed, float64(b.N)/elapsed.Seconds())
}
