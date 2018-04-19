package app

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/awesomenix/keypropstore/core"
	"github.com/julienschmidt/httprouter"
)

func (ctx *Context) registerRoutes() {
	ctx.appRoutes = []Route{
		Route{"POST", "/store/:store/query", ctx.queryStore},
		Route{"POST", "/store/:store/update", ctx.updateStore},
		Route{"GET", "/store/:store/backup", ctx.backupStore},
		Route{"POST", "/store/:store/restore", ctx.restoreStore},
	}
}

func (ctx *Context) queryStore(w http.ResponseWriter, r *http.Request, httpParams httprouter.Params) {
	storeName := httpParams.ByName("store")
	store, ok := ctx.stores[storeName]

	if !ok {
		err := fmt.Sprintf("invalid or store %s not found", storeName)
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	propQuery, err := ioutil.ReadAll(r.Body)
	if err != nil {

	}

	if err := r.Body.Close(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	jsRes, err := core.QueryStore(store, propQuery)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, jsRes)
}

func (ctx *Context) updateStore(w http.ResponseWriter, r *http.Request, httpParams httprouter.Params) {
	storeName := httpParams.ByName("store")
	store, ok := ctx.stores[storeName]

	if !ok {
		err := fmt.Sprintf("invalid or store %s not found", storeName)
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	jsReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.Body.Close(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = core.UpdateStore(store, jsReq)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondOK(w, "ok")
}

func (ctx *Context) backupStore(w http.ResponseWriter, r *http.Request, httpParams httprouter.Params) {
	storeName := httpParams.ByName("store")
	store, ok := ctx.stores[storeName]

	if !ok {
		err := fmt.Sprintf("invalid or store %s not found", storeName)
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	jsRes, err := core.SerializeStore(store)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, jsRes)
}

func (ctx *Context) restoreStore(w http.ResponseWriter, r *http.Request, httpParams httprouter.Params) {
	storeName := httpParams.ByName("store")
	store, ok := ctx.stores[storeName]

	if !ok {
		err := fmt.Sprintf("invalid or store %s not found", storeName)
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	jsReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.Body.Close(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = core.DeSerializeStore(store, jsReq)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondOK(w, "ok")
}
