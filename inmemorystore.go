package keypropstore

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// Concurrent friendly Store
// Map of Property and List of Keys associated with that property
type InMemoryStore struct {
	store map[string][]string
	lock  sync.RWMutex
}

// Initialize Store with custom configuration
func (s *InMemoryStore) Initialize(cfg Config) error {
	s.store = make(map[string][]string)
	return nil
}

// Not much to do while shutdown
func (s *InMemoryStore) Shutdown() error {
	return nil
}

func (s *InMemoryStore) Update(key, value string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	storeValue, ok := s.store[key]
	if !ok {
		s.store[key] = make([]string, 1)
	}

	value = strings.ToLower(value)

	if !ContainsKey(value, storeValue) {
		s.store[key] = append(storeValue, value)
	}

	return nil
}

func (s *InMemoryStore) Query(key string) ([]string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	keyList, ok := s.store[key]

	if !ok {
		// property may not, thats ok since we just have to return empty
		// but at same time we should stop continuing the search since the intersection would be empty
		return nil, errors.New(fmt.Sprintf("Error querying property %s", key))
	}

	return keyList, nil
}
