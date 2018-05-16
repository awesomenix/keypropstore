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
		Addr:         "127.0.0.1:" + ctx.config.Port,
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
func CreateContext(configName, configDir string) (*Context, error) {
	ctx := &Context{}
	// Initialize configuration
	if err := ctx.config.Initialize(configName, configDir); err != nil {
		return nil, err
	}
	// Initialize any stores, primary and backup
	if err := ctx.InitializeStores(); err != nil {
		return nil, err
	}
	// Create context
	if err := ctx.Create(); err != nil {
		return nil, err
	}

	return ctx, nil
}

// CreateDefaultContext creates and sets up default context, stores and starts HTTP Server
func CreateDefaultContext() (*Context, error) {
	return CreateContext("config", "./config")
}

// DeleteContext app context and HTTP Server
func DeleteContext(ctx *Context) error {
	// Shutdown HTTP server
	// even if there is an error shutting down HTTP its ok to ignore
	ctx.srv.Shutdown(context.TODO())
	// Shutdown of all the stores
	return ctx.ShutdownStores()
}

func waitForCtrlC() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	<-sigc
}

// Execute starts application and waits for ctrl+c
func Execute() error {
	ctx, err := CreateDefaultContext()
	if err != nil {
		return err
	}
	waitForCtrlC()
	return DeleteContext(ctx)
}
