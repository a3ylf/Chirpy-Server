package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/a3ylf/web-servers/database"
	"github.com/joho/godotenv"
)
type apiconfig struct{
    fileserverhits int
    db *database.DB
    secret string
}



func main() {
    
    dbg := flag.Bool("debug",false,"enable debug mode")
    flag.Parse()
    
    if *dbg {
        os.Remove("database/database.json")
    }
    db, err := database.NewDB("database/database.json")

    if err != nil {
        log.Println(err)
    }

     
    err = godotenv.Load()
    if err != nil {
        log.Println(err)
        return
    }

    jwtSecret := os.Getenv("JWT_SECRET")


    apicfg := apiconfig{
        fileserverhits: 0,
        db: db,
        secret: jwtSecret,
    }


    mux := http.NewServeMux()
    mux.Handle("/app/*",apicfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
    mux.HandleFunc("GET /api/healthz",handlerReadiness)
    mux.HandleFunc("/api/reset",apicfg.handlerReset)
    mux.HandleFunc("GET /admin/metrics",apicfg.handlerMetrics)
    mux.HandleFunc("POST /api/chirps",apicfg.handlePost)
    mux.HandleFunc("GET /api/chirps",apicfg.handleGet)
    mux.HandleFunc("GET /api/chirps/{ID}",apicfg.handleGetChirp)
    mux.HandleFunc("POST /api/users",apicfg.handleUserPost)
    mux.HandleFunc("POST /api/login",apicfg.handleUserLogin)
    mux.HandleFunc("PUT /api/users",apicfg.handleUserPut)
    
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


