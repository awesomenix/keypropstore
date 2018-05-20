package core

import (
	"encoding/json"
	"strings"

	"github.com/dgraph-io/badger"
)

// Store (Key, Value), value in JSON format

// BadgerStore store for db
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
		opts.Dir = "./badgerdb"
		opts.ValueDir = "./badgerdb"
	}

	db, err := badger.Open(opts)
	s.db = db
	return err
}

// Shutdown db, by closing all the open handles
func (s *BadgerStore) Shutdown() error {
	return s.db.Close()
}

// Update db with key value pair
func (s *BadgerStore) Update(key, value string) error {
	var jsStoreValue []byte
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

	var storeList []string

	// Deserialize JSON value to string array
	if jsStoreValue != nil {
		if err := json.Unmarshal(jsStoreValue, &storeList); err != nil {
			return err
		}
	}

	storeValue := make(map[string]bool)

	for _, value := range storeList {
		storeValue[value] = true
	}

	value = strings.ToLower(value)

	// append the current value to JSON
	if _, ok := storeValue[value]; !ok {
		storeList = append(storeList, value)
		jsValue, err := json.Marshal(storeList)
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

// Query for key, return value would be a list of keys associated with the property
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

// Serialize store to backup, could be optionally compressed
func (s *BadgerStore) Serialize() (map[string][]string, error) {
	store := make(map[string][]string)
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()
			jsStoreValue, err := item.Value()
			if err != nil {
				return err
			}
			// Deserialize the JSON value to string array
			var keyList []string
			if err := json.Unmarshal(jsStoreValue, &keyList); err != nil {
				return err
			}
			store[string(key)] = keyList
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return store, nil
}
