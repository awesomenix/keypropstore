package core

import (
	"os"
	"testing"
)

func TestBoltStoreSingleKey(t *testing.T) {
	directory := "./boltdb"
	os.RemoveAll(directory)
	defer os.RemoveAll(directory)

	opts := &BoltStoreConfig{directory, 600, nil}

	boltStore := new(BoltStore)
	InitializeStore(boltStore, opts)
	defer ShutdownStore(boltStore)
	err := UpdateStore(boltStore, byt)
	if err != nil {
		t.Error(err)
		return
	}

	testStoreSingleKeyReturn(boltStore, t)
}

func TestBoltStoreMultipleKey(t *testing.T) {
	directory := "./boltdb"
	os.RemoveAll(directory)
	defer os.RemoveAll(directory)

	opts := &BoltStoreConfig{directory, 600, nil}

	boltStore := new(BoltStore)
	InitializeStore(boltStore, opts)
	defer ShutdownStore(boltStore)
	err := UpdateStore(boltStore, byt)
	if err != nil {
		t.Error(err)
		return
	}

	testStoreMultipleKeyReturn(boltStore, t)
}

func TestBoltStoreSerializeDeSerialize(t *testing.T) {
	directory := "./boltdb"
	os.RemoveAll(directory)
	defer os.RemoveAll(directory)

	opts := &BoltStoreConfig{directory, 600, nil}

	boltStore := new(BoltStore)
	InitializeStore(boltStore, opts)
	defer ShutdownStore(boltStore)
	err := UpdateStore(boltStore, byt)
	if err != nil {
		t.Error(err)
		return
	}

	directoryNew := "./boltdbNew"
	os.RemoveAll(directoryNew)
	defer os.RemoveAll(directoryNew)

	opts.Path = directoryNew

	boltStoreNew := new(BoltStore)
	InitializeStore(boltStoreNew, opts)
	defer ShutdownStore(boltStoreNew)

	testStoreSerializeDeSerialize(boltStore, boltStoreNew, t)
}
