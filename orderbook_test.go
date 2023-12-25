package main

import (
	"fmt"
	"testing"
)

func TestOrderbook(t *testing.T) {
}

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrder := newOrder(true, 5)
	l.AddOrder(buyOrder)
	fmt.Println(l)
}
