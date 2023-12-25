package main

import (
	"fmt"
	"testing"
)

func TestOrderbook(t *testing.T) {
}

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrderA := newOrder(true, 5)
	buyOrderB := newOrder(true, 8)
	buyOrderC := newOrder(true, 10)
	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)

	l.deleteOrder(buyOrderB)

	fmt.Println(l)

}
