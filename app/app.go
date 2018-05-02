package app

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// APIVERSION current supported version
const APIVERSION string = "/v1"

// Context stores local and aggregate stores
type Context struct {
	config    Config
	stores    map[string]*CoreStores
	appRoutes []Route
	srv       *http.Server
}

// Create App context creating router handling multiple REST API
func (ctx *Context) Create() error {
	// register default and app routes
	ctx.registerRoutes()
	appRouter := NewAppRouter(ctx)
	ctx.srv = &http.Server{
		Addr:         ":" + ctx.config.Port,
		Handler:      appRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second}
	go func() {
		if err := ctx.srv.ListenAndServe(); err != nil {
			if flag.Lookup("test.v") == nil {
				log.Fatal(err)
			} else {
				log.Println(err)
			}
		}
	}()
	return nil
}

// CreateContext creates and sets up context, stores and starts HTTP Server
func CreateContext() *Context {
	ctx := &Context{}
	// Initialize configuration
	ctx.config.Initialize("config", "./config")
	// Initialize any stores, primary and backup
	ctx.InitializeStores()
	// Create context
	ctx.Create()

	return ctx
}

// DeleteContext app context and HTTP Server
func DeleteContext(ctx *Context) {
	// Shutdown of all the stores
	ctx.ShutdownStores()
	// Shutdown HTTP server
	ctx.srv.Shutdown(context.TODO())
}

func waitForCtrlC() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	<-sigc
}

// Execute starts application and waits for ctrl+c
func Execute() {
	ctx := CreateContext()
	waitForCtrlC()
	DeleteContext(ctx)
}
