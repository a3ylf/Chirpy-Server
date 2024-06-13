package main

import (
	"log"
	"net/http"
)

func main() {
    mux := http.NewServeMux()
    mux.Handle("/app/",http.StripPrefix("/app",http.FileServer(http.Dir("."))))
    mux.HandleFunc("/healthz",handlerReadiness)

    server := &http.Server{
		Addr:    ":8080",  // Define the address and port to listen on
		Handler: mux,      // Use the new ServeMux as the server's handler
	} 
	log.Printf("Serving on port: 8080\n")
	log.Fatal(server.ListenAndServe())
} 
func handlerReadiness(w http.ResponseWriter, r *http.Request){
    w.Header().Add("Content-Type","text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte("OK"))
}
