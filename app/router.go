package app

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Route defines registration of different routes supported by app
type Route struct {
	Method   string
	RouteURI string
	Handler  httprouter.Handle
}

// NewAppRouter registers multiple logged routes
func NewAppRouter(ctx *Context) http.Handler {
	router := httprouter.New()
	// Register default routes, which typically provides healtcheck and prof
	for _, route := range defaultRoutes {
		router.Handle(route.Method, APIVERSION+route.RouteURI, LoggerMiddleware(route.Handler))
	}

	for _, route := range ctx.appRoutes {
		router.Handle(route.Method, APIVERSION+route.RouteURI, LoggerMiddleware(route.Handler))
	}

	return http.Handler(router)
}
