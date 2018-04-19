package app

import (
	"log"
	"net/http"

	"github.com/awesomenix/keypropstore/core"
)

// APIVERSION current supported version
const APIVERSION string = "/v1"

// Context stores local and aggregate stores
type Context struct {
	stores    map[string]core.Store
	appRoutes []Route
}

// CreateStore specified in Configuration
func createStore(storeType string) core.Store {
	switch storeType {
	case "InMemory":
		return new(core.InMemoryStore)
	case "BadgerDB":
		return new(core.BadgerStore)
	}
	return new(core.InMemoryStore)
}

// InitializeStores initializes all the predefined store in configuration
func (ctx *Context) InitializeStores() error {
	ctx.stores = make(map[string]core.Store)
	ctx.stores["local"] = createStore("InMemory")
	return core.InitializeStore(ctx.stores["local"], nil)
}

// ShutdownStores shutsdown all the predefined store in configuration
func (ctx *Context) ShutdownStores() error {
	return core.ShutdownStore(ctx.stores["local"])
}

// Execute App context creating router handling multiple REST API
func (ctx *Context) Execute() {
	ctx.registerRoutes()
	appRouter := NewAppRouter(ctx)
	log.Fatal(http.ListenAndServe(":8080", appRouter))
}

// Execute creates a local app context and Executes the context
func Execute() {
	ctx := Context{}
	ctx.InitializeStores()
	defer ctx.ShutdownStores()
	ctx.Execute()
}
