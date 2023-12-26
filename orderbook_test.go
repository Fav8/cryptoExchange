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
	assert(t, len(matches), 1)
	assert(t, len(ob.asks), 1)
	assert(t, len(ob.bids), 0)
	assert(t, ob.asks[0].TotalVolume, 3.0)
	assert(t, matches[0].Price, 20_000.0)
	assert(t, matches[0].SizeFilled, 2.0)
	assert(t, matches[0].Ask, sellOrderA)
	assert(t, matches[0].Bid, buyOrderA)
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
	sellOrderA := newOrder(false, 5)
	sellOrderB := newOrder(false, 5)
	ob.PlaceLimitOrder(20_000, sellOrderA)
	ob.PlaceLimitOrder(19_000, sellOrderB)
	buyOrderA := newOrder(true, 6)
	matches := ob.PlaceMarketOrder(buyOrderA)
	assert(t, len(matches), 2)
	assert(t, len(ob.asks), 1)
	assert(t, len(ob.bids), 0)
	assert(t, ob.asks[0].TotalVolume, 4.0)
	assert(t, matches[0].Price, 20_000.0)
	assert(t, matches[0].SizeFilled, 5.0)
	assert(t, matches[0].Ask, sellOrderA)
	assert(t, matches[0].Bid, buyOrderA)
	assert(t, matches[1].Price, 19_000.0)
	assert(t, matches[1].SizeFilled, 1.0)
	assert(t, matches[1].Ask, sellOrderB)
	assert(t, matches[1].Bid, buyOrderA)

	fmt.Printf("%+v", matches)
}
