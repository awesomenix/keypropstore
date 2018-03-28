package keypropstore

import (
    "fmt"
    "encoding/json"
    ) 

type Store struct {
    store map[string][]string
}

func (s *Store) FmtKeyVal(key, val string) string{
    return fmt.Sprintf("%s:%s", key, val)
}

func (s *Store) InitializeStore(byt []byte) error {
    var dat map[string]map[string]string    

   if err := json.Unmarshal(byt, &dat); err != nil {
        panic(err)
    }

    s.store = make(map[string][]string)

    for key, value := range dat {
        for valkey, valval := range value {
            keyval := s.FmtKeyVal(valkey, valval)
            val, ok := s.store[keyval]
            if !ok {
                s.store[keyval] = make([]string, 1)
            }
            s.store[keyval] = append(val, key)
        }
    }

    return nil
}

func (s *Store) Intersect(a, b []string) []string {
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

func (s *Store) QueryStore(jsQuery []byte) []string {
    keys := make([]string, 0)

    var query map[string]string

    if err := json.Unmarshal(jsQuery, &query); err != nil {
        panic(err)
    }

    q := make([]string, 0)

    for key, val := range query {
        q = append(q, s.FmtKeyVal(key, val))
    }

    for _, val := range q {
        keyList, ok := s.store[val]

        if !ok {
            continue
        }

        keys = s.Intersect(keys, keyList)
    }

    return keys
}
