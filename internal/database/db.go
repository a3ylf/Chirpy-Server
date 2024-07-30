package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)
var ErrAlreadyExists = errors.New("already exists")
var ErrNotExist = errors.New("Do not exist")


type DB struct {
    path string
    mux *sync.RWMutex
}

type DBstructure struct {
    Chirps map[int]Chirp `json:"chirps"`
    Users map[int]User `json:"users"`
    Tokens map[string]RefreshToken `json:"refresh_tokens"`
}

type User struct {
    Id int `json:"id"`
    Email string `json:"email"`
    Password string `json:"password"`
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
func (db *DB) GetChirp(id int) (Chirp,error) {
    dbs,err := db.loadDB()
    if err != nil {
        return Chirp{},err
    }
    if dbs.Chirps[id].Id == 0 {
        return Chirp{}, errors.New("Unchirporpable")
    }
    return dbs.Chirps[id], nil
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


func (db *DB) CreateUser(email, hashedPassword string) (User, error) {
	 _, err := db.GetUserByEmail(email) 

	if err != ErrNotExist {
	    return User{}, err
    }

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		Id:             id,
		Email:          email,
		Password: hashedPassword,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, errors.New("Not exist")
	}

	return user, nil
}
func (db *DB) UpdateUser(id int, email, hashedPassword string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, errors.New("user do not exist")	
	}

	user.Email = email
	user.Password  = hashedPassword
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
    for _, user := range dbStructure.Users{
        if user.Email == email {
            return user,nil
        }
    }

	return User{}, ErrNotExist
}
func(db *DB) createDB() error {
    DBstructure:= DBstructure{
        Chirps: map[int]Chirp{},
        Users: map[int]User{},
        Tokens: map[string]RefreshToken{},
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

