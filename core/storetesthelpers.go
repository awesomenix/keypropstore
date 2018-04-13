package core

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/dgraph-io/badger"
)

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
