package app

import (
	"io/ioutil"
	"net/http"

	"github.com/awesomenix/keypropstore/core"
	"github.com/julienschmidt/httprouter"
)

func (ctx *Context) registerRoutes() {
	ctx.appRoutes = []Route{
		Route{"GET", "/store/:store/query", ctx.queryStore},
		Route{"POST", "/store/:store/update", ctx.updateStore},
	}
}

func (ctx *Context) queryStore(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryParams := r.URL.Query()

	propQuery := queryParams.Get("prop")
	if propQuery == "" {
		respondWithError(w, http.StatusBadRequest, "prop param not found")
		return
	}

	jsRes, err := core.QueryStore(ctx.inMemLocalStore, []byte(propQuery))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, jsRes)
}

func (ctx *Context) updateStore(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jsReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.Body.Close(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = core.UpdateStore(ctx.inMemLocalStore, jsReq)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondOK(w, "ok")
}
