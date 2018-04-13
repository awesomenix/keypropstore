package app

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// NewAppRouter registers multiple logged routes
func NewAppRouter(ctx *Context) http.Handler {
	router := httprouter.New()
	router.GET("/api/v1/:store/query", LoggerMiddleware(ctx.queryStore))
	router.POST("/api/v1/:store/update", LoggerMiddleware(ctx.updateStore))
	return http.Handler(router)
}
