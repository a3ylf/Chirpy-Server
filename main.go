package main

import (
	"fmt"
	"log"
	"net/http"
)
type apiconfig struct{
    fileserverhits int
}



func main() {
    var apicfg apiconfig
    mux := http.NewServeMux()
    mux.Handle("/app/*",apicfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
    mux.HandleFunc("/healthz",handlerReadiness)
    mux.HandleFunc("/reset",apicfg.handlerReset)
    mux.HandleFunc("/metrics",apicfg.handlerMetrics)
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
func (cfg *apiconfig)handlerMetrics(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte(fmt.Sprintf("Hits: %d",cfg.fileserverhits)))
}

func (cfg *apiconfig) handlerReset(w http.ResponseWriter, r*http.Request){
    cfg.fileserverhits = 0
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hits is now 0"))
}

func handlerReadiness(w http.ResponseWriter, r *http.Request){
    w.Header().Add("Content-Type","text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte("OK"))
}
