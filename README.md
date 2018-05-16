![](logo.png)
------
[![Go Report Card](https://goreportcard.com/badge/github.com/awesomenix/keypropstore)](https://goreportcard.com/report/github.com/awesomenix/keypropstore) [![Build Status](https://travis-ci.org/awesomenix/keypropstore.svg?branch=master)](https://travis-ci.org/awesomenix/keypropstore) [![codecov](https://codecov.io/gh/awesomenix/keypropstore/branch/master/graph/badge.svg)](https://codecov.io/gh/awesomenix/keypropstore)

Simulates a column store, where a key could be tagged with multiple key value properties.

- Multiple services could report keys with associated properties, it will appended the local store.
- Query for property/properties would return the keys associated with those property/properties from local store (default)
- Keypropstore is hosted in a region consists of local and optional aggregate stores, with configurable backends (default InMemoryStore).
- Aggregate stores are configured to sync remote local stores into a separate aggregate store instance.
- Queries could be made to aggregate store, could also potentially host multiple aggregate store instances.

## Architecture

![Architecture](architecture.png)

## Store Core

Keypropstore consists of multiple [store core](docs/CORE.md), local and/or aggregate.
