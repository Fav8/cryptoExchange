package main

import (
	"fmt"
	"sort"
	"time"
)

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}

type Order struct {
	Bid       bool
	Size      float64
	Limit     *Limit
	Timestamp int64
}

type Orders []*Order

func (o Orders) Len() int           { return len(o) }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].Timestamp < o[j].Timestamp }

func newOrder(bid bool, size float64) *Order {
	return &Order{
		Bid:       bid,
		Size:      size,
		Timestamp: time.Now().UnixNano(),
	}
}

func (o *Order) EditOrder(size float64) {
	sizeDiff := size - o.Size
	o.Size = size
	o.Limit.TotalVolume += sizeDiff
}

func (o *Order) String() string {
	return fmt.Sprintf("[Size: %.2f, Time: %v]", o.Size, o.Timestamp)
}

type Limit struct {
	Price       float64
	Orders      Orders
	TotalVolume float64
}

type Limits []*Limit

type ByBestAsk struct{ Limits }

func (a ByBestAsk) Len() int           { return len(a.Limits) }
func (a ByBestAsk) Swap(i, j int)      { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }
func (a ByBestAsk) Less(i, j int) bool { return a.Limits[i].Price < a.Limits[j].Price }

type ByBestBid struct{ Limits }

func (a ByBestBid) Len() int           { return len(a.Limits) }
func (a ByBestBid) Swap(i, j int)      { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }
func (a ByBestBid) Less(i, j int) bool { return a.Limits[i].Price > a.Limits[j].Price }

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

func (l *Limit) DeleteOrder(o *Order) {
	for i, order := range l.Orders {
		if order == o {
			l.Orders = append(l.Orders[:i], l.Orders[i+1:]...)
			l.TotalVolume -= o.Size
			return
		}
	}
	o.Limit = nil
	l.TotalVolume -= o.Size

	sort.Sort(l.Orders)
}

type Orderbook struct {
	Asks []*Limit
	Bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		Asks:      []*Limit{},
		Bids:      []*Limit{},
		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

func (ob *Orderbook) PlaceOrder(price float64, o *Order) []Match {
	if o.Bid {
		if limit, ok := ob.AskLimits[price]; ok {
			for orderIndex := 0; orderIndex < len(limit.Orders); orderIndex++ {
				if remainerSize := limit.Orders[orderIndex].Size - o.Size; remainerSize > 0 {
					limit.Orders[orderIndex].Size = remainerSize
					limit.TotalVolume -= o.Size
					o.Size = 0
				} else if remainerSize := limit.Orders[orderIndex].Size - o.Size; remainerSize < 0 {
					limit.DeleteOrder(limit.Orders[orderIndex])
					//check if i need to move on to other orders or add a pending order
				}
			}

		}
	} else {
		if limit, ok := ob.BidLimits[price]; ok {
			for orderIndex := 0; orderIndex < len(limit.Orders) && o.Size > 0; orderIndex++ {
				if remainerSize := limit.Orders[orderIndex].Size - o.Size; remainerSize > 0 {
					limit.Orders[orderIndex].Size = remainerSize
					limit.TotalVolume = limit.TotalVolume - o.Size
					o.Size = 0
				} else if remainerSize := limit.Orders[orderIndex].Size - o.Size; remainerSize < 0 {
					limit.DeleteOrder(limit.Orders[orderIndex])
					//check if i need to move on to other orders or add a pending order
				}
			}

		}
	}
	if o.Size > 0.0 {
		ob.add(price, o)
	}
	return []Match{}
}

func (ob *Orderbook) add(price float64, o *Order) {
	if o.Bid {
		if limit, ok := ob.BidLimits[price]; ok {
			limit.AddOrder(o)
		} else {
			limit := NewLimit(price)
			limit.AddOrder(o)
			ob.BidLimits[price] = limit
			ob.Bids = append(ob.Bids, limit)
		}
	} else {
		if limit, ok := ob.AskLimits[price]; ok {
			limit.AddOrder(o)
		} else {
			limit := NewLimit(price)
			limit.AddOrder(o)
			ob.AskLimits[price] = limit
			ob.Asks = append(ob.Asks, limit)
		}
	}
}
