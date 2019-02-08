package orderbook

import (
	"strconv"

	"github.com/emirpasic/gods/examples/redblacktreeextended"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/shopspring/decimal"
)

type Comparator func(a, b interface{}) int

func decimalComparator(a, b interface{}) int {
	aAsserted := a.(decimal.Decimal)
	bAsserted := b.(decimal.Decimal)
	switch {
	case aAsserted.GreaterThan(bAsserted):
		return 1
	case aAsserted.LessThan(bAsserted):
		return -1
	default:
		return 0
	}
}

type OrderTree struct {
	priceTree *redblacktreeextended.RedBlackTreeExtended
	priceMap  map[string]*OrderList // Dictionary containing price : OrderList object
	orderMap  map[string]*Order     // Dictionary containing order_id : Order object
	volume    decimal.Decimal       // Contains total quantity from all Orders in tree
	numOrders int                   // Contains count of Orders in tree
	depth     int                   // Number of different prices in tree (http://en.wikipedia.org/wiki/Order_book_(trading)#Book_depth)
}

func NewOrderTree() *OrderTree {
	priceTree := &redblacktreeextended.RedBlackTreeExtended{Tree: rbt.NewWith(decimalComparator)}
	priceMap := make(map[string]*OrderList)
	orderMap := make(map[string]*Order)
	return &OrderTree{priceTree, priceMap, orderMap, decimal.Zero, 0, 0}
}

func (ordertree *OrderTree) Length() int {
	return len(ordertree.orderMap)
}

func (ordertree *OrderTree) Order(orderID string) *Order {
	return ordertree.orderMap[orderID]
}

func (ordertree *OrderTree) PriceList(price decimal.Decimal) *OrderList {
	return ordertree.priceMap[price.String()]
}

func (ordertree *OrderTree) CreatePrice(price decimal.Decimal) {
	ordertree.depth = ordertree.depth + 1
	newList := NewOrderList(price)
	ordertree.priceTree.Put(price, newList)
	ordertree.priceMap[price.String()] = newList
}

func (ordertree *OrderTree) RemovePrice(price decimal.Decimal) {
	ordertree.depth = ordertree.depth - 1
	ordertree.priceTree.Remove(price)
	delete(ordertree.priceMap, price.String())
}

func (ordertree *OrderTree) PriceExist(price decimal.Decimal) bool {
	if _, ok := ordertree.priceMap[price.String()]; ok {
		return true
	}
	return false
}

func (ordertree *OrderTree) OrderExist(orderID string) bool {
	if _, ok := ordertree.orderMap[orderID]; ok {
		return true
	}
	return false
}

func (ordertree *OrderTree) RemoveOrderById(orderID string) {
	ordertree.numOrders = ordertree.numOrders - 1
	order := ordertree.orderMap[orderID]
	ordertree.volume = ordertree.volume.Sub(order.quantity)
	order.orderList.RemoveOrder(order)
	if order.orderList.Length() == 0 {
		ordertree.RemovePrice(order.price)
	}
	delete(ordertree.orderMap, orderID)
}

func (ordertree *OrderTree) MaxPrice() decimal.Decimal {
	if ordertree.depth > 0 {
		value, found := ordertree.priceTree.GetMax()
		if found {
			return value.(*OrderList).price
		}
		return decimal.Zero

	} else {
		return decimal.Zero
	}
}

func (ordertree *OrderTree) MinPrice() decimal.Decimal {
	if ordertree.depth > 0 {
		value, found := ordertree.priceTree.GetMin()
		if found {
			return value.(*OrderList).price
		} else {
			return decimal.Zero
		}

	} else {
		return decimal.Zero
	}
}

func (ordertree *OrderTree) MaxPriceList() *OrderList {
	if ordertree.depth > 0 {
		price := ordertree.MaxPrice()
		return ordertree.priceMap[price.String()]
	}
	return nil

}

func (ordertree *OrderTree) MinPriceList() *OrderList {
	if ordertree.depth > 0 {
		price := ordertree.MinPrice()
		return ordertree.priceMap[price.String()]
	}
	return nil
}

func (ordertree *OrderTree) InsertOrder(quote map[string]string) {
	orderID := quote["order_id"]

	if ordertree.OrderExist(orderID) {
		ordertree.RemoveOrderById(orderID)
	}
	ordertree.numOrders++

	price, _ := decimal.NewFromString(quote["price"])

	if !ordertree.PriceExist(price) {
		ordertree.CreatePrice(price)
	}

	order := NewOrder(quote, ordertree.priceMap[price.String()])
	ordertree.priceMap[price.String()].AppendOrder(order)
	ordertree.orderMap[order.orderID] = order
	ordertree.volume = ordertree.volume.Add(order.quantity)
}

func (ordertree *OrderTree) UpdateOrder(quote map[string]string) {
	order := ordertree.orderMap[quote["order_id"]]
	originalQuantity := order.quantity
	price, _ := decimal.NewFromString(quote["price"])

	if !price.Equal(order.price) {
		// Price changed. Remove order and update tree.
		orderList := ordertree.priceMap[order.price.String()]
		orderList.RemoveOrder(order)
		if orderList.Length() == 0 {
			ordertree.RemovePrice(price)
		}
		ordertree.InsertOrder(quote)
	} else {
		quantity, _ := decimal.NewFromString(quote["quantity"])
		timestamp, _ := strconv.Atoi(quote["timestamp"])
		order.UpdateQuantity(quantity, timestamp)
	}
	ordertree.volume = ordertree.volume.Add(order.quantity.Sub(originalQuantity))
}
