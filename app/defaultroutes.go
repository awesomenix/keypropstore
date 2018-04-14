package app

import (
	"net/http"
	"net/http/pprof"

	"github.com/julienschmidt/httprouter"
)

var defaultRoutes = []Route{
	Route{"GET", "/health", HealthCheckHandler},
	Route{"GET", "/debug/pprof/", IndexHandler},
	Route{"GET", "/debug/pprof/heap", HeapHandler},
	Route{"GET", "/debug/pprof/goroutine", GoroutineHandler},
	Route{"GET", "/debug/pprof/block", BlockHandler},
	Route{"GET", "/debug/pprof/threadcreate", ThreadCreateHandler},
	Route{"GET", "/debug/pprof/cmdline", CmdlineHandler},
	Route{"GET", "/debug/pprof/profile", ProfileHandler},
	Route{"GET", "/debug/pprof/symbol", SymbolHandler},
	Route{"POST", "/debug/pprof/symbol", SymbolHandler},
	Route{"GET", "/debug/pprof/trace", TraceHandler},
	Route{"GET", "/debug/pprof/mutex", MutexHandler},
}

// HealthCheckHandler provides health check for external montioring applications
func HealthCheckHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	respondOK(w, "ok")
}

// IndexHandler will pass the call from /debug/pprof to pprof
func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pprof.Index(w, r)
}

// HeapHandler will pass the call from /debug/pprof/heap to pprof
func HeapHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pprof.Handler("heap").ServeHTTP(w, r)
}

// GoroutineHandler will pass the call from /debug/pprof/goroutine to pprof
func GoroutineHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pprof.Handler("goroutine").ServeHTTP(w, r)
}

// BlockHandler will pass the call from /debug/pprof/block to pprof
func BlockHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pprof.Handler("block").ServeHTTP(w, r)
}

// ThreadCreateHandler will pass the call from /debug/pprof/threadcreate to pprof
func ThreadCreateHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pprof.Handler("threadcreate").ServeHTTP(w, r)
}

// CmdlineHandler will pass the call from /debug/pprof/cmdline to pprof
func CmdlineHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pprof.Cmdline(w, r)
}

// ProfileHandler will pass the call from /debug/pprof/profile to pprof
func ProfileHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pprof.Profile(w, r)
}

// SymbolHandler will pass the call from /debug/pprof/symbol to pprof
func SymbolHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pprof.Symbol(w, r)
}

// TraceHandler will pass the call from /debug/pprof/trace to pprof
func TraceHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pprof.Trace(w, r)
}

// MutexHandler will pass the call from /debug/pprof/mutex to pprof
func MutexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pprof.Handler("mutex").ServeHTTP(w, r)
}
