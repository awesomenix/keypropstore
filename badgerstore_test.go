package keypropstore

import (
	"github.com/dgraph-io/badger"
	"os"
	"testing"
)

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
