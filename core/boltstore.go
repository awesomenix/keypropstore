package core

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/boltdb/bolt"
)

const boltBucket string = "keypropstore"

// BoltStore (Key, Value), value in JSON format
type BoltStore struct {
	db *bolt.DB
}

// BoltStoreConfig configuration for path, filemode and options
type BoltStoreConfig struct {
	path    string
	mode    os.FileMode
	options *bolt.Options
}

// Initialize Store with custom configuration
func (s *BoltStore) Initialize(cfg Config) error {
	var opts *BoltStoreConfig
	if cfg != nil {
		opts = cfg.(*BoltStoreConfig)
	} else {
		opts = &BoltStoreConfig{"./boltdb", 600, nil}
	}

	db, err := bolt.Open(opts.path, opts.mode, opts.options)
	s.db = db
	return err
}

func (s *BoltStore) Shutdown() error {
	return s.db.Close()
}

func (s *BoltStore) Update(key, value string) error {
	var jsStoreValue []byte = nil

	s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(boltBucket))
		return err
	})

	// Check if the key exists in db
	// if its does return the JSON value
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(boltBucket))
		byteVal := bucket.Get([]byte(key))
		jsStoreValue = byteVal
		return nil
	})

	// if there was an error other than key not found, bubble up the error
	if err != nil {
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
		err = s.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(boltBucket))
			err := b.Put([]byte(key), jsValue)
			return err
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *BoltStore) Query(key string) ([]string, error) {
	var jsStoreValue []byte
	// Get the JSON value for this key
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(boltBucket))
		byteVal := bucket.Get([]byte(key))
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

func (s *BoltStore) Serialize() (map[string][]string, error) {
	store := make(map[string][]string)
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte([]byte(boltBucket)))

		err := bucket.ForEach(func(key, jsStoreValue []byte) error {
			// Deserialize the JSON value to string array
			var keyList []string
			if err := json.Unmarshal(jsStoreValue, &keyList); err != nil {
				return err
			}
			store[string(key)] = keyList
			return nil
		})

		return err
	})

	if err != nil {
		return nil, err
	}

	return store, nil
}