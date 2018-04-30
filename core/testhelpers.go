package core

import (
	"encoding/json"
	"fmt"
)

var byt = []byte(`{
                    "m1": {"num": "6.13","strs": "a","key1": "b"}, 
                    "m2": {"num": "6.13","key1": "bddd"}, 
                    "m3": {"strs": "a","key1": "b"}, 
                    "m4": {"key1": "asdasdb"}
                }`)

// CheckResults current vs expected to error in case there is a mismatch
func CheckResults(jsres, jsexpected []byte) error {
	var res, expected []string

	if err := json.Unmarshal(jsres, &res); err != nil {
		return err
	}

	if err := json.Unmarshal(jsexpected, &expected); err != nil {
		return err
	}

	for _, expKey := range expected {
		found := false
		for _, key := range res {
			if key == expKey {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Expected %v", expKey)
		}
	}

	return nil
}
