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
	buyOrderB := newOrder(true, 10)
	matches := ob.PlaceMarketOrder(buyOrderA)
	ob.PlaceMarketOrder(buyOrderB)
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
