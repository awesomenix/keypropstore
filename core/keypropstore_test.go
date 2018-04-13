package core

import (
	"testing"
)

type DummyEchoStore struct {
	store map[string]string
}

// Initialize Store with custom configuration
func (s *DummyEchoStore) Initialize(cfg Config) error {
	s.store = make(map[string]string)
	return nil
}

func (s *DummyEchoStore) Shutdown() error {
	return nil
}

func (s *DummyEchoStore) Update(key, value string) error {
	s.store[key] = value
	return nil
}

func (s *DummyEchoStore) Query(key string) ([]string, error) {
	ret := s.store[key]
	var res []string
	res = append(res, ret)
	return res, nil
}

func (s *DummyEchoStore) Serialize() (map[string][]string, error) {
	return nil, nil
}

func TestDummyEchoStore(t *testing.T) {
	dummyEchoStore := &DummyEchoStore{}
	InitializeStore(dummyEchoStore, nil)
	defer ShutdownStore(dummyEchoStore)
	err := UpdateStore(dummyEchoStore, byt)
	if err != nil {
		t.Error(err)
		return
	}

	query := []byte(`{"key1": "asdasdb"}`)
	expected := []byte(`["m4"]`)

	res, err := QueryStore(dummyEchoStore, query)

	if err != nil {
		t.Error(err)
		return
	}

	t.Log("Store returned", string(res), "Expect", string(expected))

	if err := CheckResults(res, expected); err != nil {
		t.Error(err)
		return
	}
}
