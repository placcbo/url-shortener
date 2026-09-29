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
	links map[string]*Link
}

func NewLinkStore() *LinkStore {
	return &LinkStore{
		links: map[string]*Link{},
	}
}

func (s *LinkStore) Create(code, url string) *Link {
	l := &Link{
		Code: code,
		URL:  url,
	}
	s.links[code] = l
	return l
}

func main() {
	store := NewLinkStore()
	store.Create("abc", "www.jumia.com/computers/windows/laptop/hp/hp15s")

	fmt.Println(store)
}
