package main

import (
	"crypto/rand"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Link struct {
	Code   string `json:"code"`
	URL    string `json:"url"`
	Clicks int    `json:"clicks"`
}

type LinkStore struct {
	mu    sync.RWMutex
	links map[string]*Link
}

func NewLinkStore() *LinkStore {
	return &LinkStore{
		links: map[string]*Link{},
	}
}

func (s *LinkStore) Create(code, url string) (*Link, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.links[code]; exists {
		return nil, false
	}

	l := &Link{
		Code: code,
		URL:  url,
	}

	s.links[code] = l

	return l, true
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

func generateCode(length int) string {
	b := make([]byte, length)

	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}

	return string(b)
}

func main() {
	store := NewLinkStore()

	r := chi.NewRouter()

	// Create a shortened link
	r.Post("/links", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			URL string `json:"url"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if input.URL == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		var code string

		for {
			code = generateCode(6)

			if code == "" {
				http.Error(w, "Failed to generate code", http.StatusInternalServerError)
				return
			}

			if _, exists := store.Get(code); !exists {
				break
			}
		}

		link, ok := store.Create(code, input.URL)
		if !ok {
			http.Error(w, "Code already exists", http.StatusConflict)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(link)
	})

	// Redirect to the original URL
	r.Get("/link/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")

		link, ok := store.IncrementClicks(code)
		if !ok {
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, link.URL, http.StatusFound)
	})

	// Get link information and statistics
	r.Get("/links/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")

		link, ok := store.Get(code)
		if !ok {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(link)
	})

	http.ListenAndServe(":8080", r)
}