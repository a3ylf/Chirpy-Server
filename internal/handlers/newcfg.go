package handlers

import "github.com/a3ylf/web-servers/internal/database"

func Newcfg (db *database.DB, secret,key string)  (Apiconfig){
    return Apiconfig{
        fileserverhits: 0,
        db: db,
        secret: secret,
        key: key,
        current:1,
    }
}
