package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *Apiconfig) HandleUserPost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
    psw, err := 	bcrypt.GenerateFromPassword([]byte(params.Password),10)
    
    if err!= nil {
        respondWithError(w,http.StatusInternalServerError,"Couldn't hash it")
    }
   
	user, err := cfg.db.CreateUser(params.Email,string(psw))
	if err != nil {
		log.Printf("Couldn't define user %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't define user")
		return
	}

	respondWithJSON(w, 201, User{
		Id:   user.Id,
		Email: user.Email,
	})
}
func (cfg *Apiconfig) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
		Expires int `json:"expires_in_seconds"`
	}

	type response struct {
	    User
	    Token string `json:"token"`
    }
	params := parameters{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
    
	user, err := cfg.db.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't Find this user")
		return
	}
    err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(params.Password))
    if err != nil {
        respondWithError(w,401,"Password incorrect!")
        return
    }

    plus := 60*60*24
    if params.Expires != 0 {
        plus = params.Expires 
    }
    claims := jwt.RegisteredClaims{
        Issuer: "chirpy",
        IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
        ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(plus)*time.Second)),
        Subject: fmt.Sprintf("%d",user.Id),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

    assignedtoken, err := token.SignedString([]byte(cfg.secret))

    if err != nil {
        respondWithError(w,http.StatusInternalServerError,"Couldn't create JMT")
    }


	respondWithJSON(w, 200, response{
        User: User{
		Id:   user.Id,
		Email: user.Email,
	},
	Token: assignedtoken,
	    }) 

}
func (cfg *Apiconfig) HandleUserPut(w http.ResponseWriter, r *http.Request) {
    tkn := r.Header.Get("Authorization")
    if  tkn == "" {
        respondWithError(w,http.StatusUnauthorized,"Couldn't find token")
    }
   
    splitAuth := strings.Split(tkn, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		respondWithError(w,http.StatusUnauthorized,"malformed authorization header")
	}
    
    subject , err := ValidateJWT(splitAuth[1],cfg.secret)

    if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}


	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}

	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
    psw, err := bcrypt.GenerateFromPassword([]byte(params.Password),10)
    
    if err!= nil {
        respondWithError(w,http.StatusInternalServerError,"Couldn't hash it")
    }

    userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}
   
	user, err := cfg.db.UpdateUser(userIDInt,params.Email,string(psw))
	if err != nil {
		log.Printf("Couldn't define user %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't define user")
		return
	}

	respondWithJSON(w, 200, User{
		Id:   user.Id,
		Email: user.Email,
	})
}

