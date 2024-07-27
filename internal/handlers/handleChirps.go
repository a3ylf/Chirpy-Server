package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)
func cleaner(pog string, bad []string) string {
    splitted := strings.Split(pog," ")
    for _, badword := range bad {
        for k, split := range splitted {
            if strings.ToLower(split) == badword {
                splitted[k] = "****"
            }
        }
    }
    toSend := "" 
    for _, split := range splitted {
        if len(toSend) == 0 {
            toSend = fmt.Sprintf(split)
            continue
        }
        toSend = fmt.Sprintf(toSend+ " " + split )
    }
    return toSend
    
}
func (cfg *Apiconfig) HandlePostChirp(w http.ResponseWriter, r *http.Request) {
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

func (cfg *Apiconfig) HandleGetChirps(w http.ResponseWriter, r *http.Request) {

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
func (cfg *Apiconfig) HandleGetChirp(w http.ResponseWriter, r *http.Request) {
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

