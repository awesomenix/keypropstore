package core

import (
	"fmt"
	"strings"
	"sync"
)

// InMemoryStore is Concurrent friendly Store
// Map of Property and List of Keys associated with that property
type InMemoryStore struct {
	store map[string]map[string]bool
	lock  sync.RWMutex
}

// Initialize Store with custom configuration
func (s *InMemoryStore) Initialize(cfg Config) error {
	s.store = make(map[string]map[string]bool)
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
	key = strings.ToLower(key)

	if _, ok := s.store[key]; !ok {
		s.store[key] = make(map[string]bool)
		s.store[key][strings.ToLower(value)] = true
		return nil
	}

	s.store[key][strings.ToLower(value)] = true

	return nil
}

// Query for key, return value would be a list of keys associated with the property
func (s *InMemoryStore) Query(key string) ([]string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	key = strings.ToLower(key)
	keySet, ok := s.store[key]

	if !ok {
		// property may not, thats ok since we just have to return empty
		// but at same time we should stop continuing the search since the intersection would be empty
		return nil, fmt.Errorf("Error querying property %s", key)
	}

	keyList := make([]string, 0)

	for key := range keySet {
		keyList = append(keyList, key)
	}

	return keyList, nil
}

// Serialize store to backup, could be optionally compressed
func (s *InMemoryStore) Serialize() (map[string][]string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	store := make(map[string][]string)

	for key, keySet := range s.store {
		keyList := make([]string, 0)
		for key := range keySet {
			keyList = append(keyList, key)
		}
		store[key] = keyList
	}

	return store, nil
}
