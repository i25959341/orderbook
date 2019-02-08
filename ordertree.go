package orderbook

import (
	"time"

	"github.com/emirpasic/gods/examples/redblacktreeextended"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/shopspring/decimal"
)

type OrderTree struct {
	priceTree *redblacktreeextended.RedBlackTreeExtended
	prices    map[string]*OrderQueue // Dictionary containing price : OrderList object
	orders    map[string]*Order      // Dictionary containing order_id : Order object
	volume    decimal.Decimal        // Contains total quantity from all Orders in tree
	numOrders int                    // Contains count of Orders in tree
	depth     int                    // Number of different prices in tree (http://en.wikipedia.org/wiki/Order_book_(trading)#Book_depth)
}

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

func (ot *OrderTree) Length() int {
	return len(ot.orders)
}

func (ot *OrderTree) Order(orderID string) *Order {
	return ot.orders[orderID]
}

func (ot *OrderTree) PriceQueue(price decimal.Decimal) *OrderQueue {
	return ot.prices[price.String()]
}

func (ot *OrderTree) CreatePrice(price decimal.Decimal) {
	ot.depth = ot.depth + 1
	newList := NewOrderQueue(price)
	ot.priceTree.Put(price, newList)
	ot.prices[price.String()] = newList
}

func (ot *OrderTree) RemovePrice(price decimal.Decimal) {
	ot.depth = ot.depth - 1
	ot.priceTree.Remove(price)
	delete(ot.prices, price.String())
}

func (ot *OrderTree) PriceExist(price decimal.Decimal) bool {
	if _, ok := ot.prices[price.String()]; ok {
		return true
	}
	return false
}

func (ot *OrderTree) OrderExist(orderID string) bool {
	if _, ok := ot.orders[orderID]; ok {
		return true
	}
	return false
}

func (ot *OrderTree) RemoveOrder(orderID string) {
	ot.numOrders = ot.numOrders - 1
	order := ot.orders[orderID]
	ot.volume = ot.volume.Sub(order.quantity)
	order.orderList.Remove(order)
	if order.orderList.Length() == 0 {
		ot.RemovePrice(order.price)
	}
	delete(ot.orders, orderID)
}

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

func (ot *OrderTree) MaxPriceQueue() *OrderQueue {
	if ot.depth > 0 {
		price := ot.MaxPrice()
		return ot.prices[price.String()]
	}
	return nil
}

func (ot *OrderTree) MinPriceQueue() *OrderQueue {
	if ot.depth > 0 {
		price := ot.MinPrice()
		return ot.prices[price.String()]
	}
	return nil
}

func (ot *OrderTree) InsertOrderFromMap(quote map[string]string) {
	orderID := quote["order_id"]

	if ot.OrderExist(orderID) {
		ot.RemoveOrder(orderID)
	}
	ot.numOrders++

	price, _ := decimal.NewFromString(quote["price"])

	if !ot.PriceExist(price) {
		ot.CreatePrice(price)
	}

	order := NewOrderFromMap(quote, ot.prices[price.String()])
	ot.prices[price.String()].Append(order)
	ot.orders[order.orderID] = order
	ot.volume = ot.volume.Add(order.quantity)
}

func (ot *OrderTree) InsertOrder(orderID string, quantity, price decimal.Decimal, timestamp time.Time) {
	if ot.OrderExist(orderID) {
		ot.RemoveOrder(orderID)
	}
	ot.numOrders++

	if !ot.PriceExist(price) {
		ot.CreatePrice(price)
	}

	priceStr := price.String()
	order := NewOrder(ot.prices[priceStr], orderID, quantity, price, timestamp)

	ot.prices[priceStr].Append(order)
	ot.orders[order.orderID] = order
	ot.volume = ot.volume.Add(order.quantity)
}
