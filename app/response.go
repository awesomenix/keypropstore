package app

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	response, _ := json.Marshal(map[string]string{"status": "error", "message": message})
	respondJSON(w, code, response)
}

func respondOK(w http.ResponseWriter, message string) {
	response, _ := json.Marshal(map[string]string{"status": "success", "message": message})
	respondJSON(w, http.StatusOK, response)
}

func respondJSON(w http.ResponseWriter, code int, response []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
