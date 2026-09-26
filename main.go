package main

import (
	"fmt"
	"sync"
)

type Link struct {
	Code   string
	URL    string
	Clicks int
}
type LinkStore struct {
	mu    sync.Mutex
	links map[int]*Link
}

func NewLinkStore() *LinkStore {
	return &LinkStore{
		links: map[int]*Link{},
	}
}
func main() {

	store := NewLinkStore()

	fmt.Println(store)
}
