package keypropstore

import (
	"testing"
	"os"
	"fmt"
	)

var store Store

func CheckResults(res, expected []string) error{
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

func TestSingleKeyReturn(t *testing.T) {
    query := []byte(`{"num": "6.13","strs": "a"}`)
    expected := []string{"m1"} 
    t.Log("Querying Store for", string(query))

    res := store.QueryStore(query)

    t.Log("Store returned", res, "Expect", expected)

    if err := CheckResults(res, expected); err != nil {
    	t.Error(err)
    }
}

func TestMultipleKeyReturn(t *testing.T) {
    query := []byte(`{"strs": "a"}`)
    expected := []string{"m1", "m3"} 
    t.Log("Querying Store for", string(query))

    res := store.QueryStore(query)

    t.Log("Store returned", res, "Expect", expected)

    if err := CheckResults(res, expected); err != nil {
    	t.Error(err)
    }
}

func BenchmarkQuery(b *testing.B) {  
	query := []byte(`{"strs": "a"}`)
    for n := 0; n < b.N; n++ {
        store.QueryStore(query)
    }
}

func TestMain(m *testing.M) {
    byt := []byte(`{
                    "m1": {"num": "6.13","strs": "a","key1": "b"}, 
                    "m2": {"num": "6.13","key1": "bddd"}, 
                    "m3": {"strs": "a","key1": "b"}, 
                    "m4": {"key1": "asdasdb"}
                }`)

    store.InitializeStore(byt)

    os.Exit(m.Run())
}