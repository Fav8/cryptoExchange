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

func (o *Order) isFilled() bool {
	return o.Size == 0
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

func (l *Limit) Fill(o *Order) []Match {
	matches := []Match{}
	for _, order := range l.Orders {
		match := l.fillOrder(order, o)
		matches = append(matches, match)
		if o.isFilled() {
			break
		}
	}
	return matches
}

func (l *Limit) fillOrder(a, b *Order) Match {
	var (
		bid        *Order
		ask        *Order
		sizeFilled float64
	)
	if a.Bid {
		bid = a
		ask = b
	} else {
		bid = b
		ask = a
	}
	if a.Size >= b.Size {
		sizeFilled = b.Size
		a.Size -= b.Size
		b.Size = 0
	} else {
		sizeFilled = a.Size
		b.Size -= a.Size
		a.Size = 0
	}
	return Match{
		Ask:        ask,
		Bid:        bid,
		SizeFilled: sizeFilled,
		Price:      l.Price,
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
	asks []*Limit
	bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		asks:      []*Limit{},
		bids:      []*Limit{},
		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

func (ob *Orderbook) PlaceMarketOrder(o *Order) []Match {
	matches := []Match{}
	if o.Bid {
		if o.Size > ob.AskTotalVolume() {
			panic(fmt.Errorf("not enough volume  [order size: %.2f, ask volume: %.2f]", o.Size, ob.AskTotalVolume()))
		}
		for _, limit := range ob.Asks() {
			limitMatches := limit.Fill(o)
			matches = append(matches, limitMatches...)
		}
	} else {
		if o.Size > ob.BidTotalVolume() {
			panic("not enough volume")
		}
		for _, limit := range ob.Bids() {
			limitMatches := limit.Fill(o)
			matches = append(matches, limitMatches...)
		}
	}

	return matches
}

func (ob *Orderbook) PlaceLimitOrder(price float64, o *Order) []Match {
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
			ob.bids = append(ob.bids, limit)
		}
	} else {
		if limit, ok := ob.AskLimits[price]; ok {
			limit.AddOrder(o)
		} else {
			limit := NewLimit(price)
			limit.AddOrder(o)
			ob.AskLimits[price] = limit
			ob.asks = append(ob.asks, limit)
		}
	}
}

func (ob *Orderbook) BidTotalVolume() float64 {
	var total float64
	for _, limit := range ob.bids {
		total += limit.TotalVolume
	}
	return total
}

func (ob *Orderbook) AskTotalVolume() float64 {
	var total float64
	for _, limit := range ob.asks {
		total += limit.TotalVolume
	}
	return total
}

func (ob *Orderbook) Asks() []*Limit {
	sort.Sort(ByBestAsk{ob.asks})
	return ob.asks
}

func (ob *Orderbook) Bids() []*Limit {
	sort.Sort(ByBestBid{ob.bids})
	return ob.bids
}
