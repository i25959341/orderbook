package orderbook

import (
	"encoding/json"
	"reflect"
)

// Side of the order
type Side int

// Sell (asks) or Buy (bids)
const (
	Sell Side = iota
	Buy
)

// String implements fmt.Stringer interface
func (s Side) String() string {
	if s == Buy {
		return "buy"
	}

	return "sell"
}

// MarshalJSON implements json.Marshaler interface
func (s Side) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (s *Side) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"buy"`:
		*s = Buy
	case `"sell"`:
		*s = Sell
	default:
		return &json.UnsupportedValueError{
			Value: reflect.New(reflect.TypeOf(data)),
			Str:   string(data),
		}
	}

	return nil
}
