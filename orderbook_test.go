package main

import (
	"fmt"
	"testing"
)

func TestOrderbook(t *testing.T) {
	ob := NewOrderbook()

	buyOrderA := newOrder(true, 10)
	buyOrderB := newOrder(true, 10)
	sellOrderA := newOrder(false, 5)
	ob.PlaceOrder(18_000, buyOrderA)
	ob.PlaceOrder(18_000, buyOrderB)
	ob.PlaceOrder(18_000, sellOrderA)
	fmt.Println(ob.Bids)
}

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrderA := newOrder(true, 5)
	buyOrderB := newOrder(true, 8)
	buyOrderC := newOrder(true, 10)
	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)

	l.DeleteOrder(buyOrderB)
	buyOrderA.EditOrder(10)

	fmt.Println(l)

}
