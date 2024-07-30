package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included")

func GetTokenBearer (h http.Header) (string, error) {
    authHeader := h.Get("Authorization")
    if authHeader == "" {
        return "", ErrNoAuthHeaderIncluded
    }

    splitAuth := strings.Split(authHeader," ")
    if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
        return "", errors.New("wrongly formed header")

    }

    return splitAuth[1],nil


    
}
func MakeJWT(
	userID int,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {
	signingKey := []byte(tokenSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   fmt.Sprintf("%d", userID),
})	
    return token.SignedString(signingKey)
}
