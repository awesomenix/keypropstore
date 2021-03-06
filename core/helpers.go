package core

import (
	"fmt"
)

// ArrayIntersect Performs Intersection of two string array
// intersection of (m1, m3) (m1, m4) = (m1)
// any empty array results in no results
func ArrayIntersect(a, b []string) []string {
	if len(a) == 0 || len(b) == 0 {
		return nil
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

// GenerateKey returns hash of key value used as a Store Key
func GenerateKey(key, val string) string {
	return fmt.Sprintf("%s:%s", key, val)
}
