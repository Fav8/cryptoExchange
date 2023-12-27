package main

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%v != %v", a, b)
	}
}

func TestPlaceLimitOrder(t *testing.T) {
	ob := NewOrderbook()
	sellOrderA := newOrder(false, 5)
	ob.PlaceLimitOrder(10_000, sellOrderA)
	assert(t, len(ob.asks), 1)
	sellOrderB := newOrder(false, 7)
	ob.PlaceLimitOrder(9_000, sellOrderB)
	assert(t, len(ob.asks), 2)

}

func TestPlaceMarketOrder(t *testing.T) {
	ob := NewOrderbook()
	sellOrderA := newOrder(false, 5)
	ob.PlaceLimitOrder(20_000, sellOrderA)
	buyOrderA := newOrder(true, 2)
	matches := ob.PlaceMarketOrder(buyOrderA)
	assert(t, len(ob.bids), 0)
	assert(t, ob.asks[0].TotalVolume, 3.0)

	fmt.Printf("%+v", matches)
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

func TestPlaceMarketOrderMultiFill(t *testing.T) {
	ob := NewOrderbook()

	buyOrderA := newOrder(true, 5)
	buyOrderB := newOrder(true, 8)
	buyOrderC := newOrder(true, 10)
	buyOrderD := newOrder(true, 1)

	ob.PlaceLimitOrder(10_000, buyOrderA)
	ob.PlaceLimitOrder(9_000, buyOrderB)
	ob.PlaceLimitOrder(5_000, buyOrderC)
	ob.PlaceLimitOrder(5_000, buyOrderD)

	assert(t, ob.BidTotalVolume(), 24.0)

	sellOrderA := newOrder(false, 20)
	matches := ob.PlaceMarketOrder(sellOrderA)
	assert(t, ob.BidTotalVolume(), 4.0)
	fmt.Printf("%+v", matches)
	assert(t, len(matches), 3)
	assert(t, len(ob.bids), 1)
}

func TestCancelOrder(t *testing.T) {
	ob := NewOrderbook()

	buyOrderA := newOrder(true, 5)

	ob.PlaceLimitOrder(10_000, buyOrderA)

	assert(t, len(ob.bids), 1)
	assert(t, ob.BidTotalVolume(), 5.0)
	ob.CancelOrder(buyOrderA)
	assert(t, len(ob.bids), 0)
	assert(t, ob.BidTotalVolume(), 0.0)
}
