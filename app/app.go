package app

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/awesomenix/keypropstore/core"
	"github.com/dgraph-io/badger"
)

// APIVERSION current supported version
const APIVERSION string = "/v1"

// CoreStores contains primary and optional backup store
type CoreStores struct {
	primary core.Store
	backup  core.Store
}

// Context stores local and aggregate stores
type Context struct {
	config    Config
	stores    map[string]*CoreStores
	appRoutes []Route
	srv       *http.Server
}

// CreateStore specified in Configuration
func createStore(storeType, storeDir string) (core.Store, error) {
	switch storeType {
	case "InMemory":
		store := new(core.InMemoryStore)
		return store, core.InitializeStore(store, nil)
	case "BadgerDB":
		opts := badger.DefaultOptions
		opts.Dir = storeDir
		opts.ValueDir = storeDir
		store := new(core.BadgerStore)
		return store, core.InitializeStore(store, opts)
	case "BoltDB":
		opts := &core.BoltStoreConfig{Path: storeDir, Mode: 600, Options: nil}
		store := new(core.BoltStore)
		return store, core.InitializeStore(store, opts)
	}
	store := new(core.InMemoryStore)
	return store, core.InitializeStore(store, nil)
}

// InitializeStores initializes all the predefined store in configuration
func (ctx *Context) InitializeStores() error {
	var err error
	ctx.stores = make(map[string]*CoreStores)
	for _, store := range ctx.config.Stores {
		// Initialize primary in memory store
		log.Printf("Initializing Primary InMemoryStore %s\n", store.Name)
		newstore := &CoreStores{}
		var localerr error
		if newstore.primary, localerr = createStore("InMemory", ""); localerr != nil {
			err = localerr
		}
		// Initialize backup store if defined
		if len(store.Backup) > 0 {
			log.Printf("Initializing Backup Store %s of type %s, backup directory %s\n", store.Name, store.Backup, store.Backupdir)
			var localerr error
			if newstore.backup, localerr = createStore(store.Backup, store.Backupdir); localerr != nil {
				err = localerr
			} else {
				// Once initialized we need to restore the primary store from backup store
				jsStore, serr := core.SerializeStore(newstore.backup)
				if serr != nil {
					err = serr
				} else {
					if dserr := core.DeSerializeStore(newstore.primary, jsStore); dserr != nil {
						err = dserr
					}
				}
			}
		}
		ctx.stores[store.Name] = newstore
	}
	return err
}

// ShutdownStores shutsdown all the predefined store in configuration
func (ctx *Context) ShutdownStores() error {
	var err error
	for _, store := range ctx.stores {
		// shutdown primary store
		if localerr := core.ShutdownStore(store.primary); localerr != nil {
			err = localerr
		}
		// shutdown backup stores if any
		if store.backup != nil {
			if localerr := core.ShutdownStore(store.backup); localerr != nil {
				err = localerr
			}
		}
	}
	return err
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
	var exit sync.WaitGroup
	exit.Add(1)

	var sigc chan os.Signal
	sigc = make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	go func() {
		<-sigc
		exit.Done()
	}()
	exit.Wait()
}

// Execute starts application and waits for ctrl+c
func Execute() {
	ctx := CreateContext()
	waitForCtrlC()
	DeleteContext(ctx)
}
