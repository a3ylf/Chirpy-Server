package main

import (
	"log"
	"net/http"
)

func main() {
    mux := http.NewServeMux()
    mux.Handle("/assets/logo.png",http.FileServer(http.Dir(".")))
    server := &http.Server{
		Addr:    ":8080",  // Define the address and port to listen on
		Handler: mux,      // Use the new ServeMux as the server's handler
	} 
	log.Printf("Serving on port: 8080\n")
	log.Fatal(server.ListenAndServe())
} 
