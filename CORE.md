# Store Core

Store core consists of (property, [key1, key2, ...]) pairs, property is represented by "propertykey:propertyvalue". InMemorystore represents the cache layer and serves as the primary store. Secondary stores can be configured, currently supports [badgerdb](https://github.com/dgraph-io/badger)

**Store Core Usage:**

- Initialize the InMemory store with default configuration

```golang
    inMemStore := &InMemoryStore{}
    InitializeStore(inMemStore, nil)
```

- Optionally can use alternate badger db store

```golang
    badgerStore := &BadgerStore{}
    InitializeStore(badgerStore, nil)
```

- UpdateStore with Key and its Properties using JSON format

```golang
	byt := []byte(`{
                    "m1": {"num": "6.13","strs": "a","key1": "b"}, 
                    "m2": {"num": "6.13","key1": "bddd"}, 
                    "m3": {"strs": "a","key1": "b"}, 
                    "m4": {"key1": "asdasdb"}
                }`)
    UpdateStore(inMemStore, byt)
```

- Querying the Store using JSON, optional multiple key value property (always AND, query multiple times for OR), return keys string array

```golang
    query := []byte(`{"num": "6.13","strs": "a"}`)
    res, err := QueryStore(inMemStore, query)
```
- Serialize the Store to JSON, use JSON to Deserialize to other store
```golang
    res, err := SerializeStore(inMemStore)
    err := DeSerializeStore(badgerStore, res)
```
