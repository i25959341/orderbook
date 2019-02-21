package orderbook

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestNewOrder(t *testing.T) {
	t.Log(NewOrder("order-1", Sell, decimal.New(100, 0), decimal.New(100, 0), time.Now().UTC()))
}

func TestOrderJSON(t *testing.T) {
	data := []*Order{
		NewOrder("one", Buy, decimal.New(11, -1), decimal.New(11, 1), time.Now().UTC()),
		NewOrder("two", Buy, decimal.New(22, -1), decimal.New(22, 1), time.Now().UTC()),
		NewOrder("three", Sell, decimal.New(33, -1), decimal.New(33, 1), time.Now().UTC()),
		NewOrder("four", Sell, decimal.New(44, -1), decimal.New(44, 1), time.Now().UTC()),
	}

	result, _ := json.Marshal(data)
	t.Log(string(result))

	data = []*Order{}

	_ = json.Unmarshal(result, &data)
	t.Log(data)

	err := json.Unmarshal([]byte(`[{"side":"fake"}]`), &data)
	if err == nil {
		t.Fatal("can unmarshal unsupported value")
	}
}
