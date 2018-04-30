package core

import (
	"fmt"
	"strings"
	"sync"
)

// InMemoryStore is Concurrent friendly Store
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

// Shutdown -Not much to do since its inmemory
func (s *InMemoryStore) Shutdown() error {
	return nil
}

// Update key value pair
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

// Query for key, return value would be a list of keys associated with the property
func (s *InMemoryStore) Query(key string) ([]string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	keyList, ok := s.store[key]

	if !ok {
		// property may not, thats ok since we just have to return empty
		// but at same time we should stop continuing the search since the intersection would be empty
		return nil, fmt.Errorf("Error querying property %s", key)
	}

	return keyList, nil
}

// Serialize store to backup, could be optionally compressed
func (s *InMemoryStore) Serialize() (map[string][]string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.store, nil
}
