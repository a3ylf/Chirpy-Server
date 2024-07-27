package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
    Id int `json:"id"`
    Email string `json:"email"`
}
type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func (cfg *apiconfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverhits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits is now 0"))
}

func (cfg *apiconfig) handlePost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	params := parameters{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	bad_words := []string{"kerfuffle", "sharbert", "fornax"}

	cleaned := cleaner(params.Body, bad_words)

	chirp, err := cfg.db.CreateChirp(cleaned)
	if err != nil {
		log.Printf("Couldn't chirp %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't chirp")
		return
	}

	respondWithJSON(w, 201, Chirp{
		Id:   chirp.Id,
		Body: chirp.Body,
	})
}
func (cfg *apiconfig) handleUserPost(w http.ResponseWriter, r *http.Request) {
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
func (cfg *apiconfig) handleUserLogin(w http.ResponseWriter, r *http.Request) {
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
func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return "", err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != string("chirpy") {
		return "", errors.New("invalid issuer")
	}

	return userIDString, nil
}
func (cfg *apiconfig) handleUserPut(w http.ResponseWriter, r *http.Request) {
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


func (cfg *apiconfig) handleGet(w http.ResponseWriter, r *http.Request) {

	dbChirps, err := cfg.db.GetChirps()
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couln't chirp for you")
	}
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			Id:   dbChirp.Id,
			Body: dbChirp.Body,
		})
	}


	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	respondWithJSON(w, http.StatusOK, chirps)
}
func (cfg *apiconfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("ID"))
    if err != nil {
        respondWithError(w,http.StatusBadRequest,"Invalid chirp id")
        return
    }
    dbChirp, err := cfg.db.GetChirp(id)
    if err != nil {
        respondWithError(w,http.StatusNotFound,"Unchirpopable")
        return
    }
    respondWithJSON(w,http.StatusOK,Chirp{
        Id: dbChirp.Id,
        Body: dbChirp.Body,
    })
}

