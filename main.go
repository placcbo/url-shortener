package main

import (
	"sync"
)

type Link struct {
	Code   string
	URL    string
	Clicks int
}

type LinkStore struct {
	links map[string]*Link
	mu    sync.RWMutex
}

func NewLinkStore() *LinkStore {
	return &LinkStore{
		links: map[string]*Link{},
	}
}

func (s *LinkStore) Create(code, url string) *Link {
	s.mu.Lock()
	defer s.mu.Unlock()
	l := &Link{
		Code: code,
		URL:  url,
	}
	s.links[code] = l
	return l
}

func (s *LinkStore) Get(code string) (*Link, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.links[code]
	if !ok {
		return nil, false
	}

	return l, true

}

func (s *LinkStore) IncrementClicks(code string) (*Link, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	l, ok := s.links[code]
	if !ok {
		return nil, false
	}
	l.Clicks++
	return l, true
}

func main() {

}
