package core

import (
	"testing"
)

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

/* func TestInMemStoreSerializeDeSerialize(t *testing.T) {
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
} */

func BenchmarkInMemStoreUpdateQuery(b *testing.B) {
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
		_, err = QueryStore(inMemStore, query)
		if err != nil {
			b.Error(err)
			return
		}
	}
}
