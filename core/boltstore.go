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
	Path    string
	Mode    os.FileMode
	Options *bolt.Options
}

// Initialize Store with custom configuration
func (s *BoltStore) Initialize(cfg Config) error {
	var opts *BoltStoreConfig
	if cfg != nil {
		opts = cfg.(*BoltStoreConfig)
	} else {
		opts = &BoltStoreConfig{"./boltdb", 600, nil}
	}

	var err error
	s.db, err = bolt.Open(opts.Path, opts.Mode, opts.Options)
	return err
}

// Shutdown db, by closing all the open handles
func (s *BoltStore) Shutdown() error {
	return s.db.Close()
}

// Update db with key value pair
func (s *BoltStore) Update(key, value string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(boltBucket))
		return err
	})

	// if there was an error creating the bucket
	if err != nil {
		return err
	}

	var jsStoreValue []byte

	// Check if the key exists in db
	// if its does return the JSON value
	err = s.db.View(func(tx *bolt.Tx) error {
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

// Query for key, return value would be a list of keys associated with the property
func (s *BoltStore) Query(key string) ([]string, error) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(boltBucket))
		return err
	})

	if err != nil {
		return nil, err
	}

	var jsStoreValue []byte

	// Get the JSON value for this key
	err = s.db.View(func(tx *bolt.Tx) error {
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

// Serialize store to backup, could be optionally compressed
func (s *BoltStore) Serialize() (map[string][]string, error) {

	err := s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(boltBucket))
		return err
	})

	if err != nil {
		return nil, err
	}

	store := make(map[string][]string)

	err = s.db.View(func(tx *bolt.Tx) error {
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
