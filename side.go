package orderbook

// Side of the order
type Side int

// Sell (asks) or Buy (bids)
const (
	Sell Side = iota
	Buy
)

func (s Side) String() string {
	if s == Buy {
		return "buy"
	}

	return "sell"
}
