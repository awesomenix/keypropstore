package core

import (
	"encoding/json"
)

type Config interface {
}

type Store interface {
	Initialize(cfg Config) error
	Shutdown() error
	Update(key, value string) error
	Query(key string) ([]string, error)
	Serialize() (map[string][]string, error)
}

func InitializeStore(s Store, cfg Config) error {
	return s.Initialize(cfg)
}

func ShutdownStore(s Store) error {
	return s.Shutdown()
}

// Update Store called with list of Key and Associated properties
// Check for consistncy of input json
// Converts
// {"m1": {"num": "6.13","strs": "a","key1": "b"}, "m2": {"num": "6.13","key1": "bddd"}}
// To
// {"num:6.13" : ["m1", m2"], "strs:a" : ["m1"], "key1:b" : ["m1"], "key1:bddd" : ["m2"]}
func UpdateStore(s Store, byt []byte) error {
	var dat map[string]map[string]string

	if err := json.Unmarshal(byt, &dat); err != nil {
		return err
	}

	for key, value := range dat {
		for valkey, valval := range value {
			keyval := GenerateKey(valkey, valval)
			if err := s.Update(keyval, key); err != nil {
				return err
			}
		}
	}

	return nil
}

// Query Store with single/multiple properties
// Properties are AND only, if an OR is required, query multiple times
// OR could be supported, but keeping it simple for now
func QueryStore(s Store, jsQuery []byte) ([]byte, error) {
	keys := make([]string, 0)

	var query map[string]string

	if err := json.Unmarshal(jsQuery, &query); err != nil {
		return nil, err
	}

	q := make([]string, 0)

	for key, val := range query {
		q = append(q, GenerateKey(key, val))
	}

	for _, val := range q {
		keyList, err := s.Query(val)

		if err != nil {
			return nil, err
		}

		if len(keys) == 0 {
			keys = keyList
			continue
		}

		keys = ArrayIntersect(keys, keyList)
	}

	b, err := json.Marshal(keys)

	return b, err
}

// Serialize store to JSON
// useful to backup store or for syncing to other stores
// {"propkey1:propvalue1" : ["key1", "key2"], "propkey2:propvalue2" : ["key2"] ...}
func SerializeStore(s Store) ([]byte, error) {
	keyPropStore, err := s.Serialize()

	if err != nil {
		return nil, err
	}

	return json.Marshal(keyPropStore)
}

// DeSerialize JSON and Update current store
// useful to restore store or for updating alternate store
// {"key1" : {"propkey1" : "propvalue1", "propkey2" : "propvalue2"}, "key2" ...}
func DeSerializeStore(s Store, jsBuffer []byte) error {
	var keyPropStore map[string][]string

	if err := json.Unmarshal(jsBuffer, &keyPropStore); err != nil {
		return err
	}

	for key, valueArray := range keyPropStore {
		for _, value := range valueArray {
			if err := s.Update(key, value); err != nil {
				return err
			}
		}
	}

	return nil
}
