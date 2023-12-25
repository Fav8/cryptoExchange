package main

import (
	"fmt"
	"time"
)

type Order struct {
	Bid       bool
	Size      float64
	Limit     *Limit
	Timestamp int64
}

func newOrder(bid bool, size float64) *Order {
	return &Order{
		Bid:       bid,
		Size:      size,
		Timestamp: time.Now().UnixNano(),
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("[Size: %.2f]", o.Size)
}

type Limit struct {
	Price       float64
	Orders      []*Order
	TotalVolume float64
}

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

func (l *Limit) deleteOrder(o *Order) {
	for i, order := range l.Orders {
		if order == o {
			l.Orders = append(l.Orders[:i], l.Orders[i+1:]...)
			l.TotalVolume -= o.Size
			return
		}
	}
	o.Limit = nil
	l.TotalVolume -= o.Size
}

type Orderbook struct {
	Asks []*Limit
	Bids []*Limit
}
