package orderbook

import (
	"strconv"
	"testing"

	"github.com/shopspring/decimal"
)

func TestNewOrderTree(t *testing.T) {
	orderTree := NewOrderTree()

	dummyOrder := make(map[string]string)
	dummyOrder["timestamp"] = strconv.Itoa(testTimestamp)
	dummyOrder["quantity"] = testQuanity.String()
	dummyOrder["price"] = testPrice.String()
	dummyOrder["order_id"] = strconv.Itoa(testOrderId)
	dummyOrder["trade_id"] = strconv.Itoa(testTradeId)

	dummyOrder1 := make(map[string]string)
	dummyOrder1["timestamp"] = strconv.Itoa(testTimestamp1)
	dummyOrder1["quantity"] = testQuanity1.String()
	dummyOrder1["price"] = testPrice1.String()
	dummyOrder1["order_id"] = strconv.Itoa(testOrderId1)
	dummyOrder1["trade_id"] = strconv.Itoa(testTradeId1)

	if !(orderTree.volume.Equal(decimal.Zero)) {
		t.Errorf("orderTree.volume incorrect, got: %d, want: %d.", orderTree.volume, decimal.Zero)
	}

	if !(orderTree.Length() == 0) {
		t.Errorf("orderTree.Length() incorrect, got: %d, want: %d.", orderTree.Length(), 0)
	}

	orderTree.InsertOrder(dummyOrder)
	orderTree.InsertOrder(dummyOrder1)

	if !(orderTree.PriceExist(testPrice)) {
		t.Errorf("orderTree.numOrders incorrect, got: %d, want: %d.", orderTree.numOrders, 2)
	}

	if !(orderTree.PriceExist(testPrice1)) {
		t.Errorf("orderTree.numOrders incorrect, got: %d, want: %d.", orderTree.numOrders, 2)
	}

	if !(orderTree.Length() == 2) {
		t.Errorf("orderTree.numOrders incorrect, got: %d, want: %d.", orderTree.numOrders, 2)
	}

	orderTree.RemoveOrderById(dummyOrder1["order_id"])
	orderTree.RemoveOrderById(dummyOrder["order_id"])

	if !(orderTree.Length() == 0) {
		t.Errorf("orderTree.numOrders incorrect, got: %d, want: %d.", orderTree.numOrders, 2)
	}
}
