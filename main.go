package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
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
	store := NewLinkStore()

	r := chi.NewRouter()
	r.Post("/links", func(w http.ResponseWriter, r *http.Request) {
		var l Link

		err := json.NewDecoder(r.Body).Decode(&l)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		link := store.Create(l.Code, l.URL)

		json.NewEncoder(w).Encode(link)
	})

	// get code

	r.Get("/link/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		link, ok := store.IncrementClicks(code)
		if !ok {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, link.URL, http.StatusNotFound)

	})
	http.ListenAndServe(":8080", r)
}
