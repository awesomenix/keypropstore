package app

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/awesomenix/keypropstore/core"
	"github.com/dgraph-io/badger"
)

// CoreStores contains primary and optional backup store
type CoreStores struct {
	primary  core.Store
	backup   core.Store
	shutdown chan bool
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
		os.Mkdir(storeDir, os.ModePerm)
		storePath := filepath.Join(storeDir, "boltdbstore")
		opts := &core.BoltStoreConfig{Path: storePath, Mode: 600, Options: nil}
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
		if len(store.AggregateURLs) > 0 {
			newstore.shutdown = make(chan bool)
			go ctx.SyncAggregateURLs(newstore, store.AggregateURLs, time.Duration(store.SyncIntervalSec)*time.Second)
		}
		ctx.stores[store.Name] = newstore
	}
	return err
}

// ShutdownStores shutsdown all the predefined store in configuration
func (ctx *Context) ShutdownStores() error {
	var err error
	for _, store := range ctx.stores {
		if store.shutdown != nil {
			store.shutdown <- true
		}

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

// SyncAggregateURLs performs sync of key values from remote stores
// Requires Extensive error checking, metrics, alerting
func (ctx *Context) SyncAggregateURLs(store *CoreStores, aggregateURLs []string, syncIntervalsecs time.Duration) {
	defer close(store.shutdown)
	ticker := time.NewTicker(syncIntervalsecs)
	for {
		select {
		case <-ticker.C:
			for _, aggregateURL := range aggregateURLs {
				backupURL := aggregateURL + "/backup"
				httpResp, httpErr := http.Get(backupURL)

				if httpErr != nil {
					continue
				}

				bodyBytes, rerr := ioutil.ReadAll(httpResp.Body)
				httpResp.Body.Close()

				if rerr != nil {
					continue
				}
				if store.backup != nil {
					if err := core.DeSerializeStore(store.backup, bodyBytes); err != nil {
						// Error deserializing the backup store
					}
				}
				if err := core.DeSerializeStore(store.primary, bodyBytes); err != nil {
					// Error deserializing the primary store
				}
			}
		case <-store.shutdown:
			ticker.Stop()
			return
		}
	}

}
