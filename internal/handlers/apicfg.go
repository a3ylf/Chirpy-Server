package handlers

import (
	"errors"
	"github.com/a3ylf/web-servers/internal/database"
	"github.com/golang-jwt/jwt/v5"
)
type Apiconfig struct{
    fileserverhits int
    db *database.DB
    secret string
}

type User struct {
    Id int `json:"id"`
    Email string `json:"email"`
}
type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
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



