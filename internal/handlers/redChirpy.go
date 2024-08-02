package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/a3ylf/web-servers/internal/auth"
)

func (cfg *Apiconfig) HandleRedChirpy(w http.ResponseWriter, r *http.Request) {

    apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find api key")
		return
	}
	if apiKey != cfg.key {
		respondWithError(w, http.StatusUnauthorized, "API key is invalid")
		return
	}

	type data struct {
		User_id int `json:"user_id"`
	}

	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}
	params := parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	
	if params.Event == "user.upgraded" {
		user, err := cfg.db.GetUser(params.Data.User_id)
		if err != nil {
			respondWithError(w, 404, "couldn't get user")
			return
		}
		user, err = cfg.db.UpdateUser(user.Id, user.Email, user.Password, true)
		
	}

    w.WriteHeader(http.StatusNoContent)
}
