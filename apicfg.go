package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"

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
		Expires int `json:"expires_in_seconds"`
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

