package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/a3ylf/web-servers/internal/auth"
	"github.com/a3ylf/web-servers/internal/database"
)

func cleaner(pog string, bad []string) string {
	splitted := strings.Split(pog, " ")
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
		toSend = fmt.Sprintf(toSend + " " + split)
	}
	return toSend

}
func (cfg *Apiconfig) HandlePostChirp(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetTokenBearer(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	subject, err := ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	author_id, err := cfg.db.GetAuthor(subject)

	type parameters struct {
		Body string `json:"body"`
	}
	params := parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
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

	var chirp = database.Chirp{}

	if author_id == 0 {
		chirp, err = cfg.db.CreateChirp(cleaned, cfg.current)
		cfg.db.CreateAuthor(subject, cfg.current)
		cfg.current++
	} else {
		chirp, err = cfg.db.CreateChirp(cleaned, author_id)
	}

	if err != nil {
		log.Printf("Couldn't chirp %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't chirp")
		return
	}

	respondWithJSON(w, 201, chirp)
}

func (cfg *Apiconfig) HandleGetChirps(w http.ResponseWriter, r *http.Request) {

	dbChirps, err := cfg.db.GetChirps()
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couln't chirp for you")
	}
	asc := true
	ascordesc := r.URL.Query().Get("sort")

	if ascordesc == "desc" {
		asc = false
	}

	s := r.URL.Query().Get("author_id")
	if s != "" {
		author_id, err := strconv.Atoi(s)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author id")
			return
		}

		chirps, err := cfg.db.GetChirpsByAuthor(author_id)

		if err != nil {
			respondWithError(w, http.StatusNotFound, "couldn't get chirps")
			return
		}
		if !asc {
			for i, j := 0, len(chirps)-1; i < j; i, j = i+1, j-1 {
				chirps[i], chirps[j] = chirps[j], chirps[i]
			}
		}
		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	chirps := []database.Chirp{}

	for _, dbChirp := range dbChirps {
		chirps = append(chirps, database.Chirp{
			Id:        dbChirp.Id,
			Body:      dbChirp.Body,
			Author_id: dbChirp.Author_id,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})
	if !asc {
		for i, j := 0, len(chirps)-1; i < j; i, j = i+1, j-1 {
			chirps[i], chirps[j] = chirps[j], chirps[i]
		}
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
func (cfg *Apiconfig) HandleGetChirp(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("ID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id")
		return
	}
	chirp, err := cfg.db.GetChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unchirpopable")
		return
	}
	respondWithJSON(w, http.StatusOK, chirp)
}
func (cfg *Apiconfig) HandleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetTokenBearer(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	subject, err := ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	author_id, err := cfg.db.GetAuthor(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get author_id")
		return
	}

	id, err := strconv.Atoi(r.PathValue("ID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id")
		return
	}
	chirp, err := cfg.db.GetChirp(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp")
		return
	}
	if chirp.Author_id == author_id {
		err = cfg.db.DeleteChirp(id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
			return
		}
		respondWithJSON(w, 204, "Chirp deleted sucessfully")
		return
	}
	respondWithError(w, 403, "Chirp could not be deleted")
}
