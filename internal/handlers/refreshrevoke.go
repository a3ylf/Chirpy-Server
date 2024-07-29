package handlers

import "net/http"

func (cfg *Apiconfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
    
    type response struct {
        Token string `json:"token"`
    }

    


}


