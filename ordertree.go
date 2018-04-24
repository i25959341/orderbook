package orderbook

import (
	"strconv"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/overlord/go-version/redblacktreeextended"
	"github.com/shopspring/decimal"
)

type Comparator func(a, b interface{}) int

func DecimalComparator(a, b interface{}) int {
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
	price_tree *redblacktreeextended.RedBlackTreeExtended
	price_map  map[decimal.Decimal]*OrderList // Dictionary containing price : OrderList object
	order_map  map[string]*Order              // Dictionary containing order_id : Order object
	volume     decimal.Decimal                // Contains total quantity from all Orders in tree
	num_orders int                            // Contains count of Orders in tree
	depth      int                            // Number of different prices in tree (http://en.wikipedia.org/wiki/Order_book_(trading)#Book_depth)
}

func NewOrderTree() *OrderTree {
	price_tree := &redblacktreeextended.RedBlackTreeExtended{rbt.NewWith(DecimalComparator)}
	price_map := make(map[decimal.Decimal]*OrderList)
	order_map := make(map[string]*Order)
	dec, _ := decimal.NewFromString("0.0")
	return &OrderTree{price_tree, price_map, order_map, dec, 0, 0}
}

func (ordertree *OrderTree) Length() int {
	return len(ordertree.order_map)
}

func (ordertree *OrderTree) Order(order_id string) *Order {
	return ordertree.order_map[order_id]
}

func (ordertree *OrderTree) PriceList(price decimal.Decimal) *OrderList {
	return ordertree.price_map[price]
}

func (ordertree *OrderTree) CreatePrice(price decimal.Decimal) {
	ordertree.depth = ordertree.depth + 1
	new_list := NewOrderList(price)
	ordertree.price_tree.Put(price, new_list)
	ordertree.price_map[price] = new_list
}

func (ordertree *OrderTree) RemovePrice(price decimal.Decimal) {
	ordertree.depth = ordertree.depth - 1
	ordertree.price_tree.Remove(price)
	delete(ordertree.price_map, price)
}

func (ordertree *OrderTree) PriceExist(price decimal.Decimal) bool {
	if _, ok := ordertree.price_map[price]; ok {
		return true
	}
	return false
}

func (ordertree *OrderTree) OrderExist(order_id string) bool {
	if _, ok := ordertree.order_map[order_id]; ok {
		return true
	}
	return false
}

func (ordertree *OrderTree) RemoveOrderById(order_id string) {
	ordertree.num_orders = ordertree.num_orders - 1
	order := ordertree.order_map[order_id]
	ordertree.volume = ordertree.volume.Sub(order.quantity)
	order.order_list.RemoveOrder(order)
	if order.order_list.Length() == 0 {
		ordertree.RemovePrice(order.price)
	}
	delete(ordertree.order_map, order_id)
}

func (ordertree *OrderTree) MaxPrice() (value interface{}, found bool) {
	if ordertree.depth > 0 {
		value, found := ordertree.price_tree.GetMax()
		return value, found
	} else {
		return nil, false
	}
}

func (ordertree *OrderTree) MinPrice() (value interface{}, found bool) {
	if ordertree.depth > 0 {
		value, found := ordertree.price_tree.GetMin()
		return value, found
	} else {
		return nil, false
	}
}

func (ordertree *OrderTree) MaxPriceList() *OrderList {
	if ordertree.depth > 0 {
		price, err := ordertree.MaxPrice()
		if err {
			return nil
		}
		return ordertree.price_map[price.(decimal.Decimal)]
	} else {
		return nil
	}
}

func (ordertree *OrderTree) MinPriceList() *OrderList {
	if ordertree.depth > 0 {
		price, err := ordertree.MinPrice()
		if err {
			return nil
		}
		return ordertree.price_map[price.(decimal.Decimal)]
	}
	return nil
}

func (ordertree *OrderTree) InsertOrder(quote map[string]string) {
	order_id := quote["order_id"]
	if ordertree.OrderExist(order_id) {
		ordertree.RemoveOrderById(order_id)
	}
	ordertree.num_orders++

	price, _ := decimal.NewFromString(quote["price"])

	if !ordertree.PriceExist(price) {
		ordertree.CreatePrice(price)
	}

	order := NewOrder(quote, ordertree.price_map[price])
	ordertree.price_map[price].AppendOrder(order)
	ordertree.order_map[order.order_id] = order
	ordertree.volume = ordertree.volume.Add(order.quantity)
}

func (ordertree *OrderTree) UpdateOrder(quote map[string]string) {
	order := ordertree.order_map[quote["order_id"]]
	original_quantity := order.quantity
	price, _ := decimal.NewFromString(quote["price"])

	if !price.Equal(order.price) {
		// Price changed. Remove order and update tree.
		order_list := ordertree.price_map[order.price]
		order_list.RemoveOrder(order)
		if order_list.Length() == 0 {
			ordertree.RemovePrice(price)
		}
		ordertree.InsertOrder(quote)
	} else {
		quantity, _ := decimal.NewFromString(quote["quantity"])
		timestamp, _ := strconv.Atoi(quote["timestamp"])
		order.UpdateQuantity(quantity, timestamp)
	}
	ordertree.volume = ordertree.volume.Add(order.quantity.Sub(original_quantity))
}
