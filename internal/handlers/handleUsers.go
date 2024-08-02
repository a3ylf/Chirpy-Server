package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/a3ylf/web-servers/internal/auth"
	"github.com/a3ylf/web-servers/internal/database"
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
        return
    }
   
	user, err := cfg.db.CreateUser(params.Email,string(psw)) 
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			respondWithError(w, http.StatusConflict, "User already exists")
			return
		}
        log.Printf(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, 201, user)
}
func (cfg *Apiconfig) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
	    database.User 	    
	    Token string `json:"token"`
	    RefreshToken string `json:"refresh_token"`
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

    
    claims := jwt.RegisteredClaims{
        Issuer: "chirpy",
        IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
        ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
        Subject: fmt.Sprintf("%d",user.Id),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

    assignedtoken, err := token.SignedString([]byte(cfg.secret))

    if err != nil {
        respondWithError(w,http.StatusInternalServerError,"Couldn't create JMT")
    }
    
    rndbyte := make([]byte,32)

    _, err = rand.Read([]byte(cfg.secret))

    if err != nil {
        respondWithError(w,http.StatusInternalServerError,"Couldn't Make refresh token")
    }

    refreshtoken := hex.EncodeToString(rndbyte)

    err = cfg.db.SaveRFtokens(user.Id, refreshtoken)
    if err != nil {
        respondWithError(w,http.StatusInternalServerError,"couldn't save rf token")
    }


	respondWithJSON(w, 200, response{
    User:  user,
	Token: assignedtoken,
	RefreshToken: refreshtoken,    
	}) 

}
func (cfg *Apiconfig) HandleUserPut(w http.ResponseWriter, r *http.Request) {
    
    token, err := auth.GetTokenBearer(r.Header)
    if err != nil {
        respondWithError(w,http.StatusUnauthorized,err.Error())
    }
    
    subject , err := ValidateJWT(token,cfg.secret)
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

	respondWithJSON(w, 200, user)
	
}

