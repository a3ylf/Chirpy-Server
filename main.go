package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)
type apiconfig struct{
    fileserverhits int
}



func main() {
    var apicfg apiconfig
    mux := http.NewServeMux()
    mux.Handle("/app/*",apicfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
    mux.HandleFunc("GET /api/healthz",handlerReadiness)
    mux.HandleFunc("/api/reset",apicfg.handlerReset)
    mux.HandleFunc("/api/validate_chirp",handleValidation)
    mux.HandleFunc("GET /admin/metrics",apicfg.handlerMetrics)
    server := &http.Server{
		Addr:    ":8080",  	
		Handler: mux,      

	} 
	log.Printf("Serving on port: 8080\n")
	log.Fatal(server.ListenAndServe())
} 

func (cfg *apiconfig) middlewareMetricsInc(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverhits++
        next.ServeHTTP(w,r)
    })
}
func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Cache-control","no-cache")
    w.WriteHeader(200)
}


func (cfg *apiconfig) handlerReset(w http.ResponseWriter, r*http.Request){
    cfg.fileserverhits = 0
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hits is now 0"))
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
        Error string `json:"error`
    }

    respondWithJSON(w,code,errorResponse{Error: msg})
}
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

func cleaner(pog string, bad []string) string {
    splitted := strings.Split(pog," ")
    for _, badword := range bad {
        for k, split := range splitted {
            if strings.ToLower(split) == badword {
                splitted[k] = "****"
            }
        }
    }
    toSend := "" 
    for _, split := range splitted {
        if len(toSend) == 0 {
            toSend = fmt.Sprintf(split)
            continue
        }
        toSend = fmt.Sprintf(toSend+ " " + split )
    }
    return toSend
    
}


