package orderbook

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestNewOrderTree(t *testing.T) {

	orderTree := NewOrderTree()

	if !(orderTree.volume.Equal(decimal.Zero)) {
		t.Errorf("orderTree.volume incorrect, got: %d, want: %d.", orderTree.volume, decimal.Zero)
	}

	if !(orderTree.Length() == 0) {
		t.Errorf("orderTree.Length() incorrect, got: %d, want: %d.", orderTree.Length(), 0)
	}
}
