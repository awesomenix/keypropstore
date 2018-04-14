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
	inMemLocalStore *core.InMemoryStore
	appRoutes       []Route
}

// Execute App context creating router handling multiple REST API
func (ctx *Context) Execute() {
	ctx.registerRoutes()
	appRouter := NewAppRouter(ctx)
	log.Fatal(http.ListenAndServe(":8080", appRouter))
}

// Execute creates a local app context and Executes the context
func Execute() {
	ctx := Context{inMemLocalStore: &core.InMemoryStore{}}
	core.InitializeStore(ctx.inMemLocalStore, nil)
	defer core.ShutdownStore(ctx.inMemLocalStore)
	ctx.Execute()
}
