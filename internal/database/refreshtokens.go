package database

import "time"

type RefreshToken struct {
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (db *DB) SaveRFtokens(id int, token string) error {
	dbs, err := db.loadDB()
	if err != nil {
		return err
	}
	dbs.Tokens[token] = RefreshToken{
		UserID:    id,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	}
	err = db.writeDB(dbs)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) RevokeRefreshToken(token string) error {
	dbs, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbs.Tokens, token)

	err = db.writeDB(dbs)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UserForRefreshToken(token string) (User, error) {
    dbs, err := db.loadDB()
    if err != nil {
        return User{},err
    }

    rt, ok := dbs.Tokens[token]

    if !ok {
        return User{}, errnotexist
    }

    if rt.ExpiresAt.Before(time.Now()) {
        return User{}, errnotexist
    }

    user, err := db.GetUser(rt.UserID)
    if err != nil {
        return User{},err
    }

    return user,nil 
   
}
