package main

import (
	"github.com/golang-jwt/jwt/v5"
)

var secretKey string = "key"

func GenJWT(user *UserToken) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)

	tokenString, err := token.SignedString([]byte(secretKey))
	return tokenString, err
}

func ParseJWT(tokenString string) (*UserToken, error) {

	rawToken, err := jwt.ParseWithClaims(tokenString, &UserToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !rawToken.Valid {
		return nil, err
	}

	token, ok := rawToken.Claims.(*UserToken)
	if !ok {
		return nil, err
	}
	return token, nil

}
