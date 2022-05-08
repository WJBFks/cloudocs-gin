package common

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var validTime time.Duration = (7 * 24) // hours

var jwtKey = []byte("NEUcseDocs2203271556")

type Claims struct {
	Id string
	jwt.StandardClaims
}

func ReleaseToken(id string) (string, error) {
	expirationTime := time.Now().Add(validTime * time.Hour)
	claims := &Claims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "docs",
			Subject:   "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParaseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, err
	})

	return token, claims, err
}
