package app

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// LoggerMiddleware Handler
func LoggerMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// log app name, method and URI
		log.Printf("[keypropstore] Method: %s, URI: %s\n", r.Method, r.RequestURI)
		next(w, r, ps)
	}
}
