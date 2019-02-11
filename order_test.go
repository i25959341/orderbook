package orderbook

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestNewOrder(t *testing.T) {
	t.Log(NewOrder("order1", decimal.New(100, 0), decimal.New(100, 0), time.Now().UTC()))
}
