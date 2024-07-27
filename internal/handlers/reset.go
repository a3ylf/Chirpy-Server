package handlers

import "net/http"

func (cfg *Apiconfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverhits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits is now 0"))
}
