package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

type DB struct {
    path string
    mux *sync.RWMutex
}

type DBstructure struct {
    Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
    Id int `json:"id"`
    Body string `json:"body"`
}

func NewDB(path string) (*DB, error) {
    db := DB{
        path: path,
        mux: &sync.RWMutex{},
    }
    err := db.ensureDB()



    return &db,err
    
}

func (db *DB) CreateChirp(body string,) (Chirp, error) {
    DBS, err :=  db.loadDB()
    if err != nil {
        return Chirp{},err
    }
    chirp := Chirp{
        Id: len(DBS.Chirps)+1,
        Body: body,
    }
    DBS.Chirps[len(DBS.Chirps)+1] = chirp
    err = db.writeDB(DBS)
    if err != nil {
        return Chirp{},err
    }
    return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
    dbs, err := db.loadDB()
    if err != nil {
        return []Chirp{}, err
    }
    chirps := make([]Chirp,0,len(dbs.Chirps))
    for _, chirp := range dbs.Chirps{
       chirps = append(chirps, chirp) 
    }
    return chirps, nil

}
func(db *DB) createDB() error {
    DBstructure:= DBstructure{
        Chirps: map[int]Chirp{},
    }
    return db.writeDB(DBstructure)
}

func (db *DB) ensureDB() error {
    _, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func(db *DB) loadDB() (DBstructure, error) {
    db.mux.RLock()
    defer db.mux.RUnlock()
    bytes, err := os.ReadFile(db.path)

    if err != nil {
        return DBstructure{}, err
    }
    var dbs DBstructure
    err = json.Unmarshal(bytes, &dbs)
    if err != nil {
        return DBstructure{}, err
    }
    return dbs, nil
}

func(db *DB) writeDB(dbs DBstructure) error {
    toSend, err := json.Marshal(dbs)
    if err != nil {
        log.Printf("Error Marshaling %s",err)
        return err
    }
    db.mux.Lock()
    defer db.mux.Unlock()

    err = os.WriteFile(db.path,toSend,0600)
    if err != nil {
        return err
    }
    
    return nil 
}

