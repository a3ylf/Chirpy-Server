package handlers

import (
	"net/http"
	"time"

	"github.com/a3ylf/web-servers/internal/auth"
)

func (cfg *Apiconfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
    
    type response struct {
        Token string `json:"token"`
    }

    rt, err := auth.GetTokenBearer(r.Header)

    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Couldn't find token")
        return
    }

    user, err := cfg.db.UserForRefreshToken(rt)
    if err != nil {
        respondWithError(w,http.StatusUnauthorized,"Couldn't get user for refresh token")
        return
    }
    accesToken, err := auth.MakeJWT(
        user.Id,
        cfg.secret,
        time.Hour,
        )
    if err != nil {
        respondWithError(w,http.StatusUnauthorized,"Coulndn't validade token")
    }

    respondWithJSON(w,http.StatusOK,response{
        Token: accesToken,
    })  
}


func (cfg *Apiconfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {

    rt, err := auth.GetTokenBearer(r.Header)

    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Couldn't find token")
        return
    }
    err = cfg.db.RevokeRefreshToken(rt)
    if err != nil {
        respondWithError(w,http.StatusInternalServerError,"Couldn't revoke session")
        return
    }
    
    w.WriteHeader(http.StatusNoContent)

}


