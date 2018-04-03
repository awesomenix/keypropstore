package keypropstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger"
	"os"
	"strings"
	"testing"
)

var byt = []byte(`{
                    "m1": {"num": "6.13","strs": "a","key1": "b"}, 
                    "m2": {"num": "6.13","key1": "bddd"}, 
                    "m3": {"strs": "a","key1": "b"}, 
                    "m4": {"key1": "asdasdb"}
                }`)

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
	inMemStore := &InMemoryStore{}
	InitializeStore(inMemStore, nil)
	defer ShutdownStore(inMemStore)
	err := UpdateStore(inMemStore, byt)
	if err != nil {
		t.Error(err)
		return
	}

	testStoreSingleKeyReturn(inMemStore, t)
}

func TestInMemStoreMultiple(t *testing.T) {
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

func testStoreSerializeDeSerialize(oldStore Store, newStore Store, t *testing.T) error {
	oldRes, err := SerializeStore(oldStore)

	if err != nil {
		t.Error(err)
		return err
	}

	t.Log("Store returned", string(oldRes))

	if err := DeSerializeStore(newStore, oldRes); err != nil {
		t.Error(err)
		return err
	}

	newRes, err := SerializeStore(newStore)

	if err != nil {
		t.Error(err)
		return err
	}

	if strings.Compare(string(oldRes), string(newRes)) != 0 {
		err := errors.New(fmt.Sprintf("Expected Old Store %s to be same as New Store %s", string(oldRes), string(newRes)))
		t.Error(err)
		return err
	}

	return nil
}

func TestInMemStoreSerializeDeSerialize(t *testing.T) {
	inMemStore := &InMemoryStore{}
	InitializeStore(inMemStore, nil)
	defer ShutdownStore(inMemStore)
	err := UpdateStore(inMemStore, byt)
	if err != nil {
		t.Error(err)
		return
	}

	inMemStoreNew := &InMemoryStore{}
	InitializeStore(inMemStoreNew, nil)
	defer ShutdownStore(inMemStoreNew)

	testStoreSerializeDeSerialize(inMemStore, inMemStoreNew, t)
}

func TestBadgerStoreSerializeDeSerialize(t *testing.T) {
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

	directoryNew := "./badgerdbNew"
	os.RemoveAll(directoryNew)
	defer os.RemoveAll(directoryNew)

	opts.Dir = directoryNew
	opts.ValueDir = directoryNew

	badgerStoreNew := new(BadgerStore)
	InitializeStore(badgerStoreNew, opts)
	defer ShutdownStore(badgerStoreNew)

	testStoreSerializeDeSerialize(badgerStore, badgerStoreNew, t)
}

func TestCrossStoreSerializeDeSerialize(t *testing.T) {
	inMemStore := &InMemoryStore{}
	InitializeStore(inMemStore, nil)
	defer ShutdownStore(inMemStore)
	err := UpdateStore(inMemStore, byt)
	if err != nil {
		t.Error(err)
		return
	}

	directory := "./badgerdb"
	os.RemoveAll(directory)
	defer os.RemoveAll(directory)

	opts := badger.DefaultOptions
	opts.Dir = directory
	opts.ValueDir = directory

	badgerStore := new(BadgerStore)
	InitializeStore(badgerStore, opts)
	defer ShutdownStore(badgerStore)

	testStoreSerializeDeSerialize(inMemStore, badgerStore, t)
}

func TestMoreCrossStoreSerializeDeSerialize(t *testing.T) {
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

	inMemStore := &InMemoryStore{}
	InitializeStore(inMemStore, nil)
	defer ShutdownStore(inMemStore)

	testStoreSerializeDeSerialize(badgerStore, inMemStore, t)
}

func BenchmarkInMemStoreQuery(b *testing.B) {
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
