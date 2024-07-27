package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)
func handleValidation(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Body string `json:"body"`
    }
   type returnvals struct {
        Cleaned string `json:"cleaned_body"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    
    if err != nil {
        respondWithError(w,http.StatusInternalServerError,"Couldn't decode parameters")
        return
    }
    if len(params.Body) > 140 {
       respondWithError(w, http.StatusBadRequest,"Chirp is too long") 
       return
    }

    bad_words := []string{"kerfuffle","sharbert","fornax"}
    respondWithJSON(w,http.StatusOK,returnvals{Cleaned: cleaner(params.Body,bad_words)})
       
}

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
        Error string `json:"error"`
    }

    respondWithJSON(w,code,errorResponse{Error: msg})
}
