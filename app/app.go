package app

import (
	"log"
	"net/http"

	"github.com/awesomenix/keypropstore/core"
)

// Context stores local and aggregate stores
type Context struct {
	inMemLocalStore *core.InMemoryStore
}

// Execute App context creating router handling multiple REST API
func (ctx *Context) Execute() {
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
