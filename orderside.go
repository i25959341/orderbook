package orderbook

import (
	"container/list"
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

// NewOrderSide creates new OrderSide manager
func NewOrderSide() *OrderSide {
	return &OrderSide{
		priceTree: &rbtx.RedBlackTreeExtended{
			Tree: rbt.NewWith(func(a, b interface{}) int {
				return a.(decimal.Decimal).Cmp(b.(decimal.Decimal))
			}),
		},
		prices: map[string]*OrderQueue{},
	}
}

// Len returns amount of orders
func (ot *OrderSide) Len() int {
	return ot.numOrders
}

// Depth returns depth of market
func (ot *OrderSide) Depth() int {
	return ot.depth
}

// Append appends order to definite price level
func (ot *OrderSide) Append(o *Order) *list.Element {
	price := o.Price()
	strPrice := price.String()

	priceQueue, ok := ot.prices[strPrice]
	if !ok {
		priceQueue = NewOrderQueue(o.Price())
		ot.prices[strPrice] = priceQueue
		ot.priceTree.Put(price, priceQueue)
		ot.depth++
	}
	ot.numOrders++
	return priceQueue.Append(o)
}

// Remove removes order from definite price level
func (ot *OrderSide) Remove(e *list.Element) *Order {
	price := e.Value.(*Order).Price()
	strPrice := price.String()

	priceQueue := ot.prices[strPrice]
	o := priceQueue.Remove(e)

	if priceQueue.Len() == 0 {
		delete(ot.prices, strPrice)
		ot.priceTree.Remove(price)
		ot.depth--
	}

	ot.numOrders--
	return o
}

// MaxPriceQueue returns maximal level of price
func (ot *OrderSide) MaxPriceQueue() *OrderQueue {
	if ot.depth > 0 {
		if value, found := ot.priceTree.GetMax(); found {
			return value.(*OrderQueue)
		}
	}
	return nil
}

// MinPriceQueue returns maximal level of price
func (ot *OrderSide) MinPriceQueue() *OrderQueue {
	if ot.depth > 0 {
		if value, found := ot.priceTree.GetMin(); found {
			return value.(*OrderQueue)
		}
	}
	return nil
}

func (ot *OrderSide) String() string {
	sb := strings.Builder{}

	prices := []decimal.Decimal{}
	for k := range ot.prices {
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
		sb.WriteString(fmt.Sprintf("\n%s -> %s", strings.Repeat(" ", maxLen-len(price))+price, ot.prices[price].Volume()))
	}

	return sb.String()
}
