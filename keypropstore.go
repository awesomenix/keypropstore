package keypropstore

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
			err := s.Update(keyval, key)
			if err != nil {
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
