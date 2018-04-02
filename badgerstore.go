package keypropstore

import (
	"encoding/json"
	"github.com/dgraph-io/badger"
	"strings"
)

// Store (Key, Value), value in JSON format

type BadgerStore struct {
	db *badger.DB
}

// Initialize Store with custom configuration
func (s *BadgerStore) Initialize(cfg Config) error {
	var opts badger.Options
	if cfg != nil {
		opts = cfg.(badger.Options)
	} else {
		opts = badger.DefaultOptions
	}

	db, err := badger.Open(opts)
	s.db = db
	return err
}

func (s *BadgerStore) Shutdown() error {
	return s.db.Close()
}

func (s *BadgerStore) Update(key, value string) error {
	var jsStoreValue []byte = nil
	// Check if the key exists in db
	// if its does return the JSON value
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		byteVal, err := item.Value()
		if err != nil {
			return err
		}
		jsStoreValue = byteVal
		return nil
	})

	// if there was an error other than key not found, bubble up the error
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	var storeValue []string

	// Deserialize JSON value to string array
	if jsStoreValue != nil {
		if err := json.Unmarshal(jsStoreValue, &storeValue); err != nil {
			return err
		}
	}

	value = strings.ToLower(value)

	// append the current value to JSON
	if !ContainsKey(value, storeValue) {
		storeValue = append(storeValue, value)
		jsValue, err := json.Marshal(storeValue)
		if err != nil {
			return err
		}

		// Update DB
		err = s.db.Update(func(txn *badger.Txn) error {
			err := txn.Set([]byte(key), jsValue)
			return err
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *BadgerStore) Query(key string) ([]string, error) {
	var jsStoreValue []byte
	// Get the JSON value for this key
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		byteVal, err := item.Value()
		if err != nil {
			return err
		}
		jsStoreValue = byteVal
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Deserialize the JSON value to string array
	var keyList []string
	if err := json.Unmarshal(jsStoreValue, &keyList); err != nil {
		return nil, err
	}

	return keyList, nil
}
