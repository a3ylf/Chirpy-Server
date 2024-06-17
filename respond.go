package main

import (
	"encoding/json"
	"log"
	"net/http"
)


func respondWithJSON(w http.ResponseWriter, code int, toSend interface{}) {
    w.Header().Set("Content-Type", "application/json")
    dat, err := json.Marshal(toSend)
    if err != nil {
        log.Printf("Error Marshaling: %s", err)
        w.WriteHeader(500)
        return
    }
    w.WriteHeader(code)
    w.Write(dat)
}
func respondWithError(w http.ResponseWriter,code int, msg string) {
    if code > 499 {
        log.Printf("Responding with error 5XX : %s",msg)
    }
    type errorResponse struct {
        Error string `json:"error`
    }

    respondWithJSON(w,code,errorResponse{Error: msg})
}
