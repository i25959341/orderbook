package orderbook

import (
	"strconv"
	"testing"
)

func TestNewOrder(t *testing.T) {
	var order_list OrderList
	dummyOrder := make(map[string]string)
	dummyOrder["timestamp"] = strconv.Itoa(testTimestamp)
	dummyOrder["quantity"] = testQuanity.String()
	dummyOrder["price"] = testPrice.String()
	dummyOrder["order_id"] = strconv.Itoa(testOrderId)
	dummyOrder["trade_id"] = strconv.Itoa(testTradeId)

	order := NewOrder(dummyOrder, &order_list)

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
	orderList := NewOrderList(testPrice)

	dummyOrder := make(map[string]string)
	dummyOrder["timestamp"] = strconv.Itoa(testTimestamp)
	dummyOrder["quantity"] = testQuanity.String()
	dummyOrder["price"] = testPrice.String()
	dummyOrder["order_id"] = strconv.Itoa(testOrderId)
	dummyOrder["trade_id"] = strconv.Itoa(testTradeId)

	order := NewOrder(dummyOrder, orderList)

	orderList.AppendOrder(order)

	order.UpdateQuantity(testQuanity1, testTimestamp1)

	if !(order.quantity.Equal(testQuanity1)) {
		t.Errorf("order id incorrect, got: %s, want: %d.", order.orderID, testOrderId)
	}

	if !(order.timestamp == testTimestamp1) {
		t.Errorf("trade id incorrect, got: %s, want: %d.", order.tradeID, testTradeId)
	}
}
