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
    Authors map[string]int `json:"authors"`
}

type User struct {
    Id int `json:"id"`
    Email string `json:"email"`
    Password string `json:"password"`
    IsChirpyRed bool `json:"is_chirpy_red"`
}
type Chirp struct {
    Id int `json:"id"`
    Author_id int `json:"author_id"`
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

func (db *DB) GetAuthor(token string) (int,error) {
    dbs,err := db.loadDB()
    if err != nil {
        return 0, err

    }
    return dbs.Authors[token], nil
}
func (db *DB) CreateAuthor(token string,author int) (error) {
    dbs,err := db.loadDB()
    if err != nil {
        return err

    }
    dbs.Authors[token] = author
    err = db.writeDB(dbs)
    if err != nil {
        return err
    }

    return  nil
}
func (db *DB) CreateChirp(body string, current int) (Chirp, error) {
    DBS, err :=  db.loadDB()
    if err != nil {
        return Chirp{},err
    }
    chirp := Chirp{
        Id: len(DBS.Chirps)+1,
        Author_id: current,
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
func (db *DB) GetChirpsByAuthor(author int) ([]Chirp, error) {
    dbs, err := db.loadDB()
    if err != nil {
        return []Chirp{}, err
    }
    chirps := make([]Chirp,0,len(dbs.Chirps))
    for _, chirp := range dbs.Chirps{
        if chirp.Author_id == author{
       chirps = append(chirps, chirp) 
        }
    }
    return chirps, nil
}


func (db *DB) DeleteChirp(id int ) error {
    dbs, err := db.loadDB()
    if err != nil {
        return err
    }
    delete(dbs.Chirps,id)
    err = db.writeDB(dbs)
    if err != nil {
        return err
    }
    return nil
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
		IsChirpyRed: false,
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
func (db *DB) UpdateUser(id int, email, hashedPassword string, isRed ... bool) (User, error) {
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
	if len(isRed) > 0 {
        user.IsChirpyRed = isRed[0]
    }
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
        Authors: map[string]int{},
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

