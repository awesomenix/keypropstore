package keypropstore

import (
    "fmt"
    "encoding/json"
    "strings"
    "sync"
    ) 

// Performs Intersection of two string array
// intersection of (m1, m3) (m1, m4) = (m1)
func intersect(a, b []string) []string {
    if len(a) == 0 {
        return b
    }

    if len(b) == 0 {
        return a
    }

    hashKey := make(map[string]struct{})

    for _, key := range a {
        hashKey[key] = struct{}{}
    }

    ret := make([]string, 0)

    for _, key := range b {
        if _, ok := hashKey[key]; ok {
            ret = append(ret, key)
        }
    }    
    return ret
}

// Returns hash of key value used as a Store Key
func fmtKeyVal(key, val string) string{
    return fmt.Sprintf("%s:%s", key, val)
}

func contains(key string, arr []string) bool {
    for _, str := range arr {
        if strings.Compare(str, key) == 0 {
            return true
        }
    }
    return false
}

// Concurrent friendly Store
// Map of Property and List of Keys associated with that property
type Store struct {
    store map[string][]string
    lock sync.RWMutex
}

// Initialize Store with custom configuration
func (s *Store) InitializeStore() error {
    s.store = make(map[string][]string)
    return nil
}

// Update Store called with list of Key and Associated properties
// Check for consistncy of input json
// Converts
// {"m1": {"num": "6.13","strs": "a","key1": "b"}, "m2": {"num": "6.13","key1": "bddd"}}
// To
// {"num:6.13" : ["m1", m2"], "strs:a" : ["m1"], "key1:b" : ["m1"], "key1:bddd" : ["m2"]}
func (s *Store) UpdateStore(byt []byte) error {
    var dat map[string]map[string]string    

    if err := json.Unmarshal(byt, &dat); err != nil {
        return err
    }

    for key, value := range dat {
        for valkey, valval := range value {
            keyval := fmtKeyVal(valkey, valval)
            s.lock.Lock()
            val, ok := s.store[keyval]
            if !ok {
                s.store[keyval] = make([]string, 1)
            }

            if !contains(key, val) {
                s.store[keyval] = append(val, strings.ToLower(key))
            }

            s.lock.Unlock()
        }
    }

    return nil
}

// Query Store with single/multiple properties
// Properties are AND only, if an OR is required, query multiple times
// OR could be supported, but keeping it simple for now
func (s *Store) QueryStore(jsQuery []byte) ([]byte, error) {
    keys := make([]string, 0)

    var query map[string]string

    if err := json.Unmarshal(jsQuery, &query); err != nil {
        return nil, err
    }

    q := make([]string, 0)

    for key, val := range query {
        q = append(q, fmtKeyVal(key, val))
    }

    for _, val := range q {

        s.lock.RLock()
        keyList, ok := s.store[val]
        s.lock.RUnlock()

        if !ok {
            continue
        }

        keys = intersect(keys, keyList)
    }

    b, err:= json.Marshal(keys)

    return b, err
}
