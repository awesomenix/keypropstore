# keypropstore

Simulates column store, where a key could be tagged with multiple key value properties. Query to property/properties would return the keys associated with the property/properties

Usage:

- Initialize the store with default configuration

```golang
   byt := []byte(`{
                    "m1": {"num": "6.13","strs": "a","key1": "b"}, 
                    "m2": {"num": "6.13","key1": "bddd"}, 
                    "m3": {"strs": "a","key1": "b"}, 
                    "m4": {"key1": "asdasdb"}
                }`)
    store.InitializeStore()
```

- UpdateStore with Key and its Properties using JSON format

```golang
	store.UpdateStore(byt)
```

- Querying the Store using JSON, optional multiple key value property (always AND, query multiple times for OR), return keys string array

```golang
    query := []byte(`{"num": "6.13","strs": "a"}`)
    res, err := store.QueryStore(query)
```
