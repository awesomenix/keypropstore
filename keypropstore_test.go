package keypropstore

import (
	"encoding/json"
	"fmt"
    "os"
	"testing"
    "github.com/dgraph-io/badger"
)

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

func testStoreSingleKeyReturn(s Store, t *testing.T) error {
    query := []byte(`{"num": "6.13","strs": "a"}`)
    expected := []byte(`["m1"]`)
    t.Log("Querying Store for", string(query))

    res, err := QueryStore(s, query)

    if err != nil {
        t.Error(err)
        return err
    }

    t.Log("Store returned", string(res), "Expect", string(expected))

    if err := CheckResults(res, expected); err != nil {
        t.Error(err)
        return err
    }

    return nil
}

func testStoreMultipleKeyReturn(s Store, t *testing.T) error {
    query := []byte(`{"strs": "a"}`)
    expected := []byte(`["m1","m3"]`)
    t.Log("Querying Store for", string(query))

    res, err := QueryStore(s, query)

    if err != nil {
        t.Error(err)
        return err
    }

    t.Log("Store returned", string(res), "Expect", string(expected))

    if err := CheckResults(res, expected); err != nil {
        t.Error(err)
        return err
    }

    return nil
}

func TestInMemStoreSingleKey(t *testing.T) {
    byt := []byte(`{
                    "m1": {"num": "6.13","strs": "a","key1": "b"}, 
                    "m2": {"num": "6.13","key1": "bddd"}, 
                    "m3": {"strs": "a","key1": "b"}, 
                    "m4": {"key1": "asdasdb"}
                }`)

    inMemStore := &InMemoryStore{}
    InitializeStore(inMemStore, nil)
    defer ShutdownStore(inMemStore)
    err:= UpdateStore(inMemStore, byt)
    if err != nil {
        t.Error(err)
        return
    }

    testStoreSingleKeyReturn(inMemStore, t)
}

func TestInMemStoreMultiple(t *testing.T) {
    byt := []byte(`{
                    "m1": {"num": "6.13","strs": "a","key1": "b"}, 
                    "m2": {"num": "6.13","key1": "bddd"}, 
                    "m3": {"strs": "a","key1": "b"}, 
                    "m4": {"key1": "asdasdb"}
                }`)

    inMemStore := &InMemoryStore{}
    InitializeStore(inMemStore, nil)
    defer ShutdownStore(inMemStore)
    err := UpdateStore(inMemStore, byt)
    if err != nil {
        t.Error(err)
        return
    }

	testStoreMultipleKeyReturn(inMemStore, t)
}

func TestBadgerStoreSingleKey(t *testing.T) {
    byt := []byte(`{
                    "m1": {"num": "6.13","strs": "a","key1": "b"}, 
                    "m2": {"num": "6.13","key1": "bddd"}, 
                    "m3": {"strs": "a","key1": "b"}, 
                    "m4": {"key1": "asdasdb"}
                }`)

    directory := "./badgerdb"
    os.RemoveAll(directory)
    defer os.RemoveAll(directory)

    opts := badger.DefaultOptions
    opts.Dir = directory
    opts.ValueDir = directory

    badgerStore := new(BadgerStore)
    InitializeStore(badgerStore, opts)
    defer ShutdownStore(badgerStore)
    err := UpdateStore(badgerStore, byt)
    if err != nil {
        t.Error(err)
        return
    }

    testStoreSingleKeyReturn(badgerStore, t)
}

func TestBadgerStoreMultipleKey(t *testing.T) {
    byt := []byte(`{
                    "m1": {"num": "6.13","strs": "a","key1": "b"}, 
                    "m2": {"num": "6.13","key1": "bddd"}, 
                    "m3": {"strs": "a","key1": "b"}, 
                    "m4": {"key1": "asdasdb"}
                }`)

    directory := "./badgerdb"
    os.RemoveAll(directory)
    defer os.RemoveAll(directory)

    opts := badger.DefaultOptions
    opts.Dir = directory
    opts.ValueDir = directory

    badgerStore := new(BadgerStore)
    InitializeStore(badgerStore, opts)
    defer ShutdownStore(badgerStore)
    err := UpdateStore(badgerStore, byt)
    if err != nil {
        t.Error(err)
        return
    }

    testStoreMultipleKeyReturn(badgerStore, t)
}

func BenchmarkInMemStoreQuery(b *testing.B) {
    byt := []byte(`{
                    "m1": {"num": "6.13","strs": "a","key1": "b"}, 
                    "m2": {"num": "6.13","key1": "bddd"}, 
                    "m3": {"strs": "a","key1": "b"}, 
                    "m4": {"key1": "asdasdb"}
                }`)

    inMemStore := &InMemoryStore{}
    InitializeStore(inMemStore, nil)
    defer ShutdownStore(inMemStore)
    err := UpdateStore(inMemStore, byt)
    if err != nil {
        b.Error(err)
        return
    }

	query := []byte(`{"strs": "a"}`)
	for n := 0; n < b.N; n++ {
		QueryStore(inMemStore, query)
	}
}