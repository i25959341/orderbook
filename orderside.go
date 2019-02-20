package orderbook

import (
	"container/list"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	rbtx "github.com/emirpasic/gods/examples/redblacktreeextended"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/shopspring/decimal"
)

// OrderSide implements facade to operations with order queue
type OrderSide struct {
	priceTree *rbtx.RedBlackTreeExtended
	prices    map[string]*OrderQueue

	numOrders int
	depth     int
}

func rbtComparator(a, b interface{}) int {
	return a.(decimal.Decimal).Cmp(b.(decimal.Decimal))
}

// NewOrderSide creates new OrderSide manager
func NewOrderSide() *OrderSide {
	return &OrderSide{
		priceTree: &rbtx.RedBlackTreeExtended{
			Tree: rbt.NewWith(rbtComparator),
		},
		prices: map[string]*OrderQueue{},
	}
}

// Len returns amount of orders
func (os *OrderSide) Len() int {
	return os.numOrders
}

// Depth returns depth of market
func (os *OrderSide) Depth() int {
	return os.depth
}

// Append appends order to definite price level
func (os *OrderSide) Append(o *Order) *list.Element {
	price := o.Price()
	strPrice := price.String()

	priceQueue, ok := os.prices[strPrice]
	if !ok {
		priceQueue = NewOrderQueue(o.Price())
		os.prices[strPrice] = priceQueue
		os.priceTree.Put(price, priceQueue)
		os.depth++
	}
	os.numOrders++
	return priceQueue.Append(o)
}

// Remove removes order from definite price level
func (os *OrderSide) Remove(e *list.Element) *Order {
	price := e.Value.(*Order).Price()
	strPrice := price.String()

	priceQueue := os.prices[strPrice]
	o := priceQueue.Remove(e)

	if priceQueue.Len() == 0 {
		delete(os.prices, strPrice)
		os.priceTree.Remove(price)
		os.depth--
	}

	os.numOrders--
	return o
}

// MaxPriceQueue returns maximal level of price
func (os *OrderSide) MaxPriceQueue() *OrderQueue {
	if os.depth > 0 {
		if value, found := os.priceTree.GetMax(); found {
			return value.(*OrderQueue)
		}
	}
	return nil
}

// MinPriceQueue returns maximal level of price
func (os *OrderSide) MinPriceQueue() *OrderQueue {
	if os.depth > 0 {
		if value, found := os.priceTree.GetMin(); found {
			return value.(*OrderQueue)
		}
	}
	return nil
}

// Orders returns all of *list.Element orders
func (os *OrderSide) Orders() (orders []*list.Element) {
	for _, price := range os.prices {
		iter := price.Head()
		for iter != nil {
			orders = append(orders, iter)
			iter = iter.Next()
		}
	}
	return
}

func (os *OrderSide) String() string {
	sb := strings.Builder{}

	prices := []decimal.Decimal{}
	for k := range os.prices {
		num, _ := decimal.NewFromString(k)
		prices = append(prices, num)
	}

	sort.Slice(prices, func(i, j int) bool {
		return prices[i].GreaterThan(prices[j])
	})

	var (
		strPrices   []string
		maxLen      int
		strPrice    string
		strPriceLen int
	)
	for _, price := range prices {
		strPrice = price.String()
		strPriceLen = len(strPrice)
		if strPriceLen > maxLen {
			maxLen = strPriceLen
		}
		strPrices = append(strPrices, price.String())
	}

	for _, price := range strPrices {
		sb.WriteString(fmt.Sprintf("\n%s -> %s", strings.Repeat(" ", maxLen-len(price))+price, os.prices[price].Volume()))
	}

	return sb.String()
}

func (os *OrderSide) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			NumOrders int                    `json:"numOrders"`
			Depth     int                    `json:"depth"`
			Prices    map[string]*OrderQueue `json:"prices"`
		}{
			NumOrders: os.numOrders,
			Depth:     os.depth,
			Prices:    os.prices,
		},
	)
}

func (os *OrderSide) UnmarshalJSON(data []byte) error {
	obj := struct {
		NumOrders int                    `json:"numOrders"`
		Depth     int                    `json:"depth"`
		Prices    map[string]*OrderQueue `json:"prices"`
	}{}

	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	os.numOrders = obj.NumOrders
	os.depth = obj.Depth
	os.prices = obj.Prices
	os.priceTree = &rbtx.RedBlackTreeExtended{
		Tree: rbt.NewWith(rbtComparator),
	}

	for price, queue := range os.prices {
		os.priceTree.Put(decimal.RequireFromString(price), queue)
	}

	return nil
}
