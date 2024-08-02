package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/a3ylf/web-servers/internal/database"
	"github.com/a3ylf/web-servers/internal/handlers"
	"github.com/joho/godotenv"
)




func main() {
    
    dbg := flag.Bool("debug",false,"enable debug mode")
    flag.Parse()
    
    if *dbg {
        os.Remove("internal/database/database.json")
    }
    db, err := database.NewDB("internal/database/database.json")

    if err != nil {
        log.Println(err)
    }

     
    err = godotenv.Load()
    if err != nil {
        log.Println(err)
        return
    }

    jwtSecret := os.Getenv("JWT_SECRET")
    apikey := os.Getenv("POLKA_KEY")
    apicfg := handlers.Newcfg(db,jwtSecret,apikey) 



    mux := http.NewServeMux()
    mux.Handle("/app/*",apicfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
    mux.HandleFunc("/api/reset",apicfg.HandlerReset)

    mux.HandleFunc("GET /api/healthz",handlers.HandlerReadiness)
    mux.HandleFunc("GET /admin/metrics",apicfg.HandlerMetrics)
    mux.HandleFunc("GET /api/chirps",apicfg.HandleGetChirps)
    mux.HandleFunc("GET /api/chirps/{ID}",apicfg.HandleGetChirp)

    mux.HandleFunc("POST /api/chirps",apicfg.HandlePostChirp)
    mux.HandleFunc("POST /api/users",apicfg.HandleUserPost)
    mux.HandleFunc("POST /api/login",apicfg.HandleUserLogin)
    mux.HandleFunc("POST /api/refresh",apicfg.HandlerRefresh)
    mux.HandleFunc("POST /api/revoke",apicfg.HandlerRevoke)
    
    mux.HandleFunc("POST /api/polka/webhooks",apicfg.HandleRedChirpy)

    mux.HandleFunc("DELETE /api/chirps/{ID}",apicfg.HandleDeleteChirp)
    
    server := &http.Server{
		Addr:    ":8080",  	
		Handler: mux,      

	} 
	log.Printf("Serving on port: 8080\n")
	log.Fatal(server.ListenAndServe())
} 










