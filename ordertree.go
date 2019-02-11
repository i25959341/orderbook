package orderbook

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/emirpasic/gods/examples/redblacktreeextended"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/shopspring/decimal"
)

// OrderTree implements facade to operations with order queue
type OrderTree struct {
	priceTree *redblacktreeextended.RedBlackTreeExtended
	prices    map[string]*OrderQueue // Dictionary containing price : OrderList object
	orders    map[string]*Order      // Dictionary containing order_id : Order object
	volume    decimal.Decimal        // Contains total quantity from all Orders in tree
	numOrders int                    // Contains count of Orders in tree
	depth     int                    // Number of different prices in tree (http://en.wikipedia.org/wiki/Order_book_(trading)#Book_depth)
}

// NewOrderTree creates new OrderTree manager
func NewOrderTree() *OrderTree {
	return &OrderTree{
		priceTree: &redblacktreeextended.RedBlackTreeExtended{
			Tree: rbt.NewWith(func(a, b interface{}) int {
				return a.(decimal.Decimal).Cmp(b.(decimal.Decimal))
			}),
		},
		prices: map[string]*OrderQueue{},
		orders: map[string]*Order{},
		volume: decimal.Zero,
	}
}

// Length returns total amount of orders in tree
func (ot *OrderTree) Length() int {
	return len(ot.orders)
}

// CreateOrder creates new order and inserts it to the tree
func (ot *OrderTree) CreateOrder(orderID string, quantity, price decimal.Decimal, timestamp time.Time) error {
	if _, ok := ot.orders[orderID]; ok {
		return ErrOrderExists
	}
	if quantity.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidQuantity
	}
	if price.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidPrice
	}

	priceQueue := ot.getOrCreatePrice(price)
	order := NewOrder(orderID, quantity, price, timestamp)

	priceQueue.Append(order)
	ot.orders[orderID] = order
	ot.volume = ot.volume.Add(order.quantity)
	ot.numOrders++
	return nil
}

// RemoveOrder removes definite order by ID from tree
func (ot *OrderTree) RemoveOrder(orderID string) error {
	order, ok := ot.orders[orderID]
	if !ok {
		return ErrOrderNotExists
	}

	priceStr := order.price.String()

	priceQueue := ot.prices[priceStr]
	priceQueue.Remove(order)

	if priceQueue.Len() == 0 {
		ot.depth--
		ot.priceTree.Remove(order.price)
		delete(ot.prices, priceStr)
	}

	delete(ot.orders, orderID)
	ot.volume = ot.volume.Sub(order.quantity)
	ot.numOrders--
	return nil
}

// MaxPrice returns maximal level of price
func (ot *OrderTree) MaxPrice() decimal.Decimal {
	if ot.depth > 0 {
		value, found := ot.priceTree.GetMax()
		if found {
			return value.(*OrderQueue).price
		}
		return decimal.Zero
	} else {
		return decimal.Zero
	}
}

// MaxPrice returns minimal level of price
func (ot *OrderTree) MinPrice() decimal.Decimal {
	if ot.depth > 0 {
		value, found := ot.priceTree.GetMin()
		if found {
			return value.(*OrderQueue).price
		} else {
			return decimal.Zero
		}
	} else {
		return decimal.Zero
	}
}

func (ot *OrderTree) getOrCreatePrice(price decimal.Decimal) *OrderQueue {
	priceStr := price.String()

	if queue, ok := ot.prices[priceStr]; ok {
		return queue
	}

	ot.depth++
	newQueue := NewOrderQueue(price)
	ot.priceTree.Put(price, newQueue)
	ot.prices[priceStr] = newQueue
	return newQueue
}

func (ot *OrderTree) String() string {
	sb := strings.Builder{}
	prices := []decimal.Decimal{}
	for k, _ := range ot.prices {
		num, _ := decimal.NewFromString(k)
		prices = append(prices, num)
	}

	sort.Slice(prices, func(i, j int) bool {
		return prices[i].LessThan(prices[j])
	})

	pricesStr := []string{}
	for _, price := range prices {
		pricesStr = append(pricesStr, price.String())
	}

	maxLen := 0
	for _, price := range pricesStr {
		if len(price) > maxLen {
			maxLen = len(price)
		}
	}

	for _, price := range pricesStr {
		sb.WriteString(fmt.Sprintf("\n%s -> %s", strings.Repeat(" ", maxLen-len(price))+price, ot.prices[price].volume))
	}

	return sb.String()
}
