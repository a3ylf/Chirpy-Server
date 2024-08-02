package handlers

import "github.com/a3ylf/web-servers/internal/database"

func Newcfg (db *database.DB, secret string)  (Apiconfig){
    return Apiconfig{
        fileserverhits: 0,
        db: db,
        secret: secret,
        current:1,
    }
}
